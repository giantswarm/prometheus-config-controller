package prometheus

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/pkg/relabel"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
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
			expectedJobName: "workload-cluster-bar-cadvisor",
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
		provider             string

		expectedScrapeConfigs []config.ScrapeConfig
	}{
		// 0. Test that when there are no services available,
		// no scrape configs are returned.
		{
			services:             nil,
			certificateDirectory: "/certs",
			provider:             "aws-test",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// 1. Test that a non-annotated service does not create a scrape config.
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
			provider:             "aws-test",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// 2. Test that a service that specifies the cluster annotation creates a scrape config.
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
			provider:             "aws-test",

			expectedScrapeConfigs: []config.ScrapeConfig{
				TestConfigOneApiserver,
				TestConfigOneAWSNode,
				TestConfigOneCadvisor,
				TestConfigOneCalicoNode,
				TestConfigOneDocker,
				TestConfigOneIngress,
				TestConfigOneKubeProxy,
				TestConfigOneKubeStateManagedApp,
				TestConfigOneKubelet,
				TestConfigOneManagedApp,
				TestConfigOneNodeExporter,
				TestConfigOneWorkload,
			},
		},

		// 3. Test that two services that specify different clusters create separate configs.
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
			provider:             "aws-test",

			expectedScrapeConfigs: []config.ScrapeConfig{
				TestConfigTwoApiserver,
				TestConfigTwoAWSNode,
				TestConfigTwoCadvisor,
				TestConfigTwoCalicoNode,
				TestConfigTwoDocker,
				TestConfigTwoIngress,
				TestConfigTwoKubeProxy,
				TestConfigTwoKubeStateManagedApp,
				TestConfigTwoKubelet,
				TestConfigTwoManagedApp,
				TestConfigTwoNodeExporter,
				TestConfigTwoWorkload,

				TestConfigOneApiserver,
				TestConfigOneAWSNode,
				TestConfigOneCadvisor,
				TestConfigOneCalicoNode,
				TestConfigOneDocker,
				TestConfigOneIngress,
				TestConfigOneKubeProxy,
				TestConfigOneKubeStateManagedApp,
				TestConfigOneKubelet,
				TestConfigOneManagedApp,
				TestConfigOneNodeExporter,
				TestConfigOneWorkload,
			},
		},
	}

	for index, test := range tests {
		metaConfig := Config{
			CertDirectory: test.certificateDirectory,
			Provider:      test.provider,
		}
		scrapeConfigs, err := GetScrapeConfigs(test.services, metaConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !cmp.Equal(test.expectedScrapeConfigs, scrapeConfigs, cmpopts.IgnoreUnexported(relabel.Regexp{}, regexp.Regexp{})) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs.\ndiff: %s\n",
				index,
				cmp.Diff(test.expectedScrapeConfigs, scrapeConfigs, cmpopts.IgnoreUnexported(relabel.Regexp{}, regexp.Regexp{})),
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
		TestConfigTwoAWSNode,
		TestConfigTwoCadvisor,
		TestConfigTwoCalicoNode,
		TestConfigTwoDocker,
		TestConfigTwoIngress,
		TestConfigTwoKubeProxy,
		TestConfigTwoKubeStateManagedApp,
		TestConfigTwoKubelet,
		TestConfigTwoManagedApp,
		TestConfigTwoNodeExporter,
		TestConfigTwoWorkload,

		TestConfigOneApiserver,
		TestConfigOneAWSNode,
		TestConfigOneCadvisor,
		TestConfigOneCalicoNode,
		TestConfigOneDocker,
		TestConfigOneIngress,
		TestConfigOneKubeProxy,
		TestConfigOneKubeStateManagedApp,
		TestConfigOneKubelet,
		TestConfigOneManagedApp,
		TestConfigOneNodeExporter,
		TestConfigOneWorkload,
	}

	for index := 0; index < 50; index++ {
		metaConfig := Config{
			CertDirectory: "/certs",
			Provider:      "aws-test",
		}
		scrapeConfigs, err := GetScrapeConfigs(services, metaConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !cmp.Equal(expectedScrapeConfigs, scrapeConfigs, cmpopts.IgnoreUnexported(relabel.Regexp{}, regexp.Regexp{})) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs. GetScrapeConfigs not deterministic.\ndiff: %s\n",
				index,
				cmp.Diff(spew.Sdump(expectedScrapeConfigs), spew.Sdump(scrapeConfigs)),
			)
		}
	}
}

var update = flag.Bool("update", false, "update .golden scrapeconfig file")

// Test_Prometheus_YamlMarshal tests that Prometheus marshals YAML correctly.
//
// It uses golden file as reference template and when changes to template are
// intentional, they can be updated by providing -update flag for go test.
//
//  go test ./service/controller/v1/prometheus -run Test_Prometheus_YamlMarshal -update
//
func Test_Prometheus_YamlMarshal(t *testing.T) {
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apiserver",
			Namespace: "xa5ly",
			Annotations: map[string]string{
				ClusterAnnotation: "xa5ly",
			},
		},
	}

	metaConfig := Config{
		CertDirectory: "/certs",
		Provider:      "aws-test",
	}
	scrapeConfigs, err := GetScrapeConfigs([]v1.Service{service}, metaConfig)
	if err != nil {
		t.Fatalf("error returned creating scrape configs: %s\n", err)
	}

	data, err := yaml.Marshal(scrapeConfigs)
	if err != nil {
		t.Fatalf("error occurred marshaling yaml: %s\n", err)
	}

	p := filepath.Join("testdata", "scrapeconfig.golden")

	if *update {
		err := ioutil.WriteFile(p, data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	goldenFile, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(data, goldenFile) {
		t.Fatalf("\n\n%s\n", cmp.Diff(goldenFile, data))
	}
}
