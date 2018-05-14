package prometheus

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test_Prometheus_getJobName tests the getJobName function.
func Test_Prometheus_getJobName(t *testing.T) {
	tests := []struct {
		service         v1.Service
		name            string
		expectedJobName string
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			name:            "cadvisor",
			expectedJobName: "guest-cluster-bar-cadvisor",
		},
	}

	for index, test := range tests {
		jobName := getJobName(test.service, test.name)

		if test.expectedJobName != jobName {
			t.Fatalf(
				"%d: expected job name does not match returned job name.\nexpected: %s\nreturned: %s\n",
				index,
				test.expectedJobName,
				jobName,
			)
		}
	}
}

// Test_Prometheus_getTargetHost tests the getTargetHost function.
func Test_Prometheus_getTargetHost(t *testing.T) {
	tests := []struct {
		service            v1.Service
		expectedTargetHost string
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedTargetHost: "foo.bar",
		},
	}

	for index, test := range tests {
		targetHost := getTargetHost(test.service)

		if test.expectedTargetHost != targetHost {
			t.Fatalf(
				"%d: expected target host does not match returned target host.\nexpected: %s\nreturned: %s\n",
				index,
				test.expectedTargetHost,
				targetHost,
			)
		}
	}
}

// Test_Prometheus_getTarget tests the getTarget function.
func Test_Prometheus_getTarget(t *testing.T) {
	tests := []struct {
		service        v1.Service
		expectedTarget model.LabelSet
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedTarget: model.LabelSet{model.AddressLabel: "foo.bar"},
		},
	}

	for index, test := range tests {
		target := getTarget(test.service)

		if !reflect.DeepEqual(test.expectedTarget, target) {
			t.Fatalf(
				"%d: expected target does not match returned target.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedTarget),
				spew.Sdump(target),
			)
		}
	}
}

// Test_Prometheus_GetScrapeConfigs tests the GetScrapeConfigs function.
func Test_Prometheus_GetScrapeConfigs(t *testing.T) {
	tests := []struct {
		services             []v1.Service
		certificateDirectory string

		expectedScrapeConfigs []config.ScrapeConfig
	}{
		// Test that when there are no services available,
		// no scrape configs are returned.
		{
			services:             nil,
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a non-annotated service does not create a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a service that specifies the cluster annotation creates a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{
				TestConfigOneApiserver,
				TestConfigOneCadvisor,
				TestConfigOneKubelet,
				TestConfigOneNodeExporter,
				TestConfigOneWorkload,
			},
		},

		// Test that two services that specify different clusters create separate configs.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "0ba9v",
						Annotations: map[string]string{
							ClusterAnnotation: "0ba9v",
						},
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{
				TestConfigTwoApiserver,
				TestConfigTwoCadvisor,
				TestConfigTwoKubelet,
				TestConfigTwoNodeExporter,
				TestConfigTwoWorkload,

				TestConfigOneApiserver,
				TestConfigOneCadvisor,
				TestConfigOneKubelet,
				TestConfigOneNodeExporter,
				TestConfigOneWorkload,
			},
		},
	}

	for index, test := range tests {
		scrapeConfigs, err := GetScrapeConfigs(test.services, test.certificateDirectory)
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !reflect.DeepEqual(test.expectedScrapeConfigs, scrapeConfigs) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedScrapeConfigs),
				spew.Sdump(scrapeConfigs),
			)
		}
	}
}

