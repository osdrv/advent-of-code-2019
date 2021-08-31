package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	BASE_PATTERN = []int{0, 1, 0, -1}
)

type Digit uint8

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func NewDigit(v int) Digit {
	return Digit(abs(v) % 10)
}

type Number struct {
	digits []Digit
}

func NewNumber(s string) *Number {
	digits := make([]Digit, 0, len(s))
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			panic(fmt.Sprintf("Unexpected digit: %c", ch))
		}
		digits = append(digits, Digit(ch-'0'))
	}
	return &Number{
		digits: digits,
	}
}

func (n *Number) NextGenWithOff(off int) *Number {
	digits := make([]Digit, 0, len(n.digits))
	for gen := off; gen <= off+len(n.digits)-1; gen++ {
		//next := genPattern(BASE_PATTERN, off, gen)
		var d int
		for _, digit := range n.digits {
			//d += next() * int(digit)
			d += int(digit)
		}
		digits = append(digits, NewDigit(d))
	}
	return &Number{
		digits: digits,
	}
}

func (n *Number) NextGen() *Number {
	return n.NextGenWithOff(0)
}

func (n *Number) String() string {
	var buf bytes.Buffer
	for _, digit := range n.digits {
		buf.WriteByte(byte('0' + digit))
	}
	return buf.String()
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func genPattern(base []int, off int, gen int) func() int {
	ix := (off / (gen + 1)) % len(base)
	n := off % (gen + 1)
	n += 1
	return func() int {
		if n > gen {
			n = 0
			ix++
		}
		if ix >= len(base) {
			ix = 0
		}
		n++
		return base[ix]
	}
}

const (
	REP = 10_000
	OFF = 7
	GEN = 100
)

func main() {
	f, err := os.Open("INPUT-TST6")
	noerr(err)
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	noerr(err)
	s := strings.Trim(string(bs), "\n\r\t")

	var buf bytes.Buffer
	for i := 0; i < REP; i++ {
		buf.WriteString(s)
	}
	s = buf.String()

	off, err := strconv.Atoi(s[:OFF])
	noerr(err)

	num := NewNumber(s)

	for i := 0; i < GEN; i++ {
		log.Printf("Generation: %d", i)
		num = num.NextGen()
	}

	log.Printf("Result: %s", num.String()[off:off+8])
}
