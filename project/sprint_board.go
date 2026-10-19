package project

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Sprint is one iteration: its name, goal and inclusive dates.
type Sprint struct {
	Name       string
	Goal       string
	Start, End core.Date
}

// SprintBoardOptions configure a SprintBoard.
type SprintBoardOptions struct {
	Board KanbanBoardOptions
	// Today counts the days left; zero uses the system clock.
	Today core.Date
}

// SprintBoard combines iteration information, completion and a task board.
func SprintBoard(c *ui.Context, sprint Sprint, tasks *[]Task, opts SprintBoardOptions) KanbanResult {
	core.DatetimeCheckDate("SprintBoard Start", sprint.Start)
	core.DatetimeCheckDate("SprintBoard End", sprint.End)
	if sprint.Start.IsZero() || sprint.End.Compare(sprint.Start) < 0 {
		panic(fmt.Sprintf("mujica: SprintBoard dates %v..%v are not a span", sprint.Start, sprint.End))
	}
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	var r KanbanResult
	col := ui.Column(c).Gap(d.Unit() * 4).Grow(1).MinHeight(0).Role(ui.RoleGroup).Label(sprint.Name)
	cache := ui.Local(col, "sprint-counts", func() sprintCache { return sprintCache{} })
	var first *Task
	if len(*tasks) > 0 {
		first = &(*tasks)[0]
	}
	if !cache.ready || cache.first != first || cache.n != len(*tasks) || cache.version != opts.Board.Version || !slices.Equal(cache.statuses, opts.Board.Statuses) {
		*cache = sprintCache{first: first, n: len(*tasks), version: opts.Board.Version, statuses: slices.Clone(opts.Board.Statuses), ready: true}
		for _, t := range *tasks {
			if g := projectStatus(opts.Board.Statuses, t.Status).Group; g != StatusCanceled {
				cache.total++
				if g == StatusCompleted {
					cache.done++
				}
			}
		}
	}
	done, total := cache.done, cache.total
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 6).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Column(c).Gap(d.Unit()).Grow(1).Children(func() {
				ui.Text(c, sprint.Name).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).TextColor(k.Text)
				if sprint.Goal != "" {
					ui.Text(c, sprint.Goal).TextColor(k.TextMuted)
				}
				ui.Text(c, datetime.WeekViewTitle(c, sprint.Start, sprint.End)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			})
			sprintStat(c, sprintDaysLeft(c, sprint, today), core.Msg(c, "sprint_board.time"))
			frac := 0.0
			if total > 0 {
				frac = float64(done) / float64(total)
			}
			ui.Column(c).Width(d.ControlHeight() * 5).Gap(d.Unit()).Children(func() {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "sprint_board.done"), done, total)).TextColor(k.Text).FontFeatures("tnum")
				feedback.Progress(c, frac, feedback.ProgressOptions{Label: core.Msg(c, "sprint_board.progress")})
			})
		})
		ui.Box(c).Height(12).FillWidth().Draw(func(p *ui.Painter, rr ui.Rect) { internal.RuleDiamond(p, rr, 8, 6, k.Ornament) })
		r = KanbanBoard(c, tasks, opts.Board)
	})
	if _, moved := r.Moved(); moved {
		cache.ready = false
		c.Invalidate()
	}
	r.Element = col
	return r
}

// sprintDaysLeft formats remaining time relative to the iteration.
func sprintDaysLeft(c *ui.Context, s Sprint, today core.Date) string {
	switch {
	case today.Compare(s.Start) < 0:
		return core.CalendarCount(c, "sprint_board.starts_in", int(s.Start.Civil().Sub(today.Civil()).Hours()/24))
	case today.Compare(s.End) > 0:
		return core.Msg(c, "sprint_board.ended")
	}
	return core.CalendarCount(c, "sprint_board.left", int(s.End.Civil().Sub(today.Civil()).Hours()/24)+1)
}

// sprintStat is a figure above its caption.
func sprintStat(c *ui.Context, value, caption string) {
	k := core.Tokens(c)
	ui.Column(c).Gap(2).Role(ui.RoleStatus).Label(caption + ": " + value).Children(func() {
		ui.Text(c, value).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.AccentText)
		ui.Text(c, caption).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
}

type sprintCache struct {
	first                   *Task
	n, version, done, total int
	statuses                []WorkflowStatus
	ready                   bool
}
