package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/fuato1/huffman/bitwise"
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
	codes := encoder.BuildLookupTable(tree, []byte{}, map[rune][]byte{})

	// writing bits buffer to file
	var stream bytes.Buffer
	bw := bitwise.NewWriter(&stream, bitwise.MSB)

	for _, c := range data {
		if c != '\n' {
			if codes[rune(c)] != nil {
				for _, b := range codes[rune(c)] {
					if b == '0' {
						bw.WriteBits(0x00, 1)
					} else {
						bw.WriteBits(0x01, 1)
					}
				}
			}
		}
	}

	if err := bw.Close(); err != nil {
		panic(err)
	}
	fmt.Printf("encoded stream: %08b\n", stream.Bytes())

	err = ioutil.WriteFile("./test_files/data", stream.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	bitwise.ReaderTest()
}
