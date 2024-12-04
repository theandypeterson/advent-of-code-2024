package main

import (
	"bufio"
	"fmt"
	"os"
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

	resCh := make(chan int)
	var wg sync.WaitGroup
	wordsearch := WordSearch{}
	coors := [][]int{}
	height := 0
	for scanner.Scan() {
		line := scanner.Text()
		wordsearch = append(wordsearch, line)
		for x, c := range line {
			if c == 'X' {
				coors = append(coors, []int{ x,height })
			}
		}
		height++
	}

	// down
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ 0, 1 }
		wordsearch.search(coors, dir, resCh)
	}()

	// up
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ 0, -1 }
		wordsearch.search(coors, dir, resCh)
	}()

	// left
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ -1, 0 }
		wordsearch.search(coors, dir, resCh)
	}()

	// right
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ 1, 0 }
		wordsearch.search(coors, dir, resCh)
	}()

	// up right
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ 1, -1 }
		wordsearch.search(coors, dir, resCh)
	}()

	// up left
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ -1, -1 }
		wordsearch.search(coors, dir, resCh)
	}()

	// down right
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ 1, 1 }
		wordsearch.search(coors, dir, resCh)
	}()

	// down left
	wg.Add(1)
	go func() {
		defer wg.Done()

		dir := [2]int{ -1, 1 }
		wordsearch.search(coors, dir, resCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res := 0
	for x := range resCh {
		res += x
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

type WordSearch []string

func (ws WordSearch) search(coors [][]int, direction [2]int, resCh chan int) {
	height := len(ws)
	width := len(ws[0])

	obCheck := func (x int, y int) bool {
		return y >= height || x>= width || x < 0 || y < 0
	}
	
	for _, coor := range coors {

		x := coor[0]
		y := coor[1]

		x += direction[0]
		y += direction[1]
		if obCheck(x,y) {
			continue
		} else if ws[y][x] == 'M' {
			x += direction[0]
			y += direction[1]
			if obCheck(x,y) {
				continue
			} else if ws[y][x] == 'A' {
				x += direction[0]
				y += direction[1]
				if obCheck(x,y) {
					continue
				} else if ws[y][x] == 'S' {
					resCh<-1
				}
			}
		}
	}
}