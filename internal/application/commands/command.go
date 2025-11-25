package commands

import "context"

// Command represents a write operation that modifies state
// Following CQRS pattern - commands should not return data (except identifiers)
type Command interface {
	// Validate performs command-level validation
	Validate() error
}

// CommandHandler handles a specific command
// Generic interface for all command handlers following CQRS pattern
type CommandHandler[TCommand Command, TResult any] interface {
	// Handle executes the command and returns a result
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

// CommandBus dispatches commands to their handlers
// Provides a centralized way to execute commands
type CommandBus interface {
	// Execute dispatches a command to its handler
	Execute(ctx context.Context, cmd Command) (interface{}, error)
}
