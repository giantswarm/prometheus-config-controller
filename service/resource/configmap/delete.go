package configmap

import (
	"context"

	"github.com/giantswarm/operatorkit/framework"
)

// NewDeletePatch is a no-op.
// We do not want to delete the configmap, as the running prometheus relies on it.
func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

// ApplyDeleteChange is a no-op.
func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	return nil
}
