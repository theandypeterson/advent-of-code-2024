package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

	labmap := make([][]rune, 0)
	currPosX := -1
	currPosY := -1
	for scanner.Scan() {
		line := scanner.Text()
		if currPosX == -1 {
			x := strings.Index(line, "^")
			if x != -1 {
				currPosX = x
				currPosY = len(labmap)
			}
		}
		runes := make([]rune, 0)
		for _,r := range line {
			runes = append(runes, r)
		}
		labmap = append(labmap, runes)
	}

  positions, isLoop := navigateMap(labmap, currPosX, currPosY)
	fmt.Printf("isLoop: %v\n", isLoop)

	res := 0
	resCh := make(chan [2]int)
	var wg sync.WaitGroup
	for hash := range positions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parts := strings.Split(hash, ",")

			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			dir := parts[2]
			nx := x
			ny := y
			if dir == "up" {
				ny--
			} else if dir == "down" {
				ny++
			} else if dir == "left" {
				nx--
			} else if dir == "right" {
				nx++
			}
			if obCheck(labmap, nx, ny) {
				return
			} else {
				newMap := make([][]rune, len(labmap))
				for i := range labmap {
					newMap[i] = make([]rune, len(labmap[i]))
					copy(newMap[i], labmap[i])
				}
				newMap[ny][nx] = 'A'
			
				_, isLoop := navigateMap(newMap, currPosX, currPosY)
				if isLoop {
					resCh<-[2]int{nx,ny}
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	hash := make(map[string]struct{})
	for coor := range resCh {
		key := strconv.Itoa(coor[0])+"|"+strconv.Itoa(coor[1])
		if _, ok := hash[key]; !ok {
			res++
		}
		hash[key] = struct{}{}
	}

	fmt.Printf("res: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}
func obCheck(labmap [][]rune, x int,y int) bool {
	width := len(labmap[0])
	height := len(labmap)

	return x < 0 || y < 0 || x >= width || y >= height
}

func navigateMap(labmap [][]rune, x int, y int) (map[string]struct{}, bool) {
	currDirection := "up"
	currPosX := x
	currPosY := y

	positions := make(map[string]struct{}, 0)
	for {
		if obCheck(labmap, currPosX, currPosY) {
			return positions, false
		}
		 
		_, ok := positions[strconv.Itoa(currPosX)+","+strconv.Itoa(currPosY)+","+currDirection]
		if ok {
			return positions, true
		}
		positions[strconv.Itoa(currPosX)+","+strconv.Itoa(currPosY)+","+currDirection] = struct{}{}
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
		if obCheck(labmap, nextPosX, nextPosY) {
			currPosX = nextPosX
			currPosY = nextPosY
			continue
		} else if labmap[nextPosY][nextPosX] == '#' || labmap[nextPosY][nextPosX] == 'A' {
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
}