package resource

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/resource/certificate"
	"github.com/giantswarm/prometheus-config-controller/flag/service/resource/configmap"
)

type Resource struct {
	Certificate certificate.Certificate
	ConfigMap   configmap.ConfigMap
	Retries     string
}
