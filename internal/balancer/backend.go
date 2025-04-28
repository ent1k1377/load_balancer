package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(rawURL string) (*Backend, error) {
	logger.Infof("Creating backend for URL: %s", rawURL)

	serverUrl, err := url.Parse(rawURL)
	if err != nil {
		logger.Errorf("Failed to parse url: %v", err)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	newBackend := &Backend{
		URL:          serverUrl,
		Alive:        true,
		ReverseProxy: proxy,
	}
	newBackend.ReverseProxy.ErrorHandler = ErrorHandler(newBackend)

	logger.Infof("Successfully created backend for URL: %s", rawURL)
	return newBackend, nil
}

func ErrorHandler(b *Backend) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Errorf("Backend %s is down: %v", b.URL, err)

		b.Alive = false

		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.Alive != alive {
		logger.Infof("Backend %s changed alive status to: %v", b.URL, alive)
	}

	b.Alive = alive
}

func (b *Backend) checkAlive() bool {
	timeout := time.Second * 1
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		logger.Warnf("Backend %s is down: %v", b.URL, err)
		return false
	}
	defer conn.Close()

	logger.Infof("Backend %s is alive", b.URL)
	return true
}
