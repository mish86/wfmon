package main

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	mode "wfmon/pkg"
	log "wfmon/pkg/logger"
	"wfmon/pkg/radio"
	"wfmon/pkg/serv"
	"wfmon/pkg/view/wifitable"
	"wfmon/pkg/wifi"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

const (
	envMode          = "MODE"
	defaultGSTimeout = time.Second * 15
)

type Application struct {
	log *zap.SugaredLogger

	servs       []serv.Serv
	starters    []serv.Starter
	shutdowners []serv.Shutdowner
	program     *tea.Program

	mode          mode.Mode
	gsTimeout     time.Duration
	chHopInterval time.Duration
	ifaceName     string
	iface         *net.Interface
}

func loadApplication() *Application {
	var err error

	app := &Application{
		// TODO
		ifaceName: strings.TrimSpace("en0"),
	}

	app.mode = mode.FromString(os.Getenv(envMode))

	if app.gsTimeout, err = time.ParseDuration(os.Getenv("GRACEFUL_SHUTDOWN_TIMEOUT")); err != nil {
		app.gsTimeout = defaultGSTimeout
	}

	if app.chHopInterval, err = time.ParseDuration(os.Getenv("CHANNEL_HOP_INTERVAL")); err != nil {
		app.chHopInterval = radio.DefaultHopInterval
	}

	return app
}

func (app *Application) initLogger() {
	logger := log.NewLogger(app.mode)

	app.log = logger.Sugar()

	zap.ReplaceGlobals(logger)
}

func (app *Application) closeLogger() {
	err := app.log.Sync()
	// https://github.com/uber-go/zap/issues/991#issuecomment-962098428
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		panic(err)
	}
}

func (app *Application) findInterface() error {
	var err error

	if app.iface, err = net.InterfaceByName(app.ifaceName); err != nil {
		return err
	}

	return nil
}

// Stops and closes services within @app.gsTimeout timeout.
func (app *Application) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), app.gsTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(app.servs))

	for _, shutdowner := range app.shutdowners {
		go func(shutdowner serv.Shutdowner) {
			// indicate completion of service shutdown
			defer wg.Done()
			// close handles
			defer shutdowner.Close()

			// stop processing
			if e := shutdowner.Stop(); e != nil {
				log.Error(e)
			}
		}(shutdowner)
	}

	go func() {
		wg.Wait()
		// all services completed their work gracefully
		cancel()
	}()

	// block until all services completed their work or timeout
	<-ctx.Done()

	if e := ctx.Err(); e != nil && !errors.Is(e, context.Canceled) {
		log.Error(e)
	}
}

// Creates and initializes services and tea program.
func (app *Application) init(ctx context.Context) {
	// find interface by name
	if err := app.findInterface(); err != nil {
		log.Fatal(err)
	}

	// create wifi monitor
	mon := wifi.NewMonitor(&wifi.Config{
		IFace: app.iface,
	})

	// create channel hopper
	hopper := radio.NewChannelHopperServ(&radio.ChannelHopperConfig{
		IFace:       app.iface,
		HopInterval: app.chHopInterval,
	})

	// create table controller
	table := wifitable.NewTableCtrl(mon.GetFrames())

	// setup services
	app.servs = []serv.Serv{mon, hopper}
	app.starters = []serv.Starter{mon, hopper, table}
	app.shutdowners = []serv.Shutdowner{mon, hopper}

	// run configurations
	for _, configer := range app.servs {
		serv := configer
		if err := serv.Configure(); err != nil {
			log.Fatal(err)
		}
	}

	// create tea program
	app.program = tea.NewProgram(
		table,
		tea.WithContext(ctx),
		tea.WithInputTTY(),
	)
}

// Runs services in seprate goroutings and blocks main with tea program.
func (app *Application) start(ctx context.Context) {
	// run services
	for _, starter := range app.starters {
		go func(starter serv.Starter) {
			if err := starter.Start(ctx); err != nil {
				log.Fatal(err)
			}
		}(starter)
	}

	// Blocks application execution until SIGINT (Ctrl+C) and SIGTERM (Ctrl+/)
	if _, err := app.program.Run(); err != nil {
		log.Fatal(err)
	}

	log.Info("shutting down")
	app.shutdown()
}

func main() {
	app := loadApplication()

	app.initLogger()
	defer app.closeLogger()

	log.Info("🚀 starting")
	log.Debugf("app mode %s", app.mode)

	ctx := context.Background()
	app.init(ctx)
	app.start(ctx)
}
