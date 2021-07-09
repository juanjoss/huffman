package bitwise

import (
	"fmt"
	"io"
	"log"
	"os"
)

func ReaderTest() {
	var file io.Reader = (*os.File)(nil)
	file, err := os.Open("./test_files/data")
	if err != nil {
		panic(err)
	}

	br := NewReader(file, MSB)
	var result []byte
	for {
		v, err := br.ReadBits(8)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		result = append(result, byte(v))
	}

	fmt.Printf("bitwise reader test: \nRead from ./test_files/data back as %08b\n\n", result)
}
