package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestDrawingToolsPlaceMoveDeleteUndo(t *testing.T) {
	cs := sampleCandles(60)
	v := MarketView{Count: 60}
	d := DrawingTools{Tool: ToolHorizontal}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Drawing"}, Drawings: &d})
		})
	}, 640, 340)
	tt.ClickAt(300, 100)
	if len(d.Items) != 1 || d.Items[0].Tool != ToolHorizontal || d.Selected() != 0 {
		t.Fatalf("drawings %+v", d)
	}
	before := d.Items[0].A.Price
	tt.Key(ui.Alt, ui.KeyUp)
	if d.Items[0].A.Price <= before {
		t.Fatal("keyboard move failed")
	}
	tt.Key(0, ui.KeyDelete)
	if len(d.Items) != 0 {
		t.Fatal("delete failed")
	}
	tt.Key(ui.Cmd, ui.KeyZ)
	if len(d.Items) != 1 {
		t.Fatal("undo failed")
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyDelete)
	tt.ClickAt(200, 150)
	if len(d.Items) != 1 {
		t.Fatal("disabled drawing changed")
	}
}
func TestDrawingToolsSnapshot(t *testing.T) {
	cs := sampleCandles(60)
	v := MarketView{Count: 60}
	d := DrawingTools{Items: []Drawing{{Tool: ToolTrendLine, A: DrawingPoint{cs[10].Time, 98}, B: DrawingPoint{cs[45].Time, 105}}, {Tool: ToolHorizontal, A: DrawingPoint{cs[20].Time, 100}}}, sel: 1}
	snapshot(t, "drawing_tools", 640, 320, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Drawings", Height: 260}, Drawings: &d})
		})
	})
}

func TestDrawingToolsTrendLinePaint(t *testing.T) {
	cs := sampleCandles(2)
	img := ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		m := &mktCore{c: c, k: core.Tokens(c), candles: cs, from: 0, to: 2}
		g := mktGeom{plot: ui.Rect{W: 200, H: 200}, band: internal.Band{N: 2, R0: 0, R1: 200}, ys: internal.Linear{D0: 0, D1: 10, R0: 180, R1: 20}}
		d := DrawingTools{Items: []Drawing{{Tool: ToolTrendLine, A: DrawingPoint{cs[0].Time, 0}, B: DrawingPoint{cs[1].Time, 10}}}}
		ui.Box(c).Size(200, 200).Draw(func(p *ui.Painter, _ ui.Rect) { d.paint(p, g, m, 12) })
	}, 200, 200, 1)
	bg := theme.Light().Background
	pixel := img.RGBAAt(100, 100)
	if int(bg.R)-int(pixel.R) < 40 || int(bg.G)-int(pixel.G) < 20 {
		t.Fatalf("diagonal midpoint lacks the stroke: %v", pixel)
	}
}

func TestDrawingToolsAllTools(t *testing.T) {
	for _, tool := range []DrawingTool{ToolTrendLine, ToolHorizontal, ToolFibonacci, ToolRectangle, ToolText} {
		t.Run(string(rune('0'+tool)), func(t *testing.T) {
			cs := sampleCandles(30)
			v := MarketView{Count: 30}
			d := DrawingTools{Tool: tool, Text: "Note"}
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Drawing"}, Drawings: &d})
			}, 640, 340)
			tt.ClickAt(200, 120)
			if tool.twoPoint() {
				if d.pending == nil || len(d.Items) != 0 {
					t.Fatal("missing first anchor")
				}
				tt.ClickAt(400, 160)
			}
			if len(d.Items) != 1 || d.Items[0].Tool != tool || d.Selected() != 0 {
				t.Fatalf("created %+v", d.Items)
			}
			if tool == ToolText && d.Items[0].Text != "Note" {
				t.Fatal("annotation text lost")
			}
			if !d.Clear() || len(d.Items) != 0 || !d.Undo() || len(d.Items) != 1 {
				t.Fatal("clear or undo failed")
			}
		})
	}
}
