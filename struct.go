package main

import (
	"container/heap"
	"fmt"
)

type Tree interface {
	Freq() int
}

type Leaf struct {
	freq  int
	value rune
}

type Node struct {
	freq        int
	left, right Tree
}

func (leaf Leaf) Freq() int {
	return leaf.freq
}

func (node Node) Freq() int {
	return node.freq
}

// queue

type PriorityQueue []Tree

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq() < pq[j].Freq()
}

func (pq *PriorityQueue) Push(e interface{}) {
	*pq = append(*pq, e.(Tree))
}

func (pq *PriorityQueue) Pop() (e interface{}) {
	e = (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return
}

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

// encoding

func encode(symFreq map[rune]int) Tree {
	var tree PriorityQueue

	for c, freq := range symFreq {
		tree = append(tree, Leaf{freq, c})
	}
	heap.Init(&tree)

	for tree.Len() > 1 {
		x := heap.Pop(&tree).(Tree)
		y := heap.Pop(&tree).(Tree)

		heap.Push(&tree, Node{x.Freq() + y.Freq(), x, y})
	}

	return heap.Pop(&tree).(Tree)
}

// decodig

func decode(tree Tree) map[rune]int {
	return nil
}

// build lookup table

func printCodes(tree Tree, code []byte) {
	switch node := tree.(type) {
	case Leaf:
		fmt.Printf("%c\t%d\t%s\n", node.value, node.freq, string(code))

	case Node:
		// left
		code = append(code, '0')
		printCodes(node, code)
		code = code[:len(code)-1]

		// right
		code = append(code, '1')
		printCodes(node, code)
		code = code[:len(code)-1]
	}
}
