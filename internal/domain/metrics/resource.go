package metrics

import (
	"errors"
	"time"
)

var (
	ErrInvalidResourceName = errors.New("resource name cannot be empty")
	ErrInvalidCreatedTime  = errors.New("resource created at cannot be zero")
)

type Resource struct {
	name      string
	seconds   float64
	createdAt time.Time
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
	case r.createdAt.IsZero():
		return ErrInvalidCreatedTime
	default:
		return nil
	}
}

func NewResource(name string, seconds float64, createdAt time.Time) (*Resource, error) {
	r := Resource{name: name, seconds: seconds, createdAt: createdAt}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return &r, nil
}
