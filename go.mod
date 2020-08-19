module github.com/giantswarm/prometheus-config-controller

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/giantswarm/apiextensions/v2 v2.1.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/k8sclient/v4 v4.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v1.2.0
	github.com/giantswarm/operatorkit/v2 v2.0.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-cmp v0.5.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.11.1
	github.com/prometheus/prometheus v2.20.1+incompatible
	github.com/spf13/afero v1.2.2
	github.com/spf13/viper v1.6.2
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
)
