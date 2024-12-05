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

	rules := make(map[int][]int)
	listOfPages := make([][]int, 0)
	mode := "rules"
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mode = "pages"
			continue
		}

		if mode == "rules" {
			parts := strings.Split(line, "|")
			first, _ := strconv.Atoi(parts[0])
			second, _ := strconv.Atoi(parts[1])
			entry, ok := rules[first]
			if !ok {
				entry = []int{}
			}
			entry = append(entry, second)
			rules[first] = entry
		} else {
			parts := strings.Split(line, ",")
			pages := make([]int, len(parts))
			for i, p := range parts {
				x, _ := strconv.Atoi(p)
				pages[i] = x
			}
			listOfPages = append(listOfPages, pages)
		}
	}

	resCh := make(chan int)
	var wg sync.WaitGroup

	for _, pages := range listOfPages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			legit, _ := IsSorted(pages, rules)
			if legit {
				return
			}
			sortedPages := SortPages(pages, rules)
			mid := sortedPages[(len(sortedPages)-1)/2]
			resCh<-mid
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res := 0
	for x := range resCh {
		res += x
	}

	fmt.Printf("Result: %v\n", res)
	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func Intersect(a []int, b []int) []int {
	res := make([]int, 0)
	hash := make(map[int]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			res = append(res, v)
		}
	}

	return res
}

/*

75 97 47 61 53

if no rule, move on

if there is a rule, check if that number is behind us
  if so move, me to the front of the list
*/
func SortPages(pages []int, rules map[int][]int) []int {
	isSorted, offendingPageIndex := IsSorted(pages, rules)
	if isSorted {
		return pages
	} else {
		newPages := []int{pages[offendingPageIndex]}
		newPages = append(newPages, pages[:offendingPageIndex]...)
		newPages = append(newPages, pages[offendingPageIndex+1:]...)
		return SortPages(newPages, rules)
	}
}

func IsSorted(pages []int, rules map[int][]int) (bool, int) {
	for i, page := range pages {
		if i == 0 {
			continue
		}
		rule, ok := rules[page]
		if !ok {
			continue
		}
		if len(Intersect(rule, pages[:i])) == 0 {
			continue
		}
		return false, i
	}
	return true, -1
}