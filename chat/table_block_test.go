package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var (
	tableBlockHeader = []string{"Guest", "Title", "Seat"}
	tableBlockRows   = [][]string{{"Marie", "Queen", "1"}, {"Louis", "King", "2"}, {"Axel", "Count\tof Fersen", "14"}}
)

func TestTableBlockCopiesTSVAndScrolls(t *testing.T) {
	copies := 0
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(160).Children(func() {
			r := TableBlock(c, tableBlockHeader, tableBlockRows, TableBlockOptions{Align: []ui.Align{ui.Start, ui.Start, ui.End}})
			w = r.Element.Bounds().W
			if r.Copied() {
				copies++
			}
		})
	}, 400, 300)
	tt.Frame()
	if w > 160 {
		t.Fatalf("table widened its message to %v", w)
	}
	if err := tt.Click("Copy as TSV"); err != nil {
		t.Fatal(err)
	}
	want := "Guest\tTitle\tSeat\nMarie\tQueen\t1\nLouis\tKing\t2\nAxel\tCount of Fersen\t14\n"
	if tt.Clipboard() != want || copies != 1 {
		t.Fatalf("copies %d clipboard %q", copies, tt.Clipboard())
	}
}

func TestTableBlockBadRows(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a short row")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TableBlock(c, []string{"a", "b"}, [][]string{{"1"}}, TableBlockOptions{})
	}, 100, 100)
}

func TestTableBlockSnapshot(t *testing.T) {
	snapshot(t, "table_block", 400, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			TableBlock(c, tableBlockHeader, tableBlockRows, TableBlockOptions{Align: []ui.Align{ui.Start, ui.Start, ui.End}})
		})
	})
}
