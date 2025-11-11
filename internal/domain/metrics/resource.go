package metrics

import (
	"errors"
	"time"
)

var (
	ErrInvalidResourceName   = errors.New("resource name cannot be empty")
	ErrInvalidResourceStatus = errors.New("resource status cannot be empty")
	ErrInvalidCreatedTime    = errors.New("resource created at cannot be zero")
)

type Resource struct {
	name, status string
	seconds      float64
	createdAt    time.Time
}

func (r Resource) Status() string {
	return r.status
}

func (r Resource) Name() string {
	return r.name
}

func (r Resource) Seconds() float64 {
	return r.seconds
}

func (r Resource) CreatedAt() time.Time {
	return r.createdAt
}

func (r Resource) validate() error {
	switch {
	case r.name == "":
		return ErrInvalidResourceName
	case r.status == "":
		return ErrInvalidResourceStatus
	case r.createdAt.IsZero():
		return ErrInvalidCreatedTime
	default:
		return nil
	}
}

func NewResource(name, status string, seconds float64, createdAt time.Time) (*Resource, error) {
	r := Resource{name: name, status: status, seconds: seconds, createdAt: createdAt}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return &r, nil
}
