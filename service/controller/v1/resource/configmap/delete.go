package configmap

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"
)

// NewDeletePatch calls NewUpdatePatch as the ConfigMap must be updated with
// removed cluster. This is important when the last cluster in the
// installation is removed.
func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	p, err := r.NewUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return p, nil
}

// ApplyDeleteChange calls ApplyUpdateChange as the ConfigMap must be updated
// with removed cluster. This is important when the last cluster in the
// installation is removed.
func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	err := r.ApplyUpdateChange(ctx, obj, deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