// Test_Prometheus_GetScrapeConfigs_Deterministic tests that the GetScrapeConfigs function is deterministic,
// and that scrape configs are returned in alphabetical order by the job name.
func Test_Prometheus_GetScrapeConfigs_Deterministic(t *testing.T) {
	services := []v1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apiserver",
				Namespace: "xa5ly",
				Annotations: map[string]string{
					ClusterAnnotation: "xa5ly",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apiserver",
				Namespace: "0ba9v",
				Annotations: map[string]string{
					ClusterAnnotation: "0ba9v",
				},
			},
		},
	}

	expectedScrapeConfigs := []config.ScrapeConfig{
		TestConfigTwoApiserver,
		TestConfigTwoCadvisor,
		TestConfigTwoKubelet,
		TestConfigTwoNodeExporter,
		TestConfigTwoWorkload,

		TestConfigOneApiserver,
		TestConfigOneCadvisor,
		TestConfigOneKubelet,
		TestConfigOneNodeExporter,
		TestConfigOneWorkload,
	}

	for index := 0; index < 50; index++ {
		scrapeConfigs, err := GetScrapeConfigs(services, "/certs")
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !reflect.DeepEqual(expectedScrapeConfigs, scrapeConfigs) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs. GetScrapeConfigs not deterministic.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(expectedScrapeConfigs),
				spew.Sdump(scrapeConfigs),
			)
		}
	}
}

