package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	CODE_ADD  = 1
	CODE_MULT = 2

	CODE_INPUT  = 3
	CODE_OUTPUT = 4

	CODE_JMPNZ = 5
	CODE_JMPZ  = 6
	CODE_JMPLT = 7
	CODE_JMPEQ = 8

	CODE_TERM = 99

	MODE_POSITION  = 0
	MODE_IMMEDIATE = 1
)

func readval(program []int, pos, mode int) int {
	switch mode {
	case MODE_IMMEDIATE:
		return program[pos]
	case MODE_POSITION:
		return program[program[pos]]
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func setval(program []int, pos, val, mode int) {
	switch mode {
	case MODE_IMMEDIATE:
		program[pos] = val
		return
	case MODE_POSITION:
		program[program[pos]] = val
		return
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func opOnly(op int) int {
	return op % 100
}

func opMono(op int) (int, int) {
	mode1 := (op / 100) % 10
	return op % 100, mode1
}

func opDuo(op int) (int, int, int) {
	mode2 := (op / 1000) % 10
	op, mode1 := opMono(op)
	return op, mode1, mode2
}

func opTrio(op int) (int, int, int, int) {
	mode3 := (op / 10_000) % 10
	op, mode1, mode2 := opDuo(op)
	return op, mode1, mode2, mode3
}

func opQuatro(op int) (int, int, int, int, int) {
	mode4 := (op / 100_000) % 10
	op, mode1, mode2, mode3 := opTrio(op)
	return op, mode1, mode2, mode3, mode4
}

func CopyProgram(program []int) []int {
	res := make([]int, len(program))
	copy(res, program)
	return res
}

func Compute(program []int, input <-chan int, output chan<- int) {
	pcnt := 0
	for {
		log.Printf("program: %+v", program)
		log.Printf("program counter: %d", pcnt)
		if pcnt < 0 || pcnt >= len(program) {
			panic("segmentation fault")
		}
		op := opOnly(program[pcnt])
		log.Printf("Interpret opcode: %d[%d]", op, program[pcnt])
		switch op {
		case CODE_ADD:
			_, mode1, mode2, mode3 := opTrio(program[pcnt])
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2)
			log.Printf("Computing %d + %d", a, b)
			setval(program, pcnt+3, a+b, mode3)
			pcnt += 4
		case CODE_MULT:
			_, mode1, mode2, mode3 := opTrio(program[pcnt])
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2)
			log.Printf("Computing %d * %d", a, b)
			setval(program, pcnt+3, a*b, mode3)
			pcnt += 4
		case CODE_INPUT:
			_, mode1 := opMono(program[pcnt])
			log.Printf("mono mode: %d", mode1)
			val := <-input
			log.Printf("Reading %d from input", val)
			setval(program, pcnt+1, val, mode1)
			pcnt += 2
		case CODE_OUTPUT:
			_, mode1 := opMono(program[pcnt])
			log.Printf("mono mode: %d", mode1)
			log.Printf("Reading from pos %d", pcnt+1)
			val := readval(program, pcnt+1, mode1)
			log.Printf("Writing %d to the output", val)
			output <- val
			pcnt += 2
		case CODE_JMPNZ:
			_, mode1, mode2 := opDuo(program[pcnt])
			log.Printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2)
			if cond > 0 {
				log.Printf("cond %d > 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				log.Printf("cond %d is not > 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPZ:
			_, mode1, mode2 := opDuo(program[pcnt])
			log.Printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2)
			if cond == 0 {
				log.Printf("cond %d = 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				log.Printf("cond %d is not 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPLT:
			_, mode1, mode2, mode3 := opTrio(program[pcnt])
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2), readval(program, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if left < right {
				log.Printf("left %d is less than right %d, writing flag 1 to %d", left, right, pos)
				setval(program, pos, 1, MODE_IMMEDIATE)
			} else {
				log.Printf("left %d is not less than right %d, writing flag 0 to %d", left, right, pos)
				setval(program, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_JMPEQ:
			_, mode1, mode2, mode3 := opTrio(program[pcnt])
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := readval(program, pcnt+1, mode1), readval(program, pcnt+2, mode2), readval(program, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if left == right {
				log.Printf("left %d equals to right %d, writing flag 1 to %d", left, right, pos)
				setval(program, pos, 1, MODE_IMMEDIATE)
			} else {
				log.Printf("left %d is not equal to right %d, writing flag 0 to %d", left, right, pos)
				setval(program, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_TERM:
			pcnt += 1
			log.Printf("Successfully terminated program")
			log.Printf("Program: %+v", program)
			return
		default:
			log.Fatalf("Unknown opcode: %d", op)
		}
	}
}

func ParseProgram(s string) ([]int, error) {
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
