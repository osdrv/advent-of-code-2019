package main

import (
	"fmt"
	"log"
	"strings"
)

type Vector3 struct {
	x, y, z int
}

func (v Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{
		x: v.x + v2.x,
		y: v.y + v2.y,
		z: v.z + v2.z,
	}
}

type Planet struct {
	pos Vector3
	vel Vector3
}

func NewPlanet(pos Vector3) *Planet {
	return &Planet{
		pos: pos,
		vel: Vector3{0, 0, 0},
	}
}

func (p *Planet) String() string {
	return fmt.Sprintf("pos=<x=%d, y=%d, z=%d>, vel=<x=%d, y=%d, z=%d>",
		p.pos.x, p.pos.y, p.pos.z, p.vel.x, p.vel.y, p.vel.z)
}

func (p *Planet) Pot() int {
	return abs(p.pos.x) + abs(p.pos.y) + abs(p.pos.z)
}

func (p *Planet) Kin() int {
	return abs(p.vel.x) + abs(p.vel.y) + abs(p.vel.z)
}

func (p *Planet) EqualsTo(p2 *Planet) bool {
	return p.pos == p2.pos && p.vel == p2.vel
}

type PlanetSystem struct {
	planets []*Planet
	time    int
}

func NewPlanetSystem(planets []*Planet) *PlanetSystem {
	return &PlanetSystem{
		planets: planets,
	}
}

func (ps *PlanetSystem) adjustGravity() {
	for _, planet := range ps.planets {
		var dv Vector3
		for _, another := range ps.planets {
			if planet == another {
				continue
			}
			dv = dv.Add(Vector3{
				posToGrav(planet.pos.x, another.pos.x),
				posToGrav(planet.pos.y, another.pos.y),
				posToGrav(planet.pos.z, another.pos.z),
			})
		}
		planet.vel = planet.vel.Add(dv)
	}
}

func (ps *PlanetSystem) adjustVelocity() {
	for _, planet := range ps.planets {
		planet.pos = planet.pos.Add(planet.vel)
	}
}

func (ps *PlanetSystem) Tick() {
	ps.time++
	ps.adjustGravity()
	ps.adjustVelocity()
}

func (ps *PlanetSystem) Energy() int {
	energy := 0
	for _, planet := range ps.planets {
		energy += planet.Pot() * planet.Kin()
	}
	return energy
}

func (ps *PlanetSystem) Snapshot() *PlanetSystem {
	return &PlanetSystem{
		planets: cpPlanets(ps.planets),
	}
}

func (ps *PlanetSystem) EqualsTo(s *PlanetSystem) bool {
	total := len(s.planets)
	xMatch, yMatch, zMatch := 0, 0, 0
	for ix := range ps.planets {
		if ps.planets[ix].pos.x == s.planets[ix].pos.x && ps.planets[ix].vel.x == s.planets[ix].vel.x {
			xMatch++
		}
		if ps.planets[ix].pos.y == s.planets[ix].pos.y && ps.planets[ix].vel.y == s.planets[ix].vel.y {
			yMatch++
		}
		if ps.planets[ix].pos.z == s.planets[ix].pos.z && ps.planets[ix].vel.z == s.planets[ix].vel.z {
			zMatch++
		}
	}
	if xMatch == total {
		log.Printf("X match on %d", ps.time)
	}
	if yMatch == total {
		log.Printf("Y match on %d", ps.time)
	}
	if zMatch == total {
		log.Printf("Z match on %d", ps.time)
	}
	return total == xMatch && total == yMatch && total == zMatch
}

func (ps *PlanetSystem) String() string {
	chunks := make([]string, 0, len(ps.planets))
	for _, planet := range ps.planets {
		chunks = append(chunks, planet.String())
	}
	return strings.Join(chunks, "\n")
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func posToGrav(v1, v2 int) int {
	if v1 < v2 {
		return 1
	} else if v1 > v2 {
		return -1
	}
	return 0
}

func cpPlanets(planets []*Planet) []*Planet {
	res := make([]*Planet, 0, len(planets))
	for _, planet := range planets {
		res = append(res, &Planet{
			pos: planet.pos,
			vel: planet.vel,
		})
	}
	return res
}

func main() {
	scan := []Vector3{
		{-13, -13, -13},
		{5, -8, 3},
		{-6, -10, -3},
		{0, 5, -5},
	}
	//scan := []Vector3{
	//	{-1, 0, 2},
	//	{2, -10, -7},
	//	{4, -8, 8},
	//	{3, 5, -1},
	//}
	//scan := []Vector3{
	//	{-8, -10, 0},
	//	{5, 5, 10},
	//	{2, -7, 3},
	//	{9, -8, -3},
	//}
	/*
		<x=-13, y=-13, z=-13>
		<x=5, y=-8, z=3>
		<x=-6, y=-10, z=-3>
		<x=0, y=5, z=-5>
	*/

	planets := make([]*Planet, 0, len(scan))
	for _, s := range scan {
		planets = append(planets, NewPlanet(s))
	}

	ps := NewPlanetSystem(planets)
	snap := ps.Snapshot()

	//log.Printf("initial state:\n%s", ps)

	//T := 1000
	//for i := 0; i < T; i++ {
	//	ps.Tick()
	//	//log.Printf("tick %d:\n%s", i, ps)
	//}

	//energy := ps.Energy()

	//log.Printf("Total energy: %d", energy)

	total := len(ps.planets)
	var periods [3]int
	set := 0
	for {
		ps.Tick()
		xMatch, yMatch, zMatch := 0, 0, 0
		for ix := range ps.planets {
			if ps.planets[ix].pos.x == snap.planets[ix].pos.x && ps.planets[ix].vel.x == snap.planets[ix].vel.x {
				xMatch++
			}
			if ps.planets[ix].pos.y == snap.planets[ix].pos.y && ps.planets[ix].vel.y == snap.planets[ix].vel.y {
				yMatch++
			}
			if ps.planets[ix].pos.z == snap.planets[ix].pos.z && ps.planets[ix].vel.z == snap.planets[ix].vel.z {
				zMatch++
			}
		}
		if xMatch == total {
			if periods[0] == 0 {
				set++
				periods[0] = ps.time
			}
		}
		if yMatch == total {
			if periods[1] == 0 {
				set++
				periods[1] = ps.time
			}
		}
		if zMatch == total {
			if periods[2] == 0 {
				set++
				periods[2] = ps.time
			}
		}
		if set == 3 {
			break
		}
	}

	res := lcm(periods[0], lcm(periods[1], periods[2]))

	log.Printf("res: %d", res)

}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}
