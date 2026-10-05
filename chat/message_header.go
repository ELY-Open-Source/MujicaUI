package chat

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MessageStatus is where a chat message is in its delivery.
type MessageStatus int

const (
	// MessageNoStatus shows no status.
	MessageNoStatus MessageStatus = iota
	MessageSending
	MessageStreaming
	MessageSent
	MessageFailed
)

// MessageHeaderOptions configure a MessageHeader; zero fields are hidden.
type MessageHeaderOptions struct {
	Model  string
	Time   time.Time
	Status MessageStatus
}

// MessageHeader shows a message's model, time and status as icon and text.
func MessageHeader(c *ui.Context, opts MessageHeaderOptions) *ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	row := ui.Row(c).Gap(u * 2).AlignItems(ui.Center).FontSize(small).TextColor(k.TextMuted)
	row.Children(func() {
		if opts.Model != "" {
			ui.Text(c, opts.Model).FontWeight(600).TextColor(k.Text).SingleLine()
		}
		if !opts.Time.IsZero() {
			ui.Text(c, opts.Time.Format(core.Msg(c, "message_header.time"))).FontFeatures("tnum").SingleLine()
		}
		if opts.Status == MessageNoStatus {
			return
		}
		key, icon, color := "", "", k.TextMuted
		switch opts.Status {
		case MessageSending:
			key, icon = "message_header.sending", "clock"
		case MessageStreaming:
			key = "message_header.streaming"
		case MessageSent:
			key, icon = "message_header.sent", "check"
		case MessageFailed:
			key, icon, color = "message_header.failed", "circle-alert", k.Danger
		default:
			panic("mujica: unknown MessageStatus")
		}
		ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(color).Children(func() {
			if icon == "" {
				core.Spinner(c, core.SpinnerOptions{Size: 12, Label: core.Msg(c, key)})
			} else {
				ui.Icon(c, icons.Must(icon)).Size(12, 12)
			}
			ui.Text(c, core.Msg(c, key)).SingleLine()
		})
	})
	return row
}
