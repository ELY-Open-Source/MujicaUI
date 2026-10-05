package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// LogStreamState is a LogStream's pause, level filter and place.
type LogStreamState struct {
	Paused bool
	// MinLevel hides lines below it.
	MinLevel LogLevel
	// frozen is the line count at pause; lines after it wait.
	frozen  int
	shown   []int
	scanned int
	min     LogLevel
	rows    ui.ListState
}

// LogStreamOptions configure a LogStream.
type LogStreamOptions struct {
	// Height of the line well; zero means 320.
	Height   float32
	Disabled bool
}

// LogStream shows live ANSI log lines that follow the end, pause and filter by level, virtualized.
func LogStream(c *ui.Context, lines []LogLine, s *LogStreamState, opts LogStreamOptions) ChoiceResult {
	if s == nil {
		panic("mujica: LogStream needs a state")
	}
	if opts.Height < 0 || s.MinLevel < LogDebug || s.MinLevel > LogError {
		panic("mujica: LogStream needs Height >= 0 and a known MinLevel")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	var r ChoiceResult
	r.Element = codePanel(c).Disabled(opts.Disabled).Height(h)
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		codeHeader(c, func() {
			was := s.Paused
			label := msg(c, "log_stream.pause")
			if s.Paused {
				label = msg(c, "log_stream.resume")
			}
			if Toggle(c, &s.Paused, label, ToggleOptions{Disabled: locked}).Changed() && !locked {
				r.changed = true
			} else {
				s.Paused = was
			}
			levels := make([]Segment, 4)
			for lv := LogDebug; lv <= LogError; lv++ {
				levels[lv] = Segment{Label: logViewerName(c, lv)}
			}
			pick := int(s.MinLevel)
			SegmentedControl(c, &pick, levels, SegmentedControlOptions{Label: msg(c, "log_stream.level"), Disabled: locked})
			if LogLevel(pick) != s.MinLevel && !locked {
				s.MinLevel, r.changed = LogLevel(pick), true
			}
			ui.Box(c).Grow(1)
			if s.Paused && len(lines) > s.frozen {
				codeCaption(c, fmt.Sprintf(msg(c, "log_stream.waiting"), len(lines)-s.frozen)).Role(ui.RoleStatus)
			}
		})
		if !s.Paused {
			s.frozen = len(lines)
		}
		logStreamFilter(lines[:min(s.frozen, len(lines))], s)
		if len(s.shown) == 0 {
			codeCaption(c, msg(c, "log_stream.empty")).Padding(Density(c).Unit() * 4)
			return
		}
		s.rows.FollowEnd = !s.Paused
		s.rows.Key = func(i int) any { return s.shown[i] }
		d := Density(c)
		ui.List(c, &s.rows, len(s.shown), func(i int) { logStreamRow(c, lines[s.shown[i]]) }).Grow(1).Padding(d.Unit(), 0)
	})
	if r.changed {
		slog.Debug("mujica: log_stream change", "paused", s.Paused, "min_level", s.MinLevel)
	}
	return r
}

// logStreamFilter extends the shown lines, or redoes them after a filter change.
func logStreamFilter(lines []LogLine, s *LogStreamState) {
	if s.min != s.MinLevel || s.scanned > len(lines) {
		s.min, s.scanned, s.shown = s.MinLevel, 0, s.shown[:0]
	}
	for i := s.scanned; i < len(lines); i++ {
		lv := lines[i].Level
		if lv < LogDebug || lv > LogError {
			panic("mujica: LogLine level out of range")
		}
		if lv >= s.MinLevel {
			s.shown = append(s.shown, i)
		}
	}
	s.scanned = len(lines)
}

// logStreamRow shows one line: time, level word and ANSI text.
func logStreamRow(c *ui.Context, l LogLine) {
	k := Tokens(c)
	d := Density(c)
	col := [...]ui.Color{k.TextMuted, k.Info, k.Warning, k.Danger}[l.Level]
	ui.Row(c).MinHeight(d.ControlHeight() - 8).Padding(0, d.Unit()*3).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
		if !l.Time.IsZero() {
			codeMono(c, l.Time.Format("15:04:05.000")).TextColor(k.TextMuted).Shrink(0)
		}
		ui.Text(c, fmt.Sprintf("%-5s", logViewerName(c, l.Level))).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).
			FontWeight(700).TextColor(col).SingleLine().Shrink(0)
		ui.RichText(c, codeAnsiSpans(c, l.Text)...).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).Grow(1).MinWidth(0)
	})
}
