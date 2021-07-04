package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	EMPTY_TILE  = 0
	WALL_TILE   = 1
	BLOCK_TILE  = 2
	PADDLE_TILE = 3
	BALL_TILE   = 4
)

const (
	CONTROL_LEFT    = -1
	CONTROL_NEUTRAL = 0
	CONTROL_RIGHT   = 1
)

type Coord struct {
	x, y int
}

type Game struct {
	sync.Mutex
	score   int
	control int
	program *Program
	field   map[Coord]int
}

func NewGame(program *Program) *Game {
	return &Game{
		program: program,
		field:   make(map[Coord]int),
	}
}

func (g *Game) Field() map[Coord]int {
	return g.field
}

func (g *Game) Control(control int) {
	g.control = control
}

func (g *Game) SetScore(score int) {
	g.score = score
}

func (g *Game) Run(in <-chan int64) {
	out := make(chan int64, 3)
	g.control = 0
	go func() {
		for {
			if g.program.Status() == PROGRAM_TERM {
				break
			}
			x, y, c := <-out, <-out, <-out
			// score signal
			if x == -1 && y == 0 {
				g.SetScore(int(c))
				continue
			}
			g.Lock()
			g.field[Coord{int(x), int(y)}] = int(c)
			g.Unlock()
		}
	}()
	Compute(g.program, in, out)

	close(out)
}

func (g *Game) Frame() image.Image {
	g.Lock()
	defer g.Unlock()
	maxX, maxY := 0, 0
	minX, minY := 0, 0
	for coord := range g.field {
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

	pane := image.NewRGBA(image.Rect(0, 0, width, height))
	var c color.Color
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			switch g.field[Coord{i - adj.x, j - adj.y}] {
			case EMPTY_TILE:
				//empty tile
				c = color.White
			case WALL_TILE:
				// wall
				// grey
				c = color.RGBA{0xD0, 0xD0, 0xD0, 0xFF}
			case BLOCK_TILE:
				// block
				c = color.RGBA{0, 0, 0xFF, 0xFF}
			case PADDLE_TILE:
				// horizontal paddle
				c = color.RGBA{0, 0xFF, 0, 0xFF}
			case BALL_TILE:
				// ball
				c = color.RGBA{0xFF, 0, 0, 0xFF}
			default:
				c = color.Transparent
			}
			pane.Set(i, j, c)
		}
	}
	return pane
}

func main() {
	Debug = 0
	file, err := os.Open("INPUT")
	if err != nil {
		panic(fmt.Sprintf("Failed to open input file: %s", err))
	}
	defer file.Close()
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

	// put 2 quarters
	program.SetVal(0, 0, 2, MODE_IMMEDIATE)

	if err := ui.Init(); err != nil {
		panic(fmt.Sprintf("failed to initialize termui: %v", err))
	}
	defer ui.Close()

	game := NewGame(program)
	img := widgets.NewImage(nil)
	img.SetRect(0, 0, 44, 22) // I took these from a pre-run

	control := make(chan int64, 128)

	go game.Run(control)

	render := func() {
		img.Image = game.Frame()
		img.Title = "Score: " + strconv.Itoa(game.score)
		ui.Render(img)
	}

	go func() {
		for {
			render()
			<-time.After(100 * time.Millisecond)
		}
	}()

	uiEvents := ui.PollEvents()
EVENT_LOOP:
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				break EVENT_LOOP
			case "<Left>":
				control <- CONTROL_LEFT
			case "<Right>":
				control <- CONTROL_RIGHT
			case "<Space>", "<Enter>":
				control <- CONTROL_NEUTRAL
			}
			//case <-time.After(500 * time.Millisecond):
			//	control <- CONTROL_NEUTRAL
		}
	}
}
