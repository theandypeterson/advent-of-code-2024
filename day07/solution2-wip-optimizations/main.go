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

type Operation struct {
	opStack []string
	key string
	prev *Operation
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
	cache := make(map[int][]Operation)
	// prime the cache
	generateOpStacks(maxLength, cache)
	for _, equation := range equations {
		wg.Add(1)
		calcCache := make(map[string]int)
		go func() {
			defer wg.Done()
			operations := generateOpStacks(len(equation.values), cache)
			for _, operation := range operations {
				x := calcStack2(equation.values, &operation, calcCache)
				if x == equation.testValue  {
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

func generateOpStacks(x int, cache map[int][]Operation) []Operation {
	if x == 2 {
		res := []Operation{ 
			{ opStack: []string{"*"}, key: "*", prev: nil }, 
			{ opStack: []string{"+"}, key: "+", prev: nil },
			{ opStack: []string{"|"}, key: "|", prev: nil },
		}
		return res;
	} else {
		if y, ok := cache[x]; ok {
			return y;
		}
		prev := generateOpStacks(x-1, cache)

		res := make([]Operation, len(prev)*3)
		for i,o := range prev {
			o1 := make([]string, x-1)
			o2 := make([]string, x-1)
			o3 := make([]string, x-1)
			for ii,oo := range o.opStack {
				o1[ii] = oo
				o2[ii] = oo
				o3[ii] = oo
			}
			o1[x-2] = "*"
			o2[x-2] = "+"
			o3[x-2] = "|"
			key1 := o.key+"*"
			key2 := o.key+"+"
			key3 := o.key+"|"
			res[i*3] = Operation{
				opStack: o1,
				key: key1,
				prev: &o,
			} 
			res[i*3+1] = Operation{
				opStack: o2,
				key: key2,
				prev: &o,
			}
			res[i*3+2] = Operation{
				opStack: o3,
				key: key3,
				prev: &o,
			} 
		}
		cache[x] = res
		return res
	}
}

// func calcStack(values []int, opStack []string) int  {
// 	x := values[0]
// 	for i, op := range opStack {
// 		if op == "*" {
// 			x = x * values[i+1]
// 		} else if op == "+" {
// 			x = x + values[i+1]
// 		} else {
// 			y := values[i+1]
// 			z := 1
// 			for z < y {
// 				z *=10
// 			}
// 			x = x*z+y
// 		}
// 	}
// 	return x
// }

func calcStack2(values []int, operation *Operation, cache map[string]int) int {
	if len(values) == 1 {
		return values[0]
	}
	
	c, ok := cache[operation.key] 
	if ok {
		fmt.Printf("cache hit %v\n", operation.key)
		return c
	} else {
		fmt.Println("cache miss")
	}

	prev := calcStack2(values[:len(values)-1], operation.prev, cache) 
	op := operation.opStack[len(operation.opStack)-1]
	curr := values[len(values)-1]
	res := 0
	if op == "*" {
		res = prev * curr
	} else if op == "+" {
		res = prev + curr
	} else {
		r, _ := strconv.Atoi(strconv.Itoa(prev) + strconv.Itoa(curr))
		res = r
	}
	cache[operation.key] = res
	return res
}