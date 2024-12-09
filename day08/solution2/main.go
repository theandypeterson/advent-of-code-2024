package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Antenna struct {
	x int
	y int
	freq string
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

	grid := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		grid = append(grid, row)
	}

	antennas := []Antenna{}

	for y, row := range grid {
		for x, spot := range row {
			if spot == "." {
				continue
			}
			a := Antenna{
				x: x,
				y: y,
				freq: spot,
			}
			antennas = append(antennas, a)
		}
	}

	antinodes := map[string]struct{}{}

	obCheck := func(x int, y int) bool {
		return x<0 || y<0 || y>=len(grid) || x>=len(grid[0])
	}

	for _, antenna := range antennas {
		matchingFreq := []Antenna{}
		for _, a := range antennas {
			if a.freq == antenna.freq && !(a.x == antenna.x && a.y == antenna.y) {
				matchingFreq = append(matchingFreq, a)
			}
		}

		for _, a := range matchingFreq {
			rise := a.y - antenna.y
			run := a.x - antenna.x

			nodex := a.x
			nodey := a.y
			for {
				if obCheck(nodex,nodey) {
					break
				}
				antinodes[strconv.Itoa(nodex)+","+strconv.Itoa(nodey)]=struct{}{}
				nodex = nodex + run
				nodey = nodey + rise
			}

			nodex = nodex - run
			nodey = nodey - rise
			for {
				if obCheck(nodex,nodey) {
					break
				}
				antinodes[strconv.Itoa(nodex)+","+strconv.Itoa(nodey)]=struct{}{}
				nodex = nodex - run
				nodey = nodey - rise
			}
		}
	}

	res := 0
	for range antinodes {
		res++
	}

	fmt.Printf("Result: %v\n", res)
	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}