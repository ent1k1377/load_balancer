package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

var emptyFunc = func() {}

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

type ServerPool struct {
	backends []*Backend
	strategy Strategy
}

type Strategy interface {
	NextBackend(backends []*Backend) (*Backend, func(), error)
}
