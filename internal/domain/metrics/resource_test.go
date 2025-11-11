package metrics_test

import (
	"testing"
	"time"

	"github.com/bruli/pinger-exporter/internal/domain/metrics"
	"github.com/stretchr/testify/require"
)

func TestNewResource(t *testing.T) {
	type args struct {
		name, status string
		seconds      float64
		createdAt    time.Time
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name:        "with an invalid name, then it returns an invalid resource name error",
			args:        args{},
			expectedErr: metrics.ErrInvalidResourceName,
		},
		{
			name: "with an invalid status, then it returns an invalid resource status error",
			args: args{
				name: "resource",
			},
			expectedErr: metrics.ErrInvalidResourceStatus,
		},
		{
			name: "with an invalid created at, then it returns an invalid created time error",
			args: args{
				name:      "resource",
				status:    "ok",
				seconds:   10,
				createdAt: time.Time{},
			},
			expectedErr: metrics.ErrInvalidCreatedTime,
		},
		{
			name: "with valid data, then it returns a valid struct",
			args: args{
				name:      "resource",
				status:    "ok",
				seconds:   10,
				createdAt: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a NewResource constructor,
		when it's called '`+tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := metrics.NewResource(tt.args.name, tt.args.status, tt.args.seconds, tt.args.createdAt)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.Equal(t, tt.args.name, got.Name())
			require.Equal(t, tt.args.seconds, got.Seconds())
			require.Equal(t, tt.args.createdAt, got.CreatedAt())
		})
	}
}
