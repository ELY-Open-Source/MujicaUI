package mujica

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"log/slog"

	"mujicaui/internal"
	"mujicaui/theme"
)

// IssueCardOptions configure an IssueCard.
type IssueCardOptions struct {
	Statuses []WorkflowStatus
	Members  []ProjectMember
	Labels   []TaskLabel
	Selected bool
	Disabled bool
}

// IssueCard displays task fields in a selectable card.
func IssueCard(c *ui.Context, task Task, opts IssueCardOptions) *ui.Element {
	projectCheckPriority(task.Priority)
	st := projectStatus(opts.Statuses, task.Status)
	k := Tokens(c)
	d := Density(c)
	small := FontSize(c, theme.CaptionSize)
	b := ui.ButtonBase(c).Column().AlignItems(ui.Stretch).Justify(ui.Start).Padding(d.Unit()*3).Gap(d.Unit()*2).
		Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).Background(k.Surface).
		Label(task.ID + " " + task.Title).Description(fmt.Sprintf(msg(c, "issue_card.description"), st.Name, PriorityName(c, task.Priority))).Disabled(opts.Disabled)
	internal.Ring(b, k.Focus, theme.CardRadius)
	switch {
	case opts.Selected:
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover).Border(theme.BorderWidth, k.ControlBorder)
	}
	selected := opts.Selected
	b.Transition(colorFade).Draw(func(p *ui.Painter, r ui.Rect) {
		if selected {
			datetimeSelectedMark(p, r, k.AccentText)
		}
	})
	b.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			statusGlyphBox(c, st.Group)
			ui.Text(c, task.ID).Font(theme.MonoFont).FontSize(small).TextColor(k.TextMuted).Grow(1)
			PriorityIndicator(c, task.Priority, PriorityIndicatorOptions{IconOnly: true})
		})
		t := ui.Text(c, task.Title).TextColor(k.Text).MaxLines(2).Ellipsis("…")
		if st.Group == StatusCanceled {
			t.Strikethrough().TextColor(k.TextMuted)
		}
		if len(task.Labels) == 0 && len(task.Assignees) == 0 {
			return
		}
		ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Wrap().Children(func() {
			for _, id := range task.Labels {
				issueLabelChip(c, projectLabel(opts.Labels, id))
			}
			ui.Spacer(c)
			if len(task.Assignees) > 0 {
				items := make([]AvatarItem, len(task.Assignees))
				for i, id := range task.Assignees {
					items[i] = AvatarItem{Name: projectMember(opts.Members, id).Name}
				}
				AvatarGroup(c, items, AvatarGroupOptions{Max: 3, Size: 22})
			}
		})
	})
	if b.Clicked() && !b.IsDisabled() {
		slog.Debug("mujica: issue_card open", "id", task.ID)
	}
	return b
}

// issueLabelChip displays a label with its color dot.
func issueLabelChip(c *ui.Context, l TaskLabel) {
	k := Tokens(c)
	d := Density(c)
	col := eventColor(k, l.Color)
	ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Padding(1, d.Unit()*1.5).Radius(theme.ControlRadius).
		Border(theme.BorderWidth, k.Border).Children(func() {
		ui.Box(c).Size(8, 8).Radius(4).Background(col)
		ui.Text(c, l.Name).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Text).SingleLine()
	})
}
