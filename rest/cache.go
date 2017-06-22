package rest

import (
	"crypto/md5"
	"github.com/bradfitz/gomemcache/memcache"
	"strings"
)

// Memcached connection settings
type MC struct {
	Enabled bool     `json:"enabled"`
	Hosts   []string `json:"hosts"`
	Prefix  string   `json:"prefix"`
	Client  *memcache.Client
}

// InitConnection to memcached cluster
func (mc MC) InitConnection() *memcache.Client {
	client := memcache.New(mc.Hosts...)
	return client
}

// StoreToCache save data to memecache
func (mc MC) StoreToCache(key string, data []byte) error {
	prefixed_key := strings.Join([]string{mc.Prefix, key}, "_")

	var item memcache.Item
	// item.Key = hashKey(prefixed_key)
	item.Key = prefixed_key
	item.Value = data

	// store
	err := mc.Client.Set(&item)

	return err
}

// (mc MC) GetFromCache ...
func (mc MC) GetFromCache(key string) ([]byte, error) {
	var data *memcache.Item
	prefixed_key := strings.Join([]string{mc.Prefix, key}, "_")

	// data, err := mc.Client.Get(hashKey(prefixed_key))
	data, err := mc.Client.Get(prefixed_key)
	if err != nil {
		return []byte(""), err
	}

	return data.Value, err
}

// hashKey return md5sum of prefixed_key
func hashKey(key string) string {
	h := md5.New()
	sum := h.Sum([]byte(key))

	return string(sum)
}
