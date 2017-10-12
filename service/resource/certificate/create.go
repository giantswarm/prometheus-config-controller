package certificate

import (
	"context"
)

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	return nil, nil
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	return nil
}
