package model

type Tree interface {
	Freq() int
}

type Leaf struct {
	Frequency int
	Value     rune
}

type Node struct {
	Frequency   int
	Left, Right Tree
}

func (leaf Leaf) Freq() int {
	return leaf.Frequency
}

func (node Node) Freq() int {
	return node.Frequency
}
