package fuzzyfs

import (
	//"fmt"
	"testing"
)

func TestDirList(t *testing.T) {

	llista := NewDirList()

	a1 := llista.Add("a1", nil)
	a2 := llista.Add("a2", a1)
	a3 := llista.Add("a3", a1)
	_ = llista.Add("a4", a3)
	a5 := llista.Add("a5", a3)

	b := llista.List[0]
	//fmt.Println(llista)
	if b != *a1 {
		t.Errorf("%#v, %#v", b, a3)
	}

	r := llista.Get("a2", -1)
	if *r != *a2 {
		t.Errorf("%s, %s", *r, *a2)
	}

	r1 := llista.Get("a3", 1)
	if r1 != a3 {
		t.Errorf("%s<>%s", &r1, &a3)
	}

	if r.Name != "a2" {
		t.Errorf("%s!=a2", r.Name)
	}

	if r.Parent != a1 {
		t.Errorf("%s!=%s", r.Parent, a1)
	}
	m := a5.Parents()
	if m[0] != a5 {
		t.Errorf("%s!=parents", m)
	}

	if m[2] != a1 {
		t.Errorf("%s!=parents", m)
	}

	if len(m) != 3 {
		t.Errorf("%s!=parents", m)
	}

	if a5.Path() != "a1/a3/a5" {
		t.Errorf("%s!=a1/a3/a5", a5.Path())
	}

	if a1.Path() != "a1/" {
		t.Errorf("%#s!=a1/", a1.Path())
	}

	if a1.calcDepth() != 0 {
		t.Errorf("Wrong depth a1 %s", a1.calcDepth())
	}

	if a2.calcDepth() != 1 {
		t.Errorf("Wrong depth a2 %s", a2.calcDepth())
	}

	if a5.calcDepth() != 2 {
		t.Errorf("Wrong depth a5 %s", a5.calcDepth())
	}

	if a5.Depth != 2 {
		t.Errorf("Wrong depth a5 %s", a5.calcDepth())
	}

	if a1.Depth != 0 {
		t.Errorf("Wrong depth a1 %s", a1.calcDepth())
	}

}
