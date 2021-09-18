package server

import (
	"net"
	"net/http"
	"s8s/staging/apimachinery/pkg/runtime/serializer"
	"s8s/staging/apiserver/pkg/server/dynamiccertificates"
)

type Config struct {
	// SecureServing is required to serve https
	SecureServing *SecureServingInfo

	// LoopbackClientConfig is a config for a privileged loopback connection to the API server
	// This is required for proper functioning of the PostStartHooks on a GenericAPIServer
	// TODO: move into SecureServing(WithLoopback) as soon as insecure serving is gone
	//LoopbackClientConfig *restclient.Config
	// BuildHandlerChainFunc allows you to build custom handler chains by decorating the apiHandler.
	BuildHandlerChainFunc func(apiHandler http.Handler, c *Config) (secure http.Handler)
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c completedConfig) New(name string, delegationTarget DelegationTarget) (*GenericAPIServer, error) {
	handlerChainBuilder := func(handler http.Handler) http.Handler {
		return c.BuildHandlerChainFunc(handler, c.Config)
	}
	apiServerHandler := NewAPIServerHandler(name, handlerChainBuilder, delegationTarget.UnprotectedHandler())
	s := &GenericAPIServer{
		delegationTarget:  delegationTarget,
		SecureServingInfo: c.SecureServing,
		Handler:           apiServerHandler,
	}
	return s, nil
}

type RecommendedConfig struct {
	Config
}

type SecureServingInfo struct {
	// Listener is the secure server network listener.
	Listener net.Listener

	// Cert is the main server cert which is used if SNI does not match. Cert must be non-nil and is
	// allowed to be in SNICerts.
	Cert dynamiccertificates.CertKeyContentProvider

	// SNICerts are the TLS certificates used for SNI.
	//SNICerts []dynamiccertificates.SNICertKeyContentProvider

	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	//ClientCA dynamiccertificates.CAContentProvider

	// MinTLSVersion optionally overrides the minimum TLS version supported.
	// Values are from tls package constants (https://golang.org/pkg/crypto/tls/#pkg-constants).
	MinTLSVersion uint16

	// CipherSuites optionally overrides the list of allowed cipher suites for the server.
	// Values are from tls package constants (https://golang.org/pkg/crypto/tls/#pkg-constants).
	CipherSuites []uint16

	// HTTP2MaxStreamsPerConnection is the limit that the api server imposes on each client.
	// A value of zero means to use the default provided by golang's HTTP/2 support.
	HTTP2MaxStreamsPerConnection int

	// DisableHTTP2 indicates that http2 should not be enabled.
	DisableHTTP2 bool
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{c}}
}

func (c *RecommendedConfig) Complete() CompletedConfig {
	return c.Config.Complete()
}

func NewConfig(codecs serializer.CodecFactory) *Config {
	return &Config{
		BuildHandlerChainFunc:       DefaultBuildHandlerChain,
	}
}

func DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) http.Handler {
	return apiHandler
}