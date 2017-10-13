package resource

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/resource/configmap"
)

type Resource struct {
	CertificateDirectory string
	ConfigMap            configmap.ConfigMap
	Retries              string
}
