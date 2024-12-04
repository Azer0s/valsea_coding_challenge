package util

import "testing"

func TestAtomicLinkedList_Add(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("size should be 3, but got %d", list.Size())
	}
}

func TestAtomicLinkedList_Contains(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if !list.Contains(1) {
		t.Errorf("list should contain 1")
	}

	if !list.Contains(2) {
		t.Errorf("list should contain 2")
	}

	if !list.Contains(3) {
		t.Errorf("list should contain 3")
	}

	if list.Contains(4) {
		t.Errorf("list should not contain 4")
	}
}

func TestAtomicLinkedList_Size(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("size should be 3, but got %d", list.Size())
	}
}

func TestAtomicLinkedList_AddString(t *testing.T) {
	list := NewAtomicLinkedList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if list.Size() != 3 {
		t.Errorf("size should be 3, but got %d", list.Size())
	}
}

func TestAtomicLinkedList_ContainsString(t *testing.T) {
	list := NewAtomicLinkedList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if !list.Contains("a") {
		t.Errorf("list should contain a")
	}

	if !list.Contains("b") {
		t.Errorf("list should contain b")
	}

	if !list.Contains("c") {
		t.Errorf("list should contain c")
	}

	if list.Contains("d") {
		t.Errorf("list should not contain d")
	}
}

func TestAtomicLinkedList_Remove(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if err := list.Remove(2); err != nil {
		t.Errorf("error should be nil, but got %v", err)
	}

	if list.Size() != 2 {
		t.Errorf("size should be 2, but got %d", list.Size())
	}

	if list.Contains(2) {
		t.Errorf("list should not contain 2")
	}

	sl := list.Slice()
	if sl[0] != 1 || sl[1] != 3 {
		t.Errorf("list should be [1, 3], but got %v", list.Slice())
	}
}

func TestAtomicLinkedList_RemoveNonExistent(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if err := list.Remove(4); err == nil {
		t.Errorf("error should not be nil")
	}

	if list.Size() != 3 {
		t.Errorf("size should be 3, but got %d", list.Size())
	}

	if !list.Contains(2) {
		t.Errorf("list should contain 2")
	}
}

func TestAtomicLinkedList_Slice(t *testing.T) {
	list := NewAtomicLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	sl := list.Slice()
	if len(sl) != 3 {
		t.Errorf("slice length should be 3, but got %d", len(sl))
	}

	if sl[0] != 1 || sl[1] != 2 || sl[2] != 3 {
		t.Errorf("slice should be [1, 2, 3], but got %v", sl)
	}
}

func TestAtomicLinkedList_ConcurrentAccess(t *testing.T) {
	doneChan1 := make(chan struct{})
	doneChan2 := make(chan struct{})

	list := NewAtomicLinkedList[int]()
	go func() {
		for i := 0; i < 1000; i++ {
			list.Add(i)
		}
		doneChan1 <- struct{}{}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			list.Add(i)
		}
		doneChan2 <- struct{}{}
	}()

	<-doneChan1
	<-doneChan2

	if list.Size() != 2000 {
		t.Errorf("size should be 2000, but got %d", list.Size())
	}
}
