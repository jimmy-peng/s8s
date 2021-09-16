package server

import (
	"time"
)

type DelegationTarget interface {
	ListedPaths() []string
	// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
	PrepareRun() preparedGenericAPIServer
}

type emptyDelegate struct {
}

// preparedGenericAPIServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedGenericAPIServer struct {
	*GenericAPIServer
}

func (s emptyDelegate) ListedPaths() []string {
	return []string{}
}

func (s emptyDelegate) PrepareRun() preparedGenericAPIServer {
	return preparedGenericAPIServer{nil}
}

func NewEmptyDelegate() DelegationTarget {
	return emptyDelegate{}
}

// GenericAPIServer contains state for a Kubernetes cluster api server.
type GenericAPIServer struct {
	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration
	// delegationTarget is the next delegate in the chain. This is never nil.
	delegationTarget DelegationTarget

	// SecureServingInfo holds configuration of the TLS server.
	SecureServingInfo *SecureServingInfo

	// "Outputs"
	// Handler holds the handlers being used by this API server
	Handler *APIServerHandler
}

func (s GenericAPIServer) ListedPaths() []string {
	return []string{}
}

// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
func (s *GenericAPIServer) PrepareRun() preparedGenericAPIServer {
	s.delegationTarget.PrepareRun()

	return preparedGenericAPIServer{s}
}

// NonBlockingRun spawns the secure http server. An error is
// returned if the secure port cannot be listened on.
// The returned channel is closed when the (asynchronous) termination is finished.
func (s preparedGenericAPIServer) NonBlockingRun() error {

	if s.SecureServingInfo != nil && s.Handler != nil {
		var err error
		err = s.SecureServingInfo.ServeWithListenerStopped(s.Handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
	err := s.NonBlockingRun()
	if err != nil {
		return err
	}
	return nil
}
