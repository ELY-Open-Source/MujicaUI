package diagram

import (
	"fmt"
	"strconv"
	"strings"
)

// numbers reads geometry attributes; a missing one is zero.
func (f *frame) numbers(names ...string) ([]float32, error) {
	out := make([]float32, len(names))
	for i, name := range names {
		v, has := f.geo[name]
		if !has {
			continue
		}
		n, ok := length(v)
		if !ok {
			return nil, fmt.Errorf("<%s> %s=%q is not a number", f.name, name, v)
		}
		out[i] = n
	}
	return out, nil
}

// corner reads rx or ry; a percentage is of side.
func (st *style) corner(v string, side float32) (float32, bool) {
	v = strings.ToLower(strings.TrimSpace(v))
	if p, ok := strings.CutSuffix(v, "%"); ok {
		k, ok := number(p)
		return side * float32(k) / 100, ok && k >= 0
	}
	n, ok := st.measure(v)
	return n, ok && n >= 0
}

func num(v float32) string { return strconv.FormatFloat(float64(v), 'f', -1, 32) }

// data is the element's outline as SVG path data.
func (f *frame) data() (string, error) {
	switch f.name {
	case "path":
		return f.geo["d"], nil
	case "polygon":
		return "M" + f.geo["points"] + "Z", nil
	case "rect":
		return f.rect()
	}
	n, err := f.numbers("cx", "cy", "r", "rx", "ry")
	if err != nil {
		return "", err
	}
	rx, ry := n[3], n[4]
	if f.name == "circle" {
		rx, ry = n[2], n[2]
	}
	arc := " A" + num(rx) + " " + num(ry) + " 0 1 0 "
	return "M" + num(n[0]-rx) + " " + num(n[1]) + arc + num(n[0]+rx) + " " + num(n[1]) + arc + num(n[0]-rx) + " " + num(n[1]) + "Z", nil
}

// rect is a box, rounded by rx and ry.
func (f *frame) rect() (string, error) {
	n, err := f.numbers("x", "y", "width", "height")
	if err != nil {
		return "", err
	}
	x, y, w, h := n[0], n[1], n[2], n[3]
	rx, hasX := f.st.corner(f.geo["rx"], w)
	ry, hasY := f.st.corner(f.geo["ry"], h)
	if !hasX {
		rx = ry
	}
	if !hasY {
		ry = rx
	}
	rx, ry = min(rx, w/2), min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		return "M" + num(x) + " " + num(y) + "h" + num(w) + "v" + num(h) + "h" + num(-w) + "Z", nil
	}
	arc := func(ex, ey float32) string {
		return "A" + num(rx) + " " + num(ry) + " 0 0 1 " + num(ex) + " " + num(ey)
	}
	return "M" + num(x+rx) + " " + num(y) + "H" + num(x+w-rx) + arc(x+w, y+ry) + "V" + num(y+h-ry) + arc(x+w-rx, y+h) +
		"H" + num(x+rx) + arc(x, y+h-ry) + "V" + num(y+ry) + arc(x+rx, y) + "Z", nil
}
