package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func NewPoint(x, y int) Point {
	return Point{x, y}
}

type Segment struct {
	p1, p2 Point
}

func NewSegment(p1, p2 Point) Segment {
	if p1.x == p2.x {
		if p1.y > p2.y {
			p1, p2 = p2, p1
		}
	}
	if p1.y == p2.y {
		if p1.x > p2.x {
			p1, p2 = p2, p1
		}
	}
	return Segment{p1: p1, p2: p2}
}

func (s Segment) Contains(point Point) bool {
	if s.IsVertical() {
		return s.p1.x == point.x && s.p1.y <= point.y && point.y <= s.p2.y
	}
	return s.p1.y == point.y && s.p1.x <= point.x && point.x <= s.p2.x
}

func (s Segment) IsVertical() bool {
	return s.p1.x == s.p2.x
}

func (s Segment) Len() int {
	if s.IsVertical() {
		return s.p2.y - s.p1.y
	}
	return s.p2.x - s.p1.x
}

type Path struct {
	segments []Segment
}

func NewPath(srl string) Path {
	chunks := strings.Split(srl, ",")
	segments := make([]Segment, 0, len(chunks))
	var point, newPoint Point
	for _, ch := range chunks {
		switch ch[0] {
		case 'R':
			newPoint = NewPoint(point.x+mustInt(ch[1:]), point.y)
		case 'L':
			newPoint = NewPoint(point.x-mustInt(ch[1:]), point.y)
		case 'U':
			newPoint = NewPoint(point.x, point.y+mustInt(ch[1:]))
		case 'D':
			newPoint = NewPoint(point.x, point.y-mustInt(ch[1:]))
		default:
			log.Fatalf("Unrecognised directive: %s", ch)
		}
		segment := NewSegment(NewPoint(point.x, point.y), NewPoint(newPoint.x, newPoint.y))
		segments = append(segments, segment)
		point = newPoint
	}
	return Path{segments: segments}
}

func mustInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to parse number: %s", err)
	}
	return val
}

func IntersectSegments(s1, s2 Segment) *Point {
	if s1 == s2 {
		return nil
	}
	var vert, hor Segment
	if s1.IsVertical() {
		if s2.IsVertical() {
			return nil
		}
		vert, hor = s1, s2
	} else {
		if !s2.IsVertical() {
			return nil
		}
		vert, hor = s2, s1
	}
	if (hor.p1.x <= vert.p1.x && vert.p1.x <= hor.p2.x) &&
		(vert.p1.y <= hor.p1.y && hor.p1.y <= vert.p2.y) {
		p := NewPoint(vert.p1.x, hor.p1.y)
		return &p
	}
	return nil
}

func IntersectPaths(p1, p2 Path) []Point {
	res := make([]Point, 0, 1)
	zero := NewPoint(0, 0)
	for _, s1 := range p1.segments {
		for _, s2 := range p2.segments {
			if point := IntersectSegments(s1, s2); point != nil {
				if *point != zero {
					res = append(res, *point)
				}
			}
		}
	}
	return res
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func DistanceToPoint(path Path, dest Point) int {
	log.Printf("Calculating distance to point %+v for path %+v", dest, path)
	point := path.segments[0].p1
	dist := 0
	found := false
	for _, segment := range path.segments {
		if segment.Contains(dest) {
			log.Printf("Segment %+v contains the point", segment)
			dist += NewSegment(point, dest).Len()
			found = true
			break
		}
		dist += segment.Len()
	Points:
		for _, nextPoint := range []Point{segment.p1, segment.p2} {
			if point != nextPoint {
				point = nextPoint
				break Points
			}
		}
	}
	if !found {
		return -1
	}
	log.Printf("The distance is %d", dist)
	return dist
}

func main() {
	file, err := os.Open("INPUT")
	if err != nil {
		log.Fatalf("Failed to open input file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	paths := make([]Path, 0, 2)
	for scanner.Scan() {
		line := scanner.Text()
		path := NewPath(line)
		paths = append(paths, path)
		log.Printf("pparsed path: %+v", path)
	}
	p1, p2 := paths[0], paths[1]
	points := IntersectPaths(p1, p2)
	log.Printf("Intersections: %+v", points)

	minDist := 0
	for _, point := range points {
		dist := DistanceToPoint(p1, point) + DistanceToPoint(p2, point)
		if minDist == 0 || minDist > dist {
			minDist = dist
		}
	}
	log.Printf("Min distance: %d", minDist)
}
