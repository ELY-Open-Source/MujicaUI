package internal

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

// diffSides rebuilds both inputs from a diff and checks the line numbers.
func diffSides(t *testing.T, d []DiffLine) (a, b []string) {
	t.Helper()
	for _, l := range d {
		if l.Op != DiffInsert {
			a = append(a, l.Text)
			if l.Old != len(a) {
				t.Fatalf("old number %d, want %d", l.Old, len(a))
			}
		}
		if l.Op != DiffDelete {
			b = append(b, l.Text)
			if l.New != len(b) {
				t.Fatalf("new number %d, want %d", l.New, len(b))
			}
		}
	}
	return a, b
}

// lcs is the reference longest common subsequence length.
func lcs(a, b []string) int {
	prev := make([]int, len(b)+1)
	for i := range a {
		cur := make([]int, len(b)+1)
		for j := range b {
			if a[i] == b[j] {
				cur[j+1] = prev[j] + 1
			} else {
				cur[j+1] = max(prev[j+1], cur[j])
			}
		}
		prev = cur
	}
	return prev[len(b)]
}

func TestLineDiffKnown(t *testing.T) {
	d := LineDiff(SplitLines("a\nb\nc\nd\n"), SplitLines("a\nc\nd\ne\n"))
	var got []string
	for _, l := range d {
		got = append(got, fmt.Sprintf("%d%s", l.Op, l.Text))
	}
	want := []string{"0a", "1b", "0c", "0d", "2e"}
	if !slices.Equal(got, want) {
		t.Fatalf("diff %v, want %v", got, want)
	}
	if add, del := DiffStats(d); add != 1 || del != 1 {
		t.Fatalf("stats +%d -%d", add, del)
	}
}

func TestLineDiffEdges(t *testing.T) {
	for _, tc := range []struct{ a, b string }{{"", ""}, {"", "x\ny"}, {"x\ny", ""}, {"same\n", "same\n"}, {"a", "b"}} {
		d := LineDiff(SplitLines(tc.a), SplitLines(tc.b))
		a, b := diffSides(t, d)
		if !slices.Equal(a, SplitLines(tc.a)) || !slices.Equal(b, SplitLines(tc.b)) {
			t.Errorf("%q -> %q: rebuilt %q, %q", tc.a, tc.b, a, b)
		}
	}
}

func TestLineDiffMinimalRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for range 300 {
		a := make([]string, r.Intn(30))
		b := make([]string, r.Intn(30))
		for i := range a {
			a[i] = string(rune('a' + r.Intn(4)))
		}
		for i := range b {
			b[i] = string(rune('a' + r.Intn(4)))
		}
		d := LineDiff(a, b)
		ra, rb := diffSides(t, d)
		if !slices.Equal(ra, a) || !slices.Equal(rb, b) {
			t.Fatalf("rebuilt %q %q from %q %q", ra, rb, a, b)
		}
		add, del := DiffStats(d)
		if keep := len(a) - del; keep != lcs(a, b) || len(b)-add != keep {
			t.Fatalf("%q -> %q: kept %d lines, LCS is %d", a, b, keep, lcs(a, b))
		}
	}
}

func TestLineDiffLargeFileSmallEdit(t *testing.T) {
	a := make([]string, 50000)
	for i := range a {
		a[i] = fmt.Sprint("line ", i)
	}
	b := slices.Clone(a)
	b[25000] = "changed"
	d := LineDiff(a, b)
	if add, del := DiffStats(d); add != 1 || del != 1 || len(d) != 50001 {
		t.Fatalf("+%d -%d over %d lines", add, del, len(d))
	}
	if !strings.HasPrefix(d[25000].Text, "line 25000") || d[25000].Op != DiffDelete {
		t.Fatalf("line 25000: %+v", d[25000])
	}
}
