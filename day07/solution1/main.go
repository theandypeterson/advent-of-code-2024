package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Equation struct {
	testValue int
	values []int
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

	equations := make([]Equation, 0)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		testValue, _ := strconv.Atoi(parts[0][:len(parts[0])-1])

		values := make([]int, len(parts)-1)
		for i, p := range parts {
			if i == 0 {
				continue
			}
			v,_ := strconv.Atoi(p)
			values[i-1] = v
		}
		
		equations = append(equations, Equation{
			testValue: testValue,
			values: values,
		})
	}

	res := 0
	for _, equation := range equations {
		opStacks := generateOpStacks(len(equation.values))
		for _, opStack := range opStacks {
			x := equation.values[0]
			for i, op := range opStack {
				if op == "*" {
					x = x * equation.values[i+1]
				} else {
					x = x + equation.values[i+1]
				}
			}
			if x == equation.testValue {
				res += x
				break
			}
		}
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func generateOpStacks(x int) [][]string {
	if x == 2 {
		res := [][]string{ {"*"}, {"+"} }
		return res;
	} else {
		prev := generateOpStacks(x-1)

		res := make([][]string, len(prev)*2)
		for i,o := range prev {
			o1 := make([]string, x-1)
			o2 := make([]string, x-1)
			for ii,oo := range o {
				o1[ii] = oo
				o2[ii] = oo
			}
			o1[x-2] = "*"
			o2[x-2] = "+"
			res[i*2] = o1
			res[i*2+1] = o2
		}
		return res
	}
}