package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type v2 struct {
	x, y int
}

type vertex struct {
	ch  byte
	pos v2
}

func (v vertex) String() string {
	return fmt.Sprintf("{ch: %c, pos: %+v}", v.ch, v.pos)
}

type transit struct {
	v vertex
	k uint32
	l int
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func isDoor(ch byte) bool {
	return ch >= 'A' && ch <= 'Z'
}

func isKey(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

func isStart(ch byte) bool {
	return ch == '@'
}

func isWall(ch byte) bool {
	return ch == '#'
}

func isSomething(ch byte) bool {
	return isKey(ch) || isDoor(ch) || isStart(ch)
}

func getAdjMatrix(input io.Reader) (map[vertex][]transit, vertex, uint32) {
	data, err := ioutil.ReadAll(input)
	data = bytes.TrimRight(data, "\n\t\r")
	noerr(err)
	ss := strings.Split(string(data), "\n")
	bs := make([][]byte, 0, len(ss))
	for _, s := range ss {
		bs = append(bs, []byte(s))
	}

	q := make([]vertex, 0, 1)
	var start vertex
	var allKeys uint32

	for y := 0; y < len(bs); y++ {
		for x := 0; x < len(bs[y]); x++ {
			ch := bs[y][x]
			v := vertex{
				ch:  ch,
				pos: v2{x, y},
			}
			if isStart(ch) {
				start = v
			}
			if isKey(ch) {
				allKeys |= 1 << int(ch-'a')
			}
			if isSomething(ch) {
				q = append(q, v)
			}
		}
	}

	matrix := make(map[vertex][]transit)
	var v vertex
	for len(q) > 0 {
		v, q = q[0], q[1:]
		matrix[v] = calcTransits(bs, v)
	}

	return matrix, start, allKeys
}

var (
	STEPS = []v2{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}
)

func calcTransits(bs [][]byte, from vertex) []transit {
	q := make([]transit, 0, 1)
	ts := make([]transit, 0, 1)
	q = append(q, transit{
		v: from,
		k: 0,
		l: 0,
	})
	vs := make(map[v2]bool)
	var h transit
	for len(q) > 0 {
		h, q = q[0], q[1:]
		vs[h.v.pos] = true
		if isDoor(h.v.ch) {
			h.k |= 1 << int(h.v.ch-'A')
		}
		if h.v != from && isKey(h.v.ch) {
			ts = append(ts, h)
		}
		for _, s := range STEPS {
			np := v2{x: h.v.pos.x + s.x, y: h.v.pos.y + s.y}
			npch := bs[np.y][np.x]
			if vs[np] {
				continue
			}
			if np.x < 0 || np.y < 0 || np.x >= len(bs[0]) || np.y >= len(bs) {
				continue
			}
			if isWall(npch) {
				continue
			}
			q = append(q, transit{
				v: vertex{
					ch:  npch,
					pos: np,
				},
				k: h.k,
				l: h.l + 1,
			})
		}
	}
	return ts
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type memoItem struct {
	p v2
	k uint32
}

func solve(adj map[vertex][]transit, start vertex, allKeys uint32) int {
	var visit func(vertex, uint32, int) int
	memo := make(map[memoItem]int)
	minSoFar := 999999

	visit = func(v vertex, keys uint32, dist int) int {
		if dist >= minSoFar {
			return -1
		}
		if isKey(v.ch) {
			keys |= 1 << int(v.ch-'a')
		}
		mi := memoItem{
			p: v.pos,
			k: keys,
		}
		if d, ok := memo[mi]; ok {
			if dist >= d {
				return -1
			}
		}
		memo[mi] = dist
		if keys == allKeys {
			log.Printf("Visiting vertex: %c, keys: %b, dist: %d\n", v.ch, keys, dist)
			log.Println("all keys found")
			minSoFar = min(minSoFar, dist)
			return dist
		}

		d := 999999
		for _, tr := range adj[v] {
			// this check guarantees the recursion depth won't exceed the amount
			// of keys
			if keys&(1<<int(tr.v.ch-'a')) > 0 {
				// we already have this key
				continue
			}
			// do we have enough keys to transit?
			if (keys & tr.k) == tr.k {
				if nd := visit(tr.v, keys, dist+tr.l); nd >= 0 {
					d = min(d, nd)
				}
			}
		}
		return d
	}
	return visit(start, 0, 0)
}

func main() {
	file, err := os.Open("INPUT")
	noerr(err)
	defer file.Close()

	adj, start, allKeys := getAdjMatrix(file)
	log.Printf("adjacency matrix: %+v", adj)
	res := solve(adj, start, allKeys)
	log.Printf("Result is: %d", res)
}
