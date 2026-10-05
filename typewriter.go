package mujica

import (
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// TypewriterOptions configure a Typewriter.
type TypewriterOptions struct {
	// Speed is characters per second; zero means 40.
	Speed float64
	// Skippable adds a Skip button that shows the rest at once.
	Skippable bool
}

// TypewriterResult is a built Typewriter and its events.
type TypewriterResult struct {
	Element *ui.Element
	done    bool
}

// Finished reports the frame the whole text came into view.
func (r TypewriterResult) Finished() bool { return r.done }

type typewriterState struct {
	text    []rune
	shown   int
	last    time.Time
	skipped bool
	done    bool
}

// Typewriter reveals text a character at a time; text that grows continues where it was.
func Typewriter(c *ui.Context, text string, opts TypewriterOptions) TypewriterResult {
	speed := opts.Speed
	if speed == 0 {
		speed = 40
	}
	if speed < 0 {
		panic("mujica: Typewriter Speed is negative")
	}
	col := ui.Column(c).AlignItems(ui.Start).Gap(Density(c).Unit() * 2)
	res := TypewriterResult{Element: col}
	st := ui.Local(col, "typing", func() *typewriterState { return &typewriterState{last: c.Now()} })
	s := *st
	runes := []rune(text)
	if !strings.HasPrefix(text, string(s.text)) {
		s.shown, s.skipped, s.done = 0, false, false
		slog.Debug("mujica: typewriter restart")
	}
	if len(runes) > len(s.text) {
		s.done = false
	}
	s.text = runes
	now := c.Now()
	if c.Preferences().ReduceMotion || s.skipped {
		s.shown = len(runes)
	} else if s.shown < len(runes) {
		step := time.Duration(float64(time.Second) / speed)
		if n := int(now.Sub(s.last) / step); n > 0 {
			s.shown = min(len(runes), s.shown+n)
			s.last = s.last.Add(time.Duration(n) * step)
		}
		if s.shown < len(runes) {
			c.After(s.last.Add(step).Sub(now))
		}
	} else {
		s.last = now
	}
	col.Children(func() {
		ui.Text(c, string(runes[:s.shown])).Label(text)
		if opts.Skippable && s.shown < len(runes) {
			if Button(c, msg(c, "typewriter.skip"), ButtonOptions{Variant: Ghost}).Clicked() {
				s.skipped, s.shown = true, len(runes)
				slog.Debug("mujica: typewriter skip")
			}
		}
	})
	if s.shown == len(runes) && !s.done {
		s.done, res.done = true, true
	}
	return res
}
