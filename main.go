package main

import (
	"fmt"
	"github.com/Prem-099/lru-cache/lru"
)



func main() {
	cache := lru.NewSharded[string,int](100,16)
	cache.Put("a", 1)
	cache.Put("b", 2)
	val, ok := cache.Get("a")
	fmt.Println("Value of key a : ",val,ok)
	cache.Put("c",3)
	_,ok = cache.Get("b")
	fmt.Println("b exists ",ok)
	val,ok = cache.Get("c")
	fmt.Println("C value exists ",val, ok)
}
