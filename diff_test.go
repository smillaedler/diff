// Copyright 2012 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"testing"
)

type testcase struct {
	name string
	a, b []int
	res  []Change
}

var tests = []testcase{
	{"shift",
		[]int{1, 2, 3},
		[]int{0, 1, 2, 3},
		[]Change{{0, 0, 0, 1}},
	},
	{"push",
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		[]Change{{3, 3, 0, 1}},
	},
	{"unshift",
		[]int{0, 1, 2, 3},
		[]int{1, 2, 3},
		[]Change{{0, 0, 1, 0}},
	},
	{"pop",
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3},
		[]Change{{3, 3, 1, 0}},
	},
	{"all changed",
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int{10, 11, 12, 13, 14},
		[]Change{
			{0, 0, 10, 5},
		},
	},
	{"all same",
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]Change{},
	},
	{"wrap",
		[]int{1},
		[]int{0, 1, 2, 3},
		[]Change{
			{0, 0, 0, 1},
			{1, 2, 0, 2},
		},
	},
	{"snake",
		[]int{0, 1, 2, 3, 4, 5},
		[]int{1, 2, 3, 4, 5, 6},
		[]Change{
			{0, 0, 1, 0},
			{6, 5, 0, 1},
		},
	},
	// note: input is ambiguous
	// first two traces differ from fig.1
	// it still is a lcs and ses path
	{"paper fig. 1",
		[]int{1, 2, 3, 1, 2, 2, 1},
		[]int{3, 2, 1, 2, 1, 3},
		[]Change{
			{0, 0, 1, 1},
			{2, 2, 1, 0},
			{5, 4, 1, 0},
			{7, 5, 0, 1},
		},
	},
}

func TestDiffAB(t *testing.T) {
	for _, test := range tests {
		res := Ints(test.a, test.b)
		if len(res) != len(test.res) {
			t.Error(test.name, "expected length", len(test.res), "for", res)
			continue
		}
		for i, c := range test.res {
			if c != res[i] {
				t.Error(test.name, "expected ", c, "got", res[i])
			}
		}
	}
}

func TestDiffBA(t *testing.T) {
	// interesting: fig.1 Diff(b, a) results in the same path as `diff -d a b`
	tests[len(tests)-1].res = []Change{
		{0, 0, 2, 0},
		{3, 1, 1, 0},
		{5, 2, 0, 1},
		{7, 5, 0, 1},
	}
	for _, test := range tests {
		res := Ints(test.b, test.a)
		if len(res) != len(test.res) {
			t.Error(test.name, "expected length", len(test.res), "for", res)
			continue
		}
		for i, c := range test.res {
			// flip change data also
			rc := Change{c.B, c.A, c.Ins, c.Del}
			if rc != res[i] {
				t.Error(test.name, "expected ", rc, "got", res[i])
			}
		}
	}
}

func TestDiffRunes(t *testing.T) {
	d := &runes{
		[]rune("brown fox jumps over the lazy dog"),
		[]rune("brwn faax junps ovver the lay dago"),
	}
	n, m := len(d.a), len(d.b)
	res := Diff(n, m, d)
	echange := []Change{
		{2, 2, 1, 0},
		{7, 6, 1, 2},
		{12, 12, 1, 1},
		{18, 18, 0, 1},
		{27, 28, 1, 0},
		{31, 31, 0, 2},
		{32, 34, 1, 0},
	}
	for i, c := range res {
		t.Log(c)
		if c != echange[i] {
			t.Error("expected", echange[i], "got", c)
		}
	}
}

func BenchmarkDiff(b *testing.B) {
	t := tests[len(tests)-1]
	d := &ints{t.a, t.b}
	n, m := len(d.a), len(d.b)
	for i := 0; i < b.N; i++ {
		Diff(n, m, d)
	}
}

func BenchmarkDiffRunes(b *testing.B) {
	d := &runes{[]rune("1231221"), []rune("321213")}
	n, m := len(d.a), len(d.b)
	for i := 0; i < b.N; i++ {
		Diff(n, m, d)
	}
}
