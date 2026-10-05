package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// CarouselOptions configure a Carousel.
type CarouselOptions struct {
	// Label names the carousel for assistive technology.
	Label string
	// Autoplay advances slides every Interval; it pauses on hover, focus or the pause button.
	Autoplay bool
	// Interval defaults to 5 seconds.
	Interval time.Duration
	// Loop wraps from the last slide to the first and back; Autoplay always wraps.
	Loop bool
}

type carouselState struct {
	paused bool
	next   time.Time
	// dotFocus leaves the arrows to the dots' focus group.
	dotFocus bool
}

// Carousel shows slide *index of n with previous, next and position controls.
func Carousel(c *ui.Context, index *int, n int, opts CarouselOptions, slide func(i int)) ChoiceResult {
	if n < 0 || opts.Interval < 0 {
		panic("mujica: Carousel with negative count or Interval")
	}
	interval := opts.Interval
	if interval == 0 {
		interval = 5 * time.Second
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var v ChoiceResult
	v.Element = ui.Column(c).Gap(u * 2).Label(opts.Label)
	st := ui.Local(v.Element, "state", func() carouselState { return carouselState{} })
	if n > 0 && (*index < 0 || *index >= n) {
		slog.Debug("mujica: carousel index clamped", "index", *index, "n", n)
		*index = max(0, min(*index, n-1))
	}
	wrap := opts.Loop || opts.Autoplay
	goTo := func(i int) {
		if n == 0 {
			return
		}
		if wrap {
			i = (i%n + n) % n
		}
		if i < 0 || i >= n || i == *index {
			return
		}
		*index = i
		v.changed = true
		c.Announce(fmt.Sprintf(msg(c, "carousel.position"), i+1, n))
		slog.Debug("mujica: carousel slide", "index", i)
	}
	now := c.Now()
	playing := opts.Autoplay && n > 1 && !st.paused && !v.Element.Hovered() && !v.Element.FocusWithin()
	if !playing {
		st.next = time.Time{}
	} else if st.next.IsZero() {
		st.next = now.Add(interval)
	} else if !now.Before(st.next) {
		goTo(*index + 1)
		st.next = now.Add(interval)
	}
	if playing {
		c.After(st.next.Sub(now))
	}
	switch {
	case st.dotFocus:
	case v.Element.Shortcut(0, ui.KeyLeft):
		goTo(*index - 1)
	case v.Element.Shortcut(0, ui.KeyRight):
		goTo(*index + 1)
	case v.Element.Shortcut(0, ui.KeyHome):
		goTo(0)
	case v.Element.Shortcut(0, ui.KeyEnd):
		goTo(n - 1)
	}
	v.Element.Children(func() {
		ui.Row(c).Grow(1).MinHeight(120).AlignItems(ui.Center).Gap(u * 2).Children(func() {
			prev := Button(c, "", ButtonOptions{Variant: Ghost, Icon: icons.Must("chevron-left"), Label: msg(c, "carousel.prev"), Disabled: n < 2 || !wrap && *index == 0})
			ui.Column(c).Grow(1).MinWidth(0).FillHeight().Clip().Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
				Role(ui.RoleGroup).Label(fmt.Sprintf(msg(c, "carousel.position"), *index+1, n)).Children(func() {
				if n > 0 {
					ui.Column(c).Key(*index).Grow(1).Children(func() { slide(*index) })
				}
			})
			next := Button(c, "", ButtonOptions{Variant: Ghost, Icon: icons.Must("chevron-right"), Label: msg(c, "carousel.next"), Disabled: n < 2 || !wrap && *index == n-1})
			if prev.Clicked() {
				goTo(*index - 1)
			}
			if next.Clicked() {
				goTo(*index + 1)
			}
		})
		ui.Row(c).Justify(ui.Center).AlignItems(ui.Center).Gap(u * 2).Children(func() {
			st.dotFocus = false
			ui.Row(c).Gap(u).FocusGroup(ui.Horizontal).Children(func() {
				for i := range n {
					dot := carouselDot(c, i, n, i == *index)
					if dot.Focused() {
						st.dotFocus = true
					}
					// The arrows move the focus among the dots; the slide follows it.
					if dot.Clicked() || dot.Focused() && i != *index {
						goTo(i)
						c.Invalidate()
					}
				}
			})
			if opts.Autoplay && carouselPause(c, st.paused).Clicked() {
				st.paused = !st.paused
				slog.Debug("mujica: carousel autoplay", "paused", st.paused)
			}
		})
	})
	return v
}

// carouselDot is a position button: a filled diamond for the current slide.
func carouselDot(c *ui.Context, i, n int, on bool) *ui.Element {
	k := Tokens(c)
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).
		Label(fmt.Sprintf(msg(c, "carousel.position"), i+1, n)).Role(ui.RoleTab).Selected(on)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	hover := b.Hovered()
	b.Background(ui.Transparent).Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
		switch {
		case on:
			p.FillPath(internal.Diamond(cx, cy, 12), k.AccentText)
		case hover:
			p.FillPath(internal.Diamond(cx, cy, 8), k.TextMuted)
		default:
			p.StrokePath(internal.Diamond(cx, cy, 8), 1.25, k.TextMuted)
		}
	})
	return b
}

// carouselPause is the pause button: two bars, or a play triangle while paused.
func carouselPause(c *ui.Context, paused bool) *ui.Element {
	k := Tokens(c)
	label := msg(c, "carousel.pause")
	if paused {
		label = msg(c, "carousel.play")
	}
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Label(label)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Pressed() {
		b.Background(k.Selection)
	} else if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	return b.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
		if !paused {
			p.Fill(ui.Rect{X: cx - 4, Y: cy - 5, W: 3, H: 10}, k.AccentText, 0)
			p.Fill(ui.Rect{X: cx + 1, Y: cy - 5, W: 3, H: 10}, k.AccentText, 0)
			return
		}
		var tri ui.Path
		tri.MoveTo(cx-3, cy-5).LineTo(cx+5, cy).LineTo(cx-3, cy+5).Close()
		p.FillPath(&tri, k.AccentText)
	})
}
