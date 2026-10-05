package finance

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func level2Sample() Level2Data {
	return Level2Data{
		Bids: []L2Level{{100.00, 500, 12}, {99.99, 300, 4}},
		Asks: []L2Level{{100.01, 200, 3}, {100.02, 900, 21}, {100.03, 50, 1}},
	}
}

func TestLevel2QuotesLevels(t *testing.T) {
	tt := tradingTester(core.GreenUp, 600, 300, func(c *ui.Context) {
		ui.Box(c).Width(560).Children(func() { Level2Quotes(c, level2Sample(), Level2QuotesOptions{Levels: 3}) })
	})
	for _, name := range []string{"B1 100.00 500 12 orders", "S2 100.02 900 21 orders"} {
		if _, ok := tt.Find(name); !ok {
			t.Fatalf("no row %q in %q", name, tt.Texts())
		}
	}
	if !tt.HasText("B3") || !tt.HasText("—") {
		t.Fatal("missing third bid not shown as empty")
	}
}

// A changed level marks its size move and flashes.
func TestLevel2QuotesHighlightsChangedLevel(t *testing.T) {
	data := level2Sample()
	var marks *map[string]*level2Mark
	tt := tradingTester(core.GreenUp, 600, 300, func(c *ui.Context) {
		e := Level2Quotes(c, data, Level2QuotesOptions{Levels: 2})
		marks = ui.Local(e, "marks", func() map[string]*level2Mark { return nil })
	})
	before := bytes.Clone(tt.Image().Pix)
	data.Bids = []L2Level{{100.00, 650, 13}, {99.99, 300, 4}}
	tt.Frame()
	m := (*marks)["1 0"]
	if m.dir != 1 || m.at.IsZero() {
		t.Fatalf("bid 1 mark %+v", m)
	}
	if (*marks)["1 1"].dir != 0 || !(*marks)["1 1"].at.IsZero() {
		t.Fatal("unchanged level marked")
	}
	if bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("nothing drawn for the change")
	}
	data.Bids = []L2Level{{100.01, 650, 13}, {99.99, 300, 4}}
	tt.Frame()
	if (*marks)["1 0"].dir != 0 {
		t.Fatal("new price read as a size move")
	}
}

func TestLevel2QuotesSnapshot(t *testing.T) {
	snapshot(t, "level2_quotes", 600, 230, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { Level2Quotes(c, level2Sample(), Level2QuotesOptions{}) })
	})
}
