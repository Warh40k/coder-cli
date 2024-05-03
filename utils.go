package coder_cli

import (
	"bytes"
	"fmt"
	"github.com/Warh40k/bw-coder/bwcoder"
	"io"
	"math"
	"os"
	"strings"
)

const CHUNK_SIZE = 8096

func TranslateSequence(seq *os.File) *bytes.Buffer {
	var chunk = make([]byte, CHUNK_SIZE) // чанк (в байтах)
	var translated bytes.Buffer
	var bitCount = int(math.Ceil(math.Log2(float64(CHUNK_SIZE))))

	for {
		var n, slen int
		slen, err := seq.Read(chunk)
		if err == io.EOF {
			break
		}
		var lcol = make([]byte, slen)
		n = bwcoder.Encode(chunk, lcol, slen)
		bnum := getBin(n, bitCount)
		translated.WriteString(bnum)
		translated.Write(lcol)
	}

	return &translated
}

func getBin(num, bitCount int) string {
	var numBit = 1
	if num != 0 {
		numBit = int(math.Log2(float64(num))) + 1
	}
	zeroCount := bitCount - numBit
	return strings.Repeat("0", zeroCount) + fmt.Sprintf("%b", num)
}
