package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Antenna struct {
	x int
	y int
	freq string
}

type Grid [][]string

type Coor struct {
	x int
	y int
}

type Region struct {
	grid Grid
	label string
	coors []Coor
}

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("pass in input, e.g. \"1\"")
		return
	}
	arg := os.Args[1]
	file, err := os.Open(fmt.Sprintf("inputs/input%v.txt", arg))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		grid = append(grid, row)
	}

	coorRegistry := make(map[string]struct{})

	regions := []Region{}
	// Build regions
	for y, row := range grid {
		for x, plant := range row {
			coor := Coor{
				x: x,
				y: y,
			}

			if coor.registered(coorRegistry) {
				continue
			}
			region := buildRegion(plant, coor, coorRegistry, grid)
			regions = append(regions, region)
		}
	}

	// for _, region := range regions {
	// 	area := region.area()
	// 	per := region.perimeter()
	// 	fmt.Printf("%v: %v * %v = %v\n", region.label, area, per, region.price())
	// }

	res := 0
	for _, region := range regions {
		res += region.price()
	}

	fmt.Printf("Result: %v\n", res)


	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}


func (grid Grid) get(c Coor) string {
	x := c.x
	y := c.y
	if x < 0 || y < 0 || y >= len(grid) || x >= len(grid[0]) {
		return "~"
	} else {
		return grid[y][x]
	}
}

func (r Region) price() int{
	return r.area() * r.sides() 
}

func (r Region) area() int {
	return len(r.coors)
}

func (region Region) perimeter() int {
	res := 0
	g := region.grid
	for _, coor := range region.coors {
		u := coor.up()
		d := coor.down()
		l := coor.left()
		r := coor.right()
		if g.get(u) != region.label {
			res++
		}
		if g.get(d) != region.label {
			res++
		}
		if g.get(l) != region.label {
			res++
		}
		if g.get(r) != region.label {
			res++
		}
	}
	return res
}

func (c Coor) key() string {
	return fmt.Sprintf("%v,%v", c.x, c.y)
}

func (c Coor) registered(registry map[string]struct{}) bool {
	_, ok := registry[c.key()]
	return ok
}

func buildRegion(plant string, coor Coor, registry map[string]struct{}, grid Grid) Region {
	coors := findNeighbors(plant, coor, registry , grid)
	
	return Region{
		label: plant,
		coors: coors,
		grid: grid,
	}
}

func (c Coor) up() Coor {
	return Coor{
		x: c.x,
		y: c.y - 1,
	}
}

func (c Coor) down() Coor {
	return Coor{
		x: c.x,
		y: c.y + 1,
	}
}

func (c Coor) left() Coor {
	return Coor{
		x: c.x - 1,
		y: c.y,
	}
}

func (c Coor) right() Coor {
	return Coor{
		x: c.x + 1,
		y: c.y,
	}
}

func findNeighbors(plant string, coor Coor, registry map[string]struct{}, grid Grid) []Coor {
	_,ok := registry[coor.key()]
	if ok {
		return []Coor{}
	}

	height := len(grid)
	width := len(grid[0])
	if coor.y >= height || coor.x>= width || coor.x < 0 || coor.y < 0 {
		return []Coor{}
	}

	if grid[coor.y][coor.x] != plant {
		return []Coor{}
	}

	res := []Coor{ coor }

	registry[coor.key()] = struct{}{}
	// up
	u := findNeighbors(plant, coor.up(), registry, grid)
	res = append(res, u...)
	// down
	d := findNeighbors(plant, coor.down(), registry, grid)
	res = append(res, d...)
	// left
	l := findNeighbors(plant, coor.left(), registry, grid)
	res = append(res, l...)
	// right
	r := findNeighbors(plant, coor.right(), registry, grid)
	res = append(res, r...)

	return res
}

func (region Region) sides() int {
	res := 0
	for _, coor := range region.coors {
		if region.grid.get(coor.up()) != region.label && region.grid.get(coor.right()) != region.label {
			res++
		}
		if region.grid.get(coor.right()) != region.label && region.grid.get(coor.down()) != region.label {
			res++
		}
		if region.grid.get(coor.down()) != region.label && region.grid.get(coor.left()) != region.label {
			res++
		}
		if region.grid.get(coor.left()) != region.label && region.grid.get(coor.up()) != region.label {
			res++
		}

		if region.grid.get(coor.up()) == region.label && region.grid.get(coor.right()) == region.label && region.grid.get(coor.up().right()) != region.label {
			res++
		}
		if region.grid.get(coor.right()) == region.label && region.grid.get(coor.down()) == region.label && region.grid.get(coor.right().down()) != region.label {
			res++
		}
		if region.grid.get(coor.down()) == region.label && region.grid.get(coor.left()) == region.label && region.grid.get(coor.down().left()) != region.label {
			res++
		}
		if region.grid.get(coor.left()) == region.label && region.grid.get(coor.up()) == region.label && region.grid.get(coor.left().up()) != region.label {
			res++
		}

	}

	return res

}