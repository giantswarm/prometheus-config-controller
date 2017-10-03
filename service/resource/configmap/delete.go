package configmap

import (
	"context"
)

// GetDeleteState is a no-op.
// We do not want to delete the configmap, as the running prometheus relies on it.
func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	return nil, nil
}

// ProcessDeleteState is a no-op.
func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	return nil
}
