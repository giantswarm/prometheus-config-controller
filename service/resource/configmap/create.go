package configmap

import (
	"context"
)

// GetCreateState is a no-op.
// The controller is given a configuration, which it then controls,
// we do not create a prometheus configuration from scratch.
func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	return nil, nil
}

// ProcessCreateState is a no-op.
func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	return nil
}
