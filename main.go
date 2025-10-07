package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("Please paste the directory to traverse: ")
	dir := bufio.NewScanner(os.Stdin)
	dir.Scan()
	entryDir := strings.TrimSpace(dir.Text())
	fmt.Println("Start to traverse filebase: ", entryDir)
	pid := os.Getpid()
	fmt.Println("PId: ", pid)

	start := time.Now()
	result := make(chan string)
	queue := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		queue <- entryDir
	}()

	go func() {
		for f := range queue {
			fCopy := f
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				if err := traverse(&fCopy, result, queue, &wg); err != nil {
					panic(err)
				}
			}(f)
		}
	}()

	go func() {
		wg.Wait()
		close(queue)
		close(result)
	}()
	// File gets truncated if it already exists
	f, err := os.Create(entryDir + "/.fileBase.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int = 0
	for entry := range result {
		f.Write([]byte((entry + "\n")))
		count++
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	dst, err := os.Create(entryDir + "/.fileBase.txt.gz")
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	zip, err := gzip.NewWriterLevel(dst, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer zip.Close()

	if _, err := io.Copy(zip, f); err != nil {
		panic(err)
	}

	end := time.Now()
	fmt.Println("Finished traversing filebase in ", end.Sub(start).Minutes(), " minutes")
	fmt.Println("Filebase objects: ", count)
}

// Scans the content of the directory f and sends it to channel
// queue and result if its a directory, to result if its a file
func traverse(f *string, result chan string, queue chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	dir, err := os.ReadDir(*f)
	if err != nil {
		return err
	}
	for _, d := range dir {
		name := *f + "/" + d.Name()
		result <- name
		if d.IsDir() {
			wg.Add(1)
			queue <- name
		}
	}
	return nil
}
