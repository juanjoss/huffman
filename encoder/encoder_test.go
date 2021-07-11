package encoder

import (
	"bytes"
	"io/ioutil"
	"sort"
	"testing"

	model "github.com/fuato1/huffman/model"
)

/*
	This test just takes care of comparing the encoded data
	with the expected theoretical result, and doesn't tests
	the generated e.Stream (HUFF_NUMBER + header + body + PSEUDO_EOF)
	of bits which are written to the file after encoding.
*/
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
		enc := NewEncoder()
		tree := enc.Encode(pairs)
		codes := enc.BuildLookupTable(tree)
		enc.SaveEncodedData("./data", data, codes)

		// encoding data
		var stream []byte
		for _, c := range data {
			if c != '\n' {
				if codes[rune(c)] != nil {
					for _, b := range codes[rune(c)] {
						if b == '0' {
							stream = append(stream, '0')
						} else {
							stream = append(stream, '1')
						}
					}
				}
			}
		}

		// assert
		if bytes.Compare(expected, stream) != 0 {
			t.Errorf("test %s failed: \n\tresult: %s\n\texpected: %s", file, stream, expected)
		}
	}
}
