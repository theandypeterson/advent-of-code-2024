package main

import (
	"bufio"
	"fmt"
	"math"
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

	resCh := make(chan bool)
	var wg sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		wg.Add(1)
		go func() {
			defer wg.Done()
			numParts := []int{}
			for _,x := range parts {
				n, _ := strconv.Atoi(x)
				numParts = append(numParts, n)
			}
			resCh<-isSafe(numParts)
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res:=0
	for x := range resCh {
		if x {
			res++
		}
	}

	fmt.Printf("Result: %v\n", res)
	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func isSafe(level []int) bool {
	var levelType string
	if level[0] < level[1] {
		levelType = "decreasing"
	} else {
		levelType = "increasing"
	}

	if levelType == "decreasing" {
		for i := range len(level)-1 {
			if level[i] >= level[i+1] {
				return false
			} else if math.Abs(float64(level[i] - level[i+1])) > 3 {
				return false
			}
		}
		return true
	} else {
		for i := range len(level)-1 {
			if level[i] <= level[i+1] {
				return false
			} else if math.Abs(float64(level[i] - level[i+1])) > 3 {
				return false
			}
		}
		return true
	}
}
