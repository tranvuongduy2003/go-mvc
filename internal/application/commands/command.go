package commands

import "context"

type Command interface {
	Validate() error
}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

type CommandBus interface {
	Execute(ctx context.Context, cmd Command) (interface{}, error)
}
