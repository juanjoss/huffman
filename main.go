package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/fuato1/huffman/encoder"
	"github.com/fuato1/huffman/model"
)

func main() {
	// reading file
	data, err := ioutil.ReadFile("./test_files/data.txt")
	if err != nil {
		panic(err)
	}

	// symbol's frequencies
	symFreq := make(map[rune]int)
	pairs := make(model.Pairs, 0)

	for _, c := range data {
		if c != '\n' {
			symFreq[rune(c)]++
		}
	}

	for k, v := range symFreq {
		pairs = append(pairs, model.Pair{k, v})
	}
	sort.Sort(pairs)

	// encoding
	tree := encoder.Encode(pairs)
	stream := encoder.BuildStream(tree, []byte{})

	// expected = 1 1 1 01 01 00
	fmt.Println(string(stream))
}
