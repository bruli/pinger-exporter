package app

import (
	"context"

	"github.com/rs/zerolog"
)

type CommandHandlerMiddleware func(h CommandHandler) CommandHandler

type CommandHandlerFunc func(ctx context.Context, cmd Command) error

func (c CommandHandlerFunc) Handle(ctx context.Context, cmd Command) error {
	return c(ctx, cmd)
}

func NewLogCommandHandlerMiddleware(log *zerolog.Logger) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) error {
			err := h.Handle(ctx, cmd)
			if err != nil {
				log.Err(err).Msgf("failed to handle command %q", cmd.Name())
			}
			return err
		})
	}
}
