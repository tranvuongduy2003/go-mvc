package queries

import "context"

// Query represents a read operation that doesn't modify state
// Following CQRS pattern - queries should only read data
type Query interface {
	// Validate performs query-level validation
	Validate() error
}

// QueryHandler handles a specific query
// Generic interface for all query handlers following CQRS pattern
type QueryHandler[TQuery Query, TResult any] interface {
	// Handle executes the query and returns the result
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

// QueryBus dispatches queries to their handlers
// Provides a centralized way to execute queries
type QueryBus interface {
	// Execute dispatches a query to its handler
	Execute(ctx context.Context, query Query) (interface{}, error)
}
