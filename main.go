package main

import (
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/server"
	"github.com/giantswarm/prometheus-config-controller/service"
)

var (
	f *flag.Flag = flag.New()

	description string = "The prometheus-config-controller provides Prometheus service discovery for Kubernetes clusters on Kubernetes."
	gitCommit   string = "n/a"
	name        string = "prometheus-config-controller"
	source      string = "https://github.com/giantswarm/prometheus-config-controller"
)

func panicOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainError() error {
	var err error

	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			c := service.Config{
				Logger: newLogger,

				Description: description,
				Flag:        f,
				GitCommit:   gitCommit,
				Name:        name,
				Source:      source,
				Viper:       v,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
			go newService.Boot()
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:  newLogger,
				Service: newService,
				Viper:   v,

				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description:    description,
			GitCommit:      gitCommit,
			Name:           name,
			Source:         source,
			VersionBundles: service.NewVersionBundles(),
		}

		newCommand, err = command.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.Address, "http://127.0.0.1:6443", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
	daemonCommand.PersistentFlags().Bool(f.Service.Kubernetes.InCluster, false, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.KubeConfig, "", "KubeConfig used to connect to Kubernetes. When empty other settings are used.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CAFile, "", "Certificate authority file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CrtFile, "", "Certificate file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.KeyFile, "", "Key file path to use to authenticate with Kubernetes.")

	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Address, "http://127.0.0.1:9090", "Address of Prometheus to reload.")

	daemonCommand.PersistentFlags().Int(f.Service.Resource.Retries, 3, "Number of times to retry resources.")

	daemonCommand.PersistentFlags().String(f.Service.Resource.Certificate.ComponentName, "prometheus", "Component name label for certificates.")
	daemonCommand.PersistentFlags().String(f.Service.Resource.Certificate.Directory, "/certs", "Directory in which to store certificates.")
	daemonCommand.PersistentFlags().String(f.Service.Resource.Certificate.Namespace, "default", "Namespace for certificates.")
	daemonCommand.PersistentFlags().Int(f.Service.Resource.Certificate.Permission, 0600, "File permission for certificates.")

	daemonCommand.PersistentFlags().String(f.Service.Resource.ConfigMap.Key, "prometheus.yml", "Key in configmap under which prometheus configuration is held.")
	daemonCommand.PersistentFlags().String(f.Service.Resource.ConfigMap.Name, "prometheus", "Name of prometheus configmap to control.")
	daemonCommand.PersistentFlags().String(f.Service.Resource.ConfigMap.Namespace, "monitoring", "Namespace of prometheus configmap to control.")
	daemonCommand.PersistentFlags().Duration(f.Service.Resource.ConfigMap.MinimumReloadTime, 2*time.Minute, "Minimum time between reloads of Prometheus.")

	daemonCommand.PersistentFlags().Duration(f.Service.Controller.ControllerBackOffDuration, time.Minute*5, "Maximum backoff duration for controller.")
	daemonCommand.PersistentFlags().Duration(f.Service.Controller.FrameworkBackOffDuration, time.Minute*5, "Maximum backoff duration for operator framework.")
	daemonCommand.PersistentFlags().Duration(f.Service.Controller.ResyncPeriod, time.Minute*1, "Controller resync period.")

	newCommand.CobraCommand().Execute()

	return nil
}
