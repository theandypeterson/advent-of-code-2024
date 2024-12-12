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
	threshold, _ := strconv.Atoi(os.Args[2])
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



	res := 0
	cache := make(map[string]int)
	for _, blinkStone := range stones {
		res += countStones(blinkStone, threshold, cache)
	}


	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func countStones(stone int, blink int, cache map[string]int) int {
	if blink == 0 {
		return 1
	}
	key := strconv.Itoa(stone)+","+strconv.Itoa(blink)
	v, ok := cache[key]
	if ok {
		return v
	}

	newStones := processStone(stone)
	res := 0
	for _,stone := range newStones {
			res += countStones(stone, blink-1, cache)
	}
	cache[key] = res
	return res
}

func processStone(stone int) []int {
	stoneStr := strconv.Itoa(stone)
	if stone == 0 {
		return []int{1}
	} else if len(stoneStr)%2==0 {
		stone1, _ := strconv.Atoi(stoneStr[:len(stoneStr)/2])
		stone2, _ := strconv.Atoi(stoneStr[len(stoneStr)/2:])
		return []int{stone1, stone2}
	} else {
		return []int{stone*2024}
	}
}