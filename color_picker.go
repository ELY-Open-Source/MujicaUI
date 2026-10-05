package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ColorPickerOptions configure a ColorPicker.
type ColorPickerOptions struct {
	// Presets are offered as swatches; nil uses the built-in court palette.
	Presets []ui.Color
	// Label names the trigger; empty uses the built-in "Color".
	Label    string
	Disabled bool
}

// ColorPickerParts are the parts of a ColorPicker.
type ColorPickerParts struct {
	Trigger *ui.Element
	changed bool
}

// Changed reports whether the user picked a new color this frame.
func (p ColorPickerParts) Changed() bool { return p.changed }

// colorPickerPresets are the default swatches: court tones and neutrals.
var colorPickerPresets = []ui.Color{
	ui.Hex("#7D2034"), ui.Hex("#A12635"), ui.Hex("#AD9164"), ui.Hex("#805714"),
	ui.Hex("#326247"), ui.Hex("#355E80"), ui.Hex("#4B3A5A"), ui.Hex("#211C1D"),
	ui.Hex("#6A5B55"), ui.Hex("#D2C3AF"), ui.Hex("#F2EBDD"), ui.Hex("#FFFFFF"),
}

// colorPickerHSVA keeps hue for grays, which RGB loses.
type colorPickerHSVA struct{ h, s, v, a float64 }

type colorPickerState struct {
	open, wasOpen bool
	hsv           colorPickerHSVA
	last          ui.Color
	hex, r, g, b  string
	alpha         string
	bad           string
}

