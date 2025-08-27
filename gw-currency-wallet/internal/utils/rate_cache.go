package utils

import (
	"fmt"
	"sync"
	"time"
)

type RateCacheInterface interface {
	UpdatedRates(rates map[string]float64)
	GetRate(from, to string) (float64, bool)
	GetAllRates() (map[string]float64, bool)
}

type RateCache struct {
	Rates     map[string]float64
	UpdatedAt time.Time
	TTL       time.Duration
	mu        sync.RWMutex
}

func NewRateCache(ttl time.Duration) *RateCache {
	return &RateCache{
		Rates:     make(map[string]float64),
		UpdatedAt: time.Time{},
		TTL:       ttl,
		mu:        sync.RWMutex{},
	}
}

func (c *RateCache) UpdatedRates(rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Rates = rates
	c.UpdatedAt = time.Now()
}

func (c *RateCache) GetRate(from, to string) (float64, bool) {
	key := fmt.Sprintf("From %s to %s", from, to)

	c.mu.RLock()
	defer c.mu.RUnlock()

	if rate, ok := c.Rates[key]; ok {
		if time.Since(c.UpdatedAt) <= c.TTL {
			return rate, true
		}
	}

	revKey := fmt.Sprintf("From %s to %s", to, from)
	if rate, ok := c.Rates[revKey]; ok {
		if time.Since(c.UpdatedAt) <= c.TTL {
			return 1 / rate, true
		}
	}

	return 0, false
}

func (c *RateCache) GetAllRates() (map[string]float64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.UpdatedAt) > c.TTL || len(c.Rates) == 0 {
		return nil, false
	}

	return c.Rates, true
}
