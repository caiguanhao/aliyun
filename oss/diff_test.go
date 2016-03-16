package main

import "testing"

var testNo int

type TestDiffCase struct {
	Local        []OSSFile
	Remote       []OSSFile
	L2R_Expected []OSSFile
	R2L_Expected []OSSFile
}

func (c TestDiffCase) test(t *testing.T) {
	testNo++

	l2r, r2l := getDiff(c.Local, c.Remote, 0, 0, false)

	if len(l2r) != len(c.L2R_Expected) {
		t.Errorf("#%d l2r: diff should return %d files instead of %d", testNo, len(c.L2R_Expected), len(l2r))
	} else {
		for i := range l2r {
			if l2r[i].Name != c.L2R_Expected[i].Name {
				t.Errorf("#%d l2r: %s should be %s", testNo, l2r[i].Name, c.L2R_Expected[i].Name)
			}
		}
	}

	if len(r2l) != len(c.R2L_Expected) {
		t.Errorf("#%d r2l: diff should return %d files instead of %d", testNo, len(c.R2L_Expected), len(r2l))
	} else {
		for i := range r2l {
			if r2l[i].Name != c.R2L_Expected[i].Name {
				t.Errorf("#%d r2l: %s should be %s", testNo, r2l[i].Name, c.R2L_Expected[i].Name)
			}
		}
	}
}

func TestDiff(t *testing.T) {
	TestDiffCase{
		Local:        []OSSFile{OSSFile{Name: "/a"}},
		Remote:       []OSSFile{OSSFile{Name: "/a"}},
		L2R_Expected: []OSSFile{},
		R2L_Expected: []OSSFile{},
	}.test(t)

	TestDiffCase{
		Local:        []OSSFile{OSSFile{Name: "/a"}},
		Remote:       []OSSFile{OSSFile{Name: "/b"}},
		L2R_Expected: []OSSFile{OSSFile{Name: "/a"}},
		R2L_Expected: []OSSFile{OSSFile{Name: "/b"}},
	}.test(t)

	TestDiffCase{
		Local:        []OSSFile{OSSFile{Name: "/b"}, OSSFile{Name: "/a"}},
		Remote:       []OSSFile{OSSFile{Name: "/d"}, OSSFile{Name: "/c"}},
		L2R_Expected: []OSSFile{OSSFile{Name: "/b"}, OSSFile{Name: "/a"}},
		R2L_Expected: []OSSFile{OSSFile{Name: "/d"}, OSSFile{Name: "/c"}},
	}.test(t)

	TestDiffCase{
		Local:        []OSSFile{OSSFile{Name: "/b"}, OSSFile{Name: "/c"}, OSSFile{Name: "/a"}},
		Remote:       []OSSFile{OSSFile{Name: "/d"}, OSSFile{Name: "/e"}, OSSFile{Name: "/c"}},
		L2R_Expected: []OSSFile{OSSFile{Name: "/b"}, OSSFile{Name: "/a"}},
		R2L_Expected: []OSSFile{OSSFile{Name: "/d"}, OSSFile{Name: "/e"}},
	}.test(t)
}
