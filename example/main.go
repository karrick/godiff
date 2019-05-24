package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/godiff"
)

func main() {
	if err := cmdMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func cmdMain() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("usage: %s file1 file2", filepath.Base(os.Args[0]))
	}

	first, err := linesFromFile(os.Args[1])
	if err != nil {
		return err
	}

	second, err := linesFromFile(os.Args[2])
	if err != nil {
		return err
	}

	if false {
		for _, s := range first {
			fmt.Fprintf(os.Stderr, "FIRST: %q\n", s)
		}
		for _, s := range second {
			fmt.Fprintf(os.Stderr, "SECOND: %q\n", s)
		}
	}

	for _, s := range godiff.Strings(first, second) {
		fmt.Println(s)
	}

	return nil
}

func linesFromFile(pathname string) (lines []string, err error) {
	var f *os.File

	f, err = os.Open(pathname)
	if err != nil {
		return
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	err = s.Err()

	if err2 := f.Close(); err == nil {
		err = err2
	}

	return
}
