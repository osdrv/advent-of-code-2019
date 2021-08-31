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

const (
	COMMA = ','
	SCAFF = '#'
	SPACE = '.'

	UP    = '^'
	DOWN  = 'v'
	LEFT  = '<'
	RIGHT = '>'

	TURN_LEFT  = 'L'
	TURN_RIGHT = 'R'

	YES = "y"
	NO  = "n"
)

var (
	PATH_A = []string{"L", "10", "R", "8", "R", "8"}
	PATH_B = []string{"L", "10", "L", "12", "R", "8", "R", "10"}
	PATH_C = []string{"R", "10", "L", "12", "R", "10"}
)

type Coord struct {
	x, y int
}

func parseField(input []byte) [][]byte {
	return bytes.Split(bytes.Trim(input, "\n\r\t"), []byte{10})
}

var (
	STEP_LEFT  = Coord{-1, 0}
	STEP_RIGHT = Coord{1, 0}
	STEP_UP    = Coord{0, -1}
	STEP_DOWN  = Coord{0, 1}

	STEPS = []Coord{
		STEP_LEFT,
		STEP_RIGHT,
		STEP_DOWN,
		STEP_UP,
	}

	STEP_DIR = map[Coord]byte{
		STEP_LEFT:  LEFT,
		STEP_RIGHT: RIGHT,
		STEP_UP:    UP,
		STEP_DOWN:  DOWN,
	}
)

func getRoutineA() string {
	return strings.Join(PATH_A, ",")
}

func getRoutineB() string {
	return strings.Join(PATH_B, ",")
}

func getRoutineC() string {
	return strings.Join(PATH_C, ",")
}

func getMainRoutine(path string) string {
	a := getRoutineA()
	b := getRoutineB()
	c := getRoutineC()

	chunks := make([]string, 0, 1)

	ptr := path
	var ch string
	for len(ptr) > 0 {
		if strings.HasPrefix(ptr, ",") {
			ptr = ptr[1:]
		}
		if strings.HasPrefix(ptr, a) {
			ch = "A"
			ptr = ptr[len(a):]
		} else if strings.HasPrefix(ptr, b) {
			ch = "B"
			ptr = ptr[len(b):]
		} else if strings.HasPrefix(ptr, c) {
			ch = "C"
			ptr = ptr[len(c):]
		} else {
			log.Fatalf("Unknown prefix on string: %s.\nKnown prefixes are:\n\t%s\n\t%s\n\t%s", ptr, a, b, c)
		}
		chunks = append(chunks, ch)
	}
	return strings.Join(chunks, ",")
}

func findIntersects(field [][]byte) []Coord {
	res := make([]Coord, 0, 1)
	for i := 1; i < len(field)-1; i++ {
	Next:
		for j := 1; j < len(field[0])-1; j++ {
			if field[i][j] != SCAFF {
				continue Next
			}
			for _, s := range STEPS {
				if field[i+s.x][j+s.y] != SCAFF {
					continue Next
				}
			}
			res = append(res, Coord{i, j})
		}
	}
	return res
}

func withinRange(field [][]byte, pos Coord) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < len(field[0]) && pos.y < len(field)
}

func traverseField(field [][]byte) []byte {
	start, dir := findRobot(field)
	visited := make(map[Coord]struct{})
	path := make([]byte, 0, 1)

	var visit func(Coord)
	visit = func(pos Coord) {
		visited[pos] = struct{}{}
		var turn byte
		var nextStep Coord
		for _, step := range STEPS {
			npos := Coord{pos.x + step.x, pos.y + step.y}
			if !withinRange(field, npos) {
				continue
			}
			if field[npos.y][npos.x] != SCAFF {
				continue
			}
			if _, ok := visited[npos]; ok {
				continue
			}
			turn = getTurn(step, dir)
			if turn == 0 {
				log.Fatalf("Wrong turn from %+v with cur dir: %c", step, dir)
			}
			dir = STEP_DIR[step]
			nextStep = step
			break
		}
		if turn == 0 {
			return
		}
		pathlen := 0
		for {
			newpos := Coord{pos.x + nextStep.x, pos.y + nextStep.y}
			if !withinRange(field, newpos) || field[newpos.y][newpos.x] != SCAFF {
				break
			}
			visited[newpos] = struct{}{}
			pathlen++
			pos = newpos
		}
		if len(path) != 0 {
			path = append(path, COMMA)
		}
		path = append(path, turn)
		path = append(path, COMMA)
		path = append(path, []byte(strconv.Itoa(pathlen))...)

		visit(pos)
	}

	visit(start)

	return path
}

