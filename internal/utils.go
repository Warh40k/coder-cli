package internal

import (
	"bytes"
	"fmt"
	bwcoder "github.com/Warh40k/bw-coder/coder"
	"io"
	"math"
	"os"
	"strings"
)

const CHUNK_SIZE = 256

func GetSequence(inputPath string) []byte {

	var chunk = make([]byte, CHUNK_SIZE) // чанк (в байтах)
	var result bytes.Buffer
	var bitCount = int(math.Ceil(math.Log2(float64(CHUNK_SIZE))))

	input, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("error opening input file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()

	for {
		var n, slen int
		slen, err = input.Read(chunk)
		if err == io.EOF {
			break
		}
		var lcol = make([]byte, slen)
		n = bwcoder.Encode(chunk, lcol, slen)
		bnum := getBin(n, bitCount)
		result.WriteString(bnum)
		result.Write(chunk)
	}

	return result.Bytes()
}

func getBin(num, bitCount int) string {
	var numBit = 1
	if num != 0 {
		numBit = int(math.Log2(float64(num))) + 1
	}
	zeroCount := bitCount - numBit
	return strings.Repeat("0", zeroCount) + fmt.Sprintf("%b", num)
}
