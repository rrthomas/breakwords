///usr/bin/env go run $0 $@ ; exit
// Adapted from https://rosettacode.org/wiki/Word_break_problem#Go

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var progname = "breakwords"
var version = "0.1 (03 Oct 2019)"
var author = "Reuben Thomas <rrt@sc3d.org>"

// Command-line arguments
var versionFlag *bool = flag.Bool("version", false, "output version information and exit")
var helpFlag *bool = flag.Bool("help", false, "display this help and exit")

func usage() {
	os.Stderr.WriteString(progname + " " + version + "\n\n" +
		"Add spaces between words in text.\n\n" +
		"Usage: " + progname + " WORDLIST-FILE\n\n")
	flag.PrintDefaults()
}

func showVersion() {
	os.Stderr.WriteString(progname + " " + version + " " + author + "\n")
}

// Read an os.File into memory and return a slice of its lines.
func readLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type dict map[string]bool

func newDict(words []string) dict {
	d := dict{}
	for _, w := range words {
		d[w] = true
	}
	return d
}

func (d dict) wordBreak(s string) (broken []string, ok bool) {
	if s == "" {
		return nil, true
	}
	type prefix struct {
		length int
		broken []string
	}
	bp := []prefix{{0, nil}}
	for end := 1; end <= len(s); end++ {
		for i := 0; i < len(bp); i++ {
			w := s[bp[i].length:end]
			if d[w] {
				b := append(bp[i].broken, w)
				if end == len(s) {
					return b, true
				}
				bp = append(bp, prefix{end, b})
				break
			}
		}
	}
	return nil, false
}

func main() {
	// Parse command-line arguments
	flag.Parse()
	if *versionFlag {
		showVersion()
		os.Exit(0)
	}
	if *helpFlag {
		usage()
		os.Exit(0)
	}
	args := flag.Args()
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
	wordlist_file := args[0]

	// Read wordlist and create dictionary
	h, err := os.Open(wordlist_file)
	if err != nil {
		panic(err)
	}
	words, err := readLines(h)
	if err != nil {
		panic(err)
	}
	d := newDict(words)

	// Read text
	lines, err := readLines(os.Stdin)
	if err != nil {
		panic(err)
	}
	if err:= h.Close(); err != nil {
		panic(err)
	}

	// Process text
	for _, s := range lines {
		if b, ok := d.wordBreak(s); ok {
			fmt.Printf("%s: %s\n", s, strings.Join(b, " "))
		} else {
			fmt.Println("can't break")
		}
	}
}
