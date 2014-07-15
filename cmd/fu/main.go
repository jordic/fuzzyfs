package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/jordic/fuzzyfs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	path       = flag.String("p", "/", "Directory to search")
	reindex    = flag.Bool("reindex", false, "Regenerate index")
	depth      = flag.Int("depth", 7, "Depth to crawl to index")
	query      = flag.String("q", "", "Query to search")
	umbral     = flag.Int("umbral", 1, "Results umbral to output")
	verbose    = flag.Bool("v", false, "Verbose")
	index_file = flag.String("index", ".index.gob", "Index file to use")
	method     = flag.Int("method", 1, "Method of indexing 1. all files. 2 only dirs")
)

func main() {

	flag.Parse()

	dir := *path

	if _, err := os.Stat(dir); err != nil {
		log.Fatalf("Directory %s not exist", dir)
	}

	if strings.HasSuffix(dir, "/") == false {
		dir = dir + "/"
	}

	if *verbose == false {
		log.SetOutput(ioutil.Discard)
	}

	index_loaded := false
	var llista *fuzzyfs.DirList
	// Check if gob index exists... or command reindex
	if _, err := os.Stat(dir + *index_file); err != nil || *reindex == true {

		startTime := time.Now()
		log.Print("Regenerating index")
		llista = fuzzyfs.NewDirList()
		llista.MaxDepth = *depth

		if *method == 1 {
			llista.PathSelect = fuzzyfs.AllFiles
		} else {
			llista.PathSelect = fuzzyfs.DirsAndSymlinksAsDirs
			log.Print("Indexing... dirs and symlinks")
		}

		llista.Populate(dir, nil)

		log.Printf("Found %d files", llista.Length)

		f, err := os.Create(dir + *index_file)
		if err != nil {
			log.Fatalf("Unable to create index %s", err)
		}
		defer f.Close()

		enc := gob.NewEncoder(f)
		enc.Encode(&llista)

		endTime := time.Now()
		log.Printf("Index generated in %s\n", endTime.Sub(startTime))
		index_loaded = true

		if *reindex == true {
			os.Exit(0)
		}

	}

	startTime := time.Now()
	if index_loaded == false {
		f, err := os.Open(dir + *index_file)
		if err != nil {
			log.Fatalf("Unable to load index %s", err)
		}
		defer f.Close()

		enc := gob.NewDecoder(f)
		err = enc.Decode(&llista)
	}

	q := *query
	u := *umbral

	res := llista.Query(q, u)

	for _, r := range res {
		fmt.Println(dir + r.Path)
	}

	//fmt.Print(res)
	endTime := time.Now()
	log.Print(fmt.Sprint("ElapsedTime:", endTime.Sub(startTime), "\n"))
	log.Printf("Total files: %d\n", llista.Length)

}
