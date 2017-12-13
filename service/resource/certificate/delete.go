package certificate

import (
	"context"

	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	return nil
}
