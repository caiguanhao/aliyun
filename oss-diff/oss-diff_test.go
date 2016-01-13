package main

import "testing"

var testNo int

type TestDiffCase struct {
	Local        []File
	Remote       []File
	L2R_Expected []File
	R2L_Expected []File
}

func (c TestDiffCase) test(t *testing.T) {
	testNo++

	l2r := diff(c.Local, c.Remote)
	if len(l2r) != len(c.L2R_Expected) {
		t.Errorf("#%d l2r: diff should return %d files instead of %d", testNo, len(c.L2R_Expected), len(l2r))
	}
	for i := range l2r {
		if l2r[i].Name != c.L2R_Expected[i].Name {
			t.Errorf("#%d l2r: %s should be %s", testNo, l2r[i].Name, c.L2R_Expected[i].Name)
		}
	}

	r2l := diff(c.Remote, c.Local)
	if len(r2l) != len(c.R2L_Expected) {
		t.Errorf("#%d r2l: diff should return %d files instead of %d", testNo, len(c.R2L_Expected), len(r2l))
	}
	for i := range r2l {
		if r2l[i].Name != c.R2L_Expected[i].Name {
			t.Errorf("#%d r2l: %s should be %s", testNo, r2l[i].Name, c.R2L_Expected[i].Name)
		}
	}
}

func TestDiff(t *testing.T) {
	TestDiffCase{
		Local:        []File{File{Name: "/a"}},
		Remote:       []File{File{Name: "/a"}},
		L2R_Expected: []File{},
		R2L_Expected: []File{},
	}.test(t)

	TestDiffCase{
		Local:        []File{File{Name: "/a"}},
		Remote:       []File{File{Name: "/b"}},
		L2R_Expected: []File{File{Name: "/a"}},
		R2L_Expected: []File{File{Name: "/b"}},
	}.test(t)

	TestDiffCase{
		Local:        []File{File{Name: "/b"}, File{Name: "/a"}},
		Remote:       []File{File{Name: "/d"}, File{Name: "/c"}},
		L2R_Expected: []File{File{Name: "/b"}, File{Name: "/a"}},
		R2L_Expected: []File{File{Name: "/d"}, File{Name: "/c"}},
	}.test(t)

	TestDiffCase{
		Local:        []File{File{Name: "/b"}, File{Name: "/c"}, File{Name: "/a"}},
		Remote:       []File{File{Name: "/d"}, File{Name: "/e"}, File{Name: "/c"}},
		L2R_Expected: []File{File{Name: "/b"}, File{Name: "/a"}},
		R2L_Expected: []File{File{Name: "/d"}, File{Name: "/e"}},
	}.test(t)
}
