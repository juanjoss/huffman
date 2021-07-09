package encoder

import (
	"container/heap"

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

func BuildLookupTable(tree model.Tree, code []byte, codes map[rune][]byte) map[rune][]byte {
	switch node := tree.(type) {
	case model.Leaf:
		//fmt.Printf("%c\t%d\t%s\n", node.Value, node.Freq(), string(code))
		codes[node.Value] = append(codes[node.Value], code...)

		return codes

	case model.Node:
		// left
		code = append(code, '0')
		codes = BuildLookupTable(node.Left, code, codes)
		code = code[:len(code)-1]

		// right
		code = append(code, '1')
		codes = BuildLookupTable(node.Right, code, codes)
		code = code[:len(code)-1]
	}

	return codes
}
