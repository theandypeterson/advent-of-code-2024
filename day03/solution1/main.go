package main

import (
	"fmt"
	"os"
	"regexp"
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

	r, _ := regexp.Compile(`mul\(\d*,\d*\)`)
	match := r.FindAllString(content, -1)

	res := 0
	for _, command := range match {
		digits := strings.Split(command[4 : len(command)-1], ",")
		x, _ := strconv.Atoi(digits[0])
		y, _ := strconv.Atoi(digits[1])
		res += x*y
	}

	fmt.Printf("Result: %v\n", res)
	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}