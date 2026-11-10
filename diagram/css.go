package diagram

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

// palette holds Pele's theme tokens.
type palette struct {
	bg, fg, muted, line, surface, alt, border, accent, critical ui.Color
	series                                                      [8]ui.Color
	radius                                                      float32
}

// token resolves a custom property Pele names.
func (pal *palette) token(name string) (ui.Color, bool) {
	switch name {
	case "--_bg", "--pele-bg":
		return pal.bg, true
	case "--_fg", "--pele-fg":
		return pal.fg, true
	case "--_m", "--pele-muted":
		return pal.muted, true
	case "--_l", "--pele-line":
		return pal.line, true
	case "--_s", "--pele-surface":
		return pal.surface, true
	case "--_a", "--pele-surface-alt":
		return pal.alt, true
	case "--_b", "--pele-border":
		return pal.border, true
	case "--_c", "--pele-accent":
		return pal.accent, true
	case "--pele-critical":
		return pal.critical, true
	}
	if n, ok := strings.CutPrefix(name, "--pele-series-"); ok && len(n) == 1 && n[0] >= '1' && n[0] <= '8' {
		return pal.series[n[0]-'1'], true
	}
	return ui.Color{}, false
}

// weights are Pele's weight tokens; the theme has one bold.
var weights = map[string]bool{"--_w": true, "--_hw": true, "--_tw": true,
	"--pele-bold-weight": true, "--pele-heading-weight": true, "--pele-title-weight": true}

// reserved reports a name in Pele's namespaces.
func reserved(name string) bool {
	return strings.HasPrefix(name, "--_") || strings.HasPrefix(name, "--pele-")
}

// plain swaps each var() outside a color; "" is dropped.
func (w *walker) plain(v string) (string, error) {
	for {
		at := strings.Index(v, "var(")
		if at < 0 {
			return v, nil
		}
		depth, end := 0, -1
		for i := at + 3; i < len(v) && end < 0; i++ {
			switch v[i] {
			case '(':
				depth++
			case ')':
				if depth--; depth == 0 {
					end = i
				}
			}
		}
		if end < 0 {
			return "", nil
		}
		name, sub, has := strings.Cut(v[at+4:end], ",")
		switch name = strings.TrimSpace(name); {
		case name == "--_r" || name == "--pele-radius":
			sub = num(w.pal.radius)
		case name == "--_fm" || name == "--pele-font-mono":
			sub = "monospace"
		case name == "--pele-font":
			sub = w.font
		case weights[name]:
			sub = "bold"
		case reserved(name):
			return "", fmt.Errorf("unknown token %q", name)
		case !has:
			return "", nil
		}
		v = v[:at] + sub + v[end+1:]
	}
}

// color reads a CSS color; ok is false when dropped.
func (pal *palette) color(v string, current ui.Color) (c ui.Color, ok bool, err error) {
	v = strings.ToLower(strings.TrimSpace(v))
	switch {
	case v == "none" || v == "transparent":
		return ui.Transparent, true, nil
	case v == "currentcolor":
		return current, true, nil
	case strings.HasPrefix(v, "#"):
		c, ok = hex(v[1:])
		return c, ok, nil
	}
	name, body, call := strings.Cut(v, "(")
	if !call {
		c, ok = named[v]
		return c, ok, nil
	}
	body, closed := strings.CutSuffix(body, ")")
	if !closed {
		return c, false, nil
	}
	switch name {
	case "var":
		token, fallback, has := strings.Cut(body, ",")
		token = strings.TrimSpace(token)
		if c, ok = pal.token(token); ok {
			return c, true, nil
		}
		if reserved(token) {
			return c, false, fmt.Errorf("unknown token %q", token)
		}
		if has {
			return pal.color(fallback, current)
		}
		return c, false, nil
	case "color-mix":
		return pal.mix(body, current)
	case "rgb", "rgba":
		c, ok = rgb(numbers(body))
	case "hsl", "hsla":
		c, ok = hsl(numbers(body))
	}
	return c, ok, nil
}

// mix reads color-mix(in srgb, a p%, b q%).
func (pal *palette) mix(body string, current ui.Color) (ui.Color, bool, error) {
	parts := split(body)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "in srgb" {
		return ui.Color{}, false, nil
	}
	first, p, hasP, okP := share(parts[1])
	second, q, hasQ, okQ := share(parts[2])
	switch {
	case !hasP && !hasQ:
		p, q = 0.5, 0.5
	case !hasQ:
		q = 1 - p
	case !hasP:
		p = 1 - q
	}
	if !okP || !okQ || p+q <= 0 {
		return ui.Color{}, false, nil
	}
	a, ok, err := pal.color(first, current)
	if !ok || err != nil {
		return a, ok, err
	}
	b, ok, err := pal.color(second, current)
	return blend(a, b, p, q), ok, err
}

