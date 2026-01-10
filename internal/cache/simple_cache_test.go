package cache

import (
	"testing"
	"time"
)

func TestSimpleCache_SetAndGet(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	// Test Set and Get
	cache.Set("key1", "value1")
	value, found := cache.Get("key1")

	if !found {
		t.Error("Expected to find key1")
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}
}

func TestSimpleCache_GetNonExistent(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	value, found := cache.Get("nonexistent")

	if found {
		t.Error("Expected not to find nonexistent key")
	}

	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}
}

func TestSimpleCache_Expiration(t *testing.T) {
	cache := NewSimpleCache(100 * time.Millisecond)
	defer cache.Close()

	cache.Set("expiring_key", "expiring_value")

	// Should be found immediately
	_, found := cache.Get("expiring_key")
	if !found {
		t.Error("Expected to find expiring_key immediately")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not be found after expiration
	_, found = cache.Get("expiring_key")
	if found {
		t.Error("Expected expiring_key to be expired")
	}
}

func TestSimpleCache_Delete(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	cache.Set("key_to_delete", "value")
	cache.Delete("key_to_delete")

	_, found := cache.Get("key_to_delete")
	if found {
		t.Error("Expected key_to_delete to be deleted")
	}
}

func TestSimpleCache_Clear(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Size())
	}
}

func TestSimpleCache_Size(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	if cache.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", cache.Size())
	}

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}
}

func TestSimpleCache_Concurrent(t *testing.T) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()

	// Concurrent writes
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(n int) {
			cache.Set(string(rune(n)), n)
			done <- true
		}(i)
	}

	// Wait for all writes
	for i := 0; i < 100; i++ {
		<-done
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		go func(n int) {
			cache.Get(string(rune(n)))
			done <- true
		}(i)
	}

	// Wait for all reads
	for i := 0; i < 100; i++ {
		<-done
	}

	// If we reach here without panic, concurrency is handled correctly
	t.Log("Concurrent access handled correctly")
}

func TestSimpleCache_GetStats(t *testing.T) {
	cache := NewSimpleCache(100 * time.Millisecond)
	defer cache.Close()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	stats := cache.GetStats()

	totalItems, ok := stats["total_items"].(int)
	if !ok || totalItems != 2 {
		t.Errorf("Expected total_items 2, got %v", stats["total_items"])
	}

	validItems, ok := stats["valid_items"].(int)
	if !ok || validItems != 2 {
		t.Errorf("Expected valid_items 2, got %v", stats["valid_items"])
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	stats = cache.GetStats()
	expiredItems, ok := stats["expired_items"].(int)
	if !ok || expiredItems != 2 {
		t.Errorf("Expected expired_items 2, got %v", stats["expired_items"])
	}
}

func BenchmarkSimpleCache_Set(b *testing.B) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set("bench_key", "bench_value")
	}
}

func BenchmarkSimpleCache_Get(b *testing.B) {
	cache := NewSimpleCache(1 * time.Hour)
	defer cache.Close()
	cache.Set("bench_key", "bench_value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get("bench_key")
	}
}
