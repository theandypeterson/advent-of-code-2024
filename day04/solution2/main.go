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
			if c == 'A' {
				coors = append(coors, []int{ x,height })
			}
		}
		height++
	}

	width := len(wordsearch[0])

	obCheck := func (x int, y int) bool {
		return y >= height || x>= width || x < 0 || y < 0
	}
	
	for _, coor := range coors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			x := coor[0]
			y := coor[1]

			x++
			y++

			a,b := false, false
			if obCheck(x,y) {
				return
			} else if wordsearch[y][x] == 'M' {
				x-=2
				y-=2
				if obCheck(x,y) {
					return
				} else if wordsearch[y][x] == 'S' {
					a= true
				}
			} else if wordsearch[y][x] == 'S' {
				x-=2
				y-=2
				if obCheck(x,y) {
					return
				} else if wordsearch[y][x] == 'M' {
					a= true
				}
			}

			x = coor[0]
			y = coor[1]

			x++
			y--
			if obCheck(x,y) {
				return
			} else if wordsearch[y][x] == 'M' {
				x-=2
				y+=2
				if obCheck(x,y) {
					return
				} else if wordsearch[y][x] == 'S' {
					b= true
				}
			} else if wordsearch[y][x] == 'S' {
				x-=2
				y+=2
				if obCheck(x,y) {
					return
				} else if wordsearch[y][x] == 'M' {
					b= true
				}
			}

			if a && b {
				resCh<-1
			}
		}()

	}

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
