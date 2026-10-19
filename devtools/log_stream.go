package devtools

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// LogStreamState is a LogStream's pause, level filter and place.
type LogStreamState struct {
	Paused bool
	// MinLevel hides lines below it.
	MinLevel code.LogLevel
	// frozen holds the paused line count; later lines wait.
	frozen  int
	shown   []int
	scanned int
	min     code.LogLevel
	rows    ui.ListState
	src     []code.LogLine
	started bool
}

// LogStreamOptions configure a LogStream.
type LogStreamOptions struct {
	// Height of the line well; zero means 320.
	Height   float32
	Disabled bool
}

// LogStream shows virtualized ANSI logs with pause and level filters.
func LogStream(c *ui.Context, lines []code.LogLine, s *LogStreamState, opts LogStreamOptions) core.ChoiceResult {
	if s == nil {
		panic("mujica: LogStream needs a state")
	}
	if opts.Height < 0 || s.MinLevel < code.LogDebug || s.MinLevel > code.LogError {
		panic("mujica: LogStream needs Height >= 0 and a known MinLevel")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	var r core.ChoiceResult
	r.Element = code.CodePanel(c).Disabled(opts.Disabled).Height(h)
	locked := r.Element.IsDisabled()
	if !s.started {
		s.frozen, s.started = len(lines), true
	}
	r.Element.Children(func() {
		code.CodeHeader(c, func() {
			was := s.Paused
			label := core.Msg(c, "log_stream.pause")
			if s.Paused {
				label = core.Msg(c, "log_stream.resume")
			}
			if input.Toggle(c, &s.Paused, label, input.ToggleOptions{Disabled: locked}).Changed() && !locked {
				r = core.NewChoiceResult(r.Element, true)
			} else {
				s.Paused = was
			}
			levels := make([]navigation.Segment, 4)
			for lv := code.LogDebug; lv <= code.LogError; lv++ {
				levels[lv] = navigation.Segment{Label: code.LogViewerName(c, lv)}
			}
			pick := int(s.MinLevel)
			navigation.SegmentedControl(c, &pick, levels, navigation.SegmentedControlOptions{Label: core.Msg(c, "log_stream.level"), Disabled: locked})
			if code.LogLevel(pick) != s.MinLevel && !locked {
				s.MinLevel, r = code.LogLevel(pick), core.NewChoiceResult(r.Element, true)
			}
			ui.Box(c).Grow(1)
			if s.Paused && len(lines) > s.frozen {
				code.CodeCaption(c, fmt.Sprintf(core.Msg(c, "log_stream.waiting"), len(lines)-s.frozen)).Role(ui.RoleStatus)
			}
		})
		if !s.Paused {
			s.frozen = len(lines)
		}
		logStreamFilter(lines[:min(s.frozen, len(lines))], s)
		if len(s.shown) == 0 {
			code.CodeCaption(c, core.Msg(c, "log_stream.empty")).Padding(core.Density(c).Unit() * 4)
			return
		}
		s.rows.FollowEnd = !s.Paused
		s.rows.Key = func(i int) any { return s.shown[i] }
		d := core.Density(c)
		ui.List(c, &s.rows, len(s.shown), func(i int) { logStreamRow(c, lines[s.shown[i]]) }).Grow(1).Padding(d.Unit(), 0)
	})
	if r.Changed() {
		slog.Debug("mujica: log_stream change", "paused", s.Paused, "min_level", s.MinLevel)
	}
	return r
}

// logStreamFilter updates filtered indices when lines or filters change.
func logStreamFilter(lines []code.LogLine, s *LogStreamState) {
	if s.min != s.MinLevel || s.scanned > len(lines) || len(lines) > 0 && (len(s.src) == 0 || &lines[0] != &s.src[0]) {
		s.min, s.scanned, s.shown = s.MinLevel, 0, s.shown[:0]
	}
	for i := s.scanned; i < len(lines); i++ {
		lv := lines[i].Level
		if lv < code.LogDebug || lv > code.LogError {
			panic("mujica: LogLine level out of range")
		}
		if lv >= s.MinLevel {
			s.shown = append(s.shown, i)
		}
	}
	s.scanned, s.src = len(lines), lines
}

// logStreamRow shows one line: time, level word and ANSI text.
func logStreamRow(c *ui.Context, l code.LogLine) {
	k := core.Tokens(c)
	d := core.Density(c)
	col := [...]ui.Color{k.TextMuted, k.Info, k.Warning, k.Danger}[l.Level]
	ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
		if !l.Time.IsZero() {
			code.CodeMono(c, l.Time.Format("15:04:05.000")).TextColor(k.TextMuted).Shrink(0)
		}
		ui.Text(c, fmt.Sprintf("%-5s", code.LogViewerName(c, l.Level))).Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize)).
			FontWeight(700).TextColor(col).SingleLine().Shrink(0)
		ui.RichText(c, code.CodeAnsiSpans(c, l.Text)...).Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize)).Grow(1).MinWidth(0)
	})
}
