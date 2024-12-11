package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Antenna struct {
	x int
	y int
	freq string
}

type Grid [][]int

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
	trailheads := [][2]int{}
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, r := range line {
			if r == '.' {
				row[i] = -1
			} else {
				n, _ := strconv.Atoi(string(r))
				row[i] = n
				if n == 0 {
					trailheads = append(trailheads, [2]int{ i, rowCount })
				}
			}
		}
		rowCount++
		grid = append(grid, row)
	}

	res := 0
	for _, trailhead := range trailheads {
		res += trailRating(grid, trailhead[0], trailhead[1])
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}


func (grid Grid) get(x int, y int) int {
	if x < 0 || y < 0 || y >= len(grid) || x >= len(grid[0]) {
		return -1
	} else {
		return grid[y][x]
	}
}

func trailRating(grid Grid, x int, y int) int {
	spot := grid.get(x,y)
	if spot == 9 {
		return 1
	}

	rating := 0
	// up
	if grid.get(x, y-1) == spot+1 {
		rating += trailRating(grid,  x, y-1)
	}
	// down
	if grid.get(x,y+1) == spot+1 {
		rating += trailRating(grid, x, y+1)
	}
	// left
	if grid.get(x-1,y) == spot+1 {
		rating += trailRating(grid, x-1, y)
	}
	// right
	if grid.get(x+1,y) == spot+1 {
		rating += trailRating(grid, x+1, y)
	}
	return rating
}