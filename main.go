package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {

	numWorkers := runtime.NumCPU() * 8

	if numWorkers > 200 {
		numWorkers = 200
	}
	if numWorkers < 16 {
		numWorkers = 16
	}

	fmt.Printf("Using %d workers (%d CPU cores detected)\n", numWorkers, runtime.NumCPU())
	fmt.Println("Please paste the directory to traverse: ")

	dir := bufio.NewScanner(os.Stdin)
	dir.Scan()
	entryDir := strings.TrimSpace(dir.Text())
	fmt.Println("Start to traverse filebase: ", entryDir)

	pid := os.Getpid()
	fmt.Println("PId: ", pid)

	start := time.Now()

	result := make(chan string, 1000)
	queue := make(chan string, 1000)

	var wg sync.WaitGroup
	var dirWg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(queue, result, &wg, &dirWg)
	}

	//Send init dir
	dirWg.Add(1)
	queue <- entryDir

	go func() {
		dirWg.Wait()
		close(queue)
	}()

	go func() {
		wg.Wait()
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

	os.Remove(entryDir + "/.fileBase.txt")

	end := time.Now()
	fmt.Println("Finished traversing filebase in ", end.Sub(start).Minutes(), " minutes")
	fmt.Println("Filebase objects: ", count)
}

func worker(jobs chan string, results chan string, wg *sync.WaitGroup, dirWg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		if err := traverse(job, results, jobs, dirWg); err != nil {
			fmt.Fprintf(os.Stderr, "Error traversing %s: %v\n", job, err)
		}
	}
}

// Scans the content of the directory f and sends it to channel
// queue and result if its a directory, to result if its a file
func traverse(f string, result chan string, queue chan string, dirWg *sync.WaitGroup) error {
	defer dirWg.Done()
	dir, err := os.ReadDir(f)
	if err != nil {
		return err
	}
	for _, d := range dir {
		name := f + "/" + d.Name()
		result <- name
		if d.IsDir() {
			dirWg.Add(1)
			queue <- name
		}
	}
	return nil
}
