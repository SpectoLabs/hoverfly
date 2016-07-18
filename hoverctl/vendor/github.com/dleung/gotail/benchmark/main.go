package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/dleung/gotail"
)

var fname string

func main() {
	flag.StringVar(&fname, "file", "testing.log", "File to tail")
	flag.Parse()

	log.Println("Running Benchmarks")
	var concurrency int = runtime.NumCPU()
	var rowcount = 5000000          // number of rows to write
	runtime.GOMAXPROCS(concurrency) // number of writer processes
	var row int

	_ = os.Remove(fname)
	createFile("")

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	tail, err := gotail.NewTail(fname, gotail.Config{Timeout: 10})
	time.Sleep(10 * time.Millisecond)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			log.Printf("Spawned %d Workers\n", i)

			for row <= rowcount {
				body := fmt.Sprintf("%d Worker doing the write.  %d iteration is being written.\n", i, row)
				writeContents(f, body)
				row++
			}
		}(i)
	}

	count := 0
	startTime := time.Now()

	go func() {
		for {
			startTime := time.Now()
			countNow := count
			time.Sleep(5 * time.Second)
			duration := time.Since(startTime).Seconds()
			newCount := count

			log.Printf("%d processed at %0.2f rows/sec\n", count, float64(newCount-countNow)/duration)
		}
	}()

	for _ = range tail.Lines {
		count++
		if count == rowcount {
			break
		}
	}

	duration := time.Since(startTime).Seconds()
	fmt.Printf("%d rows processed in %0.4f seconds, at the rate of %0.4f rows/s\n", count, duration, float64(count)/duration)
	removeFile()
}

func writeContents(f *os.File, contents string) {
	_, err := f.WriteString(contents)
	if err != nil {
		log.Fatalln(err)
	}
}

func createFile(contents string) {
	err := ioutil.WriteFile(fname, []byte(contents), 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeFile() {
	err := os.Remove(fname)
	if err != nil {
		log.Fatalln(err)
	}
}

func renameFile() {
	oldname := fname
	newname := fname + "_new"
	err := os.Rename(oldname, newname)
	if err != nil {
		log.Fatalln(err)
	}
}

func writeFile(contents string) {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	writeContents(f, contents)
}
