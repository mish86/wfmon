package repeater

import (
	"context"
	"time"
)

type Repeater interface {
	Start(ctx context.Context, interval time.Duration, onTimer func())
}

type Fnc func(ctx context.Context, interval time.Duration, onTimer func())

func (fnc Fnc) Start(ctx context.Context, interval time.Duration, onTimer func()) {
	fnc(ctx, interval, onTimer)
}

func Default(ctx context.Context, interval time.Duration, onTimer func(), onDone func()) {
	retryTimer := time.NewTimer(interval)
	defer retryTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			onDone()
			return

		case <-retryTimer.C:
			onTimer()

			// reset timer
			retryTimer.Reset(interval)
		}
	}
}
