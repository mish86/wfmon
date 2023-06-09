package serv

import (
	"context"
)

type Starter interface {
	Start(context.Context) error
}

type Closer interface {
	Close()
}

type Configer interface {
	Configure() error
}

type Stopper interface {
	Stop() error
}

type Shutdowner interface {
	Stopper
	Closer
}

type Serv interface {
	Starter
	Closer
	Configer
	Stopper
}
