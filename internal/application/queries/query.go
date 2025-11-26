package queries

import "context"

type Query interface {
	Validate() error
}

type QueryHandler[TQuery Query, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

type QueryBus interface {
	Execute(ctx context.Context, query Query) (interface{}, error)
}
