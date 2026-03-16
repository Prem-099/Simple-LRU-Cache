package lru

import (
	"math/rand"
	"testing"
	"time"
)

var sinkInt int
var sinkOk bool

func BenchmarkTestGet(b *testing.B) {
	cache := New[int, int](1000)
	cache.Put(1, 1, 2*time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkOk = cache.Get(1)
	}
}

func BenchmarkTestGetShard(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	cache.Put(1, 1, 2*time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(1)
	}
}

func BenchmarkTestParallelGet(b *testing.B) {
	cache := New[string, int](1000)
	cache.Put("a", 1, 2*time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			sinkInt, sinkOk = cache.Get("a")
		}
	})
}


func BenchmarkParallelGetShard(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i, time.Minute)
	}
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			sinkInt, sinkOk = cache.Get(i)
			i++
			if i >= 1000 {
				i = 0
			}
		}
	})
}

func BenchmarkTestMixed(b *testing.B) {
	cache := New[int, int](1000)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			sinkInt, sinkOk = cache.Get(i)
			i++
		}
	})
}

func BenchmarkTestMixedShard(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			sinkInt, sinkOk = cache.Get(i)
			i++
		}
	})
}

func BenchmarkWriteHeavy(b *testing.B) {
	cache := New[int, int](1000)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			i++
		}
	})
}

func BenchmarkHeavyWriteShard(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			i++
		}
	})
}

func BenchmarkTestTtl(b *testing.B) {
	cache := New[int, int](1000)
	cache.Put(1, 1, 3*time.Second)
	time.Sleep(2 * time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d hitrate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(1)
	}
}

func BenchmarkCachedMixed(b *testing.B) {
	cache := New[int, int](1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i, 5*time.Second)
	}
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := i % 1500
		if i%3 == 0 {
			cache.Put(key, key, 5*time.Second)
		} else {
			sinkInt, sinkOk = cache.Get(key)
		}
	}
}

func BenchmarkZipfCache(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		zipf := rand.NewZipf(r, 1.2, 1, 1500)
		for p.Next() {
			key := int(zipf.Uint64())
			if key%4 == 0 {
				cache.Put(key, key, 2*time.Second)
			} else {
				cache.Get(key)
			}
		}
	})
}

// built in
func BenchmarkMapGet(b *testing.B) {
	m := make(map[int]int, 1000)
	for i := 0; i < 1000; i++ {
		m[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt = m[i%1000]
	}
}

func BenchmarkLruEviction(b *testing.B) {
	cache := New[int, int](1000)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	for i := 0; i < b.N; i++ {
		cache.Put(i, i, 2*time.Second)
	}
}

func BenchmarkLruParallelEviction(b *testing.B) {
	cache := NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			i++
		}
	})
}

func BenchmarkZipfGet(b *testing.B) {
	cache := NewSharded[int, int](10000, 256)
	for i := 0; i < 10000; i++ {
		cache.Put(i, i, time.Hour)
	}
	zipf := rand.NewZipf(rand.New(rand.NewSource(1)), 1.2, 1, 9999)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	keys := make([]int, b.N)
	for i := range keys {
		keys[i] = int(zipf.Uint64())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(keys[i])
	}
}

func BenchmarkZipfParallelGet(b *testing.B) {
	cache := NewSharded[int, int](10000, 256)
	zipf := rand.NewZipf(rand.New(rand.NewSource(1)), 1.5, 1, 9999)
	keys := make([]int, 10000)
	for i := range keys {
		keys[i] = int(zipf.Uint64())
	}
	for i := 0; i < 10000; i++ {
		cache.Put(i, i, 5*time.Second)
	}
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			sinkInt, sinkOk = cache.Get(keys[i%len(keys)])
			i++
		}
	})
}
