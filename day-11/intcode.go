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

	CODE_REL = 9

	CODE_TERM = 99

	MODE_POSITION  = 0
	MODE_IMMEDIATE = 1
	MODE_RELATIVE  = 2

	PROGRAM_TERM = 0
	PROGRAM_INIT = 1
	PROGRAM_RUN  = 2
	PROGRAM_ERR  = -1
)

func assert(cond bool, msg string) {
	if !cond {
		log.Fatalf("Assertion failed: %s", msg)
	}
}

type Program struct {
	memory map[int64]int64
	status int
}

func NewProgram(program []int64) *Program {
	memory := make(map[int64]int64)
	for ix, v := range program {
		memory[int64(ix)] = v
	}
	return &Program{
		memory: memory,
		status: PROGRAM_INIT,
	}
}

func (p *Program) ReadVal(rel int64, pos int64, mode int) int64 {
	switch mode {
	case MODE_IMMEDIATE:
		return p.memory[pos]
	case MODE_POSITION:
		return p.memory[p.memory[pos]]
	case MODE_RELATIVE:
		return p.memory[p.memory[pos]+rel]
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func (p *Program) SetVal(rel int64, pos int64, val int64, mode int) {
	switch mode {
	case MODE_IMMEDIATE:
		p.memory[pos] = val
		return
	case MODE_POSITION:
		p.memory[p.memory[pos]] = val
		return
	case MODE_RELATIVE:
		p.memory[p.memory[pos]+rel] = val
		return
	}
	panic(fmt.Sprintf("Unknown deref mode: %d", mode))
}

func (p *Program) OpRaw(pcnt int64) int {
	return int(p.memory[pcnt])
}

func (p *Program) OpOnly(pcnt int64) int {
	return int(p.memory[pcnt] % 100)
}

func (p *Program) OpMono(pcnt int64) (int, int) {
	mode1 := int((p.memory[pcnt] / 100) % 10)
	return p.OpOnly(pcnt), mode1
}

func (p *Program) OpDuo(pcnt int64) (int, int, int) {
	mode2 := int((p.memory[pcnt] / 1000) % 10)
	op, mode1 := p.OpMono(pcnt)
	return op, mode1, mode2
}

func (p *Program) OpTrio(pcnt int64) (int, int, int, int) {
	mode3 := int((p.memory[pcnt] / 10_000) % 10)
	op, mode1, mode2 := p.OpDuo(pcnt)
	return op, mode1, mode2, mode3
}

func (p *Program) OpQuatro(pcnt int64) (int, int, int, int, int) {
	mode4 := int((p.memory[pcnt] / 100_000) % 10)
	op, mode1, mode2, mode3 := p.OpTrio(pcnt)
	return op, mode1, mode2, mode3, mode4
}

func (p *Program) SetStatus(status int) {
	p.status = status
}

func (p *Program) Status() int {
	return p.status
}

func CopyProgram(program []int) []int {
	res := make([]int, len(program))
	copy(res, program)
	return res
}

func Compute(program *Program, input <-chan int64, output chan<- int64) {
	program.SetStatus(PROGRAM_RUN)
	var pcnt int64 = 0
	var rel int64 = 0
	for {
		//log.Printf("program: %+v", program)
		log.Printf("program counter: %d", pcnt)
		op := program.OpOnly(pcnt)
		log.Printf("Interpret opcode: %d[%d], rel: %d", op, program.OpRaw(pcnt), rel)
		switch op {
		case CODE_ADD:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			log.Printf("Computing %d + %d", a, b)
			program.SetVal(rel, pcnt+3, a+b, mode3)
			pcnt += 4
		case CODE_MULT:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			log.Printf("Computing %d * %d", a, b)
			program.SetVal(rel, pcnt+3, a*b, mode3)
			pcnt += 4
		case CODE_INPUT:
			_, mode1 := program.OpMono(pcnt)
			log.Printf("mono mode: %d", mode1)
			val := <-input
			log.Printf("Reading %d from input", val)
			program.SetVal(rel, pcnt+1, val, mode1)
			pcnt += 2
		case CODE_OUTPUT:
			_, mode1 := program.OpMono(pcnt)
			log.Printf("mono mode: %d", mode1)
			log.Printf("Reading from pos %d", pcnt+1)
			val := program.ReadVal(rel, pcnt+1, mode1)
			log.Printf("Writing %d to the output", val)
			output <- val
			pcnt += 2
		case CODE_JMPNZ:
			_, mode1, mode2 := program.OpDuo(pcnt)
			log.Printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			if cond > 0 {
				log.Printf("cond %d > 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				log.Printf("cond %d is not > 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPZ:
			_, mode1, mode2 := program.OpDuo(pcnt)
			log.Printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			if cond == 0 {
				log.Printf("cond %d = 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				log.Printf("cond %d is not 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPLT:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2), program.ReadVal(rel, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if mode3 == MODE_RELATIVE {
				pos += rel
			}
			if left < right {
				log.Printf("left %d is less than right %d, writing flag 1 to %d", left, right, pos)
				program.SetVal(rel, pos, 1, MODE_IMMEDIATE)
			} else {
				log.Printf("left %d is not less than right %d, writing flag 0 to %d", left, right, pos)
				program.SetVal(rel, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_JMPEQ:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			log.Printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2), program.ReadVal(rel, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if mode3 == MODE_RELATIVE {
				pos += rel
			}
			if left == right {
				log.Printf("left %d equals to right %d, writing flag 1 to %d", left, right, pos)
				program.SetVal(rel, pos, 1, MODE_IMMEDIATE)
			} else {
				log.Printf("left %d is not equal to right %d, writing flag 0 to %d", left, right, pos)
				program.SetVal(rel, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_REL:
			_, mode1 := program.OpMono(pcnt)
			log.Printf("mono mode: %d", mode1)
			adj := program.ReadVal(rel, pcnt+1, mode1)
			rel += adj
			log.Printf("Adj: %d, New rel: %d", adj, rel)
			pcnt += 2
		case CODE_TERM:
			pcnt += 1
			program.SetStatus(PROGRAM_TERM)
			log.Printf("Successfully terminated program")
			//log.Printf("Program: %+v", program)
			return
		default:
			log.Fatalf("Unknown opcode: %d", op)
		}
	}
}

func ParseProgram64(s string) ([]int64, error) {
	chunks := strings.Split(s, ",")
	res := make([]int64, 0, len(chunks))
	for _, ch := range chunks {
		n, err := strconv.ParseInt(ch, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}