// ColorPicker is a swatch button opening a panel to pick *color: a saturation and
// brightness square, hue and opacity sliders, HEX, RGB and opacity fields, and presets.
func ColorPicker(c *ui.Context, color *ui.Color, opts ColorPickerOptions) ColorPickerParts {
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "color_picker.label")
	}
	presets := opts.Presets
	if presets == nil {
		presets = colorPickerPresets
	}
	parts := ColorPickerParts{}
	trigger := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3, 0, d.Unit()*1.5).Gap(d.Unit() * 2).
		Radius(theme.ControlRadius).Background(k.Surface).Disabled(opts.Disabled).Label(label)
	parts.Trigger = trigger
	internal.Ring(trigger, k.Focus, theme.ControlRadius)
	st := ui.Local(trigger, "picker", func() colorPickerState {
		s := colorPickerState{last: *color, hsv: colorPickerToHSVA(*color)}
		s.sync(*color)
		return s
	})
	if *color != st.last {
		st.hsv, st.last = colorPickerToHSVA(*color), *color
		st.sync(*color)
	}
	set := func(x colorPickerHSVA, fields bool) {
		x.s, x.v, x.a = colorPickerClamp(x.s), colorPickerClamp(x.v), colorPickerClamp(x.a)
		st.hsv = x
		next := x.color()
		if next != *color {
			slog.Debug("mujica: color picker changed", "color", colorPickerHex(next))
			*color = next
			parts.changed = true
		}
		st.last = *color
		if fields {
			st.sync(*color)
		}
	}
	if trigger.Clicked() {
		st.open = !st.open
	}
	border := k.ControlBorder
	if trigger.Hovered() && !opts.Disabled {
		border = k.ControlBorder.Mix(k.Text, 0.25)
		trigger.Background(k.SurfaceHover)
	}
	trigger.Border(theme.BorderWidth, border).Description(colorPickerHex(*color))
	now := *color
	trigger.Children(func() {
		ui.Box(c).Size(d.Unit()*6, d.Unit()*6).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.Border).Clip().Shrink(0).
			Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			colorPickerCheckers(p, r, 4)
			p.Fill(r, now, 0)
		})
		ui.Text(c, strings.ToUpper(colorPickerHex(now))).Font(theme.MonoFont).SingleLine()
		ui.Icon(c, icons.Must("chevron-down")).Size(16, 16).TextColor(k.TextMuted)
	})
	opened := st.open && !st.wasOpen
	ui.PopoverBase(c, trigger, &st.open, func(panel *ui.Element) {
		panel.Margin(d.Unit(), 0, 0, 0).Padding(d.Unit()*3).Gap(d.Unit()*3).Width(260).Radius(theme.PopupRadius).
			Background(k.Surface).Border(theme.BorderWidth, k.Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.18))
		colorPickerSquare(c, st, set, opened)
		x := st.hsv
		h := x.h
		if colorPickerSlider(c, &h, 360, msg(c, "color_picker.hue"), func(p *ui.Painter, r ui.Rect) {
			hues := []ui.Color{ui.RGB(255, 0, 0), ui.RGB(255, 255, 0), ui.RGB(0, 255, 0), ui.RGB(0, 255, 255), ui.RGB(0, 0, 255), ui.RGB(255, 0, 255), ui.RGB(255, 0, 0)}
			w := r.W / 6
			for i := range 6 {
				p.FillGradient(ui.Rect{X: r.X + w*float32(i), Y: r.Y, W: w + 0.5, H: r.H}, ui.LinearGradient{From: hues[i], To: hues[i+1], Angle: 90}, 0)
			}
		}).Changed() {
			set(colorPickerHSVA{h, x.s, x.v, x.a}, true)
		}
		a := x.a
		opaque := colorPickerHSVA{x.h, x.s, x.v, 1}.color()
		if colorPickerSlider(c, &a, 1, msg(c, "color_picker.alpha"), func(p *ui.Painter, r ui.Rect) {
			colorPickerCheckers(p, r, 4)
			p.FillGradient(r, ui.LinearGradient{From: opaque.Alpha(0), To: opaque, Angle: 90}, 0)
		}).Changed() {
			set(colorPickerHSVA{x.h, x.s, x.v, a}, true)
		}
		colorPickerFields(c, st, color, set)
		ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).Gap(d.Unit() * 1.5).Children(func() {
			for _, sw := range presets {
				hex := strings.ToUpper(colorPickerHex(sw))
				b := ui.ButtonBase(c).Height(d.Unit()*6).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.Border).Label(hex).Clip()
				internal.Ring(b, k.Focus, theme.ControlRadius)
				chosen := *color == sw
				b.Selected(chosen).Draw(func(p *ui.Painter, r ui.Rect) {
					colorPickerCheckers(p, r, 4)
					p.Fill(r, sw, 0)
				})
				if chosen {
					mark := ui.RGB(255, 255, 255)
					if colorPickerLuma(sw) > 0.55 {
						mark = ui.RGB(0, 0, 0)
					}
					b.Children(func() { ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(mark) })
				}
				if b.Clicked() {
					set(colorPickerToHSVA(sw), true)
				}
			}
		})
	})
	if st.wasOpen && !st.open {
		trigger.Focus()
	}
	st.wasOpen = st.open
	return parts
}

// colorPickerSquare picks saturation across and brightness up.
func colorPickerSquare(c *ui.Context, st *colorPickerState, set func(colorPickerHSVA, bool), opened bool) {
	k := Tokens(c)
	x := st.hsv
	sq := ui.Box(c).Height(150).Radius(theme.ControlRadius).Focusable().Role(ui.RoleSlider).Cursor(ui.CursorCrosshair).
		Label(msg(c, "color_picker.square"))
	sq.FocusRing(false)
	if opened {
		sq.Focus()
	}
	if _, _, ok := sq.Dragged(); ok {
		px, py, _ := sq.PointerPosition()
		if b := sq.Bounds(); b.W > 0 && b.H > 0 {
			set(colorPickerHSVA{x.h, float64(px / b.W), 1 - float64(py/b.H), x.a}, true)
		}
	}
	switch {
	case sq.Shortcut(0, ui.KeyRight):
		set(colorPickerHSVA{x.h, x.s + 0.01, x.v, x.a}, true)
	case sq.Shortcut(0, ui.KeyLeft):
		set(colorPickerHSVA{x.h, x.s - 0.01, x.v, x.a}, true)
	case sq.Shortcut(0, ui.KeyUp):
		set(colorPickerHSVA{x.h, x.s, x.v + 0.01, x.a}, true)
	case sq.Shortcut(0, ui.KeyDown):
		set(colorPickerHSVA{x.h, x.s, x.v - 0.01, x.a}, true)
	}
	x = st.hsv
	sq.Description(fmt.Sprintf(msg(c, "color_picker.sv"), x.s*100, x.v*100))
	hue := colorPickerHSVA{x.h, 1, 1, 1}.color()
	sq.Draw(func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{From: ui.RGB(255, 255, 255), To: hue, Angle: 90}, theme.ControlRadius)
		p.FillGradient(r, ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0), To: ui.RGB(0, 0, 0), Angle: 180}, theme.ControlRadius)
		cx, cy := internal.Snap(p, r.X+r.W*float32(x.s)), internal.Snap(p, r.Y+r.H*float32(1-x.v))
		p.FillPath(internal.Diamond(cx, cy, 14), ui.RGB(0, 0, 0).Alpha(0.5))
		p.FillPath(internal.Diamond(cx, cy, 11), ui.RGB(255, 255, 255))
		p.FillPath(internal.Diamond(cx, cy, 6), x.color().Alpha(1))
		if sq.FocusVisible() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
}

