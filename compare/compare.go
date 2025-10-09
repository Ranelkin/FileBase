package compare

import (
	"bufio"
	"compress/gzip"
	"os"
)

// Compares each path entry in 2 files
// using set Operations and returns the files that
// are in the first file and not in the second
func Difference(f1 string, f2 string) (*[]string, error) {
	// Channels to receive results and errors from goroutines
	type result struct {
		set map[string]bool
		err error
	}
	ch1 := make(chan result)
	ch2 := make(chan result)

	// Run setFromFile in goroutines
	go func() {
		s, err := setFromFile(f1)
		ch1 <- result{set: s, err: err}
	}()
	go func() {
		s, err := setFromFile(f2)
		ch2 <- result{set: s, err: err}
	}()

	// Collect results from both goroutines
	r1 := <-ch1
	if r1.err != nil {
		return nil, r1.err
	}
	r2 := <-ch2
	if r2.err != nil {
		return nil, r2.err
	}

	// Compute difference
	var diff []string
	for k := range r1.set {
		if !r2.set[k] {
			diff = append(diff, k)
		}
	}

	return &diff, nil
}

// Creates a set from a compressed file which was created
// previously by the Traverse() method
func setFromFile(f string) (map[string]bool, error) {
	s := make(map[string]bool)
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		s[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return s, nil
}
