package main

type Comparator func(v1, v2 interface{}) bool

type BinHeap struct {
	elements   []interface{}
	comparator Comparator
	index      map[interface{}]int
}

func NewBinHeap(cmp Comparator) *BinHeap {
	return &BinHeap{
		elements:   make([]interface{}, 0, 1),
		comparator: cmp,
		index:      make(map[interface{}]int),
	}
}

func (h *BinHeap) Size() int {
	return len(h.elements)
}

func (h *BinHeap) Peek() interface{} {
	if h.Size() == 0 {
		return nil
	}
	return h.elements[0]
}

func (h *BinHeap) Pop() interface{} {
	if (h.Size()) == 0 {
		return nil
	}
	last := len(h.elements) - 1
	h.index[h.elements[0]] = last
	h.index[h.elements[last]] = 0
	h.elements[0], h.elements[last] = h.elements[last], h.elements[0]
	res := h.elements[last]
	delete(h.index, h.elements[last])
	h.elements = h.elements[:last]
	h.reheapDown(0)
	return res
}

func (h *BinHeap) Push(v interface{}) {
	last := len(h.elements)
	h.elements = append(h.elements, v)
	h.index[h.elements[last]] = last
	h.reheapUp(last)
}

func (h *BinHeap) ReheapAt(ix int) {
	h.reheapUp(ix)
	h.reheapDown(ix)
}

func (h *BinHeap) Find(v interface{}) int {
	ix, ok := h.index[v]
	if !ok {
		return -1
	}
	return ix
}

func (h *BinHeap) reheapUp(ix int) {
	ptr := ix
	for ptr > 0 {
		parent := (ix - 1) / 2
		if !h.compare(ptr, parent) {
			break
		}
		h.index[h.elements[parent]] = ptr
		h.index[h.elements[ptr]] = parent
		h.elements[ptr], h.elements[parent] = h.elements[parent], h.elements[ptr]
		ptr = parent
	}
}

func (h *BinHeap) reheapDown(ix int) {
	ptr := ix
	for ptr < len(h.elements) {
		left, right := ptr*2+1, ptr*2+2
		next := ptr
		if left < len(h.elements) {
			if h.compare(left, next) {
				next = left
			}
		}
		if right < len(h.elements) {
			if h.compare(right, next) {
				next = right
			}
		}
		if next == ptr {
			break
		}
		h.index[h.elements[next]] = ptr
		h.index[h.elements[ptr]] = next
		h.elements[ptr], h.elements[next] = h.elements[next], h.elements[ptr]
		ptr = next
	}
}

// compares elements at ix1 and ix2 and returns true if the one at ix1 is < than
// the one at ix2
func (h *BinHeap) compare(ix1, ix2 int) bool {
	return h.comparator(h.elements[ix1], h.elements[ix2])
}
