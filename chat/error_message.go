package chat

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ErrorMessageOptions configure an ErrorMessage.
type ErrorMessageOptions struct {
	// Title sums the failure up; empty uses the built-in "Something went wrong".
	Title string
	// RetryAt is when a rate limit lifts; until then retry waits with a countdown.
	RetryAt time.Time
}

// ErrorMessageResult is the error's element and a retry the user asked for.
type ErrorMessageResult struct {
	Element *ui.Element
	retried bool
}

// Retried reports a press of retry.
func (r ErrorMessageResult) Retried() bool { return r.retried }

// ErrorMessage shows a failed reply with the original error text, selectable, and retry.
func ErrorMessage(c *ui.Context, text string, opts ErrorMessageOptions) ErrorMessageResult {
	if text == "" {
		panic("mujica: ErrorMessage needs the error text")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	title := opts.Title
	if title == "" {
		title = core.Msg(c, "error_message.title")
	}
	wait := opts.RetryAt.Sub(c.Now())
	var r ErrorMessageResult
	r.Element = ui.Column(c).Gap(u*2).Padding(u*3, u*4).MaxWidth(560).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Danger).Role(ui.RoleGroup).Label(title)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icons.Must("circle-alert")).Size(16, 16).Shrink(0).TextColor(k.Danger)
			ui.Text(c, title).FontWeight(600).TextColor(k.Danger)
		})
		ui.Text(c, text).Selectable().Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text)
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			label := core.Msg(c, "error_message.retry")
			if wait > 0 {
				secs := int(math.Ceil(wait.Seconds()))
				label = fmt.Sprintf(core.Msg(c, "error_message.wait"), secs)
				// Wake when the shown second changes.
				c.After(wait - time.Duration(secs-1)*time.Second)
			}
			if input.Button(c, label, input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("refresh-cw"), Disabled: wait > 0}).Clicked() && !disabled {
				r.retried = true
				slog.Debug("mujica: error_message retry")
			}
			if wait > 0 {
				ui.Text(c, core.Msg(c, "error_message.limited")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			}
		})
	})
	return r
}
