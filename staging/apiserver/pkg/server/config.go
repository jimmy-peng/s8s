package server

import "s8s/staging/apimachinery/pkg/runtime/serializer"

type Config struct {

}

func NewConfig(codecs serializer.CodecFactory) *Config {
	return &Config{}
}