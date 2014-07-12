package main

import (
	"fmt"
	"github.com/sajari/fuzzy"
	//"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// A Filesystem entry, we store, the node name and
// Their parents. We only need this info, por
// later build search and comparasions..
// Also we store only the current dir name...
// this way, when we traverse the list, searching,
// we can get better results..
// Given a path lik /asdf/1/2/asdf we store each segment
//
type Dir struct {
	Name   string
	parent *Dir
}

// String
func (d *Dir) String() string {
	return d.Name
}

// Parents returns a slice of element Parents
func (d *Dir) Parents() []*Dir {

	if d.parent == nil {
		return nil
	}

	s := []*Dir{}
	s = append(s, d)
	p := d.parent
	for p.parent != nil {
		s = append(s, p)
		p = p.parent
	}
	s = append(s, p)
	return s
}

// Path rejoins the path with their ancestors
func (d *Dir) Path() string {
	if d.parent == nil {
		return d.Name + "/"
	}
	s := d.Parents()
	t := ""
	for i := len(s) - 1; i >= 0; i-- {
		t = t + s[i].Name + "/"
	}
	return t
}

// DirList Holds a dir listing
type DirList struct {
	length int
	List   []Dir
}

// NewDirList generates a new list of thirs
func NewDirList(c int) *DirList {
	return &DirList{
		length: 0,
		List:   make([]Dir, c),
	}
}

// Add a filesystem entry to the List
func (d *DirList) Add(name string, p *Dir) *Dir {
	a := Dir{
		Name:   name,
		parent: p}

	if len(d.List) == d.length {
		b := make([]Dir, len(d.List)*2)
		copy(b, d.List)
		d.List = b
	}

	d.List[d.length] = a
	d.length = d.length + 1
	return &d.List[d.length-1]
}

// Get and entry by name... Not stable, and
// and not for use... because, don't take care
// of duplicates... neither levels...
func (d *DirList) Get(name string) *Dir {
	for k := range d.List {
		if d.List[k].Name == name {
			return &d.List[k]
		}
	}
	return nil
}

// Populate acually traverses the filsystem, storing dir info
// on the list..
func (d *DirList) Populate(path string, parent *Dir) error {

	f, err := d.readDirNames(path)
	if err != nil {
		return err
	}

	for k := range f {
		if strings.HasPrefix(f[k], ".") {
			continue
		}
		res := d.Add(f[k], parent)
		err := d.Populate(path+f[k]+"/", res)
		if err != nil {
			continue
		}
	}

	return nil
}

func (d *DirList) readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	dir, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	names := []string{}
	for k := range dir {
		if dir[k].IsDir() {
			names = append(names, dir[k].Name())
		}
	}

	sort.Strings(names)
	return names, nil
}

// Result is a result for querys
type Result struct {
	Path     string
	Distance int
}

func (r Result) String() string {
	return fmt.Sprintf("%d: %s\n", r.Distance, r.Path)
}

// Results is the list of resultsets
type Results []Result

// Sort interface
func (r Results) Len() int           { return len(r) }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Results) Less(i, j int) bool { return r[i].Distance < r[j].Distance }

func (d *DirList) Query(q string, umbral int) Results {
	fmt.Printf("Query list of %d\n", d.length)
	// Levenshtein
	res := Results{}
	for r := 0; r < d.length; r++ {
		v := fuzzy.Levenshtein(&d.List[r].Name, &q)
		//fmt.Printf("Processing %s %d\n", d.List[r].Name, v)
		if v <= umbral {
			pa := Result{d.List[r].Path(), v}
			res = append(res, pa)
		}
	}

	return res

}

var llista *DirList

func main() {

	llista = NewDirList(100)
	go func() {
		startTime := time.Now()
		fmt.Printf("Building index...\n")
		err := llista.Populate(os.Args[1], nil)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		fmt.Printf("\nTotal list entries: %s\n", llista.length)
		fmt.Println("time indexing...", endTime.Sub(startTime))

	}()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRequest))
	http.ListenAndServe("localhost:8888", mux)

	//pprof.WriteHeapProfile(f)
	//fmt.Print(llista)
	//for r := 0; r < llista.length; r++ {
	//	fmt.Println(llista.List[r].Path())
	//}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	var query string
	query = r.FormValue("q")

	if len(query) < 4 {
		fmt.Fprint(w, "Query tooo short!")
		return
	}

	res := llista.Query(query, 3)

	//fmt.Println(res)

	sort.Sort(Results(res))

	endTime := time.Now()
	fmt.Fprint(w, fmt.Sprint("ElapsedTime in seconds:", endTime.Sub(startTime), "\n"))

	fmt.Fprint(w, fmt.Sprintf("%s", res))
	return
}
