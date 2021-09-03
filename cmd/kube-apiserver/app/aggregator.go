package app
import (
	genericapiserver "s8s/staging/apiserver/pkg/server"
	"s8s/cmd/kube-apiserver/app/options"
	aggregatorapiserver "s8s/staging/kube-aggregator/pkg/apiserver"
)

func createAggregatorConfig(
	kubeAPIServerConfig genericapiserver.Config,
	commandOptions *options.ServerRunOptions) (*aggregatorapiserver.Config, error) {
	// make a shallow copy to let us twiddle a few things
	// most of the config actually remains the same.  We only need to mess with a couple items related to the particulars of the aggregator
	genericConfig := kubeAPIServerConfig

	aggregatorConfig := &aggregatorapiserver.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config:                genericConfig,
		},
		ExtraConfig: aggregatorapiserver.ExtraConfig{
		},
	}
	return aggregatorConfig, nil
}

func createAggregatorServer(aggregatorConfig *aggregatorapiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*aggregatorapiserver.APIAggregator, error) {
	aggregatorServer, err := aggregatorConfig.Complete().NewWithDelegate(delegateAPIServer)
	if err != nil {
		return nil, err
	}

	return aggregatorServer, nil
}