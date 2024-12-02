package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
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
	list2 := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		n1, n2 := separateLines(line)

		list1 = append(list1, n1)
		list2 = append(list2, n2)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		slices.Sort(list1)
	}()
	go func() {
		defer wg.Done()
		slices.Sort(list2)
	}()
	wg.Wait()


	resultCh := make(chan int)
	for i := range list1 {
		wg.Add(1)
		go func () {
			defer wg.Done()
			x1 := list1[i]
			x2 := list2[i]
			dist := int(math.Abs(float64(x1 - x2)))
			resultCh<-dist
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