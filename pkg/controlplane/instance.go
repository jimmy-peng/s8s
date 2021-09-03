package controlplane

import (
	genericapiserver "s8s/staging/apiserver/pkg/server"
)

type ExtraConfig struct {
	MasterCount int
}

type Config struct {
	GenericConfig *genericapiserver.Config
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package
type CompletedConfig struct {
	*completedConfig
}

func (c *Config) Complete() CompletedConfig {
	cfg := completedConfig{
		c.GenericConfig.Complete(),
		&c.ExtraConfig,
	}
	return CompletedConfig{&cfg}
}

func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*Instance, error) {
	s, err := c.GenericConfig.New("kube-apiserver", delegationTarget)
	if err != nil {
		return nil, err
	}

	m := &Instance{
		GenericAPIServer: s,
	}

	return m, nil
}

// Instance contains state for a Kubernetes cluster api server instance.
type Instance struct {
	GenericAPIServer *genericapiserver.GenericAPIServer

	//ClusterAuthenticationInfo clusterauthenticationtrust.ClusterAuthenticationInfo
}
