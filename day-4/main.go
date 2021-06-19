package main

import (
	"fmt"
	"log"
)

const (
	MIN = 128392
	MAX = 643281
)

type Number struct {
	digits [6]int
}

func NewNumber(v int) *Number {
	if v > 999999 {
		panic(fmt.Sprintf("%d exceeds max Number value (999999)", v))
	}
	var digits [6]int
	ix := len(digits) - 1
	for v > 0 {
		digits[ix] = v % 10
		v /= 10
		ix--
	}
	return &Number{
		digits: digits,
	}
}

func (n *Number) Inc() {
	ix := len(n.digits) - 1
	carry := 1
	for carry > 0 && ix >= 0 {
		v := n.digits[ix] + 1
		n.digits[ix] = v % 10
		carry = v / 10
		ix--
	}
}

func (n *Number) Cmp(n2 *Number) int {
	ix := 0
	for ix < len(n.digits) {
		if n.digits[ix] > n2.digits[ix] {
			return 1
		} else if n.digits[ix] < n2.digits[ix] {
			return -1
		}
		ix++
	}
	return 0
}

func (n *Number) String() string {
	skip := true
	res := make([]byte, 0, len(n.digits))
	for _, v := range n.digits {
		skip = skip && v == 0
		if skip {
			continue
		}
		res = append(res, '0'+byte(v))
	}
	return string(res)
}

func hasDoubleDigit(n *Number) bool {
	ix := 1
	for ix < len(n.digits) {
		if n.digits[ix] == n.digits[ix-1] {
			return true
		}
		ix++
	}
	return false
}

func hasDoubleDigitReduced(n *Number) bool {
	ix := 1
	cnt := 1
	for ix < len(n.digits) {
		if n.digits[ix] == n.digits[ix-1] {
			cnt++
		} else {
			if cnt == 2 {
				return true
			}
			cnt = 1
		}
		ix++
	}
	return cnt == 2
}

func setNonDecr(n *Number) {
	prev := 0
	for ix := 0; ix < len(n.digits); ix++ {
		if n.digits[ix] < prev {
			for ix < len(n.digits) {
				n.digits[ix] = prev
				ix++
			}
			return
		}
		prev = n.digits[ix]
	}
}

func main() {
	max := NewNumber(MAX)
	n := NewNumber(MIN)
	setNonDecr(n)

	cnt := 0

	//log.Printf("number: %s, has2digits: %t", NewNumber(112233), hasDoubleDigitReduced(NewNumber(112233)))
	//log.Printf("number: %s, has2digits: %t", NewNumber(123444), hasDoubleDigitReduced(NewNumber(123444)))
	//log.Printf("number: %s, has2digits: %t", NewNumber(111122), hasDoubleDigitReduced(NewNumber(111122)))

	for n.Cmp(max) < 0 {
		log.Printf("n: %s", n)
		if hasDoubleDigitReduced(n) {
			cnt++
		}
		n.Inc()
		setNonDecr(n)
	}
	log.Printf("Count: %d", cnt)
}
