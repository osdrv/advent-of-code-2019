package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Chemical struct {
	quant  int
	handle string
	comps  []Chemical
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func parseReaction(s string) Chemical {
	chunks := strings.SplitN(s, " => ", 2)
	chem := parseChemicals(chunks[1])[0]
	comps := parseChemicals(chunks[0])
	chem.comps = comps
	return chem
}

func parseChemicals(s string) []Chemical {
	res := make([]Chemical, 0, 1)
	chunks := strings.Split(s, ", ")
	for _, ch := range chunks {
		res = append(res, parseChemical(ch))
	}
	return res
}

func parseChemical(s string) Chemical {
	chunks := strings.SplitN(s, " ", 2)
	quant, err := strconv.Atoi(chunks[0])
	handle := chunks[1]
	noerr(err)
	return Chemical{
		quant:  quant,
		handle: handle,
	}
}

const (
	ORE = "ORE"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func computeORE(reactions map[string]Chemical, left map[string]int, handle string, quant int) int {
	if handle == ORE {
		if quant > left[ORE] {
			return -1
		}
		left[ORE] -= quant
		return quant
	}
	rem := left[handle]
	tmp := quant - min(rem, quant)
	left[handle] = max(left[handle]-quant, 0)
	quant = tmp
	ore := 0
	if quant > 0 {
		base := reactions[handle]
		mult := int(math.Ceil(float64(quant) / float64(base.quant)))
		for _, comp := range base.comps {
			compore := computeORE(reactions, left, comp.handle, mult*comp.quant)
			if compore < 0 {
				return -1
			}
			ore += compore
		}
		left[handle] = mult*base.quant - quant
	}
	return ore
}

func main() {
	file, err := os.Open("INPUT")
	noerr(err)

	reactions := make(map[string]Chemical)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Trim(scanner.Text(), "\n\t\r")
		chem := parseReaction(s)
		reactions[chem.handle] = chem
	}

	log.Printf("Reactions: %+v", reactions)

	left := make(map[string]int)
	left[ORE] = 1000000000000
	cnt := 0
	for {
		ore := computeORE(reactions, left, "FUEL", 1)
		if ore < 0 {
			break
		}
		cnt++
	}
	log.Printf("Total count: %d", cnt)
}