func getTurn(dir Coord, cur byte) byte {
	switch cur {
	case UP:
		switch dir {
		case STEP_LEFT:
			return TURN_LEFT
		case STEP_RIGHT:
			return TURN_RIGHT
		}
	case RIGHT:
		switch dir {
		case STEP_UP:
			return TURN_LEFT
		case STEP_DOWN:
			return TURN_RIGHT
		}
	case DOWN:
		switch dir {
		case STEP_LEFT:
			return TURN_RIGHT
		case STEP_RIGHT:
			return TURN_LEFT
		}
	case LEFT:
		switch dir {
		case STEP_DOWN:
			return TURN_LEFT
		case STEP_UP:
			return TURN_RIGHT
		}
	}
	return 0
}

func findRobot(field [][]byte) (Coord, byte) {
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			switch dir := field[i][j]; dir {
			case UP, DOWN, LEFT, RIGHT:
				return Coord{j, i}, dir
			}
		}
	}
	return Coord{-1, -1}, 0
}

func main() {
	Debug = 1

	file, err := os.Open("INPUT")
	noerr(err)
	data, err := ioutil.ReadAll(file)
	noerr(err)
	file.Close()
	rawProgram := strings.Trim(string(data), "\n\r\t")
	intCode, err := ParseProgram64(rawProgram)
	noerr(err)
	program := NewProgram(intCode)

	in := make(chan int64)
	out := make(chan int64)
	var buf bytes.Buffer
	go func() {
		for v := range out {
			buf.WriteByte(byte(v))
		}
		close(in)
	}()
	Compute(program, in, out)
	close(out)
	<-in

	fmt.Print(string(buf.Bytes()))

	field := parseField(buf.Bytes())
	//intersects := findIntersects(field)

	//res := 0
	//for _, in := range intersects {
	//	res += in.x * in.y
	//}

	//log.Printf("%+v", intersects)
	//log.Printf("res: %d", res)

	path := traverseField(field)
	log.Println(string(path))

	mainRoutine := getMainRoutine(string(path))
	a, b, c := getRoutineA(), getRoutineB(), getRoutineC()
	log.Printf("Main routine: %s", mainRoutine)
	log.Printf("Routine A: %s", a)
	log.Printf("Routine B: %s", b)
	log.Printf("Routine C: %s", c)

	in2 := make(chan int64, len(mainRoutine)+1+len(a)+1+len(b)+1+len(c)+1+len(NO)+1)
	out2 := make(chan int64)

	feedInput(in2, mainRoutine)
	feedInput(in2, a)
	feedInput(in2, b)
	feedInput(in2, c)
	feedInput(in2, NO)

	//res := make([]int64, 0, 1)
	//var res int64
	res := make([]byte, 0, 1)
	program = NewProgram(intCode)
	program.SetVal(0, 0, 2, MODE_IMMEDIATE)
	go func() {
		for v := range out2 {
			res = append(res, byte(v))
		}
		close(in2)
	}()
	Compute(program, in2, out2)
	close(out2)
	<-in2

	//fmt.Printf("%+v", res)
	fmt.Println(res[len(res)-1])
}

func feedInput(ch chan<- int64, s string) {
	for _, c := range s {
		ch <- int64(c)
	}
	ch <- int64('\n')
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
