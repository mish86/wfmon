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

type Shutdowner interface {
	Shutdown() error
}

type Serv interface {
	Starter
	Closer
	Configer
	Shutdowner
}
