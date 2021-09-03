package server

import "s8s/staging/apimachinery/pkg/runtime/serializer"

type Config struct {

}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c completedConfig) New(name string, delegationTarget DelegationTarget) (*GenericAPIServer, error) {
	s := &GenericAPIServer{
		delegationTarget:           delegationTarget,
	}
	return s, nil
}

type RecommendedConfig struct {
	Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{c}}
}

func (c *RecommendedConfig) Complete() CompletedConfig {
	return c.Config.Complete()
}

func NewConfig(codecs serializer.CodecFactory) *Config {
	return &Config{}
}