package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// LiveStatus is the state of real-time sync.
type LiveStatus int

const (
	LiveSynced LiveStatus = iota
	LiveSyncing
	LiveReconnecting
	LiveOffline
	LiveError
)

// LiveIndicatorOptions configure a LiveIndicator.
type LiveIndicatorOptions struct {
	// LastSynced dates the last good sync; zero hides it.
	LastSynced time.Time
	// Now measures LastSynced; zero uses the frame's time.
	Now time.Time
	// Error explains LiveError.
	Error string
}

// LiveIndicator displays sync status, announces changes and offers retries.
func LiveIndicator(c *ui.Context, status LiveStatus, opts LiveIndicatorOptions) ChoiceResult {
	if status < LiveSynced || status > LiveError {
		panic("mujica: LiveStatus out of range")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	text := msg(c, [...]string{"live_indicator.synced", "live_indicator.syncing", "live_indicator.reconnecting", "live_indicator.offline", "live_indicator.error"}[status])
	if status == LiveError && opts.Error != "" {
		text += ": " + opts.Error
	}
	var res ChoiceResult
	res.Element = ui.Row(c).Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize)
	last := ui.Local(res.Element, "live_indicator.last", func() LiveStatus { return status })
	if *last != status {
		*last = status
		c.Announce(text)
		slog.Debug("mujica: live_indicator status", "status", status)
	}
	res.Element.Children(func() {
		box := ui.Row(c).Gap(u*2).AlignItems(ui.Center).Padding(u, u*3).Radius(theme.MinHitSize / 2).Role(ui.RoleStatus).Label(text)
		col := [...]ui.Color{k.Success, k.TextMuted, k.Warning, k.Warning, k.Danger}[status]
		box.Border(theme.BorderWidth, col.Mix(k.Surface, 0.5)).Background(k.Surface)
		box.Children(func() {
			switch status {
			case LiveSyncing, LiveReconnecting:
				Spinner(c, SpinnerOptions{Size: 14, Label: text})
			default:
				name := [...]string{"circle-check", "", "", "wifi-off", "circle-alert"}[status]
				ui.Icon(c, icons.Must(name)).Size(16, 16).Shrink(0).TextColor(col)
			}
			ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
			if !opts.LastSynced.IsZero() && status != LiveSyncing {
				ui.Text(c, liveAgo(c, now.Sub(opts.LastSynced))).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			}
		})
		if status == LiveOffline || status == LiveError {
			if Button(c, msg(c, "live_indicator.retry"), ButtonOptions{Variant: Ghost}).MinHeight(theme.MinHitSize).Clicked() && !res.Element.IsDisabled() {
				res.changed = true
				slog.Debug("mujica: live_indicator retry")
			}
		}
	})
	return res
}

// liveAgo words how long ago the last sync was.
func liveAgo(c *ui.Context, d time.Duration) string {
	switch {
	case d < time.Minute:
		return msg(c, "live_indicator.justNow")
	case d < time.Hour:
		return fmt.Sprintf(msg(c, "live_indicator.minutes"), int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf(msg(c, "live_indicator.hours"), int(d.Hours()))
	}
	return fmt.Sprintf(msg(c, "live_indicator.days"), int(d.Hours()/24))
}
