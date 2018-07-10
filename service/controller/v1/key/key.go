package key

import (
	"fmt"
	"path"
)

const (
	NginxICMetricPort   = "10254"
	KubeStaeMetricsPort = "10301"

	PrefixMaster    = "master"
	PrefixApiServer = "apiserver"

	LabelVersionBundle = "giantswarm.io/version-bundle"
	LabelEtcdDomain    = "giantswarm.io/etcd-domain"
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

func EtcdTargetUrl(etcdDomain string) string {
	return fmt.Sprintf("https://%s:443", etcdDomain)
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
