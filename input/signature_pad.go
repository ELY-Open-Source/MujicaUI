package input

import (
	"cmp"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SignaturePoint is a point of a stroke, in DIPs from the pad's top left.
type SignaturePoint struct{ X, Y float32 }

// SignaturePadOptions configure a SignaturePad.
type SignaturePadOptions struct {
	// Height is the drawing area's height; zero is 160. The pad fills its parent's width.
	Height float32
	// Label names the pad; empty uses "Signature".
	Label    string
	Disabled bool
}

// SignaturePath writes strokes as SVG path data: "M x y L x y …" per stroke.
func SignaturePath(strokes [][]SignaturePoint) string {
	var b strings.Builder
	num := func(v float32) string { return strconv.FormatFloat(float64(v), 'f', -1, 32) }
	for _, st := range strokes {
		for i, p := range st {
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			if i == 0 {
				b.WriteString("M ")
			} else {
				b.WriteString("L ")
			}
			b.WriteString(num(p.X) + " " + num(p.Y))
		}
	}
	return b.String()
}

// signaturePadDrawing is the stroke in progress.
type signaturePadDrawing struct {
	active, ended bool
}

// SignaturePad records pointer strokes into *strokes; Changed reports a finished stroke or a clear.
func SignaturePad(c *ui.Context, strokes *[][]SignaturePoint, opts SignaturePadOptions) core.ChoiceResult {
	if strokes == nil {
		panic("mujica: SignaturePad needs strokes")
	}
	if !(opts.Height >= 0) || math.IsInf(float64(opts.Height), 1) {
		panic("mujica: SignaturePad Height must be finite and >= 0")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	name := cmp.Or(opts.Label, core.Msg(c, "signature_pad.label"))
	col := ui.Column(c).FillWidth().Gap(u * 2).Disabled(opts.Disabled)
	disabled := col.IsDisabled()
	changed := false
	col.Children(func() {
		pad := ui.Box(c).FillWidth().Height(cmp.Or(opts.Height, 160)).Focusable().Radius(theme.CardRadius).Clip().
			Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Cursor(ui.CursorCrosshair).
			Role(ui.RoleImage).Label(name).Description(core.Msg(c, "signature_pad.hint"))
		internal.Ring(pad, k.Focus, theme.CardRadius)
		d := ui.Local(pad, "signature_pad", func() *signaturePadDrawing { return &signaturePadDrawing{} })
		s := *d
		pad.HandleInput(func(ev ui.InputEvent) bool {
			if pad.IsDisabled() {
				return false
			}
			pt := SignaturePoint{ev.X, ev.Y}
			switch {
			case ev.Kind == ui.InputPointerDown && ev.Button == 0:
				*strokes = append(*strokes, []SignaturePoint{pt})
				s.active = true
				pad.Focus()
			case ev.Kind == ui.InputPointerMove && s.active && len(*strokes) > 0:
				last := &(*strokes)[len(*strokes)-1]
				if q := (*last)[len(*last)-1]; math.Hypot(float64(pt.X-q.X), float64(pt.Y-q.Y)) >= 1 {
					*last = append(*last, pt)
				}
			case ev.Kind == ui.InputPointerUp && s.active:
				s.active, s.ended = false, true
			default:
				return false
			}
			return true
		})
		if s.ended {
			s.ended, changed = false, true
			slog.Debug("mujica: signature_pad stroke", "strokes", len(*strokes))
		}
		clear := func() {
			if len(*strokes) > 0 {
				*strokes, changed = nil, true
				slog.Debug("mujica: signature_pad clear")
			}
		}
		if !disabled && (pad.Shortcut(0, ui.KeyDelete) || pad.Shortcut(0, ui.KeyBackspace)) {
			clear()
		}
		ink := signaturePadClone(*strokes)
		empty := len(ink) == 0
		size := core.FontSize(c, theme.CaptionSize)
		pad.Draw(func(p *ui.Painter, r ui.Rect) {
			line := r.Y + r.H - 28
			p.Fill(ui.Rect{X: r.X + 16, Y: line, W: r.W - 32, H: 1}, k.Border, 0)
			p.Text(r.X+16, line+4, "×", size, k.TextMuted)
			if empty {
				w, _ := p.MeasureText(0, ui.Span{Text: core.Msg(c, "signature_pad.empty"), Size: size})
				p.Text(r.X+(r.W-w)/2, r.Y+r.H/2-size, core.Msg(c, "signature_pad.empty"), size, k.TextMuted)
			}
			for _, st := range ink {
				var path ui.Path
				for i, q := range st {
					x, y := internal.Snap(p, r.X+q.X), internal.Snap(p, r.Y+q.Y)
					if i == 0 {
						path.MoveTo(x, y)
					} else {
						path.LineTo(x, y)
					}
				}
				if len(st) == 1 {
					path.Circle(internal.Snap(p, r.X+st[0].X), internal.Snap(p, r.Y+st[0].Y), 1)
				}
				p.StrokePath(&path, 2, k.Text)
			}
		})
		ui.Row(c).Children(func() {
			if Button(c, core.Msg(c, "signature_pad.clear"), ButtonOptions{Variant: Ghost, Icon: icons.Must("trash-2"), Disabled: empty}).Clicked() && !disabled {
				clear()
			}
		})
	})
	return core.NewChoiceResult(col, changed)
}

// signaturePadClone copies strokes so painting never sees later edits.
func signaturePadClone(s [][]SignaturePoint) [][]SignaturePoint {
	out := make([][]SignaturePoint, len(s))
	for i, st := range s {
		out[i] = append([]SignaturePoint(nil), st...)
	}
	return out
}
