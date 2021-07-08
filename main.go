package main

import (
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("./test_files/data.txt")
	if err != nil {
		panic(err)
	}

	symFreq := make(map[rune]int)

	for _, c := range data {
		symFreq[rune(c)]++
	}

	tree := encode(symFreq)

	printCodes(tree, []byte{})
}
