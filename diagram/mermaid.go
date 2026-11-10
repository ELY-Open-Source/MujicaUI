package diagram

import (
	"errors"
	"fmt"
	"log/slog"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// MermaidOptions configure a Mermaid.
type MermaidOptions struct {
	// Label names the diagram; empty takes the source's accTitle.
	Label string
}

// inputs are what a scene is made from.
type inputs struct {
	source string
	pal    palette
	size   float32
}

// drawn is a scene, or its error, and their inputs.
type drawn struct {
	in  inputs
	sc  *scene
	err error
}

// Mermaid draws Mermaid source in the window's colors, shrunk to fit.
func Mermaid(c *ui.Context, source string, opts MermaidOptions) ui.Element {
	in := inputs{source: source, pal: tokens(core.Tokens(c)), size: core.FontSize(c, theme.BodySize)}
	root := ui.Column(c).FillWidth()
	d := ui.Local(root, "mermaid.drawn", func() drawn { return drawn{} })
	if d.in != in {
		*d = drawn{in: in}
		d.sc, d.err = render(source, in.pal, in.size, "")
		report(d, len(source))
	}
	return root.Children(func() {
		if d.err != nil {
			failed(c, d.err)
			return
		}
		sc := d.sc
		label := opts.Label
		if label == "" {
			label = sc.title
		}
		if label == "" {
			label = core.Msg(c, "mermaid.label")
		}
		ui.Box(c).FillWidth().MaxWidth(sc.w).AspectRatio(sc.w / sc.h).Role(ui.RoleImage).Label(label).
			Draw(func(p *ui.Painter, r ui.Rect) { p.Clip(r, 0, func() { sc.paint(p, r) }) })
	})
}

// report logs one render, never the source.
func report(d *drawn, bytes int) {
	var pe *pele.Error
	switch {
	case d.err == nil:
		slog.Debug("mujica: mermaid drew", "bytes", bytes, "runs", len(d.sc.runs), "width", d.sc.w, "height", d.sc.h)
	case errors.As(d.err, &pe):
		slog.Debug("mujica: mermaid failed", "bytes", bytes, "code", pe.Code, "line", pe.Line, "column", pe.Column)
	default:
		slog.Debug("mujica: mermaid failed", "bytes", bytes, "code", "walk")
	}
}

// failed shows the error; Pele's gets its line and column.
func failed(c *ui.Context, err error) {
	body := err.Error()
	var pe *pele.Error
	if errors.As(err, &pe) && pe.Line > 0 {
		body = fmt.Sprintf("%d:%d %s", pe.Line, pe.Column, body)
	}
	feedback.Alert(c, nil, body, feedback.AlertOptions{Severity: core.SeverityError, Title: core.Msg(c, "mermaid.failed")})
}

// halo is where an outlined text's stroke is stamped.
var halo = [8][2]float32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-0.7, -0.7}, {0.7, -0.7}, {-0.7, 0.7}, {0.7, 0.7}}

// paint draws the scene scaled to r's width.
func (sc *scene) paint(p *ui.Painter, r ui.Rect) {
	k := r.W / sc.w
	view := ui.Scaling(k, k).Then(ui.Translation(r.X, r.Y))
	for _, run := range sc.runs {
		p.Transform(run.t.Then(view), 0, 0, func() {
			for _, m := range run.marks {
				switch {
				case m.text != nil && m.halo > 0:
					for _, at := range halo {
						p.TextLayout(m.text, m.x+at[0]*m.halo, m.y+at[1]*m.halo, m.fill)
					}
				case m.text != nil:
					p.TextLayout(m.text, m.x, m.y, m.fill)
				default:
					if m.fill.A > 0 {
						p.FillPath(m.path, m.fill)
					}
					if m.stroke.A > 0 {
						p.StrokePath(m.path, m.width, m.stroke)
					}
				}
			}
		})
	}
}
