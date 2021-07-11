package encoder

import (
	"bytes"
	"container/heap"
	"io/ioutil"

	bits "github.com/fuato1/huffman/bitwise"
	model "github.com/fuato1/huffman/model"
)

const (
	HUFF_NUMBER = 0xface8200 // identifier for the decompressor
	PSEUDO_EOF  = 1 << 8     // identifies EOF, it's considered a leaf
)

type Encoder struct {
	w      *bits.Writer
	Stream bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{Stream: bytes.Buffer{}}
}

func (e *Encoder) Encode(symFreq model.Pairs) model.Tree {
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

/*
	Header at the beggining of file:

	- HUFF_NUMBER it's an int -> 0xface8200, so we use an int's size, 32 bits (uint32).
	HUFF_NUMBER is written at the beggining of the file as flag for the decompressor because
	it's unlikely that HUFF_NUMBER will occur at the beggining of a file.

	- Tree:
		- For leaf nodes ->
			write 1
			write 9 bits for the symbol:
				<symbol> (9 bits)

		- For internal nodes -> write 0
*/

func (e *Encoder) BuildLookupTable(tree model.Tree) map[rune][]byte {
	e.w = bits.NewWriter(&e.Stream, bits.MSB) // creating writer
	e.w.WriteBits(HUFF_NUMBER, 32)            // writing HUFF_NUMBER to e.Stream
	defer e.w.WriteBits(PSEUDO_EOF, 9)

	return e.buildTable(tree, []byte{}, map[rune][]byte{})
}

func (e *Encoder) buildTable(tree model.Tree, code []byte, codes map[rune][]byte) map[rune][]byte {
	switch node := tree.(type) {
	case model.Leaf:
		codes[node.Value] = append(codes[node.Value], code...)

		// writing header to e.Stream
		e.w.WriteBits(0x01, 1)
		e.w.WriteBits(uint32(node.Value), 8) // why not 9 bits like len(PSEUDO_EOF)?

		return codes

	case model.Node:
		e.w.WriteBits(0x00, 1)

		// left
		code = append(code, '0')
		codes = e.buildTable(node.Left, code, codes)
		code = code[:len(code)-1]

		// right
		code = append(code, '1')
		codes = e.buildTable(node.Right, code, codes)
		code = code[:len(code)-1]
	}

	return codes
}

func (e *Encoder) SaveEncodedData(file string, data []byte, codes map[rune][]byte) {
	// writing body to e.Stream
	for _, c := range data {
		if c != '\n' {
			if codes[rune(c)] != nil {
				for _, b := range codes[rune(c)] {
					if b == '0' {
						e.w.WriteBits(0x00, 1)
					} else {
						e.w.WriteBits(0x01, 1)
					}
				}
			}
		}
	}

	// writing PSEUDO_EOF to e.Stream
	e.w.WriteBits(PSEUDO_EOF, 9)

	var err error
	if err := e.w.Close(); err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(file, e.Stream.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
