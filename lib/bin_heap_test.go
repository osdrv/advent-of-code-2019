package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	MaxHeapInt Comparator = func(v1, v2 interface{}) bool {
		i1, i2 := v1.(int), v2.(int)
		return i1 > i2
	}

	MinHeapInt Comparator = func(v1, v2 interface{}) bool {
		i1, i2 := v1.(int), v2.(int)
		return i1 < i2
	}
)

func TestBinHeap_Size_Empty(t *testing.T) {
	h := NewBinHeap(MaxHeapInt)
	assert.Equal(t, 0, h.Size())
}

func TestBinHeap_MinHeap_Push(t *testing.T) {
	h := NewBinHeap(MinHeapInt)
	h.Push(1)
	h.Push(2)
	h.Push(3)
	assert.Equal(t, 1, h.Peek())
}

func TestBinHeap_MaxHeap_Push(t *testing.T) {
	h := NewBinHeap(MaxHeapInt)
	h.Push(1)
	h.Push(2)
	h.Push(3)
	assert.Equal(t, 3, h.Peek())
}

func TestBinHeap_MinHeap_Pop(t *testing.T) {
	h := NewBinHeap(MinHeapInt)
	h.Push(1)
	h.Push(2)
	h.Push(3)
	h.Pop()
	assert.Equal(t, 2, h.Peek())
}

func TestBinHeap_MaxHeap_Pop(t *testing.T) {
	h := NewBinHeap(MaxHeapInt)
	h.Push(1)
	h.Push(2)
	h.Push(3)
	h.Pop()
	assert.Equal(t, 2, h.Peek())
}
