package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
)

type Coord struct {
	x, y int
}

type Field struct {
	w, h      int
	asteroids []Coord
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func dist(c1, c2 Coord) int {
	return abs(c1.x-c2.x) + abs(c1.y-c2.y)
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func norm(a, b int) (int, int) {
	mod := gcd(abs(a), abs(b))
	return a / mod, b / mod
}

func calcConnections(f *Field) map[Coord]int {
	res := make(map[Coord]int)
	for _, ast := range f.asteroids {
		// create a map of other asteroids
		rem := make([]Coord, len(f.asteroids))
		copy(rem, f.asteroids)
		cnt := 0
		// sort by the distance
		sort.Slice(rem, func(i, j int) bool {
			return dist(ast, rem[i]) < dist(ast, rem[j])
		})
		log.Printf("Coord: %+v, sorted rem: %+v", ast, rem)
		scratch := make(map[Coord]struct{})
		scratch[ast] = struct{}{}
		for _, rr := range rem {
			if _, ok := scratch[rr]; ok {
				continue
			}
			cnt++
			x, y := ast.x, ast.y
			dx, dy := norm(rr.x-ast.x, rr.y-ast.y)
			for x >= 0 && x < f.w && y >= 0 && y < f.h {
				scratch[Coord{x, y}] = struct{}{}
				x += dx
				y += dy
			}
		}
		res[ast] = cnt
	}
	return res
}

func ReadField(in io.Reader) *Field {
	sc := bufio.NewScanner(in)
	asteroids := make([]Coord, 0, 1)
	w, h := 0, 0
	for sc.Scan() {
		line := sc.Text()
		h++
		if w <= 0 {
			w = len(line)
		}
		for x, ch := range line {
			switch ch {
			case '.':
			case '#':
				asteroids = append(asteroids, Coord{x, h - 1})
			default:
				log.Fatalf("Unknown char: %c", ch)
			}
		}
	}
	return &Field{
		w:         w,
		h:         h,
		asteroids: asteroids,
	}
}

type RadialCoord struct {
	cartesian Coord
	angle     float64
	r         int
}

func polarCoord(p1, p2 Coord) RadialCoord {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	// The angle is somewhat weird as it goes from vertical and then
	// clockwise
	angle := math.Pi/2 - math.Atan2(float64(-dy), float64(dx))
	if angle < 0 {
		angle += 2 * math.Pi
	}
	// r is not a real radius: it's a manhattan distance
	r := abs(dx) + abs(dy)
	return RadialCoord{
		cartesian: p2,
		angle:     angle,
		r:         r,
	}
}

func sortedPolarCoords(base Coord, asteroids []Coord) []RadialCoord {
	res := make([]RadialCoord, 0, len(asteroids)-1)
	for _, ast := range asteroids {
		if ast == base {
			continue
		}
		res = append(res, polarCoord(base, ast))
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].angle != res[j].angle {
			return res[i].angle < res[j].angle
		}
		return res[i].r < res[j].r
	})

	return res
}

func main() {
	file, err := os.Open("INPUT")
	noerr(err)
	defer file.Close()
	field := ReadField(file)
	log.Printf("field: %+v", field)

	conns := calcConnections(field)
	log.Printf("Connections: %+v", conns)
	maxConn := 0
	var maxCoord Coord
	for coord, conn := range conns {
		if conn > maxConn {
			maxConn = conn
			maxCoord = coord
		}
	}
	log.Printf("Max conns: %d at pos: %+v", maxConn, maxCoord)
	asteroids := sortedPolarCoords(maxCoord, field.asteroids)

	cnt := 0
	//n := 200
	var head RadialCoord
	for len(asteroids) > 0 {
		head, asteroids = asteroids[0], asteroids[1:]
		cnt++
		log.Printf("%d-th asteroid: %+v, mult: %d", cnt, head, head.cartesian.x*100+head.cartesian.y)
		//if cnt == n {
		//	break
		//}
		tail := make([]RadialCoord, 0, 1)
		for len(asteroids) > 0 && head.angle == asteroids[0].angle {
			tail = append(tail, asteroids[0])
			asteroids = asteroids[1:]
		}
		asteroids = append(asteroids, tail...)
	}
}

func noerr(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
