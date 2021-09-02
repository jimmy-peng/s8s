package server

import "time"

type DelegationTarget interface {
	ListedPaths() []string
}

type emptyDelegate struct {
}

func (s emptyDelegate) ListedPaths() []string {
	return []string{}
}

func NewEmptyDelegate() DelegationTarget {
	return emptyDelegate{}
}

// GenericAPIServer contains state for a Kubernetes cluster api server.
type GenericAPIServer struct {
	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration
}
