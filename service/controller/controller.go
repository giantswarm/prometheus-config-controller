package controller

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

const (
	// serviceLabelSelector is the label selector to match master services.
	serviceLabelSelector = "app=master"
)

type Config struct {
	BackOff           backoff.BackOff
	K8sClient         kubernetes.Interface
	Logger            micrologger.Logger
	OperatorFramework *framework.Framework

	ResyncPeriod time.Duration
}

func DefaultConfig() Config {
	return Config{
		BackOff:           nil,
		K8sClient:         nil,
		Logger:            nil,
		OperatorFramework: nil,

		ResyncPeriod: time.Duration(0),
	}
}

type Controller struct {
	backOff           backoff.BackOff
	k8sClient         kubernetes.Interface
	logger            micrologger.Logger
	operatorFramework *framework.Framework

	bootOnce     sync.Once
	mutex        sync.Mutex
	resyncPeriod time.Duration
}

func New(config Config) (*Controller, error) {
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.OperatorFramework == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.OperatorFramework must not be empty")
	}

	if config.ResyncPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResyncPeriod must not be zero")
	}

	newController := &Controller{
		backOff:           config.BackOff,
		k8sClient:         config.K8sClient,
		logger:            config.Logger,
		operatorFramework: config.OperatorFramework,

		bootOnce:     sync.Once{},
		mutex:        sync.Mutex{},
		resyncPeriod: config.ResyncPeriod,
	}

	return newController, nil
}

func (o *Controller) Boot() {
	o.bootOnce.Do(func() {
		operation := func() error {
			err := o.bootWithError()
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}

		notifier := func(err error, d time.Duration) {
			o.logger.Log("warning", fmt.Sprintf("retrying controller boot due to error: %#v", microerror.Mask(err)))
		}

		err := backoff.RetryNotify(operation, o.backOff, notifier)
		if err != nil {
			o.logger.Log("error", fmt.Sprintf("stop controller boot retries due to too many errors: %#v", microerror.Mask(err)))
			os.Exit(1)
		}
	})
}

func (o *Controller) bootWithError() error {
	o.logger.Log("debug", "starting list/watch")

	newResourceEventHandler := o.operatorFramework.NewCacheResourceEventHandler()

	listWatch := &cache.ListWatch{
		ListFunc: func(options apismetav1.ListOptions) (runtime.Object, error) {
			o.logger.Log("debug", "listing all services", "event", "list")
			options.LabelSelector = serviceLabelSelector
			return o.k8sClient.CoreV1().Services("").List(options)
		},
		WatchFunc: func(options apismetav1.ListOptions) (watch.Interface, error) {
			o.logger.Log("debug", "watching all services", "event", "watch")
			options.LabelSelector = serviceLabelSelector
			return o.k8sClient.CoreV1().Services("").Watch(options)
		},
	}

	_, informer := cache.NewInformer(
		listWatch,
		&v1.Service{},
		o.resyncPeriod,
		newResourceEventHandler,
	)
	informer.Run(nil)

	return nil
}
