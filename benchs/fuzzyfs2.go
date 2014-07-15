package main

import (
	"fmt"
	"github.com/sajari/fuzzy"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var result []Result
var q string
var ct int

func main() {

	startTime := time.Now()
	ct = 0

	var wg sync.WaitGroup
	workers := 32
	wg.Add(workers)
	worker := func(paths chan string) {
		defer wg.Done()
		for p := range paths {
			work(p)
		}
	}

	queue := make(chan string, 2*workers)
	for i := 0; i < workers; i++ {
		go worker(queue)
	}

	visitor := func(path string, info os.FileInfo, err error) error {
		if err == nil {
			ct = ct + 1
			queue <- path
		}
		return nil
	}

	result = []Result{}
	q = strings.ToLower(os.Args[1])

	for _, p := range os.Args[2:] {
		filepath.Walk(p, visitor)
	}

	close(queue)
	wg.Wait()
	sort.Sort(Results(result))

	endTime := time.Now()
	fmt.Print(fmt.Sprint("ElapsedTime in seconds:", endTime.Sub(startTime), "\n"))
	fmt.Printf("Walked, %d files\n\n", ct)
	fmt.Printf("%s", result)
	os.Exit(0)

}

func work(p string) {

	var level, v int
	//fmt.Println(p)

	inf, err := os.Stat(p)
	if err != nil {
		return
	}

	segments := strings.Split(p, "/")
	w1 := strings.ToLower(segments[len(segments)-1])

	if strings.HasPrefix(w1, ".") {
		return
	}

	if w1 == q {
		v = 0
	} else if strings.HasPrefix(w1, q) {
		v = 1
	} else if strings.Contains(w1, q) {
		v = 2
	} else {
		v = fuzzy.Levenshtein(&w1, &q)
	}
	// if umbral
	if v <= 2 {
		if inf.IsDir() {
			level = len(segments)
		} else {
			level = len(segments) - 1
		}

		pa := Result{w1, v, level, p}
		result = append(result, pa)
	}

}

type Result struct {
	Path     string `json:"path"`
	Distance int
	Depth    int `json:"depth"`
	FPath    string
}

func (r Result) String() string {
	return fmt.Sprintf("%d: %s\n", r.Distance, r.FPath)
}

// Results is the list of resultsets
type Results []Result

// Sort interface
func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Sort by distance and depth
func (r Results) Less(i, j int) bool {
	if r[i].Distance == r[j].Distance {
		return r[i].Depth > r[j].Depth
	}
	return r[i].Distance > r[j].Distance
}
