package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Coord struct {
	x, y int
}

var (
	Up    = Coord{0, 1}
	Down  = Coord{0, -1}
	Right = Coord{1, 0}
	Left  = Coord{-1, 0}
)

type Robot struct {
	program *Program
	dir     Coord
	pos     Coord
}

func NewRobot(program *Program) *Robot {
	return &Robot{
		program: program,
		dir:     Up,
		pos:     Coord{0, 0},
	}
}

func (r *Robot) Move() {
	r.pos.x += r.dir.x
	r.pos.y += r.dir.y
}

func (r *Robot) TurnCounter() {
	switch r.dir {
	case Up:
		r.dir = Left
	case Left:
		r.dir = Down
	case Down:
		r.dir = Right
	case Right:
		r.dir = Up
	}
}

func (r *Robot) TurnClock() {
	switch r.dir {
	case Up:
		r.dir = Right
	case Right:
		r.dir = Down
	case Down:
		r.dir = Left
	case Left:
		r.dir = Up
	}
}

func (r *Robot) Run(field map[Coord]int) {
	in := make(chan int64, 1)
	out := make(chan int64, 2)

	field[Coord{0, 0}] = 1

	go func() {
		for {
			if r.program.Status() == PROGRAM_TERM {
				break
			}
			in <- int64(field[r.pos])
			color := <-out
			turn := <-out
			field[r.pos] = int(color)
			switch turn {
			case 0:
				r.TurnCounter()
			case 1:
				r.TurnClock()
			default:
				log.Printf("Unexpected turn directive: %d", turn)
			}
			r.Move()
		}
	}()

	Compute(r.program, in, out)

	close(in)
	close(out)
}

func drawPicture(field map[Coord]int, w io.Writer) error {
	maxX, maxY := 0, 0
	minX, minY := 0, 0
	log.Printf("%+v", field)
	for coord, color := range field {
		log.Printf("coord: %+v, color: %d", coord, color)
		x, y := coord.x, coord.y
		if maxX < x {
			maxX = x
		}
		if maxY < y {
			maxY = y
		}
		if minX > x {
			minX = x
		}
		if minY > y {
			minY = y
		}
	}

	adj := Coord{-minX, -minY}
	width, height := maxX-minX+1, maxY-minY+1

	log.Printf("minX: %d, minY: %d, maxX: %d, maxY: %d", minX, minY, maxX, maxY)
	log.Printf("adj: %+v", adj)
	log.Printf("width: %d, height: %d", width, height)

	pane := image.NewRGBA(image.Rect(0, 0, width, height))
	var c color.Color
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			switch field[Coord{i - adj.x, j - adj.y}] {
			case 0:
				c = color.Black
			case 1:
				c = color.White
			default:
				c = color.Transparent
			}
			pane.Set(i, height-1-j, c)
		}
	}
	if err := png.Encode(w, pane); err != nil {
		return err
	}

	return nil
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
	intCode, err := ParseProgram64(rawProgram)
	if err != nil {
		log.Fatalf("Failed to parse program: %s", err)
	}
	program := NewProgram(intCode)

	robot := NewRobot(program)
	field := make(map[Coord]int)
	robot.Run(field)

	log.Printf("Painted field size: %d", len(field))

	out, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Failed to create output png file: %s", err)
	}
	defer out.Close()
	if err := drawPicture(field, out); err != nil {
		log.Fatalf("Failed to draw picture: %s", err)
	}

}
