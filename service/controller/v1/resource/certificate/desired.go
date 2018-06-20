package certificate

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

const (
	caKey  = "ca"  // CaKey is the key in the Secret that holds the CA.
	crtKey = "crt" // CrtKey is the key in the Secret that holds the certificate.
	keyKey = "key" // KeyKey is the key in the Secret that holds the key.

	// ServiceLabelSelector is the label selector to match master services.
	serviceLabelSelector = "app=master"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("debug", "fetching all services")

	servicesTimer := prometheusclient.NewTimer(kubernetesResource.WithLabelValues("services", "list"))
	services, err := r.k8sClient.CoreV1().Services("").List(metav1.ListOptions{
		LabelSelector: serviceLabelSelector,
	})
	servicesTimer.ObserveDuration()

	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred listing all services")
	}

	r.logger.Log("debug", "filtering services")
	validServices := prometheus.FilterInvalidServices(services.Items)

	r.logger.Log("debug", "fetching certificates")
	certificateFiles := []certificateFile{}

	for _, service := range validServices {
		clusterID := prometheus.GetClusterID(service)

		certificatesTimer := prometheusclient.NewTimer(kubernetesResource.WithLabelValues("secrets", "list"))
		certificates, err := r.k8sClient.CoreV1().Secrets(r.certNamespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf(
				"clusterComponent=%s, clusterID=%s",
				r.certComponentName,
				clusterID,
			),
		})
		certificatesTimer.ObserveDuration()

		if err != nil {
			return nil, microerror.Maskf(err, "an error occurred fetching certificate for cluster: %s", clusterID)
		}

		if len(certificates.Items) == 0 {
			// If the certificate can't be found, try to continue on.
			// It's possible that the certificate just hasn't been created yet.
			// If the certificate is consistently missing, we'll be notified
			// about the cluster not being scrapeable.
			r.logger.Log("warning", fmt.Sprintf("certificate for cluster '%s' is missing, continuing", clusterID))
			continue
		}
		certificate := certificates.Items[0]

		for _, certificateKey := range []string{caKey, crtKey, keyKey} {
			if data, ok := certificate.Data[certificateKey]; ok {
				var path string
				switch certificateKey {
				case caKey:
					path = key.CAPath(r.certDirectory, clusterID)
				case crtKey:
					path = key.CrtPath(r.certDirectory, clusterID)
				case keyKey:
					path = key.KeyPath(r.certDirectory, clusterID)
				}

				certificateFiles = append(certificateFiles, certificateFile{
					path: path,
					data: string(data),
				})
			}
		}
	}

	r.logger.Log("debug", "certificates fetched")

	return certificateFiles, nil
}
