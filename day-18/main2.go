package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

type v2 struct {
	x, y int
}

type v3 struct {
	x, y, z int
}

type state struct {
	pos   v3
	steps int
}

func getGrid(input io.Reader) (map[v2]bool, map[v2]uint32, map[v2]uint32, v2) {
	grid := make(map[v2]bool)
	doors := make(map[v2]uint32)
	keys := make(map[v2]uint32)
	var start v2

	scanner := bufio.NewScanner(input)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, ch := range line {
			if ch == '#' {
				continue
			}
			pos := v2{x, y}
			if ch == '@' {
				start = pos
			} else if ch >= 'a' && ch <= 'z' {
				// key
				keys[pos] = 1 << int(ch-'a')
			} else if ch >= 'A' && ch <= 'Z' {
				doors[pos] = 1 << int(ch-'A')
			}
			grid[pos] = true
		}
		y++
	}

	err := scanner.Err()
	noerr(err)

	return grid, doors, keys, start
}

func makeAllKeys(keys map[v2]uint32) uint32 {
	res := uint32(0)
	for _, k := range keys {
		res |= k
	}
	return res
}

var (
	STEPS = []v2{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}
)

func MinDistMaxKeys(i1, i2 interface{}) bool {
	v1, v2 := i1.(*state), i2.(*state)
	if v1.steps == v2.steps {
		return popcnt(uint32(v1.pos.z)) > popcnt(uint32(v2.pos.z))
	}
	return v1.steps < v2.steps
}

func solve(grid map[v2]bool, doors, keys map[v2]uint32, start v2) int {
	allkeys := makeAllKeys(keys)
	h := NewBinHeap(MinDistMaxKeys)
	h.Push(&state{
		pos:   v3{x: start.x, y: start.y, z: 0},
		steps: 0,
	})
	visited := make(map[v3]int)
	for h.Size() > 0 {
		curr := (*state)h.Pop()
		pos := v2{curr.pos.x, curr.pos.y}
		currkeys := curr.pos.z
		if curr.pos.z == allkeys {
			return curr.steps
		}
		if k, ok := keys[pos]; ok {
			currkeys |= k
		}
	}
	return -1
}

func popcnt(i uint32) int {
	cnt := 0
	for i > 0 {
		i &= i - 1
		cnt++
	}
	return cnt
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func main() {
	file, err := os.Open("INPUT")
	noerr(err)
	defer file.Close()

	grid, doors, keys, start := getGrid(file)
	res := solve(grid, doors, keys, start)

	log.Printf("The result is: %d", res)
}
