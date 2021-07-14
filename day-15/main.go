package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Coord struct {
	x, y int
}

type Direction int

const (
	NORTH Direction = 1
	SOUTH           = 2
	WEST            = 3
	EAST            = 4
)

const (
	WALL   = 0
	SPACE  = 1
	OXYGEN = 2
	PATH   = 3
)

type Robot struct {
	program *Program
	pos     Coord
	dir     Direction
	move    chan Direction
	in      chan int64
	out     chan int64
	term    chan struct{}
	done    chan struct{}
}

func NewRobot(program *Program) *Robot {
	return &Robot{
		program: program,
		pos:     Coord{0, 0},
		dir:     NORTH,
		move:    make(chan Direction),
		in:      make(chan int64, 1),
		out:     make(chan int64, 1),
		term:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (r *Robot) Pos() Coord {
	return r.pos
}

func (r *Robot) Start() {
	go Compute(r.program, r.in, r.out)
	go func() {
		for {
			select {
			case dir := <-r.move:
				r.in <- int64(dir)
			case <-r.term:
				close(r.done)
				return
			}
		}
	}()
}

func (r *Robot) Stop() {
	close(r.term)
	<-r.done
}

func (r *Robot) Move(dir Direction) int {
	r.move <- dir
	res := <-r.out
	if res > 0 {
		newPos := makeStep(r.pos, dir)
		r.dir = dir
		r.pos = newPos
	}
	return int(res)
}

func makeStep(base Coord, dir Direction) Coord {
	switch dir {
	case NORTH:
		return Coord{base.x, base.y - 1}
	case SOUTH:
		return Coord{base.x, base.y + 1}
	case WEST:
		return Coord{base.x - 1, base.y}
	case EAST:
		return Coord{base.x + 1, base.y}
	default:
		fatalf("Unexpected dir: %d", dir)
		return Coord{0, 0}
	}
}

var (
	Directions = []Direction{
		NORTH,
		EAST,
		SOUTH,
		WEST,
	}
)

type StepOption struct {
	pos   Coord
	dir   Direction
	score uint32
}

const (
	SCORE_SAME_DIR uint32 = 1 << iota
	SCORE_UNKNOWN
)

func scoreOpt(opt StepOption, field map[Coord]int, curDir Direction) uint32 {
	var score uint32
	if opt.dir == curDir {
		//score |= SCORE_SAME_DIR
	}
	if _, ok := field[opt.pos]; !ok {
		score |= SCORE_UNKNOWN
	}
	return score
}

func explore(robot *Robot, field map[Coord]int, nextStep chan struct{}) {
	unknowns := make(map[Coord]struct{})
	visits := make(map[Coord]int)
	for {
		opts := make([]StepOption, 0, 4)
	NextDir:
		for _, dir := range Directions {
			newPos := makeStep(robot.pos, dir)
			if v, ok := field[newPos]; ok {
				if v == WALL {
					continue NextDir
				}
			} else {
				unknowns[newPos] = struct{}{}
			}
			opt := StepOption{
				pos: newPos,
				dir: dir,
			}
			opt.score = scoreOpt(opt, field, robot.dir)
			opts = append(opts, opt)
		}
		sort.Slice(opts, func(i, j int) bool {
			if opts[i].score == opts[j].score {
				// least visited cell
				return visits[opts[i].pos] < visits[opts[j].pos]
			}
			return opts[i].score > opts[j].score
		})
		bestOpt := opts[0]
		field[bestOpt.pos] = robot.Move(bestOpt.dir)
		visits[robot.pos]++
		delete(unknowns, bestOpt.pos)
		if len(unknowns) == 0 {
			break
		}
		//<-nextStep
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func renderField(field map[Coord]int, robot *Robot) image.Image {
	maxX, maxY := 0, 0
	minX, minY := 0, 0
	for coord := range field {
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
			pos := Coord{i - adj.x, j - adj.y}
			v, ok := field[pos]
			if pos.x == 0 && pos.y == 0 {
				c = color.RGBA{0xFF, 0xFF, 0, 0xFF}
			} else if !ok {
				c = color.Black
			} else if robot.pos == pos {
				c = color.RGBA{0xFF, 0, 0, 0xFF}
			} else {
				switch v {
				case PATH:
					c = color.RGBA{0xFF, 0, 0xFF, 0xFF}
				case WALL:
					c = color.RGBA{0, 0xFF, 0, 0xFF}
				case SPACE:
					c = color.White
				case OXYGEN:
					c = color.RGBA{0, 0xFF, 0xFF, 0xFF}
				default:
					c = color.Transparent
				}
			}
			pane.Set(i, j, c)
		}
	}
	return pane
}

func findPath(field map[Coord]int, search int) (int, []Coord) {
	var visit func(Coord) (int, []Coord)
	visited := make(map[Coord]struct{})
	visit = func(pos Coord) (int, []Coord) {
		visited[pos] = struct{}{}
		if field[pos] == search {
			return 0, []Coord{pos}
		}
		minDist := -1
		var minChain []Coord
		for _, dir := range Directions {
			newPos := makeStep(pos, dir)
			if _, ok := visited[newPos]; ok {
				continue
			}
			if field[newPos] == WALL {
				continue
			}
			newDist, newChain := visit(newPos)
			if newDist >= 0 {
				minDist = newDist + 1
				minChain = append([]Coord{pos}, newChain...)
			}
		}
		return minDist, minChain
	}
	return visit(Coord{0, 0})
}

func findFillTime(field map[Coord]int, pos Coord) int {
	visited := make(map[Coord]int)
	var visit func(pos Coord, t int) int
	visit = func(pos Coord, t int) int {
		visited[pos] = t
		maxTime := t
	NextDir:
		for _, dir := range Directions {
			newPos := makeStep(pos, dir)
			if _, ok := visited[newPos]; ok {
				continue NextDir
			}
			if field[newPos] == WALL {
				continue NextDir
			}
			maxTime = max(maxTime, visit(newPos, t+1))
		}
		return maxTime
	}
	return visit(pos, 0)
}

func main() {
	Debug = 0

	if err := ui.Init(); err != nil {
		panic(fmt.Sprintf("failed to initialize termui: %v", err))
	}

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

	img := widgets.NewImage(nil)
	img.SetRect(0, 0, 60, 45)

	program := NewProgram(intCode)

	robot := NewRobot(program)
	robot.Start()
	defer robot.Stop()

	field := make(map[Coord]int)
	field[robot.pos] = SPACE
	control := make(chan struct{})

	//render := func() {
	//	img.Image = renderField(field, robot)
	//	img.Title = "Current cell: " + strconv.Itoa(field[robot.pos])
	//	ui.Render(img)
	//}

	explore(robot, field, control)

	minDist, minChain := findPath(field, OXYGEN)

	longestRange := findFillTime(field, minChain[len(minChain)-1])

	for _, pos := range minChain {
		field[pos] = PATH
	}

	img.Image = renderField(field, robot)
	img.Title = fmt.Sprintf("Shortest path: %d", minDist)
	ui.Render(img)

	uiEvents := ui.PollEvents()
EVENT_LOOP:
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<C-c>":
				break EVENT_LOOP
			case "<Space>":
				//render()
				//control <- struct{}{}
			}
			//case <-time.After(125 * time.Millisecond):
			//	render()
			//	control <- struct{}{}
		}
	}
	ui.Close()

	log.Printf("longestRange: %d", longestRange)
}
