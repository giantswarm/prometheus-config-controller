package service

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/kubernetes"
)

type Service struct {
	Kubernetes kubernetes.Kubernetes
}
