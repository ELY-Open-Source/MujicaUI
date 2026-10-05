package input

import (
	"cmp"
	"log/slog"
	"math"
	"sort"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// MasonryOptions configure a Masonry.
type MasonryOptions struct {
	// Key is required: a stable, unique key per item.
	Key func(i int) string
	// ColumnWidth is the narrowest column; zero is 220. Columns fill the width.
	ColumnWidth float32
	// Gap is between items; zero is 3 spacing units.
	Gap float32
	// Estimate is the height assumed until an item is measured; zero is 160.
	Estimate float32
}

// masonryState is the scroll, the measured heights and the cached placement.
type masonryState struct {
	scroll  ui.ScrollState
	heights map[string]float32
	version int
	laidOut masonryLayout
	// built are the items built last frame; their keys are checked for changed data.
	built []int
	// sized: a frame was asked for once the box had a width.
	sized bool
}

// masonryLayout places n items in cols columns; byCol lists each column's items top to bottom.
type masonryLayout struct {
	n, cols, version int
	colW, gap        float32
	keys             []string
	x, y, h          []float32
	byCol            [][]int
	total            float32
}

// Masonry lays n items in as many columns as fit, each into the shortest column, building only those in view.
func Masonry(c *ui.Context, n int, opts MasonryOptions, item func(i int)) *ui.Element {
	if opts.Key == nil || item == nil {
		panic("mujica: Masonry needs Key and item")
	}
	if n < 0 || opts.ColumnWidth < 0 || opts.Gap < 0 || opts.Estimate < 0 {
		panic("mujica: Masonry count and sizes must not be negative")
	}
	u := core.Density(c).Unit()
	colMin := cmp.Or(opts.ColumnWidth, 220)
	gap := cmp.Or(opts.Gap, 3*u)
	est := cmp.Or(opts.Estimate, 160)
	sc := ui.Scroll(c)
	st := ui.Local(sc, "masonry", func() *masonryState { return &masonryState{heights: map[string]float32{}} })
	s := *st
	sc.TrackScroll(&s.scroll)
	view := sc.Bounds()
	if view.W == 0 && !s.sized {
		// Unmeasured: lay out again once the box has its width.
		s.sized = true
		c.Invalidate()
	}
	cols := max(1, int((view.W+gap)/(colMin+gap)))
	colW := (view.W - gap*float32(cols-1)) / float32(cols)
	l := &s.laidOut
	stale := l.n != n || l.cols != cols || l.colW != colW || l.gap != gap || l.version != s.version
	for _, i := range s.built {
		stale = stale || i >= n || l.keys[i] != opts.Key(i)
	}
	if stale {
		*l = masonryPlace(n, cols, colW, gap, est, s.heights, opts.Key)
		l.version = s.version
		slog.Debug("mujica: masonry layout", "items", n, "columns", cols)
	}
	top, bottom := s.scroll.Y-view.H/2, s.scroll.Y+view.H*1.5
	var seen []int
	sc.Children(func() {
		ui.Box(c).FillWidth().Height(l.total).Children(func() { masonryBuild(c, s, l, colW, top, bottom, item, &seen) })
	})
	s.built = seen
	return sc
}

// masonryBuild builds the items between top and bottom and records measured heights.
func masonryBuild(c *ui.Context, s *masonryState, l *masonryLayout, colW, top, bottom float32, item func(int), seen *[]int) {
	for _, list := range l.byCol {
		from := sort.Search(len(list), func(j int) bool { i := list[j]; return l.y[i]+l.h[i] >= top })
		for _, i := range list[from:] {
			if l.y[i] > bottom {
				break
			}
			*seen = append(*seen, i)
			key := l.keys[i]
			box := ui.Column(c).Key(key).Absolute().Left(l.x[i]).Top(l.y[i]).Width(colW)
			box.Children(func() { item(i) })
			if h := box.Bounds().H; h > 0 && math.Abs(float64(h-l.h[i])) > 0.5 {
				s.heights[key] = h
				s.version++
				c.Invalidate()
			}
		}
	}
}

// masonryPlace puts each item into the shortest column, using measured heights or est.
func masonryPlace(n, cols int, colW, gap, est float32, heights map[string]float32, key func(int) string) masonryLayout {
	l := masonryLayout{n: n, cols: cols, colW: colW, gap: gap, keys: make([]string, n), x: make([]float32, n), y: make([]float32, n), h: make([]float32, n), byCol: make([][]int, cols)}
	colH := make([]float32, cols)
	seen := make(map[string]bool, n)
	for i := range n {
		k := key(i)
		if seen[k] {
			panic("mujica: Masonry key repeated: " + k)
		}
		seen[k] = true
		col := 0
		for j := range cols {
			if colH[j] < colH[col] {
				col = j
			}
		}
		h, ok := heights[k]
		if !ok {
			h = est
		}
		l.keys[i], l.x[i], l.y[i], l.h[i] = k, float32(col)*(colW+gap), colH[col], h
		l.byCol[col] = append(l.byCol[col], i)
		colH[col] += h + gap
	}
	for _, h := range colH {
		l.total = max(l.total, h-gap)
	}
	return l
}
