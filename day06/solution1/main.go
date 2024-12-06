package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

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

	labmap := make([]string, 0)
	currPosX := -1
	currPosY := -1
	currDirection := "up"
	for scanner.Scan() {
		line := scanner.Text()
		if currPosX == -1 {
			x := strings.Index(line, "^")
			if x != -1 {
				currPosX = x
				currPosY = len(labmap)
			}
		}
		labmap = append(labmap, line)
	}
	width := len(labmap[0])
	height := len(labmap)

	positions := make(map[int]struct{}, 0)
	obCheck := func(x int,y int) bool {
		return x < 0 || y < 0 || x >= width || y >= height
	}
	for {
		if obCheck(currPosX, currPosY) {
			break
		}
		positions[currPosX+currPosY*height] = struct{}{}
		nextPosX := currPosX
		nextPosY := currPosY
		if currDirection == "up" {
			nextPosY--
		} else if currDirection == "down" {
			nextPosY++
		} else if currDirection == "left" {
			nextPosX--
		} else if currDirection == "right" {
			nextPosX++
		}

		// need to check ob before accessing labmap
		if obCheck(nextPosX, nextPosY) {
			currPosX = nextPosX
			currPosY = nextPosY
			continue
		} else if labmap[nextPosY][nextPosX] == '#' {
			if currDirection == "up" {
				currDirection = "right"
			} else if currDirection == "down" {
				currDirection = "left"
			} else if currDirection == "left" {
				currDirection = "up"
			} else if currDirection == "right" {
				currDirection = "down"
			}
		} else {
			currPosX = nextPosX
			currPosY = nextPosY
			continue
		}
	}

	res := 0
	for range positions {
		res++
	}

	fmt.Printf("res: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}