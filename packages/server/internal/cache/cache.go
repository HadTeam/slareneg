package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) bool
	Clear()
	GetStats() CacheStats
}

type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	ItemCount   int       `json:"itemCount"`
	HitRate     float64   `json:"hitRate"`
	LastCleanup time.Time `json:"lastCleanup"`
}

type InMemoryCache struct {
	data            map[string]*CacheItem
	mutex           sync.RWMutex
	hits            int64
	misses          int64
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

func NewInMemoryCache(cleanupInterval time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		data:            make(map[string]*CacheItem),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
	}

	go cache.startCleanup()
	return cache
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.data[key] = &CacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return nil
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		c.misses++
		return nil, false
	}

	if item.IsExpired() {
		c.misses++
		go func() {
			c.mutex.Lock()
			delete(c.data, key)
			c.mutex.Unlock()
		}()
		return nil, false
	}

	c.hits++
	return item.Value, true
}

func (c *InMemoryCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, exists := c.data[key]
	if exists {
		delete(c.data, key)
	}
	return exists
}

func (c *InMemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*CacheItem)
	c.hits = 0
	c.misses = 0
}

func (c *InMemoryCache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return CacheStats{
		Hits:      c.hits,
		Misses:    c.misses,
		ItemCount: len(c.data),
		HitRate:   hitRate,
	}
}

func (c *InMemoryCache) Stop() {
	close(c.stopCleanup)
}

func (c *InMemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
		}
	}
}

type CacheService struct {
	cache Cache
}

func NewCacheService(cache Cache) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

func (s *CacheService) SetUserSession(userID string, sessionData interface{}) error {
	key := fmt.Sprintf("session:%s", userID)
	return s.cache.Set(key, sessionData, 24*time.Hour)
}

func (s *CacheService) GetUserSession(userID string) (interface{}, bool) {
	key := fmt.Sprintf("session:%s", userID)
	return s.cache.Get(key)
}

func (s *CacheService) DeleteUserSession(userID string) bool {
	key := fmt.Sprintf("session:%s", userID)
	return s.cache.Delete(key)
}

func (s *CacheService) SetGameState(gameID string, gameState interface{}) error {
	key := fmt.Sprintf("game:%s", gameID)
	return s.cache.Set(key, gameState, 2*time.Hour)
}

func (s *CacheService) GetGameState(gameID string) (interface{}, bool) {
	key := fmt.Sprintf("game:%s", gameID)
	return s.cache.Get(key)
}

func (s *CacheService) DeleteGameState(gameID string) bool {
	key := fmt.Sprintf("game:%s", gameID)
	return s.cache.Delete(key)
}

func (s *CacheService) SetPlayerConnection(playerID string, connectionInfo interface{}) error {
	key := fmt.Sprintf("connection:%s", playerID)
	return s.cache.Set(key, connectionInfo, 30*time.Minute)
}

func (s *CacheService) GetPlayerConnection(playerID string) (interface{}, bool) {
	key := fmt.Sprintf("connection:%s", playerID)
	return s.cache.Get(key)
}

func (s *CacheService) DeletePlayerConnection(playerID string) bool {
	key := fmt.Sprintf("connection:%s", playerID)
	return s.cache.Delete(key)
}

func (s *CacheService) GetCacheStats() CacheStats {
	return s.cache.GetStats()
}

func (s *CacheService) SetJSON(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.cache.Set(key, string(jsonData), ttl)
}

func (s *CacheService) GetJSON(key string, dest interface{}) (bool, error) {
	data, exists := s.cache.Get(key)
	if !exists {
		return false, nil
	}

	jsonStr, ok := data.(string)
	if !ok {
		return false, fmt.Errorf("cached data is not a JSON string")
	}

	err := json.Unmarshal([]byte(jsonStr), dest)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return true, nil
}
