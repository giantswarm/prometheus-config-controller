package certificate

import (
	"context"
	"fmt"

	prometheusclient "github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-config-controller/service/key"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
)

const (
	caKey  = "ca"  // CaKey is the key in the Secret that holds the CA.
	crtKey = "crt" // CrtKey is the key in the Secret that holds the certificate.
	keyKey = "key" // KeyKey is the key in the Secret that holds the key.
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("debug", "fetching all services")

	servicesTimer := prometheusclient.NewTimer(kubernetesResource.WithLabelValues("services", "list"))
	services, err := r.k8sClient.CoreV1().Services("").List(metav1.ListOptions{})
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
		certificates, err := r.k8sClient.CoreV1().Secrets(r.certificateNamespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf(
				"clusterComponent=%s, clusterID=%s",
				r.certificateComponentName,
				clusterID,
			),
		})
		certificatesTimer.ObserveDuration()

		if err != nil {
			return nil, microerror.Maskf(err, "an error occurred fetching certificate for cluster: %s", clusterID)
		}

		if len(certificates.Items) == 0 {
			return nil, microerror.Maskf(missingError, "certificate for cluster: %s", clusterID)
		}
		certificate := certificates.Items[0]

		for _, certificateKey := range []string{caKey, crtKey, keyKey} {
			if data, ok := certificate.Data[certificateKey]; ok {
				var path string
				switch certificateKey {
				case caKey:
					path = key.CAPath(r.certificateDirectory, clusterID)
				case crtKey:
					path = key.CrtPath(r.certificateDirectory, clusterID)
				case keyKey:
					path = key.KeyPath(r.certificateDirectory, clusterID)
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
