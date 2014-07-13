package fuzzyfs

import (
	"github.com/sajari/fuzzy"
	//"log"
	"fmt"
	"os"
	"sort"
	"strings"
)

// NewDirList generates a new list of dirs...

// For populating it with files you must
//    dirlist.Populate( path, 0)
//  path must be a valid path, and 0 indicates the starting level
//
//  MaxDepth is used for specifing the recursion level where:
//      0 indicates only the main dir
//      n indicates the level of subdirs
//  PathSelect func for filtering dir input
//
//      list.Query(q string)
func NewDirList() *DirList {
	return &DirList{
		length: 0,
		List:   make([]Dir, 100),
		// As a default we dont limit the depth crawl
		MaxDepth:   1000,
		PathSelect: OnlyDirsPathSelect,
	}
}

// DirList is the main object..
type DirList struct {
	length     int
	List       []Dir
	MaxDepth   int
	PathSelect PathSelectFn
}

// Size of index
func (d *DirList) Length() int {
	return d.length
}

// Add a filesystem entry to the List
func (d *DirList) Add(name string, p *Dir) *Dir {
	a := Dir{
		Name:   name,
		parent: p}

	a.Depth = a.calcDepth()

	if len(d.List) == d.length {
		b := make([]Dir, len(d.List)+100)
		copy(b, d.List)
		d.List = b
	}

	d.List[d.length] = a
	d.length = d.length + 1
	return &d.List[d.length-1]
}

// Traverses Dirlist, and returns a File object for a given name and depth
// if depth == -1 returns the first ocurrence of name, without matching depth
func (d *DirList) Get(name string, depth int) *Dir {
	for k := range d.List {
		if depth == -1 {
			if d.List[k].Name == name {
				return &d.List[k]
			}
		} else {
			if d.List[k].Name == name && d.List[k].Depth == depth {
				return &d.List[k]
			}
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
		res := d.Add(f[k], parent)
		// Crawl till MaxDepth
		var err error
		if d.MaxDepth > res.Depth+1 {
			err = d.Populate(path+f[k]+"/", res)
		}
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
		if d.PathSelect(dirname, dir[k]) {
			names = append(names, dir[k].Name())
		}
	}

	//sort.Strings(names)
	return names, nil
}

// Check list against Levenshtein distance and store in
// Results
func (d *DirList) Query(q string, umbral int) Results {
	// Levenshtein
	res := Results{}

	for r := 0; r < d.length; r++ {

		word1 := strings.ToLower(d.List[r].Name)
		word2 := strings.ToLower(q)

		var v int
		if word1 == word2 {
			v = 0
		} else if strings.HasPrefix(word1, word2) {
			v = 1
		} else if strings.Contains(word1, word2) {
			v = 2
		} else {
			v = fuzzy.Levenshtein(&word1, &word2)
		}

		if v <= umbral {
			pa := Result{d.List[r].Path(), v, d.List[r].Depth}
			res = append(res, pa)
		}
	}

	sort.Sort(Results(res))
	return res

}

// Path selector, is a func for determining if a given path is
// choosed to be on the index
type PathSelectFn func(path string, info os.FileInfo) bool

// Only crawl directories
func OnlyDirsPathSelect(path string, info os.FileInfo) bool {

	if info.IsDir() {
		return true
	}
	return false
}

// DirsAndSymlinksasDirs threat symlinks that are dirs, as dirs,
// and follow them.
func DirsAndSymlinksAsDirs(path string, info os.FileInfo) bool {
	// Discard hidden files and folders
	if strings.HasPrefix(info.Name(), ".") {
		return false
	}

	if info.IsDir() {
		return true
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		fx, err := os.Readlink(path + info.Name())
		if err != nil {
			return false
		}
		fxi, err := os.Stat(fx)
		if err != nil {
			return false
		}
		if fxi.IsDir() {
			return true
		}
	}

	return false
}

// AllFiles match all files..
func AllFiles(path string, info os.FileInfo) bool {
	if strings.HasPrefix(info.Name(), ".") {
		return false
	}
	return true
}

// A Filesystem entry, we store, the node name and
// Their parents. We only need this info, for
// later build search and comparasions..
// Also we store only the current dir name...
// this way, when we traverse the list, searching,
// we can get better results..
// Given a path lik /asdf/1/2/asdf we store each segment
//
type Dir struct {
	Name   string
	parent *Dir
	Depth  int
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
		t = t + s[i].Name
		if i > 0 {
			t += "/"
		}

	}
	return t
}

func (d *Dir) calcDepth() int {
	if d.parent == nil {
		return 0
	}
	depth := 1
	p := d.parent
	for p.parent != nil {
		depth += 1
		p = p.parent
	}
	return depth
}

// Result is a result for querys
type Result struct {
	Path     string `json:"path"`
	Distance int
	Depth    int `json:"depth"`
}

func (r Result) String() string {
	return fmt.Sprintf("%d: %s\n", r.Distance, r.Path)
}

// Results is the list of resultsets
type Results []Result

// Sort interface
func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Sort by distance and depth
func (r Results) Less(i, j int) bool {
	if r[i].Distance == r[j].Distance {
		return r[i].Depth < r[j].Depth
	}
	return r[i].Distance < r[j].Distance
}
