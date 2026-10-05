package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var chordNames = []string{"Nave", "Choir", "Crypt", "Cloister"}
var chordMatrix = [][]float64{{0, 5, 2, 3}, {4, 0, 1, 0}, {1, 2, 0, 6}, {2, 0, 3, 0}}

func TestChordPlace(t *testing.T) {
	l := chordPlace(chordMatrix)
	if l.total != 29 || len(l.pairs) != 5 {
		t.Fatalf("total %v pairs %v", l.total, l.pairs)
	}
	span := 2*math.Pi - 4*chordGap
	if got := l.arcs[0][1] - l.arcs[0][0]; math.Abs(got-span*10/29) > 1e-9 {
		t.Fatalf("Nave arc %v", got)
	}
	if l.sub[0][1][1]-l.sub[0][1][0] <= 0 || l.sub[1][3][1] != l.sub[1][3][0] {
		t.Fatal("sub-arcs not sized by the matrix")
	}
}

func TestChordDiagramReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ChordDiagram(c, chordNames, chordMatrix, ChordDiagramOptions{Frame: ChartFrame{Label: "Paths"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Paths")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Nave · out 10 · in 7" {
		t.Fatalf("group %q", got)
	}
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Crypt → Cloister 6 · Cloister → Crypt 3" {
		t.Fatalf("chord %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	b := h.res.Element.Bounds()
	cx, cy, rad := b.X+b.W/2, b.Y+120, float32(240)/2-30
	tt.Move(cx+(rad+6)*0.5, cy-(rad+6)*0.866) // 30°: inside the first arc
	if h.res.Point != 0 {
		t.Fatalf("hover %d", h.res.Point)
	}
}

func TestChordDiagramBad(t *testing.T) {
	for name, m := range map[string][][]float64{"ragged": {{1, 2}, {1}}, "negative": {{0, -1}, {1, 0}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ChordDiagram(c, []string{"a", "b"}, m, ChordDiagramOptions{Frame: ChartFrame{Label: "x"}})
			}, 300, 300)
		}()
	}
}

func TestChordDiagramSnapshot(t *testing.T) {
	snapshot(t, "chord_diagram", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ChordDiagram(c, chordNames, chordMatrix, ChordDiagramOptions{Frame: ChartFrame{Label: "Paths"}})
		})
	})
}

func TestChordPlaceSelfRelation(t *testing.T) {
	l := chordPlace([][]float64{{3}})
	if len(l.pairs) != 1 || l.pairs[0] != [2]int{0, 0} {
		t.Fatalf("self relation %v", l.pairs)
	}
}

func TestChordDiagramDenseDataBoundedAndDisabled(t *testing.T) {
	const n = 450
	names := make([]string, n)
	matrix := make([][]float64, n)
	for i := range matrix {
		names[i] = fmt.Sprint(i)
		matrix[i] = make([]float64, n)
		for j := range matrix[i] {
			matrix[i][j] = 1
		}
	}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ChordDiagram(c, names, matrix, ChordDiagramOptions{Frame: ChartFrame{Label: "Dense chords"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	shown := ui.Local(h.res.Element, "mujica.chord.visible", func() []int { return nil })
	if len(*shown) == 0 || len(*shown) > 360 {
		t.Fatalf("built %d chords", len(*shown))
	}
	tabTo(t, tt, "Dense chords")
	tt.Key(0, ui.KeyEnd)
	if h.res.Point != n+n*(n+1)/2-1 || len(*shown) > 361 {
		t.Fatalf("last relation %d, built %d", h.res.Point, len(*shown))
	}
	h.disabled = true
	tt.Frame()
	if h.res.Point != -1 {
		t.Fatalf("disabled chord %d", h.res.Point)
	}
}
