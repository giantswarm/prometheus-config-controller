package patchfinalizer

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	finalizerPCC = "operatorkit.giantswarm.io/prometheus-config-controller"
	labelCluster = "giantswarm.io/cluster"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	var services []corev1.Service
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "finding services to patch")

		list, err := r.k8sClient.CoreV1().Services(corev1.NamespaceAll).List(metav1.ListOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		for _, s := range list.Items {
			if hasLabel(s.Labels) {
				continue
			}

			if !hasFinalizer(s.Finalizers) {
				continue
			}

			services = append(services, s)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found %d services to patch", len(services)))
	}

	{
		for _, s := range services {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("patching service %q in namespace %q", s.Name, s.Namespace))

			s.Finalizers = withoutFinalizer(s.Finalizers)

			_, err := r.k8sClient.CoreV1().Services(s.Namespace).Update(&s)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("patched service %q in namespace %q", s.Name, s.Namespace))
		}
	}

	return nil
}

func hasFinalizer(finalizers []string) bool {
	for _, f := range finalizers {
		if f == finalizerPCC {
			return true
		}
	}

	return false
}

func hasLabel(labels map[string]string) bool {
	for _, l := range labels {
		if l == labelCluster {
			return true
		}
	}

	return false
}

func withoutFinalizer(finalizers []string) []string {
	var list []string

	for _, f := range finalizers {
		if f == finalizerPCC {
			continue
		}

		list = append(list, f)
	}

	return list
}
