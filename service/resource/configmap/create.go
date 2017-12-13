package configmap

import (
	"context"
)

// ApplyCreateChange is a no-op.
func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	return nil
}
