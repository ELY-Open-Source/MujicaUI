package project

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/datetime"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PomodoroPhase is the part of the cycle running.
type PomodoroPhase int

const (
	PomodoroFocus PomodoroPhase = iota
	PomodoroBreak
	PomodoroLongBreak
)

// PomodoroState holds the cycle, sessions and running clock.
type PomodoroState struct {
	phase   PomodoroPhase
	focused int
	running bool
	since   time.Time
	spent   time.Duration
}

// Phase returns the current phase.
func (s *PomodoroState) Phase() PomodoroPhase { return s.phase }

// Sessions returns how many focus sessions finished.
func (s *PomodoroState) Sessions() int { return s.focused }

// Running reports whether the clock runs.
func (s *PomodoroState) Running() bool { return s.running }

func (s *PomodoroState) elapsed(now time.Time) time.Duration {
	if s.running {
		return s.spent + now.Sub(s.since)
	}
	return s.spent
}

// PomodoroTimerOptions configure a PomodoroTimer.
type PomodoroTimerOptions struct {
	// Focus/Break/LongBreak default to 25/5/15 minutes.
	Focus, Break, LongBreak time.Duration
	// LongEvery counts sessions per long break; zero means four.
	LongEvery int
	// Now fixes the clock; zero uses the system clock.
	Now      time.Time
	Disabled bool
	Label    string
}

// PomodoroResult is the timer and what changed this frame.
type PomodoroResult struct {
	Element  ui.Element
	changed  bool
	finished PomodoroPhase
	done     bool
}

// Changed reports a start, pause, skip or reset.
func (r PomodoroResult) Changed() bool { return r.changed }

// Finished returns the phase that ran out this frame.
func (r PomodoroResult) Finished() (PomodoroPhase, bool) { return r.finished, r.done }

// PomodoroTimer cycles focus and breaks with pause, skip and reset.
func PomodoroTimer(c *ui.Context, s *PomodoroState, opts PomodoroTimerOptions) PomodoroResult {
	lengths := [3]time.Duration{opts.Focus, opts.Break, opts.LongBreak}
	for i, def := range []time.Duration{25 * time.Minute, 5 * time.Minute, 15 * time.Minute} {
		if lengths[i] < 0 || lengths[i]%time.Second != 0 {
			panic(fmt.Sprintf("mujica: PomodoroTimer length %v is negative or finer than a second", lengths[i]))
		}
		if lengths[i] == 0 {
			lengths[i] = def
		}
	}
	every := opts.LongEvery
	if every < 0 {
		panic("mujica: PomodoroTimer LongEvery is negative")
	}
	if every == 0 {
		every = 4
	}
	now := opts.Now
	live := now.IsZero()
	if live {
		now = c.Now()
	}
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "pomodoro_timer.label")
	}
	r := PomodoroResult{}
	next := func() {
		if s.phase == PomodoroFocus {
			s.focused++
			s.phase = PomodoroBreak
			if s.focused%every == 0 {
				s.phase = PomodoroLongBreak
			}
		} else {
			s.phase = PomodoroFocus
		}
		s.running, s.spent = false, 0
	}
	if s.running && s.elapsed(now) >= lengths[s.phase] {
		r.finished, r.done = s.phase, true
		c.Announce(fmt.Sprintf(core.Msg(c, "pomodoro_timer.done"), core.Msg(c, pomodoroKeys[s.phase])))
		slog.Debug("mujica: pomodoro_timer finished", "phase", s.phase)
		next()
	}
	col := ui.Column(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = col
	act := func(what string, fn func()) {
		if col.IsDisabled() {
			return
		}
		fn()
		r.changed = true
		slog.Debug("mujica: pomodoro_timer "+what, "phase", s.phase)
	}
	left := lengths[s.phase] - s.elapsed(now)
	secs := int(math.Ceil(left.Seconds()))
	frac := 1 - float32(secs)/float32(lengths[s.phase]/time.Second)
	col.Children(func() {
		phase := core.Msg(c, pomodoroKeys[s.phase])
		text := fmt.Sprintf("%02d:%02d", secs/60, secs%60)
		ring := ui.Box(c).Size(160, 160).Center().Role(ui.RoleStatus).Label(fmt.Sprintf(core.Msg(c, "pomodoro_timer.left"), phase, datetime.CountdownText(c, time.Duration(secs)*time.Second, datetime.CountdownWords)))
		tone := k.AccentText
		if s.phase != PomodoroFocus {
			tone = k.Success
		}
		ring.Draw(func(p *ui.Painter, rr ui.Rect) {
			cx, cy := internal.Snap(p, rr.X+rr.W/2), internal.Snap(p, rr.Y+rr.H/2)
			p.FillPath(pomodoroSector(p, cx, cy, 70, 76, 2*math.Pi), k.Border)
			if frac > 0 && !c.Preferences().ReduceMotion {
				p.FillPath(pomodoroSector(p, cx, cy, 70, 76, 2*math.Pi*float64(frac)), tone)
			}
		})
		ring.Children(func() {
			ui.Column(c).AlignItems(ui.Center).Gap(2).Children(func() {
				ui.Text(c, text).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.H1Size)).TextColor(k.Text).FontFeatures("tnum")
				ui.Text(c, phase).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(tone)
			})
		})
		ui.Text(c, fmt.Sprintf(core.Msg(c, "pomodoro_timer.sessions"), s.focused%every, every)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		ui.Row(c).Gap(d.Unit() * 2).Children(func() {
			key, icon, variant := "pomodoro_timer.start", "play", input.Primary
			if s.running {
				key, icon, variant = "pomodoro_timer.pause", "pause", input.Secondary
			}
			if input.Button(c, core.Msg(c, key), input.ButtonOptions{Variant: variant, Icon: icons.Must(icon)}).Clicked() {
				act("toggle", func() {
					if s.running {
						s.spent, s.running = s.elapsed(now), false
					} else {
						s.since, s.running = now, true
					}
				})
			}
			if input.Button(c, core.Msg(c, "pomodoro_timer.skip"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				act("skip", next)
			}
			if input.Button(c, core.Msg(c, "pomodoro_timer.reset"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("rotate-ccw")}).Clicked() {
				act("reset", func() { *s = PomodoroState{} })
			}
		})
	})
	if s.running && live {
		delay := left % time.Second
		if delay == 0 {
			delay = time.Second
		}
		c.After(delay)
	}
	return r
}

var pomodoroKeys = [...]string{"pomodoro_timer.focus", "pomodoro_timer.break", "pomodoro_timer.long_break"}

// pomodoroSector snaps every changing vertex to the cache grid.
func pomodoroSector(p *ui.Painter, cx, cy, inner, outer float32, end float64) *ui.Path {
	var path ui.Path
	n := max(2, int(math.Ceil(end/(math.Pi/48))))
	for i := 0; i <= n; i++ {
		x, y := internal.Polar(cx, cy, outer, end*float64(i)/float64(n))
		x, y = internal.Snap(p, x), internal.Snap(p, y)
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	for i := n; i >= 0; i-- {
		x, y := internal.Polar(cx, cy, inner, end*float64(i)/float64(n))
		path.LineTo(internal.Snap(p, x), internal.Snap(p, y))
	}
	path.Close()
	return &path
}
