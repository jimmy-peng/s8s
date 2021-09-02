package app

import (
	"fmt"
	"s8s/cmd/kube-apiserver/app/options"
	cliflag "s8s/component-base/cli/flag"
	"s8s/pkg/api/legacyscheme"

	"s8s/pkg/controlplane"
	genericapiserver "s8s/staging/apiserver/pkg/server"

	"github.com/spf13/cobra"
)

type completedServerRunOptions struct {
	*options.ServerRunOptions
}

func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	options.ServerRunOptions = s
	return options, nil
}

func buildGenericConfig(s *options.ServerRunOptions) (genericConfig *genericapiserver.Config) {
	genericConfig = genericapiserver.NewConfig(legacyscheme.Codecs)
	return genericConfig
}

func CreateKubeAPIServerConfig(s completedServerRunOptions) (
	*controlplane.Config,
	error) {

	genericConfig := buildGenericConfig(s.ServerRunOptions)

	config := &controlplane.Config{
		GenericConfig: genericConfig,
	}

	return config, nil

}

func CreateServerChain(completedOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	kubeAPIServerConfig, err := CreateKubeAPIServerConfig(completedOptions)
	//_, err := CreateKubeAPIServerConfig(completeOptions)
	apiExtensionsConfig, err := createAPIExtensionsConfig(*kubeAPIServerConfig.GenericConfig, completedOptions.MasterCount)

	//apiExtensionsServer, err := createAPIExtensionsServer(apiExtensionsConfig, genericapiserver.NewEmptyDelegate())
	_, err = createAPIExtensionsServer(apiExtensionsConfig, genericapiserver.NewEmptyDelegate())

	return err
}

func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	err := CreateServerChain(completeOptions, stopCh)
	if err != nil {
		return err
	}
	return nil
}

func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "kube-apiserver",
		Long: `The Kubernetes API server validates and configures data
for the api objects which include pods, services, replicationcontrollers, and
others. The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := cmd.Flags()
			cliflag.PrintFlags(fs)
			completeOptions, err := Complete(s)
			if err != nil {
				return err
			}

			if errs := completeOptions.Validate(); len(errs) != 0 {
				return nil
			}
			stopCh := make(chan struct{})
			Run(completeOptions, stopCh)

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath, args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	return cmd
}
