package mujica

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// paint draws every drawing, the one being placed, and the selected one's handles.
func (d *DrawingTools) paint(p *ui.Painter, g mktGeom, m *mktCore, size float32) {
	k := m.k
	items := d.Items
	if d.pending != nil {
		items = append(append([]Drawing(nil), items...), *d.pending)
	}
	for i, it := range items {
		col := k.Info
		width := float32(1.5)
		selected := i == d.sel-1
		if selected {
			col, width = k.AccentText, 2.5
		}
		ax, ay := d.xy(m, g, it.A)
		bx, by := d.xy(m, g, it.B)
		ax, ay, bx, by = internal.Snap(p, ax), internal.Snap(p, ay), internal.Snap(p, bx), internal.Snap(p, by)
		switch it.Tool {
		case ToolTrendLine:
			p.Line(ax, ay, bx, by, width, col)
		case ToolHorizontal:
			p.Line(g.plot.X, ay, g.plot.X+g.plot.W, ay, width, col)
			chartText(p, priceUngrouped(it.A.Price, m.dec), g.plot.X+4, ay-2, size, col, 0, 1)
		case ToolRectangle:
			r := ui.Rect{X: min(ax, bx), Y: min(ay, by), W: max(bx, ax) - min(ax, bx), H: max(by, ay) - min(ay, by)}
			p.Fill(r, col.Alpha(0.1), 0)
			p.Stroke(r, col, 0, width)
		case ToolFibonacci:
			for _, lv := range fibLevels {
				price := it.B.Price + (it.A.Price-it.B.Price)*lv
				y := internal.Snap(p, g.ys.Map(m.val(price)))
				p.Line(min(ax, bx), y, max(ax, bx), y, 1, col)
				chartText(p, strconv.FormatFloat(lv, 'f', -1, 64)+" · "+priceUngrouped(price, m.dec), max(ax, bx)+4, y, size, col, 0, 0.5)
			}
			p.Line(ax, ay, bx, by, 1, col.Alpha(0.5))
		case ToolText:
			w, h := chartText(p, it.Text, ax, ay, size, k.Text, 0.5, 0.5)
			p.Stroke(ui.Rect{X: ax - w/2 - 4, Y: ay - h/2 - 2, W: w + 8, H: h + 4}, col, 2, width/1.5)
		}
		if selected {
			for _, pt := range [][2]float32{{ax, ay}, {bx, by}} {
				p.FillPath(internal.Diamond(pt[0], pt[1], 9), k.Surface)
				p.FillPath(internal.Diamond(pt[0], pt[1], 7), col)
			}
		}
	}
}
