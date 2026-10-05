package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/egoist/mygo/ui"
)

// input is one browser event.
type input struct {
	Type   string  `json:"type"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	DX     float32 `json:"dx"`
	DY     float32 `json:"dy"`
	Button int     `json:"button"`
	Code   string  `json:"code"`
	Text   string  `json:"text"`
	Caret  int     `json:"caret"`
	Mods   []bool  `json:"mods"`
	W      int     `json:"w"`
	H      int     `json:"h"`
	Scale  float32 `json:"scale"`
	Dark   bool    `json:"dark"`
	Motion bool    `json:"reduceMotion"`
	More   bool    `json:"moreContrast"`
}

func (m *mirror) serveInput(w http.ResponseWriter, r *http.Request) {
	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.mu.Lock()
	err := m.applyWatchingClipboard(in)
	if err == nil {
		m.capture()
	}
	m.mu.Unlock()
	if err != nil {
		slog.Warn("gallery: web input rejected", "type", in.Type, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.kick()
	w.WriteHeader(http.StatusNoContent)
}

// apply sends one event to the tester; the caller holds mu.
func (m *mirror) apply(in input) error {
	mods := modifiers(in.Mods)
	switch in.Type {
	case "down":
		switch {
		case in.Button == 2:
			m.tt.RightClickAt(in.X, in.Y)
		case mods != 0:
			m.down = &in
		default:
			m.tt.Press(in.X, in.Y)
		}
	case "up":
		if d := m.down; d != nil {
			m.down = nil
			if dx, dy := in.X-d.X, in.Y-d.Y; dx*dx+dy*dy <= 16 {
				// The click lands where it is released, as a real press and release would.
				m.tt.ClickAtWith(modifiers(d.Mods), in.X, in.Y)
			} else {
				slog.Info("gallery: modified drag ignored; tester clicks only", "mods", modifiers(d.Mods))
			}
		} else if in.Button == 0 {
			m.tt.Release(in.X, in.Y)
		}
	case "move":
		m.tt.Move(in.X, in.Y)
	case "wheel":
		m.tt.Scroll(in.X, in.Y, in.DX, in.DY)
	case "key":
		k, ok := keyCodes[in.Code]
		if !ok {
			return fmt.Errorf("unmapped key code %q", in.Code)
		}
		if in.Text != "" && mods&^(ui.Shift|ui.Alt) == 0 {
			m.tt.TypeKey(mods, k, in.Text)
		} else {
			m.tt.Key(mods, k)
		}
	case "text":
		m.tt.Type(in.Text)
	case "compose":
		m.tt.Compose(in.Text, in.Caret)
	case "paste":
		m.tt.SetClipboard(in.Text)
		m.clip = in.Text
		m.tt.Command("paste")
	case "command":
		m.tt.Command(in.Text)
		m.copied = in.Text == "copy" || in.Text == "cut"
	case "resize":
		if in.W < 320 || in.H < 240 || in.Scale <= 0 {
			return fmt.Errorf("bad size %dx%d@%v", in.W, in.H, in.Scale)
		}
		m.tt.SetSize(in.W, in.H)
		m.tt.SetScale(in.Scale)
		m.tt.SetDark(in.Dark)
		m.tt.SetPreferences(ui.Preferences{ReduceMotion: in.Motion, HighContrast: in.More, TextScale: 1})
	case "focus":
		m.tt.SetFocused(in.Text == "in")
	default:
		return fmt.Errorf("unknown input type %q", in.Type)
	}
	return nil
}

func modifiers(m []bool) ui.Modifiers {
	var out ui.Modifiers
	for i, mod := range []ui.Modifiers{ui.Shift, ui.Ctrl, ui.Alt, ui.Super} {
		if i < len(m) && m[i] {
			out |= mod
		}
	}
	return out
}

// applyWatchingClipboard applies in and notes any clipboard write it caused, even of equal text.
func (m *mirror) applyWatchingClipboard(in input) error {
	if in.Type == "paste" || in.Type == "command" && in.Text == "paste" {
		return m.apply(in)
	}
	before := m.tt.Clipboard()
	m.probes++
	probe := fmt.Sprintf("\x00mujica-mirror-probe-%d", m.probes)
	m.tt.SetClipboard(probe)
	err := m.apply(in)
	if m.tt.Clipboard() == probe {
		m.tt.SetClipboard(before)
	} else {
		m.copied = true
	}
	return err
}
