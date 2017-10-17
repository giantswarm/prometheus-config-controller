package certificate

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-config-controller/service/key"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
)

const (
	CaKey  = "ca"  // CaKey is the key in the Secret that holds the CA.
	CrtKey = "crt" // CrtKey is the key in the Secret that holds the certificate.
	KeyKey = "key" // KeyKey is the key in the Secret that holds the key.
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("debug", "fetching all services")
	services, err := r.k8sClient.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred listing all services")
	}

	r.logger.Log("debug", "filtering services")
	validServices := prometheus.FilterInvalidServices(services.Items)

	r.logger.Log("debug", "fetching certificates")
	certificateFiles := []certificateFile{}
	for _, service := range validServices {
		namespace, name := prometheus.GetCertificateName(service)

		certificate, err := r.k8sClient.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return nil, microerror.Maskf(missingError, "certificate %s/%s", namespace, name)
		} else if err != nil {
			return nil, microerror.Maskf(err, "an error occurred fetching certificate %s/%s", namespace, name)
		}

		groupName := prometheus.GetGroupName(service)

		for _, certificateKey := range []string{CaKey, CrtKey, KeyKey} {
			if data, ok := certificate.Data[certificateKey]; ok {
				var path string
				switch certificateKey {
				case CaKey:
					path = key.CAPath(r.certificateDirectory, groupName)
				case CrtKey:
					path = key.CrtPath(r.certificateDirectory, groupName)
				case KeyKey:
					path = key.KeyPath(r.certificateDirectory, groupName)
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
