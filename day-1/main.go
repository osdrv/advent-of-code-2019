package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func calcFuel(mass int) int {
	return mass/3 - 2
}

func readInts(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	res := make([]int, 0, 1)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		res = append(res, num)
	}
	return res, nil
}

func main() {
	nums, err := readInts("INPUT")
	if err != nil {
		panic(fmt.Sprintf("Failed to read input data: %s", err))
	}
	res := 0
	for _, num := range nums {
		f := calcFuel(num)
		log.Printf("mass: %d, Fuel: %d\n", num, f)
		res += f
		fm := f
		for {
			extra := calcFuel(fm)
			if extra <= 0 {
				break
			}
			log.Printf("extra fuel: %d\n", extra)
			res += extra
			fm = extra
		}
	}

	fmt.Printf("Result: %d\n", res)
}
