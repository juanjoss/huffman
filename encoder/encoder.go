package encoder

import (
	"container/heap"
	"fmt"

	model "github.com/fuato1/huffman/model"
)

// encoding

func Encode(symFreq model.Pairs) model.Tree {
	var tree model.PriorityQueue

	for _, pair := range symFreq {
		tree = append(tree, model.Leaf{pair.Frequency, pair.Symbol})
	}
	heap.Init(&tree)

	for tree.Len() > 1 {
		x := heap.Pop(&tree).(model.Tree)
		y := heap.Pop(&tree).(model.Tree)

		heap.Push(&tree, model.Node{x.Freq() + y.Freq(), x, y})
	}

	return heap.Pop(&tree).(model.Tree)
}

// tree traversal

func BuildStream(tree model.Tree, stream []byte) []byte {
	switch node := tree.(type) {
	case model.Leaf:
		fmt.Printf("%c\t%d\t%s\n", node.Value, node.Freq(), string(stream))
		return stream

	case model.Node:
		// left
		stream = append(stream, '0')
		stream = append(stream, BuildStream(node.Left, stream)...)
		stream = stream[:len(stream)-1]

		// right
		stream = append(stream, '1')
		stream = append(stream, BuildStream(node.Right, stream)...)
		stream = stream[:len(stream)-1]
	}

	return stream
}
