package chart

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
)

// paint draws annotations, previews, and selected handles.
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
			drawingStroke(p, ax, ay, bx, by, width, col)
		case ToolHorizontal:
			p.Line(g.plot.X, ay, g.plot.X+g.plot.W, ay, width, col)
			ChartText(p, core.PriceUngrouped(it.A.Price, m.dec), g.plot.X+4, ay-2, size, col, 0, 1)
		case ToolRectangle:
			r := ui.Rect{X: min(ax, bx), Y: min(ay, by), W: max(bx, ax) - min(ax, bx), H: max(by, ay) - min(ay, by)}
			p.Fill(r, col.Alpha(0.1), 0)
			p.Stroke(r, col, 0, width)
		case ToolFibonacci:
			for _, lv := range fibLevels {
				price := it.B.Price + (it.A.Price-it.B.Price)*lv
				y := internal.Snap(p, g.ys.Map(m.val(price)))
				p.Line(min(ax, bx), y, max(ax, bx), y, 1, col)
				ChartText(p, strconv.FormatFloat(lv, 'f', -1, 64)+" · "+core.PriceUngrouped(price, m.dec), max(ax, bx)+4, y, size, col, 0, 0.5)
			}
			drawingStroke(p, ax, ay, bx, by, 1, col.Alpha(0.5))
		case ToolText:
			w, h := ChartText(p, it.Text, ax, ay, size, k.Text, 0.5, 0.5)
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

func drawingStroke(p *ui.Painter, ax, ay, bx, by, width float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(internal.Snap(p, ax), internal.Snap(p, ay)).LineTo(internal.Snap(p, bx), internal.Snap(p, by))
	p.StrokePath(&path, width, col)
}
