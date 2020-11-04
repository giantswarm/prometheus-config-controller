package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	Name = "configmapv1"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	CertDirectory string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string

	Provider string
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	certDirectory      string
	configMapKey       string
	configMapName      string
	configMapNamespace string
	provider           string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.CertDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertDirectory must not be empty")
	}
	if config.ConfigMapKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapKey must not be empty")
	}
	if config.ConfigMapName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapName must not be empty")
	}
	if config.ConfigMapNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapNamespace must not be empty")
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Provider must not be empty")
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		certDirectory:      config.CertDirectory,
		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,

		provider: config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ensure(ctx context.Context, obj interface{}) error {
	currentCM, err := r.getCurrentState(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	desiredCM, err := r.getDesiredState(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	cm := newConfigMapToUpdate(currentCM, desiredCM)

	{
		// currentCM is used in the log messages because cm can be nil.
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating ConfigMap %#q in namespace %#q", currentCM.GetName(), currentCM.GetNamespace()))

		if cm == nil {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ConfigMap %#q in namespace %#q is up to date", currentCM.GetName(), currentCM.GetNamespace()))

			r.logger.LogCtx(ctx, "level", "debug", "message", "cancelling resource")
			return nil
		}

		r.k8sClient.CoreV1().ConfigMaps(cm.GetNamespace()).Update(ctx, cm, metav1.UpdateOptions{})
		if apierrors.IsNotFound(err) {
			r.k8sClient.CoreV1().ConfigMaps(cm.GetNamespace()).Create(ctx, cm, metav1.CreateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		} else if err != nil {
			return microerror.Mask(err)

		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated ConfigMap %#q in namespace %#q", currentCM.GetName(), currentCM.GetNamespace()))
	}

	return nil
}
