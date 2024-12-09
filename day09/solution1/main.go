package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Node struct {
	id int
	nextNode *Node
	prevNode *Node
}

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

	mode := "file"
	id :=0
	var initial *Node 
	var curr *Node
	for _,s := range content {
		x, _ := strconv.Atoi(string(s))
		if mode == "file" {
			for range x {
				node := Node{
					id: id,
					prevNode: curr,
				}
				if initial == nil {
					initial = &node
				}
				if curr == nil {
					curr = &node
				} else {
					curr.nextNode = &node
					curr = &node
				}
			}
			id++
			mode="free space"
		} else {
			for range x {
				node := Node{
					id: -1,
					prevNode: curr,
				}
				curr.nextNode = &node
				curr = &node
			}
			mode="file"
		}
	}

	initial.defrag()
	res := 0
	c := initial.createCursor()
	for {
		res += c.node.id * c.pos
		if !c.next() || c.node.id == -1 {
			break
		}
	}

	fmt.Printf("Result: %v\n", res)

	duration := time.Since(start)
	fmt.Printf("Process took %v\n", duration)
}

func (x *Node) printAllNodes() {
	for {
		if x.id == -1 {
			fmt.Printf(".")
		} else {
			fmt.Printf("%v", x.id)
		}
		if x.nextNode == nil {
			break;
		}
		x = x.nextNode
	}
	fmt.Println()
}

type Cursor struct {
	pos int
	node *Node
}

func (x *Node) defrag() {
	cur1 := x.createCursor()
	cur2 := x.createCursor()
	for cur2.node.nextNode != nil {
		cur2.next()
	}

	for {
		// find empty space
		for cur1.node.id != -1 {
			cur1.next()
		}

		// find file
		for cur2.node.id == -1 {
			cur2.prev()
		}

		if cur1.pos >= cur2.pos {
			break
		}

		swap(cur1.node, cur2.node)
	}
}

func (x *Cursor) next() bool {
	if x.node.nextNode != nil {
		x.node = x.node.nextNode
		x.pos++
		return true
	}
	return false
}

func (x *Cursor) prev() bool {
	if x.node.prevNode != nil {
		x.node = x.node.prevNode
		x.pos--
		return true
	}
	return false
}

func (x *Node) createCursor() Cursor {
	return Cursor{
		pos: 0,
		node: x,
	}
}

func swap(x1 *Node, x2 *Node) {
	i := x2.id
	x2.id = x1.id
	x1.id = i
}