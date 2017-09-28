package kubernetes

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/kubernetes/tls"
)

type Kubernetes struct {
	Address   string
	InCluster string
	TLS       tls.TLS
}
