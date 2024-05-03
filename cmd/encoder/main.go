package main

import (
	"fmt"
	"github.com/Warh40k/bookstack-coding/bookstack"
	"github.com/Warh40k/entropy"
	"github.com/Warh40k/information_theory_lr1"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
)

var isDir bool

type tabWriter struct {
	w *tabwriter.Writer
	m sync.Mutex
}

func (tw *tabWriter) Write(p []byte) (n int, err error) {
	tw.m.Lock()
	defer tw.m.Unlock()
	return tw.w.Write(p)
}

func (tw *tabWriter) Flush() {
	tw.m.Lock()
	defer tw.m.Unlock()
	tw.w.Flush()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: encoder <input path> <output path>")
		os.Exit(1)
	}

	info, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Printf("error opening input file: %s\n", err)
		os.Exit(1)
	}
	var inFiles []string
	var dirs []string
	if info.IsDir() {
		isDir = true
		err = filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				dirs = append(dirs, path)
			} else {
				inFiles = append(inFiles, path)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error getting dir files: %s\n", err)
			os.Exit(1)
		}
	} else {
		inFiles = append(inFiles, os.Args[1])
	}

	info, err = os.Stat(os.Args[2])
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(os.Args[2]), 0777)
		if err != nil {
			fmt.Printf("error creating output directory: %s\n", err)
			os.Exit(1)
		}
	}
	for _, dir := range dirs {
		err = os.MkdirAll(filepath.Join(os.Args[2], strings.Split(dir, os.Args[1])[1]), 0777)
		if err != nil {
			fmt.Printf("error creating output directory: %s\n", err)
			os.Exit(1)
		}
	}

	wg := sync.WaitGroup{}
	var out = &tabWriter{}
	out.w = tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	defer out.Flush()

	fmt.Fprintf(out, "File\tSize\tH(X)\tH(X|X)\tH(X|XX)\tl_avg\tNewSize\n")

	for i := 0; i < len(inFiles); i++ {
		wg.Add(1)
		go processFile(inFiles[i], &wg, out)
	}
	wg.Wait()
}

func processFile(path string, wg *sync.WaitGroup, out *tabWriter) {
	defer wg.Done()

	input, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(out, "error opening input file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()

	data, err := io.ReadAll(input)

	input.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Fprintf(out, "error reading input file: %s\n", err)
		os.Exit(1)
	}

	// output initial file info
	info, err := input.Stat()
	if err != nil {
		fmt.Fprintf(out, "error stat input file: %s\n", err)
		os.Exit(1)
	}

	translatedSeq := information_theory_lr1.TranslateSequence(input)
	encodedSeq := bookstack.Encode(translatedSeq.Bytes())

	// save result
	var outPath = os.Args[2]
	if isDir {
		outPath = filepath.Join(os.Args[2], strings.Split(path, os.Args[1])[1])
	}

	encinfo, err := bookstack.SaveSequence(outPath, encodedSeq)
	if err != nil {
		fmt.Printf("error creating output file: %s\n", err)
		os.Exit(1)
	}
	entropyInfo := dumpEntropyInfo(path, data, info.Size())
	fmt.Fprintf(out, "%s\t%f\t%d\n", entropyInfo,
		getAverageCodeLen(len(encodedSeq), len([]rune(string(data)))), encinfo.Size())
}

func dumpEntropyInfo(path string, data []byte, size int64) string {
	freqs, probs := entropy.GetFreqsProbs(data)
	entr := entropy.GetEntropy(probs)
	condProbs, condFreqs := entropy.GetCondProbs(data, freqs)
	condEntr := entropy.GetCondEntropy(probs, condProbs)
	condProbsXX := entropy.GetCondProbsXX(data, condFreqs)
	condEntrXX := entropy.GetCondEntropyXX(probs, condProbs, condProbsXX)
	return fmt.Sprintf("%s\t%d\t%f\t%f\t%f",
		path, size, entr, condEntr, condEntrXX)
}

func getAverageCodeLen(seqLen, symCount int) float64 {
	return float64(seqLen) * 8 / float64(symCount)
}
