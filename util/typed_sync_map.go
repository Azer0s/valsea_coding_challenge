package util

import "sync"

type TypedSyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewTypedSyncMap[K comparable, V any]() *TypedSyncMap[K, V] {
	return &TypedSyncMap[K, V]{
		m: sync.Map{},
	}
}

func (t *TypedSyncMap[K, V]) Size() int {
	size := 0
	t.m.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

func (t *TypedSyncMap[K, V]) Range(f func(key K, value V) bool) {
	t.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (t *TypedSyncMap[K, V]) Delete(key K) {
	t.m.Delete(key)
}

func (t *TypedSyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := t.m.Load(key)
	if !ok {
		return
	}
	value = v.(V)
	return
}

func (t *TypedSyncMap[K, V]) Store(key K, value V) {
	t.m.Store(key, value)
}
