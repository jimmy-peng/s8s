package app

import (
	genericapiserver "s8s/staging/apiserver/pkg/server"
	apiextensionsapiserver "s8s/staging/apiextensions-apiserver/pkg/apiserver"
)

func createAPIExtensionsConfig(
	kubeAPIServerConfig genericapiserver.Config,
	masterCount int) (*apiextensionsapiserver.Config, error) {
	genericConfig := kubeAPIServerConfig
		apiextensionsConfig := &apiextensionsapiserver.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config:                genericConfig,
		},
		ExtraConfig: apiextensionsapiserver.ExtraConfig{
			MasterCount:          masterCount,
		},
	}

	return apiextensionsConfig, nil
}

func createAPIExtensionsServer(apiextensionsConfig *apiextensionsapiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apiextensionsapiserver.CustomResourceDefinitions, error) {
	return apiextensionsConfig.Complete().New(delegateAPIServer)
}