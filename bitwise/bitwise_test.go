package bits

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testString = "111010100"
	outputFile = "./data"
)

func TestWriter(t *testing.T) {
	var stream bytes.Buffer
	bw := NewWriter(&stream, MSB)

	// writing testString
	for _, c := range testString {
		if c != '\n' {
			if c == '0' {
				err := bw.WriteBits(0x00, 1)
				if err != nil {
					t.Errorf("error writing '0' for c = %c\n", c)
				}
			} else {
				err := bw.WriteBits(0x01, 1)
				if err != nil {
					t.Errorf("error writing '1' for c = %c", c)
				}
			}
		}
	}

	// closing writer
	if err := bw.Close(); err != nil {
		t.Errorf("error closing: %v", err)
	}

	t.Logf("writing: %08b\n", stream.Bytes())

	// writing to file
	err := ioutil.WriteFile(outputFile, stream.Bytes(), 0644)
	if err != nil {
		t.Errorf("error writing file ./data: %v", err)
	}
}

func TestReader(t *testing.T) {
	// opening file created by writer
	var file io.Reader = (*os.File)(nil)
	file, err := os.Open(outputFile)
	if err != nil {
		t.Errorf("error opening file: %s\n", outputFile)
	}

	br := NewReader(file, MSB)
	var stream []byte

	// reading bits
	for {
		v, err := br.ReadBits(8)
		if err != nil {
			if err != io.EOF {
				t.Errorf("error after reading: %e\n", err)
			}
			break
		}
		stream = append(stream, byte(v))
	}

	t.Logf("read from %s back as %08b\n\n", outputFile, stream)
}
