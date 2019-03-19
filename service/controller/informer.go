package controller

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/watch"
)

type artificialInformerConfig struct {
	ResyncPeriod time.Duration
}

// artificialInformer creates an update event with nil object on every resync
// period. It is useful when we want to use operatorkit reconciliation but not
// based on watching any object.
type artificialInformer struct {
	resyncPeriod time.Duration

	updateCh chan watch.Event
}

func newArtificialInformer(config artificialInformerConfig) (*artificialInformer, error) {
	if config.ResyncPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResyncPeriod must not be empty", config)
	}

	i := &artificialInformer{
		resyncPeriod: config.ResyncPeriod,

		updateCh: make(chan watch.Event),
	}

	return i, nil
}

func (i *artificialInformer) Boot(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return nil
		}

		e := watch.Event{
			Type: watch.Added,
		}
		i.updateCh <- e

		time.Sleep(i.resyncPeriod)
	}
}

func (i *artificialInformer) ResyncPeriod() time.Duration {
	return i.resyncPeriod
}

func (i *artificialInformer) Watch(ctx context.Context) (chan watch.Event, chan watch.Event, chan error) {
	return make(chan watch.Event), i.updateCh, make(chan error)
}
