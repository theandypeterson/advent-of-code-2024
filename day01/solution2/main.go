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

	list1 := []int{}
	list2ButNowMap := make(map[int]int)

	for scanner.Scan() {
		line := scanner.Text()
		n1, n2 := separateLines(line)

		list1 = append(list1, n1)
		_, ok := list2ButNowMap[n2]
		if !ok {
			list2ButNowMap[n2] = 1
		} else {
			list2ButNowMap[n2] += 1
		}
	}

	var wg sync.WaitGroup

	resultCh := make(chan int)
	for _, x := range list1 {
		wg.Add(1)
		go func () {
			defer wg.Done()
			entry, ok := list2ButNowMap[x]
			if !ok {
				entry = 0
			}
			simScore := x*entry
			resultCh<-simScore
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := 0
	for x := range resultCh {
		result += x
	}

	fmt.Printf("Result: %v\n", result)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func separateLines(line string) (int, int) {
	parts := strings.Split(line, " ")
	var x1, x2 int
	for i,x := range parts {
		if x == "" {
			continue
		}
		n, _ := strconv.Atoi(x)
		if i == 0 {
			x1 = n
		} else {
			x2 = n
		}
	}
	return x1, x2
}

func gatherList(in chan int, done chan bool, out chan []int) {
	l := []int{}
	for {
	select {
	case n := <-in:
		l = append(l, n)
		fmt.Println("appending")
	case <-done:
		out<-l
		fmt.Println("done!")
		break
	}
	}
}

// func sortList(in chan []int) {
//   out chan []int
// 	x := <-in
//   sort.Ints(x)
// 	out <- x
// 	return out
// }