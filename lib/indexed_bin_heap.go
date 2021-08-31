package main

type IndexedBinHeap struct {
	*BinHeap
	index map[interface{}]int
}