// colorPickerSlider is a channel slider from 0 to hi over the track track paints.
func colorPickerSlider(c *ui.Context, v *float64, hi float64, name string, track func(p *ui.Painter, r ui.Rect)) *ui.Element {
	k := Tokens(c)
	const knob = 10
	s := ui.SliderBase(c, v, 0, hi).Height(18).PaddingX(knob / 2).FocusRing(false).Label(name)
	frac := float32(*v / hi)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		bar := ui.Rect{X: r.X, Y: r.Y + 3, W: r.W, H: r.H - 6}
		p.Clip(bar, theme.ControlRadius, func() { track(p, bar) })
		p.Stroke(bar, k.Border, theme.ControlRadius, theme.BorderWidth)
		x := r.X + knob/2 + (r.W-knob)*frac
		thumb := ui.Rect{X: x - knob/2, Y: r.Y, W: knob, H: r.H}
		p.Fill(thumb, k.Surface, theme.ControlRadius)
		p.Stroke(thumb, k.Text, theme.ControlRadius, 1.5)
		if s.FocusVisible() {
			internal.FocusRing(p, thumb, theme.ControlRadius, k.Focus)
		}
	})
	return s
}

// sync shows col in every field.
func (st *colorPickerState) sync(col ui.Color) {
	st.hex = strings.TrimPrefix(strings.ToUpper(colorPickerHex(col)), "#")
	st.r, st.g, st.b = strconv.Itoa(int(col.R)), strconv.Itoa(int(col.G)), strconv.Itoa(int(col.B))
	st.alpha = strconv.Itoa(int(math.Round(float64(col.A) / 255 * 100)))
	st.bad = ""
}

// colorPickerFields edits the color as HEX, RGB and opacity text.
func colorPickerFields(c *ui.Context, st *colorPickerState, color *ui.Color, set func(colorPickerHSVA, bool)) {
	k := Tokens(c)
	d := Density(c)
	apply := func(field string, text *string, col ui.Color, ok bool) {
		if !ok {
			st.bad = field
			return
		}
		kept := *text
		set(colorPickerToHSVA(col), true)
		*text = kept
	}
	editing := false
	channel := func(s string) (uint8, bool) {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		return uint8(n), err == nil && n >= 0 && n <= 255
	}
	ui.Column(c).Gap(d.Unit() * 2).Children(func() {
		hex := InputGroup(c, &st.hex, InputGroupOptions{Prefix: "#", Label: "HEX", Invalid: st.bad == "hex"})
		editing = hex.Input.Focused()
		if hex.Changed() {
			col, err := colorPickerParseHex(st.hex)
			apply("hex", &st.hex, col, err == nil)
		}
		ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).Gap(d.Unit() * 1.5).Children(func() {
			for _, f := range []struct {
				name, label string
				text        *string
			}{{"R", msg(c, "color_picker.red"), &st.r}, {"G", msg(c, "color_picker.green"), &st.g},
				{"B", msg(c, "color_picker.blue"), &st.b}, {"A %", msg(c, "color_picker.percent"), &st.alpha}} {
				ui.Column(c).Gap(d.Unit()).Children(func() {
					ui.Text(c, f.name).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
					in := InputGroup(c, f.text, InputGroupOptions{Label: f.label, Invalid: st.bad == f.name})
					editing = editing || in.Input.Focused()
					if !in.Changed() {
						return
					}
					col := *color
					v, ok := channel(*f.text)
					switch f.name {
					case "R":
						col.R = v
					case "G":
						col.G = v
					case "B":
						col.B = v
					default:
						ok = ok && v <= 100
						col.A = uint8(math.Round(float64(v) / 100 * 255))
					}
					apply(f.name, f.text, col, ok)
				})
			}
		})
		if !editing {
			st.sync(*color)
		}
		if st.bad != "" {
			ui.Text(c, msg(c, "color_picker.invalid")).FontSize(theme.CaptionSize).TextColor(k.Danger).Role(ui.RoleStatus)
		}
	})
}

