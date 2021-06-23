package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	intCode, err := ParseProgram64(rawProgram)
	if err != nil {
		log.Fatalf("Failed to parse program: %s", err)
	}
	program := NewProgram(intCode)
	in := make(chan int64, 2)
	out := make(chan int64, 2)
	done := make(chan struct{})
	go func() {
		for v := range out {
			log.Printf("Out: %d", v)
		}
		close(done)
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			log.Print("INPUT: ")
			inp, _ := reader.ReadString('\n')
			n, err := strconv.ParseInt(strings.Trim(inp, "\n\r\t"), 10, 64)
			if err != nil {
				log.Printf("Failed to read input: %s", err)
				continue
			}
			in <- n
		}
	}()

	in <- 2
	Compute(program, in, out)
	close(out)
	<-done
}
