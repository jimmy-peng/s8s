package controlplane
import (
 	genericapiserver "s8s/staging/apiserver/pkg/server"
)

type ExtraConfig struct {
	MasterCount int
}

type Config struct {
	GenericConfig *genericapiserver.Config
	ExtraConfig ExtraConfig
}