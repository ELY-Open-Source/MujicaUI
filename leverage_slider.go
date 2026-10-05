package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// LeverageSliderOptions configure a LeverageSlider.
type LeverageSliderOptions struct {
	// Levels are the rising leverages; nil means 1× to 100×.
	Levels []float64
	// Warning and Danger start the risk zones; zero omits one.
	Warning, Danger float64
	// Label names the slider; empty uses the built-in "Leverage".
	Label    string
	Disabled bool
}

type leverageSliderState struct {
	draft float64
	seen  float64
}

// LeverageSlider drafts a leverage; Apply sets *value.
func LeverageSlider(c *ui.Context, value *float64, opts LeverageSliderOptions) ChoiceResult {
	levels := opts.Levels
	if levels == nil {
		levels = []float64{1, 2, 3, 5, 10, 20, 50, 100}
	}
	if len(levels) < 2 || !slices.IsSorted(levels) || slices.Index(levels, *value) < 0 {
		panic(fmt.Sprintf("mujica: LeverageSlider value %v must be one of at least two rising Levels", *value))
	}
	if opts.Warning < 0 || opts.Danger < 0 || opts.Danger != 0 && opts.Danger <= opts.Warning {
		panic("mujica: LeverageSlider Danger must be above Warning")
	}
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "leverage_slider.label")
	}
	var r ChoiceResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2).MinWidth(260).Disabled(opts.Disabled)
	st := ui.Local(r.Element, "lev", func() leverageSliderState { return leverageSliderState{draft: *value, seen: *value} })
	if *value != st.seen {
		st.draft, st.seen = *value, *value
	}
	disabled := r.Element.IsDisabled()
	if disabled {
		st.draft = *value
	}
	zone := func(v float64) (string, ui.Color, string) {
		switch {
		case opts.Danger > 0 && v >= opts.Danger:
			return msg(c, "leverage_slider.danger"), k.Danger, "circle-alert"
		case opts.Warning > 0 && v >= opts.Warning:
			return msg(c, "leverage_slider.warning"), k.Warning, "triangle-alert"
		}
		return "", k.TextMuted, ""
	}
	r.Element.Children(func() {
		note, col, icon := zone(st.draft)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).TextColor(k.TextMuted).Grow(1)
			ui.Text(c, leverageText(st.draft)).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).FontFeatures("tnum")
		})
		leverageTrack(c, st, levels, disabled, zone)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Children(func() {
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(col).Grow(1).Children(func() {
				if icon != "" {
					ui.Icon(c, icons.Must(icon)).Size(14, 14)
					ui.Text(c, note).FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600).Role(ui.RoleStatus)
				}
			})
			if Button(c, msg(c, "leverage_slider.apply"), ButtonOptions{Disabled: disabled || st.draft == *value}).Clicked() && !disabled {
				*value, st.seen, r.changed = st.draft, st.draft, true
				slog.Debug("mujica: leverage_slider applied", "leverage", *value)
			}
		})
	})
	return r
}

func leverageText(v float64) string { return statisticFormat(v, 0) + "×" }

// leverageTrack is the slider with zones, ticks and thumb.
func leverageTrack(c *ui.Context, st *leverageSliderState, levels []float64, disabled bool, zone func(float64) (string, ui.Color, string)) {
	k := Tokens(c)
	n := len(levels)
	at := slices.Index(levels, st.draft)
	set := func(i int) {
		i = max(0, min(n-1, i))
		if levels[i] != st.draft && !disabled {
			st.draft = levels[i]
			slog.Debug("mujica: leverage_slider draft", "leverage", st.draft)
		}
	}
	const pad = 10
	t := ui.Box(c).Height(theme.MinHitSize+16).FillWidth().Padding(0, pad).Focusable().Role(ui.RoleSlider).
		Label(msg(c, "leverage_slider.label")).Description(leverageText(st.draft))
	t.FocusRing(false)
	if !disabled {
		switch {
		case t.Shortcut(0, ui.KeyRight), t.Shortcut(0, ui.KeyUp):
			set(at + 1)
		case t.Shortcut(0, ui.KeyLeft), t.Shortcut(0, ui.KeyDown):
			set(at - 1)
		case t.Shortcut(0, ui.KeyHome):
			set(0)
		case t.Shortcut(0, ui.KeyEnd):
			set(n - 1)
		}
		if _, _, ok := t.Dragged(); ok {
			x, _, _ := t.PointerPosition()
			if w := t.Bounds().W - 2*pad; w > 0 {
				set(int(math.Round(float64((x - pad) / w * float32(n-1)))))
			}
		}
	}
	at = slices.Index(levels, st.draft)
	t.Draw(func(p *ui.Painter, r ui.Rect) {
		x0, w := r.X+pad, r.W-2*pad
		y := r.Y + 14
		xs := func(i int) float32 { return internal.Snap(p, x0+w*float32(i)/float32(n-1)) }
		// Zone bands run under the track.
		for i := 0; i < n-1; i++ {
			_, col, icon := zone(levels[i+1])
			if icon != "" {
				p.Fill(ui.Rect{X: xs(i), Y: y + 4, W: xs(i+1) - xs(i), H: 3}, col, 0)
			}
		}
		p.Line(x0, y, x0+w, y, 2, k.ControlBorder)
		p.Line(x0, y, xs(at), y, 2, k.AccentText)
		size := FontSize(c, theme.CaptionSize)
		widest, _ := p.MeasureText(0, ui.Span{Text: leverageText(levels[n-1]), Size: size})
		// Crowded labels skip levels, counted back from the last.
		every := max(1, int(math.Ceil(float64((widest+8)/(w/float32(n-1))))))
		for i, v := range levels {
			p.Line(xs(i), y+4, xs(i), y+8, 1, k.TextMuted)
			if (n-1-i)%every == 0 {
				tw, _ := p.MeasureText(0, ui.Span{Text: leverageText(v), Size: size})
				p.Text(xs(i)-tw/2, y+11, leverageText(v), size, k.TextMuted)
			}
		}
		p.FillPath(internal.Diamond(xs(at), y, 16), k.Surface)
		p.FillPath(internal.Diamond(xs(at), y, 13), k.AccentText)
		if t.FocusVisible() {
			internal.FocusRing(p, ui.Rect{X: xs(at) - 8, Y: y - 8, W: 16, H: 16}, theme.ControlRadius, k.Focus)
		}
	})
}