// Test_Prometheus_YamlMarshal tests that Prometheus marshals yaml correctly.
func Test_Prometheus_YamlMarshal(t *testing.T) {
	tests := []struct {
		service v1.Service

		expectedConfigs string
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "apiserver",
					Namespace: "xa5ly",
					Annotations: map[string]string{
						ClusterAnnotation: "xa5ly",
					},
				},
			},

			expectedConfigs: `- job_name: guest-cluster-xa5ly-apiserver
  scheme: https
  kubernetes_sd_configs:
  - api_server: https://apiserver.xa5ly
    role: endpoints
    tls_config:
      ca_file: /certs/xa5ly-ca.pem
      cert_file: /certs/xa5ly-crt.pem
      key_file: /certs/xa5ly-key.pem
      insecure_skip_verify: false
  tls_config:
    ca_file: /certs/xa5ly-ca.pem
    cert_file: /certs/xa5ly-crt.pem
    key_file: /certs/xa5ly-key.pem
    insecure_skip_verify: true
  relabel_configs:
  - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name]
    regex: default;kubernetes
    action: keep
  - source_labels: []
    target_label: app
    replacement: kubernetes
  - source_labels: []
    target_label: cluster_id
    replacement: xa5ly
  - source_labels: []
    target_label: cluster_type
    replacement: guest
- job_name: guest-cluster-xa5ly-cadvisor
  scheme: https
  kubernetes_sd_configs:
  - api_server: https://apiserver.xa5ly
    role: node
    tls_config:
      ca_file: /certs/xa5ly-ca.pem
      cert_file: /certs/xa5ly-crt.pem
      key_file: /certs/xa5ly-key.pem
      insecure_skip_verify: false
  tls_config:
    ca_file: /certs/xa5ly-ca.pem
    cert_file: /certs/xa5ly-crt.pem
    key_file: /certs/xa5ly-key.pem
    insecure_skip_verify: false
  relabel_configs:
  - source_labels: []
    target_label: __address__
    replacement: apiserver.xa5ly
  - source_labels: [__meta_kubernetes_node_name]
    target_label: __metrics_path__
    replacement: /api/v1/nodes/${1}:4194/proxy/metrics
  - source_labels: []
    target_label: app
    replacement: cadvisor
  - source_labels: []
    target_label: cluster_id
    replacement: xa5ly
  - source_labels: []
    target_label: cluster_type
    replacement: guest
  - source_labels: [__meta_kubernetes_node_address_InternalIP]
    target_label: ip
  - source_labels: [__meta_kubernetes_node_label_role]
    target_label: role
  - source_labels: [__meta_kubernetes_node_label_role]
    regex: null
    target_label: role
    replacement: worker
  metric_relabel_configs:
  - source_labels: [namespace]
    regex: (kube-system|giantswarm)
    action: keep
  - source_labels: [__name__]
    regex: container_network_.*
    action: keep
- job_name: guest-cluster-xa5ly-kubelet
  scheme: https
  kubernetes_sd_configs:
  - api_server: https://apiserver.xa5ly
    role: node
    tls_config:
      ca_file: /certs/xa5ly-ca.pem
      cert_file: /certs/xa5ly-crt.pem
      key_file: /certs/xa5ly-key.pem
      insecure_skip_verify: false
  tls_config:
    ca_file: /certs/xa5ly-ca.pem
    cert_file: /certs/xa5ly-crt.pem
    key_file: /certs/xa5ly-key.pem
    insecure_skip_verify: true
  relabel_configs:
  - source_labels: []
    target_label: app
    replacement: kubelet
  - source_labels: []
    target_label: cluster_id
    replacement: xa5ly
  - source_labels: []
    target_label: cluster_type
    replacement: guest
  - source_labels: [__meta_kubernetes_node_address_InternalIP]
    target_label: ip
  - source_labels: [__meta_kubernetes_node_label_role]
    target_label: role
  - source_labels: [__meta_kubernetes_node_label_role]
    regex: null
    target_label: role
    replacement: worker
- job_name: guest-cluster-xa5ly-node-exporter
  scheme: http
  kubernetes_sd_configs:
  - api_server: https://apiserver.xa5ly
    role: endpoints
    tls_config:
      ca_file: /certs/xa5ly-ca.pem
      cert_file: /certs/xa5ly-crt.pem
      key_file: /certs/xa5ly-key.pem
      insecure_skip_verify: false
  relabel_configs:
  - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name]
    regex: kube-system;node-exporter
    action: keep
  - source_labels: [__address__]
    regex: (.*):10250
    target_label: __address__
    replacement: ${1}:10300
  - source_labels: []
    target_label: app
    replacement: node-exporter
  - source_labels: []
    target_label: cluster_id
    replacement: xa5ly
  - source_labels: []
    target_label: cluster_type
    replacement: guest
  - source_labels: [__address__]
    regex: (.*):10300
    target_label: ip
    replacement: ${1}
  metric_relabel_configs:
  - source_labels: [fstype]
    regex: (cgroup|devpts|mqueue|nsfs|overlay|tmpfs)
    action: keep
  - source_labels: [__name__, state]
    regex: node_systemd_unit_state;(active|activating|deactivating|inactive)
    action: drop
  - source_labels: [__name__, name]
    regex: node_systemd_unit_state;(dev-disk-by|run-docker-netns|sys-devices|sys-subsystem-net|var-lib-docker-overlay2|var-lib-docker-containers|var-lib-kubelet-pods).*
    action: drop
- job_name: guest-cluster-xa5ly-workload
  scheme: http
  kubernetes_sd_configs:
  - api_server: https://apiserver.xa5ly
    role: endpoints
    tls_config:
      ca_file: /certs/xa5ly-ca.pem
      cert_file: /certs/xa5ly-crt.pem
      key_file: /certs/xa5ly-key.pem
      insecure_skip_verify: false
  relabel_configs:
  - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name]
    regex: kube-system;kube-state-metrics|node-exporter
    action: keep
  - source_labels: [__meta_kubernetes_service_name]
    target_label: app
  - source_labels: [__meta_kubernetes_namespace]
    target_label: namespace
  - source_labels: []
    target_label: cluster_id
    replacement: xa5ly
  - source_labels: []
    target_label: cluster_type
    replacement: guest
  metric_relabel_configs:
  - source_labels: [exported_namespace]
    regex: (kube-system|giantswarm)
    action: keep
`,
		},
	}

	for index, test := range tests {
		scrapeConfigs, err := GetScrapeConfigs([]v1.Service{test.service}, "/certs")
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		data, err := yaml.Marshal(scrapeConfigs)
		if err != nil {
			t.Fatalf("%d: error occurred marshaling yaml: %s\n", index, err)
		}

		expectedLines := strings.Split(test.expectedConfigs, "\n")
		returnedLines := strings.Split(string(data), "\n")

		if len(expectedLines) == len(returnedLines) {
			for i := 0; i < len(expectedLines); i++ {
				if expectedLines[i] != returnedLines[i] {
					t.Logf("\nexpected line:\n'%s'\nreturned line:\n'%s'\n", expectedLines[i], returnedLines[i])
				}
			}
		}

		if test.expectedConfigs != string(data) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs.\nexpected:\n%s\nreturned:\n%s\n",
				index,
				test.expectedConfigs,
				string(data),
			)
		}
	}
}
