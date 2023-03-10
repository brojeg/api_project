package model

import (
	"sync"
	"testing"

	"golang.org/x/time/rate"
)

func TestNewIPRateLimiter(t *testing.T) {
	// Test case 1: Verify that the ips map is initialized to an empty map.
	i := newIPRateLimiter(rate.Limit(1), 1)
	if len(i.ips) != 0 {
		t.Errorf("Test case 1: expected ips map to be empty, but got %v", i.ips)
	}

	// Test case 2: Verify that the mu mutex is initialized to a non-nil value.
	if i.mu == nil {
		t.Errorf("Test case 2: expected mu mutex to be non-nil, but got nil")
	}

	// Test case 3: Verify that the r rate and b burst size are set to the values passed in as arguments.
	if i.r != rate.Limit(1) {
		t.Errorf("Test case 3: expected r rate to be %v, but got %v", rate.Limit(1), i.r)
	}
	if i.b != 1 {
		t.Errorf("Test case 3: expected b burst size to be %v, but got %v", 1, i.b)
	}

}

func TestIPRateLimiter_AddIP(t *testing.T) {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   rate.Limit(1),
		b:   1,
	}
	limiter := i.addIP("127.0.0.1")
	if limiter == nil {
		t.Errorf("Test case 1: expected addIP() method to return a non-nil limiter, but got nil")
	}

	_, ok := i.ips["127.0.0.1"]
	if !ok {
		t.Errorf("Test case 2: expected ips map to contain key %q, but it did not", "127.0.0.1")
	}

	limiter, ok = i.ips["127.0.0.1"]
	if !ok {
		t.Fatalf("Test case 3: expected ips map to contain key %q, but it did not", "127.0.0.1")
	}
	if limiter == nil {
		t.Errorf("Test case 3: expected value associated with key %q to be non-nil, but got nil", "127.0.0.1")
	}

}

func TestIPRateLimiter_GetLimiter(t *testing.T) {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   rate.Limit(1),
		b:   1,
	}
	limiter := i.getLimiter("127.0.0.1")
	if limiter == nil {
		t.Errorf("Test case 1: expected getLimiter() method to return a non-nil limiter, but got nil")
	}

	_, ok := i.ips["127.0.0.1"]
	if !ok {
		t.Errorf("Test case 2: expected ips map to contain key %q, but it did not", "127.0.0.1")
	}

	limiter, ok = i.ips["127.0.0.1"]
	if !ok {
		t.Fatalf("Test case 3: expected ips map to contain key %q, but it did not", "127.0.0.1")
	}
	if limiter == nil {
		t.Errorf("Test case 3: expected value associated with key %q to be non-nil, but got nil", "127.0.0.1")
	}

	if limiter.Limit() != i.r {
		t.Errorf("Test case 4: expected rate.Limiter limit to be %v, but got %v", i.r, limiter.Limit())
	}
	if limiter.Burst() != i.b {
		t.Errorf("Test case 4: expected rate.Limiter burst size to be %v, but got %v", i.b, limiter.Burst())
	}
}
