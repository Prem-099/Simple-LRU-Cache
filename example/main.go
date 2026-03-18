package main

import (
	"fmt"
	"time"

	"github.com/Prem-099/lru-cache"
)

func main() {
	cache := lru.New[string, int](0)
	cache.Put("a", 1, 1*time.Second)
	cache.Put("b", 2, 1*time.Second)
	fmt.Println(cache.Len())
	val, ok := cache.Get("a")
	fmt.Println("key a value is:", val, ok)
	ook := cache.Delete("a")
	val, ok = cache.Get("a")
	if ok {
		fmt.Println("key a value is:", val)
	} else {
		fmt.Println("Key deleted", ook)
	}
	stats := cache.Stats()
	fmt.Println(stats)
	cache.Clear()
	val, ok = cache.Get("a")
	if !ok {
		fmt.Println("cache cleared")
	}

}