// colorPickerHex returns col as #rrggbb, or #rrggbbaa when translucent.
func colorPickerHex(col ui.Color) string {
	if col.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", col.R, col.G, col.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", col.R, col.G, col.B, col.A)
}

// colorPickerParseHex parses rgb, rrggbb or rrggbbaa, with or without "#".
func colorPickerParseHex(s string) (ui.Color, error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) == 6 {
		h += "ff"
	}
	if len(h) != 8 {
		return ui.Color{}, fmt.Errorf("mujica: color %q is not hex", s)
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return ui.Color{}, fmt.Errorf("mujica: color %q: %w", s, err)
	}
	return ui.Color{R: uint8(v >> 24), G: uint8(v >> 16), B: uint8(v >> 8), A: uint8(v)}, nil
}

func colorPickerToHSVA(col ui.Color) colorPickerHSVA {
	r, g, b := float64(col.R)/255, float64(col.G)/255, float64(col.B)/255
	hi, lo := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	d := hi - lo
	var h float64
	switch {
	case d == 0:
	case hi == r:
		h = 60 * math.Mod((g-b)/d, 6)
	case hi == g:
		h = 60 * ((b-r)/d + 2)
	default:
		h = 60 * ((r-g)/d + 4)
	}
	if h < 0 {
		h += 360
	}
	s := 0.0
	if hi > 0 {
		s = d / hi
	}
	return colorPickerHSVA{h, s, hi, float64(col.A) / 255}
}

func (x colorPickerHSVA) color() ui.Color {
	ch := x.v * x.s
	h := math.Mod(x.h, 360) / 60
	y := ch * (1 - math.Abs(math.Mod(h, 2)-1))
	var r, g, b float64
	switch {
	case h < 1:
		r, g = ch, y
	case h < 2:
		r, g = y, ch
	case h < 3:
		g, b = ch, y
	case h < 4:
		g, b = y, ch
	case h < 5:
		r, b = y, ch
	default:
		r, b = ch, y
	}
	m := x.v - ch
	byte8 := func(v float64) uint8 { return uint8(math.Round(colorPickerClamp(v) * 255)) }
	return ui.Color{R: byte8(r + m), G: byte8(g + m), B: byte8(b + m), A: byte8(x.a)}
}

func colorPickerClamp(v float64) float64 { return math.Max(0, math.Min(1, v)) }

// colorPickerLuma is the relative brightness of col, 0 to 1.
func colorPickerLuma(col ui.Color) float64 {
	return (0.2126*float64(col.R) + 0.7152*float64(col.G) + 0.0722*float64(col.B)) / 255
}

// colorPickerCheckers paints the grid that shows transparency.
func colorPickerCheckers(p *ui.Painter, r ui.Rect, size float32) {
	p.Fill(r, ui.RGB(255, 255, 255), 0)
	for y, row := r.Y, 0; y < r.Y+r.H; y, row = y+size, row+1 {
		for x, col := r.X, 0; x < r.X+r.W; x, col = x+size, col+1 {
			if (row+col)%2 == 0 {
				p.Fill(ui.Rect{X: x, Y: y, W: min(size, r.X+r.W-x), H: min(size, r.Y+r.H-y)}, ui.RGB(204, 204, 204), 0)
			}
		}
	}
}
