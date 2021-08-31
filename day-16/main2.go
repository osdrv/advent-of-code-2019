package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	data, err := ioutil.ReadAll(f)
	f.Close()
	data = bytes.Trim(data, "\n\t\r")
	var buf bytes.Buffer
	for i := 0; i < 10000; i++ {
		buf.Write(data)
	}
	s := buf.String()
	off, err := strconv.Atoi(s[0:7])
	noerr(err)
	s = s[off:]

	for i := 0; i < 100; i++ {
		log.Printf("%d", i)
		buf.Reset()
		ix := 0
		total := 0
		for ix < len(s) {
			if ix == 0 {
				for _, v := range s {
					total += int(v - '0')
				}
			} else {
				total -= int(s[ix-1] - '0')
			}
			buf.WriteByte('0' + byte(abs(total%10)))
			ix++
		}
		s = buf.String()
		log.Printf("%s", s)
	}
	log.Printf("Result: %s", s[0:8])
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
