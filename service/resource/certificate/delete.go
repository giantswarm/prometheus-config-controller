package certificate

import (
	"context"
)

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	return nil, nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	return nil
}
