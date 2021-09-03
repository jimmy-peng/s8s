package server

import (
	"fmt"
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
}

func (s GenericAPIServer) ListedPaths() []string {
	return []string{}
}

// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
func (s *GenericAPIServer) PrepareRun() preparedGenericAPIServer {
	//s.delegationTarget.PrepareRun()

	return preparedGenericAPIServer{s}
}

func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
	return nil
}
