package main

import (
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

var (
	Debug = 1
)

func printf(format string, v ...interface{}) {
	if Debug > 0 {
		log.Printf(format, v...)
	}
}

func fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

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
	fatalf("Unknown deref mode: %d", mode)
	return -1
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
	fatalf("Unknown deref mode: %d", mode)
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
		printf("program counter: %d", pcnt)
		op := program.OpOnly(pcnt)
		printf("Interpret opcode: %d[%d], rel: %d", op, program.OpRaw(pcnt), rel)
		switch op {
		case CODE_ADD:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			printf("Computing %d + %d", a, b)
			program.SetVal(rel, pcnt+3, a+b, mode3)
			pcnt += 4
		case CODE_MULT:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			a, b := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			printf("Computing %d * %d", a, b)
			program.SetVal(rel, pcnt+3, a*b, mode3)
			pcnt += 4
		case CODE_INPUT:
			_, mode1 := program.OpMono(pcnt)
			printf("mono mode: %d", mode1)
			printf("*** request for input")
			val := <-input
			printf("Reading %d from input", val)
			program.SetVal(rel, pcnt+1, val, mode1)
			pcnt += 2
		case CODE_OUTPUT:
			_, mode1 := program.OpMono(pcnt)
			printf("mono mode: %d", mode1)
			printf("Reading from pos %d", pcnt+1)
			val := program.ReadVal(rel, pcnt+1, mode1)
			printf("Writing %d to the output", val)
			output <- val
			pcnt += 2
		case CODE_JMPNZ:
			_, mode1, mode2 := program.OpDuo(pcnt)
			printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			if cond > 0 {
				printf("cond %d > 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				printf("cond %d is not > 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPZ:
			_, mode1, mode2 := program.OpDuo(pcnt)
			printf("duo mode: %d, %d", mode1, mode2)
			cond, dest := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2)
			if cond == 0 {
				printf("cond %d = 0, jump to %d", cond, dest)
				pcnt = dest
			} else {
				printf("cond %d is not 0, continue", cond)
				pcnt += 3
			}
		case CODE_JMPLT:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2), program.ReadVal(rel, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if mode3 == MODE_RELATIVE {
				pos += rel
			}
			if left < right {
				printf("left %d is less than right %d, writing flag 1 to %d", left, right, pos)
				program.SetVal(rel, pos, 1, MODE_IMMEDIATE)
			} else {
				printf("left %d is not less than right %d, writing flag 0 to %d", left, right, pos)
				program.SetVal(rel, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_JMPEQ:
			_, mode1, mode2, mode3 := program.OpTrio(pcnt)
			printf("trio mode: %d, %d, %d", mode1, mode2, mode3)
			left, right, pos := program.ReadVal(rel, pcnt+1, mode1), program.ReadVal(rel, pcnt+2, mode2), program.ReadVal(rel, pcnt+3, MODE_IMMEDIATE)
			//TODO: if a value like 1118 comes in, make sure to check it
			if mode3 == MODE_RELATIVE {
				pos += rel
			}
			if left == right {
				printf("left %d equals to right %d, writing flag 1 to %d", left, right, pos)
				program.SetVal(rel, pos, 1, MODE_IMMEDIATE)
			} else {
				printf("left %d is not equal to right %d, writing flag 0 to %d", left, right, pos)
				program.SetVal(rel, pos, 0, MODE_IMMEDIATE)
			}
			pcnt += 4
		case CODE_REL:
			_, mode1 := program.OpMono(pcnt)
			printf("mono mode: %d", mode1)
			adj := program.ReadVal(rel, pcnt+1, mode1)
			rel += adj
			printf("Adj: %d, New rel: %d", adj, rel)
			pcnt += 2
		case CODE_TERM:
			pcnt += 1
			program.SetStatus(PROGRAM_TERM)
			printf("Successfully terminated program")
			return
		default:
			fatalf("Unknown opcode: %d", op)
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
