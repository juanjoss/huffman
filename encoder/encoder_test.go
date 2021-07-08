package encoder

import (
	"bytes"
	"io/ioutil"
	"sort"
	"testing"

	model "github.com/fuato1/huffman/model"
)

func TestEncoder(t *testing.T) {
	// suite
	suite := map[string][]byte{
		"../test_files/data.txt": []byte("111010100"),
	}

	for file, expected := range suite {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("error reading file: %s.", file)
			t.FailNow()
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
		tree := Encode(pairs)
		stream := BuildStream(tree, []byte{})

		// assert
		// expected = 1 1 1 01 01 00
		if bytes.Compare(expected, stream) != 0 {
			t.Errorf("test %s failed: \n\tresult: %s\n\texpected: %s", file, stream, expected)
		}
	}
}
