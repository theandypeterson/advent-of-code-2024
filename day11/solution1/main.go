package main

import (
	"fmt"
	"os"
	"strconv"
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
	file, err := os.ReadFile(fmt.Sprintf("inputs/input%v.txt", arg))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	content := string(file)
	stones := []int{}
	for _,c := range strings.Split(content, " ") {
		s,_ := strconv.Atoi(c)
		stones = append(stones, s)
	}

	for range 25 {
		newStones := []int{}
		for _, stone := range stones {
			stoneStr := strconv.Itoa(stone)
			if stone == 0 {
				newStones = append(newStones, 1)
			} else if len(stoneStr)%2==0 {
				stone1, _ := strconv.Atoi(stoneStr[:len(stoneStr)/2])
				stone2, _ := strconv.Atoi(stoneStr[len(stoneStr)/2:])
				newStones = append(newStones, stone1)
				newStones = append(newStones, stone2)
			} else {
				newStones = append(newStones, stone*2024)
			}
		}
		stones = newStones
	}

	res := len(stones)
	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}