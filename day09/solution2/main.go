package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Node struct {
	id int
	size int
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
			node := Node{
				id: id,
				size: x,
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
			id++
			mode="free space"
		} else {
			if x == 0 {
				mode="file"
				continue
			}
			node := Node{
				id: -1,
				size: x,
				prevNode: curr,
			}
			curr.nextNode = &node
			curr = &node
			mode="file"
		}
	}

	initial.printAllNodes()
	initial.defrag()
	initial.printAllNodes()
	res := 0
	c := initial.createCursor()
	pos := 0
	for {
		for range c.node.size {
			if c.node.id != -1 {
				res += c.node.id * pos
			}
			pos++
		}
		if !c.next() {
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
			for range x.size {
				fmt.Printf(".")
			}
		} else {
			for range x.size {
				fmt.Printf("%v", x.id)
			}
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
	fileCursor := x.createCursor()

	// get to end
	for fileCursor.node.nextNode != nil {
		fileCursor.next()
	}
	for {

		// fmt.Println("looping")
		// fmt.Printf("fileCursor.pos: %v\n", fileCursor.pos)
		// fmt.Printf("fileCursor.node.id: %v\n", fileCursor.node.id)
		if fileCursor.node.prevNode == nil {
			// fmt.Println("prev node nil, at the beginning")
			break
		}
		// fmt.Printf("fileCursor.node.prevNode: %v\n", fileCursor.node.prevNode.id)
		// find file
		// fmt.Println("find file")
		for fileCursor.node.id == -1 {
			// fmt.Printf("before: %v\n", fileCursor.node.id)
			if !fileCursor.prev() {
				break;
			}
			// fmt.Printf("on: %v\n", fileCursor.node.id)
		}
		// fmt.Printf("fileCursor.pos: %v\n", fileCursor.pos)
		// fmt.Printf("fileCursor.node.id: %v\n", fileCursor.node.id)
		// fmt.Printf("fileCursor.node.size: %v\n", fileCursor.node.size)
		// fmt.Printf("fileCursor.node.prevNode: %v\n", fileCursor.node.prevNode.id)

		emptySpaceCursor := x.createCursor()
		// find empty space with appropiate size
		for emptySpaceCursor.node.id != -1 || emptySpaceCursor.node.size < fileCursor.node.size {
			// fmt.Printf("emptySpaceCursor.pos: %v\n", emptySpaceCursor.pos)
			// fmt.Printf("emptySpaceCursor.size: %v\n", emptySpaceCursor.node.size)
			emptySpaceCursor.next()
			// fmt.Printf("emptySpaceCursor.pos: %v\n", emptySpaceCursor.pos)
			// fmt.Printf("emptySpaceCursor.size: %v\n", emptySpaceCursor.node.size)
			if emptySpaceCursor.pos >= fileCursor.pos {
				// fmt.Println("Way passed")
				break
			}
		}

		if emptySpaceCursor.pos >= fileCursor.pos {
			fileCursor.prev()
		} else if emptySpaceCursor.node.size >= fileCursor.node.size {
			// fmt.Println("swapping")
			newEmptySpaceSize := emptySpaceCursor.node.size - fileCursor.node.size
			if newEmptySpaceSize == 0 {
				swap(emptySpaceCursor.node, fileCursor.node)
			} else {
				swap(emptySpaceCursor.node, fileCursor.node)
				emptySpaceCursor.node.size -= newEmptySpaceSize
				// fmt.Printf("newEmptySpaceSize: %v\n", newEmptySpaceSize)
				newEmptySpace1 := Node{
					id: -1,
					size: newEmptySpaceSize,
					prevNode: fileCursor.node,
					nextNode: fileCursor.node.nextNode,
				}
				newEmptySpace1.nextNode.prevNode = &newEmptySpace1
				fileCursor.node.nextNode = &newEmptySpace1
				fileCursor.pos++
			}
			fileCursor.node = emptySpaceCursor.node
		}
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
	i := x2.prevNode
	x2.prevNode = x1.prevNode
	if x2.prevNode != nil {
		x2.prevNode.nextNode = x2
	}
	x1.prevNode = i
	if x1.prevNode != nil {
		x1.prevNode.nextNode = x1
	}
	j := x2.nextNode
	x2.nextNode = x1.nextNode
	if x2.nextNode != nil {
		x2.nextNode.prevNode = x2
	}
	x1.nextNode = j
	if x1.nextNode != nil {
		x1.nextNode.prevNode=x1
	}
}