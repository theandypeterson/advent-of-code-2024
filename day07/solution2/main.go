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
	cache := make(map[int][]int)
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

func generateOpStacks(x int, cache map[int][]int) []int {
	if x == 2 {
		res := []int{ 1, 2, 3 }
		return res;
	} else {
		if y, ok := cache[x]; ok {
			return y;
		}
		prev := generateOpStacks(x-1, cache)

		res := make([]int, len(prev)*3)
		for i,o := range prev {
			o1 := smush(o,1)
			o2 := smush(o,2)
			o3 := smush(o,3)
			res[i*3] = o1
			res[i*3+1] = o2
			res[i*3+2] = o3
		}
		cache[x] = res
		return res
	}
}

func calcStack(values []int, opStack int) int  {
	if len(values) == 1 {
		return values[0]
	}

	prev := calcStack(values[:len(values)-1], opStack/10)
	res := 0
	x := values[len(values)-1]
	op := opStack - opStack/10*10
	if op == 1 {
		res = prev * x
	} else if op == 2 {
		res = prev + x
	} else {
		res = smush(prev,x)
	}

	return res
}

func smush(a int, b int) int {
	if b==0 {
		return a;
	}
	return smush(a, b/10)*10+ b%10;
}
