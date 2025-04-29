package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend представляет собой backend-сервер, на который будут перенаправляться запросы.
// Он включает URL, статус доступности (Alive), а также обратный прокси для маршрутизации запросов.
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// ServerPool управляет набором backend'ов и выбранной стратегией балансировки нагрузки
type ServerPool struct {
	backends []*Backend
	strategy Strategy
}

// Strategy — интерфейс, который определяет стратегию выбора следующего backend'а.
type Strategy interface {
	// NextBackend выбирает следующий доступный backend из списка backends.
	// Возвращает backend, функцию для выполнения дополнительной логики после обработки запроса
	// и ошибку, если backend не найден.
	NextBackend(backends []*Backend) (*Backend, func(), error)
}

// emptyFunc — функция-заглушка, используемая, когда не требуется выполнение дополнительной логики.
var emptyFunc = func() {}
