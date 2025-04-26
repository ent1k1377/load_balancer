package balancer

import (
	"sync"
	"testing"
)

func TestBackend_NewBackend(t *testing.T) {
	backend, err := NewBackend("localhost:9001")
	if backend == nil || err != nil {
		t.Fatalf("error creating backend: %v", err)
	}

	if !backend.Alive {
		t.Fatalf("backend should be alive")
	}
}

func TestBackend_NewBackend_WrongURL(t *testing.T) {
	back, err := NewBackend("http://[wrong url]")
	if err == nil {
		t.Fatalf("should have errored: %v", back.URL.String())
	}
}

func TestBackend_ConcurrentAccess(t *testing.T) {
	b, err := NewBackend("http://localhost:9001")
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}

	const (
		readers = 100
		writers = 50
	)

	var wg sync.WaitGroup

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_ = b.IsAlive()
			}
		}()
	}

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				b.SetAlive(id%2 == 0)
			}
		}(i)
	}

	wg.Wait()
}
