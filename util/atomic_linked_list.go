package util

import (
	"errors"
	"sync"
)

type NodeValue comparable

type AtomicNode[T NodeValue] struct {
	value T
	next  *AtomicNode[T]
}

type AtomicLinkedList[T NodeValue] struct {
	head *AtomicNode[T]
	size uint
	rw   *sync.RWMutex
}

func NewAtomicLinkedList[T NodeValue]() *AtomicLinkedList[T] {
	return &AtomicLinkedList[T]{rw: &sync.RWMutex{}}
}

func (l *AtomicLinkedList[T]) Add(value T) {
	l.rw.Lock()
	defer l.rw.Unlock()

	defer func() { l.size++ }()

	if l.head == nil {
		l.head = &AtomicNode[T]{value: value}
		return
	}

	current := l.head
	for current.next != nil {
		current = current.next
	}
	current.next = &AtomicNode[T]{value: value}
}

func (l *AtomicLinkedList[T]) Contains(value T) bool {
	l.rw.RLock()
	defer l.rw.RUnlock()

	current := l.head
	for current != nil {
		if current.value == value {
			return true
		}
		current = current.next
	}
	return false
}

func (l *AtomicLinkedList[T]) Size() uint {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return l.size
}

func (l *AtomicLinkedList[T]) Remove(value T) error {
	if !l.Contains(value) {
		return errors.New("value not found")
	}

	l.rw.Lock()
	defer l.rw.Unlock()

	defer func() { l.size-- }()

	if l.head.value == value {
		l.head = l.head.next
		return nil
	}

	current := l.head
	for current.next != nil {
		if current.next.value == value {
			current.next = current.next.next
			return nil
		}
		current = current.next
	}

	// should be unreachable
	return errors.New("value not found")
}

func (l *AtomicLinkedList[T]) Slice() []T {
	l.rw.RLock()
	defer l.rw.RUnlock()

	result := make([]T, 0, l.size)
	current := l.head
	for current != nil {
		result = append(result, current.value)
		current = current.next
	}
	return result
}
