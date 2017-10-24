package service

import (
	"github.com/giantswarm/prometheus-config-controller/flag/service/controller"
	"github.com/giantswarm/prometheus-config-controller/flag/service/kubernetes"
	"github.com/giantswarm/prometheus-config-controller/flag/service/prometheus"
	"github.com/giantswarm/prometheus-config-controller/flag/service/resource"
)

type Service struct {
	Controller controller.Controller
	Kubernetes kubernetes.Kubernetes
	Prometheus prometheus.Prometheus
	Resource   resource.Resource
}
