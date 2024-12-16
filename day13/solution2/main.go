package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Machine struct {
	buttonAX int
	buttonAY int
	buttonBX int
	buttonBY int
	prizeX int
	prizeY int
	minAPress int 
	minBPress int
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

	machines := make([]Machine,0)
	machineDefs := make([][]string,0)
	def := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		def = append(def, line)
		if len(def) == 3 {
			machineDefs = append(machineDefs, def)
			def = make([]string, 0)
		}
	}

	for _, def := range machineDefs {
		machine := Machine{
			buttonAX: -1,
			buttonAY: -1 ,
			buttonBX: -1,
			buttonBY: -1,
			prizeX: -1,
			prizeY: -1,
			minAPress: -1,
			minBPress: -1,
		}
		for _, line := range def {
			parts := strings.Split(line, ":")
			
			moves := strings.Split(parts[1], ",")
			if parts[0] == "Prize" {
				x, _ := strconv.Atoi(strings.Split(moves[0], "=")[1])
				y, _ := strconv.Atoi(strings.Split(moves[1], "=")[1])
				machine.prizeX = x+10000000000000
				machine.prizeY = y+10000000000000
			}
			if parts[0] == "Button A" {
				x, _ := strconv.Atoi(strings.Split(moves[0], "+")[1])
				y, _ := strconv.Atoi(strings.Split(moves[1], "+")[1])
				machine.buttonAX = x
				machine.buttonAY = y
			}
			if parts[0] == "Button B" {
				x, _ := strconv.Atoi(strings.Split(moves[0], "+")[1])
				y, _ := strconv.Atoi(strings.Split(moves[1], "+")[1])
				machine.buttonBX = x
				machine.buttonBY = y
			}
		}
		machine.minAPress = int(math.Max(float64(machine.prizeX/machine.buttonAX), float64(machine.prizeY/machine.buttonAY)))
		machine.minBPress = int(math.Max(float64(machine.prizeX/machine.buttonBX), float64(machine.prizeY/machine.buttonBY)))
		machines = append(machines, machine)
	}


	res := 0
	for _, machine := range machines {
		minCost := 0
		for apress := range machine.minAPress {
			for bpress := range machine.minBPress {
				ax := apress*machine.buttonAX
				ay := apress*machine.buttonAY

				bx := bpress*machine.buttonBX
				by := bpress*machine.buttonBY
				x := ax+bx
				y := ay+by

				if x == machine.prizeX && y == machine.prizeY {
					cost := apress*3 + bpress*1

					if minCost == 0 {
						minCost = cost
					} else if cost < minCost {
						minCost = cost
					}
				}
			}
		}
		res += minCost
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}
