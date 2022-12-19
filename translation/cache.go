package translation

import "github.com/patrickmn/go-cache"

var c *cache.Cache

func InitCache() {
	c = cache.New(cache.DefaultExpiration, cache.DefaultExpiration)
}
