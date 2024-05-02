package internal

import (
	"bytes"
	"fmt"
	"github.com/Warh40k/bw-coder/bwcoder"
	"io"
	"math"
	"os"
	"strings"
)

const CHUNK_SIZE = 4096

func GetSequence(inputPath string) (bytes.Buffer, error) {

	var chunk = make([]byte, CHUNK_SIZE) // чанк (в байтах)
	var result bytes.Buffer
	var bitCount = int(math.Ceil(math.Log2(float64(CHUNK_SIZE))))

	input, err := os.Open(inputPath)
	if err != nil {
		return result, err
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
		result.Write(lcol)
	}

	return result, nil
}

func getBin(num, bitCount int) string {
	var numBit = 1
	if num != 0 {
		numBit = int(math.Log2(float64(num))) + 1
	}
	zeroCount := bitCount - numBit
	return strings.Repeat("0", zeroCount) + fmt.Sprintf("%b", num)
}
