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

	r, _ := regexp.Compile(`mul\((\d*,\d*)\)|do\(\)|don't\(\)`)
	match := r.FindAllString(content, -1)

	res := 0
	mode := "we're so back"
	for _, command := range match {
		if (command[:3] == "mul") && mode == "we're so back" {
			digits := strings.Split(command[4 : len(command)-1], ",")
			x, _ := strconv.Atoi(digits[0])
			y, _ := strconv.Atoi(digits[1])
			res += x*y
		} else if command[:3] == "do(" {
			mode = "we're so back"
		} else if command[:5] == "don't" {
			mode = "its so over"
		}
	}

	fmt.Printf("Result: %v\n", res)
	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}