package navigation

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Step is one stage of Steps.
type Step struct {
	Title       string
	Description string
	// Error marks the step as failed, whatever its place.
	Error bool
	// Clickable lets the user go to the step.
	Clickable bool
}

// StepsOptions configure Steps.
type StepsOptions struct {
	Vertical bool
	// Label names the steps for assistive technology.
	Label string
}

type stepsState int

const (
	stepsWaiting stepsState = iota
	stepsCurrent
	stepsDone
	stepsFailed
)

// Steps shows a process whose step *current is under way: earlier steps are done.
func Steps(c *ui.Context, current *int, steps []Step, opts StepsOptions) core.ChoiceResult {
	if len(steps) == 0 {
		panic("mujica: Steps needs at least one step")
	}
	if *current < 0 || *current > len(steps) {
		panic("mujica: Steps current out of range")
	}
	u := core.Density(c).Unit()
	root := ui.Row(c).AlignItems(ui.Start).Gap(u * 2).FocusGroup(ui.Horizontal | ui.Vertical).Role(ui.RoleList)
	if opts.Vertical {
		root.Column().AlignItems(ui.Stretch).Gap(0)
	}
	if opts.Label != "" {
		root.Label(opts.Label)
	}
	res := core.ChoiceResult{Element: root}
	root.Children(func() {
		for i, s := range steps {
			if s.Title == "" {
				panic("mujica: Step needs a Title")
			}
			st := stepsWaiting
			switch {
			case s.Error:
				st = stepsFailed
			case i < *current:
				st = stepsDone
			case i == *current:
				st = stepsCurrent
			}
			if !opts.Vertical && i > 0 {
				stepsLine(c, steps, i-1, *current)
			}
			if stepsItem(c, i, len(steps), s, st, opts.Vertical, i < len(steps)-1 && opts.Vertical, i < *current) && i != *current {
				slog.Debug("mujica: steps choose", "from", *current, "to", i)
				*current, res = i, core.NewChoiceResult(res.Element, true)
			}
		}
	})
	return res
}

// stepsItem builds one step and reports a click on it.
func stepsItem(c *ui.Context, i, n int, s Step, st stepsState, vertical, tail, doneLine bool) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var e *ui.Element
	if s.Clickable {
		e = ui.ButtonBase(c).Cursor(ui.CursorPointer)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		if e.Hovered() {
			e.Background(k.SurfaceHover)
		}
		internal.Ease(e, core.Motion(c, theme.StateDuration), k.Text)
	} else {
		e = ui.Row(c).Shrink(0).Role(ui.RoleListItem)
	}
	e.AlignItems(ui.Stretch).Justify(ui.Start).Gap(u * 3).Padding(u).Radius(theme.ControlRadius).
		Label(fmt.Sprintf(core.Msg(c, "steps.item"), i+1, n, s.Title, core.Msg(c, stepsStateKey[st])))
	e.Children(func() {
		ui.Column(c).AlignItems(ui.Center).Gap(u).Shrink(0).Children(func() {
			stepsMarker(c, i, st)
			if tail {
				ui.Box(c).Width(2).Grow(1).MinHeight(u * 4).Background(stepsLineColor(k, doneLine)).Radius(1)
			}
		})
		texts := ui.Column(c).Gap(u*0.5).MinWidth(0).Padding(2, 0, 0, 0)
		if vertical {
			texts.Padding(2, 0, u*4, 0)
		}
		texts.Children(func() {
			title := ui.Text(c, s.Title).SingleLine()
			switch st {
			case stepsCurrent:
				title.FontWeight(600).TextColor(k.Text)
			case stepsFailed:
				title.FontWeight(600).TextColor(k.Danger)
			case stepsDone:
				title.TextColor(k.Text)
			default:
				title.TextColor(k.TextMuted)
			}
			if s.Description != "" {
				ui.Text(c, s.Description).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).MaxWidth(240)
			}
		})
	})
	return s.Clickable && e.Clicked()
}

var stepsStateKey = map[stepsState]string{
	stepsWaiting: "steps.waiting", stepsCurrent: "steps.current", stepsDone: "steps.done", stepsFailed: "steps.error",
}

func stepsMarker(c *ui.Context, i int, st stepsState) {
	k := core.Tokens(c)
	m := ui.Box(c).Size(24, 24).Radius(12).Center().Shrink(0).Role(ui.RoleNone)
	var icon string
	switch st {
	case stepsDone:
		m.Background(k.Accent).TextColor(k.OnAccent)
		icon = "check"
	case stepsCurrent:
		m.Background(k.Surface).Border(2, k.AccentText).TextColor(k.AccentText).FontWeight(700)
	case stepsFailed:
		m.Background(k.Surface).Border(2, k.Danger).TextColor(k.Danger)
		icon = "x"
	default:
		m.Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).TextColor(k.TextMuted)
	}
	m.Children(func() {
		if icon != "" {
			ui.Icon(c, icons.Must(icon)).Size(14, 14)
		} else {
			ui.Text(c, strconv.Itoa(i+1)).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
		}
	})
}

// stepsLine joins step i to the next in the horizontal layout.
func stepsLine(c *ui.Context, steps []Step, i, current int) {
	k := core.Tokens(c)
	ui.Box(c).Height(2).Grow(1).MinWidth(24).Margin(core.Density(c).Unit()+11, 0, 0, 0).Radius(1).
		Background(stepsLineColor(k, i < current && !steps[i].Error)).Role(ui.RoleNone)
}

func stepsLineColor(k theme.Tokens, done bool) ui.Color {
	if done {
		return k.AccentText
	}
	return k.Border
}
