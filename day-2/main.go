package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	CODE_ADD  = 1
	CODE_MULT = 2
	CODE_TERM = 99

	MODE_DIRECT  = 1
	MODE_POINTER = 2
)

func readval(program []int, pos, mode int) int {
	switch mode {
	case MODE_DIRECT:
		return program[pos]
	case MODE_POINTER:
		return program[program[pos]]
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func setval(program []int, pos, val, mode int) {
	switch mode {
	case MODE_DIRECT:
		program[pos] = val
		return
	case MODE_POINTER:
		program[program[pos]] = val
		return
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func compute(program []int) {
	pcnt := 0
	for {
		log.Printf("program counter: %d", pcnt)
		if pcnt < 0 || pcnt >= len(program) {
			panic("segmentation fault")
		}
		switch program[pcnt] {
		case CODE_ADD:
			a, b := readval(program, pcnt+1, MODE_POINTER), readval(program, pcnt+2, MODE_POINTER)
			log.Printf("Computing %d + %d", a, b)
			setval(program, pcnt+3, a+b, MODE_POINTER)
			pcnt += 4
		case CODE_MULT:
			a, b := readval(program, pcnt+1, MODE_POINTER), readval(program, pcnt+2, MODE_POINTER)
			log.Printf("Computing %d * %d", a, b)
			setval(program, pcnt+3, a*b, MODE_POINTER)
			pcnt += 4
		case CODE_TERM:
			log.Printf("Successfully terminated program")
			log.Printf("Program: %+v", program)
			return
		}
	}
}

func parseProgram(s string) ([]int, error) {
	chunks := strings.Split(s, ",")
	res := make([]int, 0, len(chunks))
	for _, ch := range chunks {
		n, err := strconv.Atoi(ch)
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func computeWithArgs(program []int, a, b int) int {
	program[1] = a
	program[2] = b
	compute(program)
	return program[0]
}

func copyIntArr(arr []int) []int {
	res := make([]int, len(arr))
	copy(res, arr)
	return res
}

func findMax(arr []int) int {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func main() {
	file, err := os.Open("INPUT")
	if err != nil {
		panic(fmt.Sprintf("Failed to open input file: %s", err))
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		program, err := parseProgram(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("Failed to parse program: %s", err))
		}
	Iterate:
		for noun := 0; noun < 100; noun++ {
			for verb := 0; verb < 100; verb++ {
				res := computeWithArgs(copyIntArr(program), noun, verb)
				log.Printf("Result: %d", res)
				if res == 19690720 {
					log.Printf("noun: %d, verb: %d", noun, verb)
					break Iterate
				}
			}
		}
	}
}
