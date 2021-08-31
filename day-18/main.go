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

var (
	STEPS = []v2{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}
)

func makeAllKeys(keys map[v2]uint32) uint32 {
	res := uint32(0)
	for _, k := range keys {
		res |= k
	}
	return res
}

func solve(grid map[v2]bool, doors, keys map[v2]uint32, start v2) int {
	allkeys := makeAllKeys(keys)
	visited := make(map[v3]bool)
	queue := []state{{pos: v3{start.x, start.y, 0}}}

	//log.Printf("grid: %+v", grid)
	//log.Printf("doors: %+v", doors)
	//log.Printf("keys: %+v", keys)
	//log.Printf("allKeys: %b", allkeys)

	var st state
	for {
		st, queue = queue[0], queue[1:]

		if uint32(st.pos.z)&allkeys == allkeys {
			return st.steps
		}

		visited[st.pos] = true

		for _, s := range STEPS {
			npos := v2{st.pos.x + s.x, st.pos.y + s.y}
			next := v3{npos.x, npos.y, st.pos.z}

			if !grid[npos] {
				continue
			}

			if visited[next] {
				continue
			}

			if d, ok := doors[npos]; ok {
				if uint32(next.z)&d != d {
					continue
				}
			}

			if k, ok := keys[npos]; ok {
				next.z = int(uint32(next.z) | k)
			}

			queue = append(queue, state{pos: next, steps: st.steps + 1})
		}
	}
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
