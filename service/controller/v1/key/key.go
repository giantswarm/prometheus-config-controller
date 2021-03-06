package key

import (
	"fmt"
	"path"
	"strings"

	"github.com/giantswarm/prometheus-config-controller/pkg/label"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	AWSNodeMetricPort                = "61678"
	ElasticLoggingMetricPort         = "9108"
	NginxIngressControllerMetricPort = "10254"
	KubeStateMetricsPort             = "10301"
	CalicoNodeMetricPort             = "9091"
	ChartOperatorMetricPort          = "8000"
	ClusterAutoscalerMetricPort      = "8085"
	CertExporterMetricPort           = "9005"
	CoreDNSMetricPort                = "9153"
	KiamMetricPort                   = "9620"
	KubeProxyMetricPort              = "10249"
	NetExporterMetricPort            = "8000"
	NicExporterMetricPort            = "10800"
	VaultExporterMetricPort          = "9410"

	AWSNodeNamespace                = "kube-system"
	ElasticLoggingNamespace         = "giantswarm-elastic-logging"
	NginxIngressControllerNamespace = "kube-system"
	KubeStateMetricsNamespace       = "kube-system"
	CalicoNodeNamespace             = "kube-system"
	ChartOperatorNamespace          = "giantswarm"
	ClusterAutoscalerNamespace      = "kube-system"
	CertExporterNamespace           = "kube-system"
	CoreDNSNamespace                = "kube-system"
	KiamNamespace                   = "kube-system"
	KubeProxyNamespace              = "kube-system"
	NetExporterNamespace            = "kube-system"
	NicExporterNamespace            = "kube-system"
	VaultExporterNamespace          = "vault-exporter"

	PrefixMaster    = "master"
	PrefixApiServer = "apiserver"

	AnnotationEtcdDomain = "giantswarm.io/etcd-domain"
)

func certPath(certificateDirectory, clusterID, suffix string) string {
	return path.Join(certificateDirectory, fmt.Sprintf("%s-%s.pem", clusterID, suffix))
}

func CAPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "ca")
}

func CrtPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "crt")
}

func KeyPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "key")
}

func APIProxyPodMetricsPath(namespace, port string) string {
	return fmt.Sprintf("/api/v1/namespaces/%s/pods/${1}:%s/proxy/metrics", namespace, port)
}

func APIServiceHost(prefix string, clusterID string) string {
	return fmt.Sprintf("%s.%s:443", prefix, clusterID)
}

func ManagedAppPodMetricsPath() string {
	return "/api/v1/namespaces/${1}/pods/${2}:${3}/proxy/${4}"
}

func LabelSelectorConfigMap() labels.Selector {
	s := fmt.Sprintf("%s=%s", label.App, "prometheus")

	selector, err := labels.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse selector %#q with error %#q", s, err))
	}

	return selector
}

func LabelSelectorService() labels.Selector {
	// Selector to check if "app=master" and if "giantswarm.io/cluster"
	// exists.
	s := fmt.Sprintf("%s=%s,%s", label.App, "master", label.Cluster)

	selector, err := labels.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse selector %#q with error %#q", s, err))
	}

	return selector
}

// PrometheusURLConfig returns the Prometheus API URL that returns the current
// configuration. It assumes that address is a valid HTTP URL.
func PrometheusURLConfig(address string) string {
	u := strings.TrimSuffix(address, "/")
	return u + "/api/v1/status/config"
}

// PrometheusURLReload returns the Prometheus API URL that reloads the
// configuration. It assumes that address is a valid HTTP URL.
func PrometheusURLReload(address string) string {
	u := strings.TrimSuffix(address, "/")
	return u + "/-/reload"
}
