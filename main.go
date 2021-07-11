package main

import (
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
	enc := encoder.NewEncoder()
	tree := enc.Encode(pairs)
	codes := enc.BuildLookupTable(tree)
	enc.SaveEncodedData("./test_files/data", data, codes)
}
