package prefixcase

import (
	"strings"
)

type PrefixCache[T any] struct {
	Sep   string   `json:"Sep"`
	Cache Cache[T] `json:"Cache"`
}

type Cache[T any] struct {
	Data  []T                 `json:"Data"`
	Child map[string]Cache[T] `json:"Child"`
}

func (c *PrefixCache[T]) WithValue(key string, value T) {
	c.Cache = setCacheValue(c.Cache, strings.Split(key, c.Sep), value)
}

func (c *PrefixCache[T]) Value(key string) []T {
	return getCacheValue(c.Cache, strings.Split(key, c.Sep))
}

func NewPrefixCase[T any](sep string) PrefixCache[T] {
	return PrefixCache[T]{
		Sep:   sep,
		Cache: Cache[T]{},
	}
}

func setCacheValue[T any](cache Cache[T], path []string, value T) Cache[T] {
	if len(path) == 0 {
		if cache.Data == nil {
			cache.Data = []T{}
		}
		cache.Data = append(cache.Data, value)
		return cache
	}

	part := path[0]
	if cache.Child == nil {
		cache.Child = make(map[string]Cache[T])
	}
	subCache, exists := cache.Child[part]
	if !exists {
		subCache = Cache[T]{}
	}
	subCache = setCacheValue(subCache, path[1:], value)
	cache.Child[part] = subCache
	return cache
}

func getCacheValue[T any](cache Cache[T], path []string) []T {
	if len(path) == 0 {
		return cache.Data
	}

	part := path[0]
	subCache, exists := cache.Child[part]
	if !exists {
		return nil
	}
	return getCacheValue(subCache, path[1:])
}
