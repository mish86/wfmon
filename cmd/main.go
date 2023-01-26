package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"wfmon/pkg"
	log "wfmon/pkg/logger"
	"wfmon/pkg/network"
	"wfmon/pkg/radio"
	"wfmon/pkg/serv"
	"wfmon/pkg/wifi"

	"go.uber.org/zap"
)

const (
	envMode          = "MODE"
	defaultGSTimeout = time.Second * 15
)

type Application struct {
	log *zap.SugaredLogger

	servs []serv.Serv

	mode          pkg.Mode
	gsTimeout     time.Duration
	chHopInterval time.Duration
	ifaceName     string
	iface         *net.Interface
}

func loadApplication() *Application {
	var err error

	app := &Application{
		ifaceName: "en0",
	}

	app.mode = pkg.FromString(os.Getenv(envMode))

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
	if app.iface, err = network.FindIFaceByName(strings.TrimSpace(app.ifaceName)); err != nil {
		return err
	}

	return nil
}

func (app *Application) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), app.gsTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(app.servs))

	for _, shutdowner := range app.servs {
		go func(shutdowner serv.Shutdowner) {
			defer wg.Done()

			if e := shutdowner.Shutdown(); e != nil {
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

func (app *Application) startWithGS(ctx context.Context) {
	c := make(chan os.Signal, 1)
	// graceful shutdown
	// when SIGINT (Ctrl+C)
	// when SIGTERM (Ctrl+/)
	// except SIGKILL, SIGQUIT will not be caught
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// run non-blocking service
	app.start(ctx)

	// block until we receive our signal.
	sig := <-c

	log.Info("received %s", sig)
	log.Info("shutting down")

	app.shutdown()
}

func (app *Application) start(ctx context.Context) {
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

	// setup services
	app.servs = []serv.Serv{
		mon, hopper,
	}

	// run configurations
	for _, configer := range app.servs {
		serv := configer
		if err := serv.Configure(); err != nil {
			log.Fatal(err)
		}
	}

	// run services
	for _, starter := range app.servs {
		serv := starter
		go func() {
			if err := serv.Start(ctx); err != nil {
				log.Fatal(err)
			}

			// close handlers
			defer mon.Close()
		}()
	}
}

func main() {
	app := loadApplication()

	app.initLogger()
	defer app.closeLogger()

	log.Info("🚀 starting")

	ctx := context.Background()
	app.startWithGS(ctx)

	// const liveTime = 10 * time.Second
	// ctx, stop := context.WithCancel(ctx)
	// app.start(ctx)
	// liveTimer := time.NewTimer(liveTime)
	// <-liveTimer.C
	// stop()
}
