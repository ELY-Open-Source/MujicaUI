package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// LogLevel is a log line's severity.
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

// LogLine is one entry of a LogViewer.
type LogLine struct {
	Time  time.Time
	Level LogLevel
	Text  string
}

// LogViewerState is a LogViewer's filters, highlight and place.
type LogViewerState struct {
	// Hidden levels are filtered out.
	Hidden [4]bool
	// Query is highlighted in the shown lines.
	Query string
	// Follow keeps the newest line in view; it starts on.
	Follow  bool
	rows    ui.ListState
	shown   []int
	scanned int
	hidden  [4]bool
	started bool
}

// LogViewerOptions configure a LogViewer.
type LogViewerOptions struct {
	Disabled bool
}

// LogViewer shows filtered, highlighted log lines, virtualized.
func LogViewer(c *ui.Context, lines []LogLine, s *LogViewerState, opts LogViewerOptions) ChoiceResult {
	if s == nil {
		panic("mujica: LogViewer needs a state")
	}
	if !s.started {
		s.started, s.Follow = true, true
	}
	var res ChoiceResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	res.Element.Children(func() {
		codeHeader(c, func() {
			for lv := LogDebug; lv <= LogError; lv++ {
				on := !s.Hidden[lv]
				if Toggle(c, &on, logViewerName(c, lv), ToggleOptions{Disabled: opts.Disabled}).Changed() {
					s.Hidden[lv] = !on
					res.changed = true
					slog.Debug("mujica: log_viewer level", "level", lv, "shown", on)
				}
			}
			SearchInput(c, &s.Query, SearchInputOptions{Disabled: opts.Disabled, Placeholder: msg(c, "log_viewer.highlight")}).Element.Grow(1).MinWidth(0)
			was := s.Follow
			if Toggle(c, &s.Follow, msg(c, "log_viewer.follow"), ToggleOptions{Disabled: opts.Disabled}).Changed() {
				res.changed = true
				slog.Debug("mujica: log_viewer follow", "on", s.Follow)
				if s.Follow && !was {
					s.rows.ScrollToEnd()
				}
			}
		})
		logViewerFilter(lines, s)
		if len(s.shown) == 0 {
			codeCaption(c, msg(c, "log_viewer.empty")).Padding(Density(c).Unit() * 4)
			return
		}
		s.rows.FollowEnd = s.Follow
		s.rows.Key = func(i int) any { return s.shown[i] }
		s.rows.Label = func(i int) string { return lines[s.shown[i]].Text }
		d := Density(c)
		ui.List(c, &s.rows, len(s.shown), func(i int) {
			logViewerRow(c, lines[s.shown[i]], s.Query)
		}).Grow(1).MinHeight(d.ControlHeight()*3).Padding(d.Unit(), 0)
	})
	return res
}

// logViewerFilter extends or redoes the shown lines.
func logViewerFilter(lines []LogLine, s *LogViewerState) {
	if s.hidden != s.Hidden || s.scanned > len(lines) {
		s.hidden, s.scanned, s.shown = s.Hidden, 0, s.shown[:0]
	}
	for i := s.scanned; i < len(lines); i++ {
		lv := lines[i].Level
		if lv < LogDebug || lv > LogError {
			panic("mujica: LogLine level out of range")
		}
		if !s.Hidden[lv] {
			s.shown = append(s.shown, i)
		}
	}
	s.scanned = len(lines)
}

func logViewerName(c *ui.Context, lv LogLevel) string {
	return msg(c, [...]string{"log_viewer.debug", "log_viewer.info", "log_viewer.warn", "log_viewer.error"}[lv])
}

func logViewerRow(c *ui.Context, l LogLine, query string) {
	k := Tokens(c)
	d := Density(c)
	col := [...]ui.Color{k.TextMuted, k.Info, k.Warning, k.Danger}[l.Level]
	ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
		if !l.Time.IsZero() {
			codeMono(c, l.Time.Format("15:04:05.000")).TextColor(k.TextMuted).Shrink(0)
		}
		ui.Text(c, fmt.Sprintf("%-5s", logViewerName(c, l.Level))).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).
			FontWeight(700).TextColor(col).SingleLine().Shrink(0)
		ui.RichText(c, codeMatchSpans(c, l.Text, query)...).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).
			SingleLine().Grow(1).MinWidth(0)
	})
}
