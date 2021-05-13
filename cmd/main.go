package main

import (
	"fmt"
	"time"

	"github.com/amartery/LRUcache/pkg/lru"
)

func main() {
	cache := lru.NewCache(3)
	cache.ShowCurrentCache() // {}

	cache.Put(1, "str1")
	cache.ShowCurrentCache() // {1: “str1"}

	cache.Put(2, "str2")
	cache.ShowCurrentCache() // {1: “str1", 2: “str2"}

	cache.Put(3, "str3")
	cache.ShowCurrentCache() // {1: “str1", 2: “str2", 3: “str3"}

	cache.Get(3)
	cache.Get(2)
	cache.Get(1)
	cache.Get(3)

	cache.Put(4, "str4")
	cache.ShowCurrentCache() // {1: “str1", 3: “str2", 4: “str4"}

	// lru with ttl
	fmt.Println("lru with ttl")

	cacheTTL := lru.NewTTLCache(3, 5*time.Second)
	cacheTTL.ShowCurrentCache() // {}

	cacheTTL.Put(1, "str1")
	cacheTTL.ShowCurrentCache() // {1: “str1"}

	cacheTTL.Put(2, "str2")
	cacheTTL.ShowCurrentCache() // {1: “str1", 2: “str2"}

	cacheTTL.Put(3, "str3")
	cacheTTL.ShowCurrentCache() // {1: “str1", 2: “str2", 3: “str3"}

	time.Sleep(6 * time.Second)
	cacheTTL.Get(3)
	cacheTTL.Get(2)
	cacheTTL.Get(1)
	cacheTTL.Get(3)

	cacheTTL.Put(4, "str4")
	cacheTTL.ShowCurrentCache() // {4: “str4"}
}
