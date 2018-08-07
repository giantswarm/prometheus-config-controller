package key

import (
	"fmt"
	"path"
)

const (
	NginxICMetricPort       = "10254"
	KubeStateMetricsPort    = "10301"
	ChartOperatorMetricPort = "8000"

	PrefixMaster    = "master"
	PrefixApiServer = "apiserver"

	AnnotationEtcdDomain = "giantswarm.io/etcd-domain"

	LabelVersionBundle = "giantswarm.io/version-bundle"
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

func APIProxyPodMetricsPath(port string) string {
	return fmt.Sprintf("/api/v1/namespaces/kube-system/pods/${1}:%s/proxy/metrics", port)
}

func APIServiceHost(prefix string, clusterID string) string {
	return fmt.Sprintf("%s.%s:443", prefix, clusterID)
}
