package core

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// overlayLayers is the window's open MujicaUI layers, oldest first.
type overlayLayers struct {
	root ui.Element // compared only: a new root means a new build pass
	pass int
	keys []any
	seen []int
}

type overlayLayersKey struct{}

// OverlayLayer marks the layer key open in this frame and reports whether it is on top.
func OverlayLayer(c *ui.Context, key any) bool {
	s := ui.Local(c.Root(), overlayLayersKey{}, func() overlayLayers { return overlayLayers{} })
	if root := c.Root(); root != s.root {
		s.root = root
		s.pass++
		n := 0
		for i, k := range s.keys {
			if s.seen[i] >= s.pass-1 { // drop layers the last pass did not build
				s.keys[n], s.seen[n] = k, s.seen[i]
				n++
			}
		}
		clear(s.keys[n:])
		s.keys, s.seen = s.keys[:n], s.seen[:n]
	}
	at := -1
	for i, k := range s.keys {
		if k == key {
			at = i
		}
	}
	if at < 0 {
		s.keys, s.seen = append(s.keys, key), append(s.seen, s.pass)
		at = len(s.keys) - 1
	}
	s.seen[at] = s.pass
	return at == len(s.keys)-1
}

// OverlayEscape reports Escape for a window-level layer, delivered only while it is on top.
func OverlayEscape(c *ui.Context, key any) bool {
	if !OverlayLayer(c, key) {
		return false
	}
	// The root is around every focus, so it takes Escape before MyGo's overlays.
	byFocus := c.Root().Shortcut(0, ui.KeyEscape)
	byWindow := c.Shortcut(0, ui.KeyEscape)
	if byFocus || byWindow {
		s := ui.Local(c.Root(), overlayLayersKey{}, func() overlayLayers { return overlayLayers{} })
		for i, k := range s.keys {
			if k == key {
				s.keys = append(s.keys[:i], s.keys[i+1:]...)
				s.seen = append(s.seen[:i], s.seen[i+1:]...)
				break
			}
		}
		c.Invalidate()
		return true
	}
	return false
}

// overlaySiteState follows the focus of a layer after it closes.
type overlaySiteState struct {
	focusIn bool // the focus was inside the layer when it last showed
	wait    int  // frames left to check the focus after a close
	holding bool // the site holds the fallback focus
}

// OverlaySite builds the invisible focus fallback at a layer's call site.
func OverlaySite(c *ui.Context) (ui.Element, *overlaySiteState) {
	site := ui.Box(c).Size(0, 0).FocusRing(false)
	st := ui.Local(site, "mujica.overlay.site", func() overlaySiteState { return overlaySiteState{} })
	st.holding = st.holding && site.FocusWithin()
	if st.holding {
		site.Focusable() // in flow, so Tab goes on from the layer's place
	} else {
		site.Absolute()
	}
	return site, st
}

// Shown records whether the focus is in the open layer's panel.
func (st *overlaySiteState) Shown(panel ui.Element) {
	st.focusIn, st.wait = panel.Valid() && panel.FocusWithin(), 0
}

// Closed moves a lost focus to the site once MyGo has had its chance to restore it.
func (st *overlaySiteState) Closed(c *ui.Context, site ui.Element) {
	if st.focusIn {
		st.focusIn, st.wait = false, 2 // MyGo prunes the panel at the end of the next frame
	}
	if st.wait == 0 {
		return
	}
	st.wait--
	if st.wait > 0 {
		c.AnimationFrame()
		return
	}
	// The site itself counts as lost: it may inherit a vanished opener's positional ID.
	if !c.Root().FocusWithin() || site.FocusWithin() {
		slog.Debug("mujica: overlay focus fallback")
		site.Focus()
		st.holding = true
	}
}
