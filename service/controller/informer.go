package controller

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		artificialObj := corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "non-existing",
				Namespace: "non-existing",
				SelfLink:  "/v1/namespace/non-existing/services/non-existing",
			},
		}

		e := watch.Event{
			Type: watch.Added,

			Object: artificialObj,
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
