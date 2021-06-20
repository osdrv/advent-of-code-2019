package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Graph struct {
	edges   map[string][]string
	reverse map[string]string
}

func NewGraph() *Graph {
	return &Graph{
		edges:   make(map[string][]string),
		reverse: make(map[string]string),
	}
}

func (g *Graph) AddEdge(from, to string) {
	g.edges[from] = append(g.edges[from], to)
	g.reverse[to] = from
}

func readGraph(file io.Reader) *Graph {
	graph := NewGraph()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\t\r\n")
		from, to := parseEdge(line)
		graph.AddEdge(from, to)
	}
	return graph
}

func parseEdge(s string) (string, string) {
	chunks := strings.Split(s, ")")
	return chunks[0], chunks[1]
}

func countOrbits(g *Graph) int {
	var count func(string, int) int
	count = func(pos string, dist int) int {
		res := dist
		for _, dest := range g.edges[pos] {
			res += count(dest, dist+1)
		}
		return res
	}
	return count("COM", 0)
}

func countTransfers(graph *Graph, from, to string) int {
	trace := make(map[string]int)
	ptr := graph.reverse[from]
	dist := 0
	for ptr != "" {
		trace[ptr] = dist
		dist++
		ptr = graph.reverse[ptr]
	}
	ptr = graph.reverse[to]
	dist = 0
	for ptr != "" {
		if fromDist, ok := trace[ptr]; ok {
			return dist + fromDist
		}
		ptr = graph.reverse[ptr]
		dist++
	}
	return -1
}

func main() {
	file, err := os.Open("INPUT")
	if err != nil {
		log.Fatalf("fail to open input file: %s", err)
	}
	graph := readGraph(file)
	log.Printf("%+v", graph)
	numOrbits := countOrbits(graph)
	log.Printf("number of orbits: %d", numOrbits)

	tr := countTransfers(graph, "YOU", "SAN")
	log.Printf("number of transfers: %d", tr)
}
