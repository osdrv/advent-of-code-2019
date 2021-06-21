package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type Amplifier struct {
	program []int
	in      <-chan int
	out     chan<- int
}

func NewAmplifier(program []int, in <-chan int, out chan<- int) *Amplifier {
	return &Amplifier{
		program: program,
		in:      in,
		out:     out,
	}
}

func (a *Amplifier) Proc() {
	Compute(a.program, a.in, a.out)
}

func rotateRight(slc []int) {
	head := slc[0]
	copy(slc, slc[1:])
	slc[len(slc)-1] = head
}

func cpArr(arr []int) []int {
	res := make([]int, len(arr))
	copy(res, arr)
	return res
}

func permutate(seq []int) <-chan []int {
	ch := make(chan []int)
	var perm func(int)
	perm = func(k int) {
		if k == 1 {
			ch <- cpArr(seq)
		} else {
			for i := 0; i < k; i++ {
				perm(k - 1)
				if k%2 > 0 {
					seq[0], seq[k-1] = seq[k-1], seq[0]
				} else {
					seq[i], seq[k-1] = seq[k-1], seq[i]
				}
			}

		}
		if k == len(seq) {
			close(ch)
		}

	}
	go perm(len(seq))
	return ch
}

func main() {
	file, err := os.Open("INPUT")
	if err != nil {
		panic(fmt.Sprintf("Failed to open input file: %s", err))
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read input file: %s", err)
	}
	rawProgram := strings.Trim(string(data), "\n\r\t")
	program, err := ParseProgram(rawProgram)
	if err != nil {
		log.Fatalf("Failed to parse program: %s", err)
	}

	maxThurst := 0
	var maxSettings []int
	for settings := range permutate([]int{5, 6, 7, 8, 9}) {
		log.Printf("settings: %+v", settings)
		ampls := make([]*Amplifier, 0, len(settings))
		var in, out chan int
		wire := make(chan int, 2)
		out = wire
		var wg sync.WaitGroup
		for ix, set := range settings {
			in = out
			if ix == len(settings)-1 {
				out = wire
			} else {
				out = make(chan int, 2)
			}
			in <- set
			ampl := NewAmplifier(CopyProgram(program), in, out)
			ampls = append(ampls, ampl)
			wg.Add(1)
			go func() {
				ampl.Proc()
				wg.Done()
			}()
		}
		wire <- 0
		wg.Wait()
		thurst := <-wire
		if thurst > maxThurst {
			maxThurst = thurst
			maxSettings = settings
		}
		log.Printf("thurst: %d", thurst)
	}

	log.Printf("Max thurst: %d", maxThurst)
	log.Printf("MaxSettings: %+v", maxSettings)
}
