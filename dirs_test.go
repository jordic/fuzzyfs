package main

import (
	//"fmt"
	"testing"
)

func TestDirList(t *testing.T) {

	llista := NewDirList(100)

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

	r := llista.Get("a2")
	if *r != *a2 {
		t.Errorf("%s, %s", *r, *a2)
	}

	r1 := llista.Get("a3")
	if r1 != a3 {
		t.Errorf("%s<>%s", &r1, &a3)
	}

	if r.Name != "a2" {
		t.Errorf("%s!=a2", r.Name)
	}

	if r.parent != a1 {
		t.Errorf("%s!=%s", r.parent, a1)
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

	if a5.Path() != "a1/a3/a5/" {
		t.Errorf("%s!=a1/a3/a5/", a5.Path())
	}

	if a1.Path() != "a1/" {
		t.Errorf("%#s!=a1/", a1.Path())
	}

}
