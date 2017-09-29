package resource

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/resource/configmap"
)

type Resource struct {
	ConfigMap configmap.ConfigMap
	Retries   string
}
