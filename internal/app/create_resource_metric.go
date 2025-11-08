package app

import (
	"context"

	"github.com/bruli/pinger-exporter/internal/domain/metrics"
)

const CreateResourceMetricCommandName = "CreateResourceMetricCommand"

type CreateResourceMetricCommand struct {
	Resource *metrics.Resource
}

func (c CreateResourceMetricCommand) Name() string {
	return CreateResourceMetricCommandName
}

type CreateResourceMetric struct {
	repo ResourceMetricRepository
}

func (c CreateResourceMetric) Handle(ctx context.Context, cmd Command) error {
	co, ok := cmd.(CreateResourceMetricCommand)
	if !ok {
		return NewInvalidCommandError(CreateResourceMetricCommandName, cmd.Name())
	}
	c.repo.Create(ctx, co.Resource)
	return nil
}

func NewCreateResourceMetric(repo ResourceMetricRepository) *CreateResourceMetric {
	return &CreateResourceMetric{repo: repo}
}
