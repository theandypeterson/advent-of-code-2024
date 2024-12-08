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
	maxLength := 0
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
		if len(values) > maxLength {
			maxLength = len(values)
		}
	}

	res := 0
	var wg sync.WaitGroup
	resCh := make(chan int)
	cache := make(map[int][][]int)
	// prime the cache
	genStart := time.Now()
	generateOpStacks(maxLength, cache)
	genDuration := time.Since(genStart)
	fmt.Printf("genDuration: %v\n", genDuration)
	for _, equation := range equations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opStacks := generateOpStacks(len(equation.values), cache)
			for _, opStack := range opStacks {
				x := calcStack(equation.values, opStack)
				if x == equation.testValue {
					resCh<-x
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for x := range resCh {
		res += x
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func generateOpStacks(x int, cache map[int][][]int) [][]int {
	if x == 2 {
		res := [][]int{ {1}, {2}, {3} }
		return res;
	} else {
		if y, ok := cache[x]; ok {
			return y;
		}
		prev := generateOpStacks(x-1, cache)

		res := make([][]int, len(prev)*3)
		for i,o := range prev {
			o1 := make([]int, x-1)
			o2 := make([]int, x-1)
			o3 := make([]int, x-1)
			for ii,oo := range o {
				o1[ii] = oo
				o2[ii] = oo
				o3[ii] = oo
			}
			o1[x-2] = 1
			o2[x-2] = 2
			o3[x-2] = 3
			res[i*3] = o1
			res[i*3+1] = o2
			res[i*3+2] = o3
		}
		cache[x] = res
		return res
	}
}

func calcStack(values []int, opStack []int) int  {
	x := values[0]
	for i, op := range opStack {
		if op == 1 {
			x = x * values[i+1]
		} else if op == 2 {
			x = x + values[i+1]
		} else {
			x = smush(x,values[i+1])
		}
	}
	return x
}

func smush(a int, b int) int {
	ystring := strconv.Itoa(b)
	xstring := strconv.Itoa(a)
	xint,_ := strconv.Atoi(xstring+ystring)
	return xint
}