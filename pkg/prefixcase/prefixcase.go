package prefixcase

import (
	"strings"
)

type PrefixCache[T any] struct {
	sep   string
	cache Cache[T]
}

type Cache[T any] struct {
	data  []T
	child map[string]Cache[T]
}

func (c *PrefixCache[T]) WithValue(key string, value T) {
	c.cache = setCacheValue(c.cache, strings.Split(key, c.sep), value)
}

func (c *PrefixCache[T]) Value(key string) []T {
	return getCacheValue(c.cache, strings.Split(key, c.sep))
}

func NewPrefixCase[T any](sep string) PrefixCache[T] {
	return PrefixCache[T]{
		sep:   sep,
		cache: Cache[T]{},
	}
}

func setCacheValue[T any](cache Cache[T], path []string, value T) Cache[T] {
	if len(path) == 0 {
		if cache.data == nil {
			cache.data = []T{}
		}
		cache.data = append(cache.data, value)
		return cache
	}

	part := path[0]
	if cache.child == nil {
		cache.child = make(map[string]Cache[T])
	}
	subCache, exists := cache.child[part]
	if !exists {
		subCache = Cache[T]{}
	}
	subCache = setCacheValue(subCache, path[1:], value)
	cache.child[part] = subCache
	return cache
}

func getCacheValue[T any](cache Cache[T], path []string) []T {
	if len(path) == 0 {
		return cache.data
	}

	part := path[0]
	subCache, exists := cache.child[part]
	if !exists {
		return nil
	}
	return getCacheValue(subCache, path[1:])
}
