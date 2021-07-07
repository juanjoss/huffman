package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	symFreq := make(map[rune]int)

	for _, c := range data {
		symFreq[rune(c)]++
	}

	tree := encode(symFreq)
	fmt.Println(tree)
}
