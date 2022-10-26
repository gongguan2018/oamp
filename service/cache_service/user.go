package cache_service

import (
	"oamp/pkg/cache"
	"strings"
)

func GetUserKeys() string {
	keys := []string{
		cache.CACHE_USER,
		"LIST",
	}
	return strings.Join(keys, "_")
}
