package mux

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

// PathRecorderMux wraps a mux object and records the registered exposedPaths.
type PathRecorderMux struct {
	// name is used for logging so you can trace requests through
	name string

	lock            sync.Mutex
	notFoundHandler http.Handler
	pathToHandler   map[string]http.Handler
	prefixToHandler map[string]http.Handler

	// mux stores a pathHandler and is used to handle the actual serving.
	// Turns out, we want to accept trailing slashes, BUT we don't care about handling
	// everything under them.  This does exactly matches only unless its explicitly requested to
	// do something different
	mux atomic.Value

	// exposedPaths is the list of paths that should be shown at /
	exposedPaths []string

	// pathStacks holds the stacks of all registered paths.  This allows us to show a more helpful message
	// before the "http: multiple registrations for %s" panic.
	pathStacks map[string]string
}

// pathHandler is an http.Handler that will satisfy requests first by exact match, then by prefix,
// then by notFoundHandler
type pathHandler struct {
	// muxName is used for logging so you can trace requests through
	muxName string

	// pathToHandler is a map of exactly matching request to its handler
	pathToHandler map[string]http.Handler

	// this has to be sorted by most slashes then by length
	prefixHandlers []prefixHandler

	// notFoundHandler is the handler to use for satisfying requests with no other match
	notFoundHandler http.Handler
}

// prefixHandler holds the prefix it should match and the handler to use
type prefixHandler struct {
	// prefix is the prefix to test for a request match
	prefix string
	// handler is used to satisfy matching requests
	handler http.Handler
}

// NewPathRecorderMux creates a new PathRecorderMux
func NewPathRecorderMux(name string) *PathRecorderMux {
	ret := &PathRecorderMux{
		name:            name,
		pathToHandler:   map[string]http.Handler{},
		prefixToHandler: map[string]http.Handler{},
		mux:             atomic.Value{},
		exposedPaths:    []string{},
		pathStacks:      map[string]string{},
	}

	ret.mux.Store(&pathHandler{notFoundHandler: http.NotFoundHandler()})
	return ret
}

// NotFoundHandler sets the handler to use if there's no match for a give path
func (m *PathRecorderMux) NotFoundHandler(notFoundHandler http.Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.notFoundHandler = notFoundHandler

	m.refreshMuxLocked()
}

// refreshMuxLocked creates a new mux and must be called while locked.  Otherwise the view of handlers may
// not be consistent
func (m *PathRecorderMux) refreshMuxLocked() {
	newMux := &pathHandler{
		muxName:         m.name,
		pathToHandler:   map[string]http.Handler{},
		prefixHandlers:  []prefixHandler{},
		notFoundHandler: http.NotFoundHandler(),
	}
	if m.notFoundHandler != nil {
		newMux.notFoundHandler = m.notFoundHandler
	}
	for path, handler := range m.pathToHandler {
		newMux.pathToHandler[path] = handler
	}
	/*
		keys := sets.StringKeySet(m.prefixToHandler).List()
		sort.Sort(sort.Reverse(byPrefixPriority(keys)))
		for _, prefix := range keys {
			newMux.prefixHandlers = append(newMux.prefixHandlers, prefixHandler{
				prefix:  prefix,
				handler: m.prefixToHandler[prefix],
			})
		}
	*/

	m.mux.Store(newMux)
}

// ServeHTTP makes it an http.Handler
func (m *PathRecorderMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.Load().(*pathHandler).ServeHTTP(w, r)
}

// ServeHTTP makes it an http.Handler
func (h *pathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if exactHandler, ok := h.pathToHandler[r.URL.Path]; ok {
		//klog.V(5).Infof("%v: %q satisfied by exact match", h.muxName, r.URL.Path)
		exactHandler.ServeHTTP(w, r)
		return
	}

	for _, prefixHandler := range h.prefixHandlers {
		if strings.HasPrefix(r.URL.Path, prefixHandler.prefix) {
			//klog.V(5).Infof("%v: %q satisfied by prefix %v", h.muxName, r.URL.Path, prefixHandler.prefix)
			prefixHandler.handler.ServeHTTP(w, r)
			return
		}
	}

	//klog.V(5).Infof("%v: %q satisfied by NotFoundHandler", h.muxName, r.URL.Path)
	h.notFoundHandler.ServeHTTP(w, r)
}
