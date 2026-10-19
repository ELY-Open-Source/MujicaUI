package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ColorSwatch is one preset of a ColorPalette; Name is read out.
type ColorSwatch struct {
	Color ui.Color
	Name  string
}

// ColorPaletteOptions configure a ColorPalette.
type ColorPaletteOptions struct {
	// Columns is the grid width in swatches; zero is 8.
	Columns int
	// Custom adds a ColorPicker after the swatches for any other color.
	Custom   bool
	Disabled bool
	// Label names the palette for assistive technology.
	Label string
}

// ColorPalette picks *value from preset swatches in a grid; arrows move in two dimensions.
func ColorPalette(c *ui.Context, value *ui.Color, swatches []ColorSwatch, opts ColorPaletteOptions) core.ChoiceResult {
	if value == nil || len(swatches) == 0 {
		panic("mujica: ColorPalette needs a value and swatches")
	}
	if opts.Columns < 0 {
		panic(fmt.Sprintf("mujica: ColorPalette Columns %d is negative", opts.Columns))
	}
	cols := opts.Columns
	if cols == 0 {
		cols = 8
	}
	for _, s := range swatches {
		if s.Name == "" {
			panic("mujica: ColorPalette swatch needs a Name")
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	g := ui.Column(c).Gap(u * 2).Disabled(opts.Disabled).Role(ui.RoleRadioGroup).Label(opts.Label)
	disabled := g.IsDisabled()
	changed := false
	g.Children(func() {
		grid := ui.Box(c).Grid().Columns(cols).Gap(u * 1.5).FocusGroup(ui.Horizontal)
		var cells []ui.Element
		grid.Children(func() {
			for i, s := range swatches {
				b := colorPaletteSwatch(c, k, s, *value == s.Color).Key(i)
				cells = append(cells, b)
				if b.Clicked() && !disabled && *value != s.Color {
					*value, changed = s.Color, true
				}
			}
		})
		// Up and Down move a whole row; FocusGroup moves by one.
		for i, b := range cells {
			switch {
			case b.Shortcut(0, ui.KeyDown) && i+cols < len(cells):
				cells[i+cols].Focus()
			case b.Shortcut(0, ui.KeyUp) && i-cols >= 0:
				cells[i-cols].Focus()
			}
		}
		if opts.Custom {
			ui.Row(c).Children(func() {
				if ColorPicker(c, value, ColorPickerOptions{Label: core.Msg(c, "color_palette.custom"), Disabled: disabled}).Changed() {
					changed = true
				}
			})
		}
	})
	if changed {
		slog.Debug("mujica: color_palette changed", "color", ColorPickerHex(*value))
	}
	return core.NewChoiceResult(g, changed)
}

// colorPaletteSwatch is one square swatch; the chosen one has an inset frame and a check.
func colorPaletteSwatch(c *ui.Context, k theme.Tokens, s ColorSwatch, on bool) ui.Element {
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Role(ui.RoleRadio).
		Label(s.Name).Tooltip(s.Name).Selected(on)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	edge := k.ControlBorder
	if b.Hovered() {
		edge = k.Text
	}
	b.Background(s.Color).Border(theme.BorderWidth, edge)
	mark := theme.Light().Text
	if theme.Contrast(s.Color, mark) < theme.Contrast(s.Color, theme.Light().OnAccent) {
		mark = theme.Light().OnAccent
	}
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !on {
			return
		}
		p.Stroke(ui.Rect{X: r.X + 3, Y: r.Y + 3, W: r.W - 6, H: r.H - 6}, mark, 1, 1.5)
		internal.Mark(p, ui.Rect{X: r.X + 7, Y: r.Y + 7, W: r.W - 14, H: r.H - 14}, "check", mark)
	})
	return b
}
