package balancer

import "testing"

var defaultBackends = []string{"localhost:9001", "localhost:9002", "localhost:9003"}

func createTestPool(t *testing.T, backends []string, strategy Strategy) *ServerPool {
	t.Helper()
	pool := NewServerPool(strategy)

	for i, rawURL := range backends {
		back, err := NewBackend(rawURL)
		if err != nil {
			t.Fatalf("failed to create backend, id - %v, url - %v: %v", i, rawURL, err)
		}

		pool.addBackend(back)
	}

	return pool
}

func TestServerPool_GetNextBackend(t *testing.T) {
	pool := createTestPool(t, defaultBackends, NewRoundRobinStrategy())

	currentBack, _, _ := pool.strategy.NextBackend(pool.backends)
	if currentBack == nil || currentBack != pool.backends[1] {
		t.Errorf("expected backend to be %v, got %v", pool.backends[1], currentBack)
	}
}

func TestServerPool_RoundRobin(t *testing.T) {
	pool := createTestPool(t, defaultBackends, NewRoundRobinStrategy())
	pool.backends[1].SetAlive(false)

	getBack1, _, _ := pool.strategy.NextBackend(pool.backends)
	getBack2, _, _ := pool.strategy.NextBackend(pool.backends)
	getBack3, _, _ := pool.strategy.NextBackend(pool.backends)

	if getBack1 == nil || getBack1 != pool.backends[2] {
		t.Errorf("expected backend to be %v, got %v", pool.backends[2], getBack1)
	}

	if getBack2 == nil || getBack2 != pool.backends[0] {
		t.Errorf("expected backend to be %v, got %v", pool.backends[0], getBack2)
	}

	if getBack3 == nil || getBack3 != pool.backends[2] {
		t.Errorf("expected backend to be %v, got %v", pool.backends[2], getBack3)
	}
}

func TestServerPool_GetNextBackend_EmptyPool(t *testing.T) {
	pool := NewServerPool(NewRoundRobinStrategy())
	currentBack, _, _ := pool.strategy.NextBackend(pool.backends)
	if currentBack != nil {
		t.Errorf("expected backend to be nil, got %v", currentBack)
	}
}

func TestServerPool_GetNextBackend_ZeroAliveBackends(t *testing.T) {
	pool := createTestPool(t, defaultBackends, NewRoundRobinStrategy())

	pool.backends[0].SetAlive(false)
	pool.backends[1].SetAlive(false)
	pool.backends[2].SetAlive(false)

	back, _, _ := pool.strategy.NextBackend(pool.backends)
	if back != nil {
		t.Errorf("expected backend to be nil, got %v", back)
	}
}

//func TestServerPool_GetNextBackend_EmptyStrategy(t *testing.T) {
//	pool := createTestPool(t, defaultBackends, nil)
//
//	back, _, _ := pool.strategy.NextBackend(pool.backends)
//	if back != nil {
//		t.Errorf("expected backend to be nil, got %v", back)
//	}
//}
