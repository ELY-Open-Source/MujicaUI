package input

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RatingOptions configure a Rating. Max is 5 and Step 1 when zero.
type RatingOptions struct {
	Max int
	// Step is the smallest change, such as 0.5 for half stars; it must divide 1.
	Step     float64
	ReadOnly bool
	Disabled bool
	// Label names the rating; empty uses the built-in "Rating".
	Label string
	// ShowValue prints "3.5 / 5" beside the stars.
	ShowValue bool
}

// Rating shows *value as stars and lets the pointer and arrows set it in steps.
func Rating(c *ui.Context, value *float64, opts RatingOptions) core.ChoiceResult {
	n, step := opts.Max, opts.Step
	if n == 0 {
		n = 5
	}
	if step == 0 {
		step = 1
	}
	if n < 0 || step < 0 || step > 1 || math.Abs(1/step-math.Round(1/step)) > 1e-9 {
		panic(fmt.Sprintf("mujica: Rating Max %d, Step %v: Max must be positive and Step divide 1", opts.Max, opts.Step))
	}
	k := core.Tokens(c)
	size := core.Density(c).ControlHeight() * 0.6
	// Each star sits centered in a hit area of at least the minimum size.
	hit := max(size, theme.MinHitSize)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "rating.label")
	}
	parts := core.ChoiceResult{}
	set := func(v float64) {
		v = math.Max(0, math.Min(float64(n), math.Round(v/step)*step))
		if v != *value {
			slog.Debug("mujica: rating changed", "from", *value, "to", v)
			*value = v
			parts = core.NewChoiceResult(parts.Element, true)
		}
	}
	e := ui.Row(c).AlignItems(ui.Center).Gap(core.Density(c).Unit() * 2).Shrink(0)
	parts.Element = e
	e.Label(label)
	e.Disabled(opts.Disabled)
	disabled := e.IsDisabled()
	interactive := !opts.ReadOnly && !disabled
	if !disabled {
		e.Focusable()
		internal.Ring(e, k.Focus, theme.ControlRadius)
	}
	if interactive {
		e.Role(ui.RoleSlider).Cursor(ui.CursorPointer)
		switch {
		case e.Shortcut(0, ui.KeyRight), e.Shortcut(0, ui.KeyUp):
			set(*value + step)
		case e.Shortcut(0, ui.KeyLeft), e.Shortcut(0, ui.KeyDown):
			set(*value - step)
		case e.Shortcut(0, ui.KeyHome):
			set(0)
		case e.Shortcut(0, ui.KeyEnd):
			set(float64(n))
		}
	}
	e.Children(func() {
		stars := ui.Box(c).Size(hit*float32(n), hit).Shrink(0)
		shown, fill := *value, k.AccentText
		if interactive {
			at := func(x float32) float64 { return math.Ceil(float64(x/hit)/step-1e-9) * step }
			if _, _, ok := stars.Dragged(); ok {
				x, _, _ := stars.PointerPosition()
				set(at(x))
				shown = *value
			} else if x, _, over := stars.PointerPosition(); over {
				shown, fill = math.Max(step, math.Min(float64(n), at(x))), k.AccentText.Alpha(0.6)
			}
		}
		stars.Draw(func(p *ui.Painter, r ui.Rect) {
			for i := range n {
				cell := ui.Rect{X: r.X + float32(i)*hit + (hit-size)/2, Y: r.Y + (hit-size)/2, W: size, H: size}
				path := ratingStar(cell)
				p.StrokePath(path, 1.25, k.ControlBorder)
				part := float32(math.Max(0, math.Min(1, shown-float64(i))))
				if part > 0 {
					p.Clip(ui.Rect{X: cell.X, Y: cell.Y, W: cell.W * part, H: cell.H}, 0, func() {
						p.FillPath(path, fill)
						p.StrokePath(path, 1.25, fill)
					})
				}
			}
		})
		if opts.ShowValue {
			ui.Text(c, fmt.Sprintf("%s / %d", strconv.FormatFloat(*value, 'f', -1, 64), n)).
				TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
		}
	})
	e.Description(fmt.Sprintf(core.Msg(c, "rating.value"), strconv.FormatFloat(*value, 'f', -1, 64), n))
	return parts
}

// ratingStar returns a five-pointed star inset in r.
func ratingStar(r ui.Rect) *ui.Path {
	var p ui.Path
	cx, cy := r.X+r.W/2, r.Y+r.H*0.54
	outer := r.W * 0.42
	inner := outer * 0.45
	for i := range 10 {
		rad := outer
		if i%2 == 1 {
			rad = inner
		}
		a := -math.Pi/2 + float64(i)*math.Pi/5
		x, y := cx+rad*float32(math.Cos(a)), cy+rad*float32(math.Sin(a))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	p.Close()
	return &p
}
