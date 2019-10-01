package project

var (
	description string = "The prometheus-config-controller provides Prometheus service discovery for Kubernetes clusters on Kubernetes."
	gitSHA      string = "n/a"
	name        string = "prometheus-config-controller"
	source      string = "https://github.com/giantswarm/prometheus-config-controller"
	version            = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
