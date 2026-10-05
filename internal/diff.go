package internal

import "strings"

// DiffOp says a line is kept, removed or added.
type DiffOp uint8

const (
	DiffEqual DiffOp = iota
	DiffDelete
	DiffInsert
)

// DiffLine is a numbered diff line; 0 means absent.
type DiffLine struct {
	Op       DiffOp
	Text     string
	Old, New int
}

// SplitLines splits text into lines; "" has none.
func SplitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

// LineDiff returns a shortest edit script from a to b.
func LineDiff(a, b []string) []DiffLine {
	// Common prefix and suffix need no search.
	pre := 0
	for pre < len(a) && pre < len(b) && a[pre] == b[pre] {
		pre++
	}
	suf := 0
	for suf < len(a)-pre && suf < len(b)-pre && a[len(a)-1-suf] == b[len(b)-1-suf] {
		suf++
	}
	out := make([]DiffLine, 0, len(a)+len(b)-pre-suf)
	for i := range pre {
		out = append(out, DiffLine{DiffEqual, a[i], i + 1, i + 1})
	}
	out = myers(out, a[pre:len(a)-suf], b[pre:len(b)-suf], pre, pre)
	for i := range suf {
		ai, bi := len(a)-suf+i, len(b)-suf+i
		out = append(out, DiffLine{DiffEqual, a[ai], ai + 1, bi + 1})
	}
	return out
}

// myers appends the edit script, numbering from ao and bo.
func myers(out []DiffLine, a, b []string, ao, bo int) []DiffLine {
	n, m := len(a), len(b)
	total := n + m
	if total == 0 {
		return out
	}
	off := total + 1
	v := make([]int, 2*total+3)
	var trace [][]int
	end := -1
	for d := 0; d <= total && end < 0; d++ {
		trace = append(trace, append([]int(nil), v[off-d:off+d+1]...)) // diagonals -d..d
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || k != d && v[off+k-1] < v[off+k+1] {
				x = v[off+k+1]
			} else {
				x = v[off+k-1] + 1
			}
			y := x - k
			for x < n && y < m && a[x] == b[y] {
				x, y = x+1, y+1
			}
			v[off+k] = x
			if x >= n && y >= m {
				end = d
				break
			}
		}
	}
	// Walk the trace back, collecting moves in reverse.
	var rev []DiffLine
	x, y := n, m
	for d := end; d >= 0; d-- {
		vd := trace[d]
		at := func(k int) int { return vd[k+d] }
		k := x - y
		var pk int
		if k == -d || k != d && at(k-1) < at(k+1) {
			pk = k + 1
		} else {
			pk = k - 1
		}
		px := 0
		if d > 0 {
			px = at(pk)
		}
		py := px - pk
		for x > px && y > py {
			x, y = x-1, y-1
			rev = append(rev, DiffLine{DiffEqual, a[x], ao + x + 1, bo + y + 1})
		}
		if d == 0 {
			break
		}
		if x == px {
			y--
			rev = append(rev, DiffLine{DiffInsert, b[y], 0, bo + y + 1})
		} else {
			x--
			rev = append(rev, DiffLine{DiffDelete, a[x], ao + x + 1, 0})
		}
	}
	for i := len(rev) - 1; i >= 0; i-- {
		out = append(out, rev[i])
	}
	return out
}

// DiffStats counts the added and removed lines of a diff.
func DiffStats(d []DiffLine) (added, removed int) {
	for _, l := range d {
		switch l.Op {
		case DiffInsert:
			added++
		case DiffDelete:
			removed++
		}
	}
	return added, removed
}
