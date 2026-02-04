package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/javaquery/unosdk/pkg/models"
)

// Cache manages cached provider information
type Cache struct {
	cachePath string
	ttl       time.Duration
}

// CacheEntry represents a cached entry
type CacheEntry struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewCache creates a new cache
func NewCache(cachePath string, ttl time.Duration) *Cache {
	return &Cache{
		cachePath: cachePath,
		ttl:       ttl,
	}
}

// Get retrieves a cached entry
func (c *Cache) Get(key string) (interface{}, bool) {
	filePath := filepath.Join(c.cachePath, key+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if time.Since(entry.Timestamp) > c.ttl {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a cache entry
func (c *Cache) Set(key string, data interface{}) error {
	os.MkdirAll(c.cachePath, 0755)

	entry := CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	filePath := filepath.Join(c.cachePath, key+".json")
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	return os.RemoveAll(c.cachePath)
}

// Provider cache methods

// GetProviderVersions retrieves cached provider versions
func (c *Cache) GetProviderVersions(provider string) ([]string, bool) {
	data, ok := c.Get("provider_" + provider)
	if !ok {
		return nil, false
	}
	versions, ok := data.([]string)
	return versions, ok
}

// SetProviderVersions caches provider versions
func (c *Cache) SetProviderVersions(provider string, versions []string) error {
	return c.Set("provider_"+provider, versions)
}

// GetProviderInfo retrieves cached provider info
func (c *Cache) GetProviderInfo(provider string) (*models.ProviderInfo, bool) {
	data, ok := c.Get("info_" + provider)
	if !ok {
		return nil, false
	}
	info, ok := data.(*models.ProviderInfo)
	return info, ok
}

// SetProviderInfo caches provider info
func (c *Cache) SetProviderInfo(provider string, info *models.ProviderInfo) error {
	return c.Set("info_"+provider, info)
}
