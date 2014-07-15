package main

import (
	"flag"
	"fmt"
	"github.com/jordic/fuzzyfs"
	"net/http"
	//"strings"
	"time"
)

var llista *fuzzyfs.DirList

var (
	depth   = flag.Int("depth", 3, "Directory depth to search")
	dirpath = flag.String("dir", ".", "Root path to start")
	mode    = flag.Int("mode", 1, "1. Only dirs are indexed. 2. All files indexed")
	umbral  = flag.Int("umbral", 3, "Umbral to discard entries..")
)

func main() {

	/*if strings.HasSuffix(*dir, "/") == false {
		*dir = *dir + "/"
	}*/
	flag.Parse()

	llista = fuzzyfs.NewDirList()
	llista.MaxDepth = *depth
	if *mode == 1 {
		llista.PathSelect = fuzzyfs.DirsAndSymlinksAsDirs
	} else {
		llista.PathSelect = fuzzyfs.AllFiles
	}

	go func() {
		startTime := time.Now()
		fmt.Printf("Building index...%s\n", *dirpath)
		err := llista.Populate(*dirpath, nil)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		fmt.Printf("Total list entries: %d\n", llista.GetLength())
		fmt.Println("time indexing...", endTime.Sub(startTime))

	}()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRequest))
	http.ListenAndServe("localhost:8888", mux)

	//pprof.WriteHeapProfile(f)
	//fmt.Print(llista)
	//for r := 0; r < llista.length; r++ {
	//  fmt.Println(llista.List[r].Path())
	//}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	var query string
	query = r.FormValue("q")

	if len(query) < 3 {
		fmt.Fprint(w, "Query tooo short!")
		return
	}

	res := llista.Query(query, *umbral)

	//fmt.Println(res)

	endTime := time.Now()
	fmt.Fprint(w, fmt.Sprint("ElapsedTime in seconds:", endTime.Sub(startTime), "\n"))

	fmt.Fprint(w, fmt.Sprintf("%s", res))
	return
}
