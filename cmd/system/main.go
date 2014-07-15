package main

import (
	"encoding/gob"
	"fmt"
	"github.com/jordic/fuzzyfs"
	"os"
)

func main() {

	/*fmt.Println("Populating list")

	llista := fuzzyfs.NewDirList()
	llista.MaxDepth = 4
	llista.PathSelect = fuzzyfs.AllFiles

	llista.Populate(os.Args[1], nil)

	fmt.Printf("Found %d\n", llista.Length)

	f, err := os.Create("dirs.gob")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	enc.Encode(&llista)*/

	f, err := os.Open("dirs.gob")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewDecoder(f)
	var llista fuzzyfs.DirList
	err = enc.Decode(&llista)

	//for _, r := range llista.List {
	//	fmt.Println(r.Name)
	//}

	res := llista.Query(os.Args[1], 2)
	fmt.Print(res)
	fmt.Printf("%#v", llista.Length)

}
