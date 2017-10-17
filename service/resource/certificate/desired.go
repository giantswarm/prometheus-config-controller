package certificate

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/certificatetpr"
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
		// TODO: Refactor groupName to clusterID.
		clusterID := prometheus.GetGroupName(service)

		// TODO: componentName should come from flags.
		assets, err := r.certificatetprService.SearchCertsForComponent(clusterID, "prometheus")
		if err != nil {
			return nil, microerror.Maskf(err, "could not fetch certificates for cluster: %s", clusterID)
		}

		for assetsKey, data := range assets {
			var path string

			switch assetsKey.Type {
			case certificatetpr.CA:
				path = key.CAPath(r.certificateDirectory, clusterID)
			}

			certificateFiles = append(certificateFiles, certificateFile{
				path: path,
				data: string(data),
			})
		}
	}

	r.logger.Log("debug", "certificates fetched")

	return certificateFiles, nil
}