// share cuts a color-mix side into color and percentage.
func share(side string) (color string, p float32, has, ok bool) {
	side = strings.TrimSpace(side)
	last, first := strings.LastIndexByte(side, ' '), strings.IndexByte(side, ' ')
	switch {
	case last >= 0 && strings.HasSuffix(side, "%"):
		color, side = side[:last], side[last+1:]
	case first >= 0 && strings.HasSuffix(side[:first], "%"):
		color, side = side[first+1:], side[:first]
	default:
		return side, 0, false, true
	}
	v, ok := number(strings.TrimSuffix(side, "%"))
	return color, float32(v) / 100, true, ok && v >= 0 && v <= 100
}

// blend mixes premultiplied, as CSS; shares under 1 fade.
func blend(a, b ui.Color, p, q float32) ui.Color {
	sum := p + q
	wa, wb := float32(a.A)/255*p/sum, float32(b.A)/255*q/sum
	alpha := wa + wb
	if alpha == 0 {
		return ui.Transparent
	}
	ch := func(x, y uint8) uint8 { return uint8((float32(x)*wa+float32(y)*wb)/alpha + 0.5) }
	return ui.Color{R: ch(a.R, b.R), G: ch(a.G, b.G), B: ch(a.B, b.B), A: uint8(alpha*min(sum, 1)*255 + 0.5)}
}

// split cuts s at commas outside parentheses.
func split(s string) []string {
	var out []string
	depth, from := 0, 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out, from = append(out, s[from:i]), i+1
			}
		}
	}
	return append(out, s[from:])
}

func hex(h string) (ui.Color, bool) {
	var v [8]uint8
	if n := len(h); n != 3 && n != 4 && n != 6 && n != 8 {
		return ui.Color{}, false
	}
	for i := range len(h) {
		d, err := strconv.ParseUint(h[i:i+1], 16, 8)
		if err != nil {
			return ui.Color{}, false
		}
		v[i] = uint8(d)
	}
	if len(h) <= 4 {
		a := uint8(15)
		if len(h) == 4 {
			a = v[3]
		}
		return ui.Color{R: v[0] * 17, G: v[1] * 17, B: v[2] * 17, A: a * 17}, true
	}
	a := uint8(255)
	if len(h) == 8 {
		a = v[6]<<4 | v[7]
	}
	return ui.Color{R: v[0]<<4 | v[1], G: v[2]<<4 | v[3], B: v[4]<<4 | v[5], A: a}, true
}

// numbers reads a color function's arguments; percent scales to 1.
func numbers(body string) (vs []float64, pct []bool) {
	for _, f := range strings.FieldsFunc(body, func(r rune) bool { return r == ',' || r == ' ' || r == '/' }) {
		f, p := strings.CutSuffix(f, "%")
		f = strings.TrimSuffix(f, "deg")
		v, ok := number(f)
		if !ok {
			return nil, nil
		}
		if p {
			v /= 100
		}
		vs, pct = append(vs, v), append(pct, p)
	}
	return vs, pct
}

func unit(v float64) uint8 { return uint8(math.Round(255 * max(0, min(v, 1)))) }

func alpha(vs []float64) uint8 {
	if len(vs) == 4 {
		return unit(vs[3])
	}
	return 255
}

func rgb(vs []float64, pct []bool) (ui.Color, bool) {
	if len(vs) != 3 && len(vs) != 4 {
		return ui.Color{}, false
	}
	var ch [3]uint8
	for i := range ch {
		if pct[i] {
			ch[i] = unit(vs[i])
		} else {
			ch[i] = unit(vs[i] / 255)
		}
	}
	return ui.Color{R: ch[0], G: ch[1], B: ch[2], A: alpha(vs)}, true
}

func hsl(vs []float64, pct []bool) (ui.Color, bool) {
	if len(vs) != 3 && len(vs) != 4 || pct[0] || !pct[1] || !pct[2] {
		return ui.Color{}, false
	}
	h, s, l := math.Mod(math.Mod(vs[0], 360)+360, 360)/30, max(0, min(vs[1], 1)), max(0, min(vs[2], 1))
	a := s * min(l, 1-l)
	f := func(n float64) uint8 {
		k := math.Mod(n+h, 12)
		return unit(l - a*max(-1, min(k-3, 9-k, 1)))
	}
	return ui.Color{R: f(0), G: f(8), B: f(4), A: alpha(vs)}, true
}

// number reads a finite CSS number.
func number(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 32)
	return f, err == nil && !math.IsInf(f, 0) && strings.Trim(s, "0123456789.+-eE") == ""
}

// length reads a number, with or without px.
func length(v string) (float32, bool) {
	f, ok := number(strings.TrimSuffix(strings.TrimSpace(v), "px"))
	return float32(f), ok
}

// lengths reads a list of numbers, as a dash array.
func lengths(v string) ([]float32, bool) {
	var out []float32
	for _, f := range strings.FieldsFunc(v, func(r rune) bool { return r == ',' || r == ' ' }) {
		n, ok := length(f)
		if !ok {
			return nil, false
		}
		out = append(out, n)
	}
	return out, true
}
