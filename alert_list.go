package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// AlertLevel is an alert's level. The zero value is AlertInfo.
type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarning
	AlertCritical
)

// AlertItem is one firing alert; the caller owns it.
type AlertItem struct {
	ID      string
	Title   string
	Level   AlertLevel
	Source  string
	At      time.Time
	Acked   bool
	Silence time.Time // silenced until; zero is not silenced
}

// AlertListOptions configure an AlertList.
type AlertListOptions struct {
	// MaxHeight caps the scrolling list; zero means 360.
	MaxHeight float32
	Disabled  bool
}

// AlertListResult reports this frame's request.
type AlertListResult struct {
	Element *ui.Element
	// Ack and Silence report requested alert IDs; empty means idle.
	Ack, Silence string
}

// AlertList shows alerts with acknowledge and silence actions.
func AlertList(c *ui.Context, alerts []AlertItem, opts AlertListOptions) AlertListResult {
	if opts.MaxHeight < 0 {
		panic("mujica: AlertList MaxHeight is negative")
	}
	seen := map[string]bool{}
	for i, a := range alerts {
		if a.ID == "" || a.Title == "" || seen[a.ID] || a.Level < AlertInfo || a.Level > AlertCritical {
			panic(fmt.Sprintf("mujica: AlertList alert %d needs a unique ID, a Title and a known Level", i))
		}
		seen[a.ID] = true
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	k := Tokens(c)
	var r AlertListResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled)
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		if len(alerts) == 0 {
			agentCaption(c, msg(c, "alert_list.empty"))
			return
		}
		agentRows(c, len(alerts), func(i int) {
			ack, silence := alertRow(c, alerts[i], locked)
			if ack && !locked {
				r.Ack = alerts[i].ID
			}
			if silence && !locked {
				r.Silence = alerts[i].ID
			}
		}).Shrink(1).MaxHeight(maxH).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Dividers(theme.BorderWidth, k.Border)
	})
	if r.Ack != "" || r.Silence != "" {
		slog.Debug("mujica: alert_list action", "ack", r.Ack, "silence", r.Silence)
	}
	return r
}

// alertRow shows an alert with acknowledge and silence buttons.
func alertRow(c *ui.Context, a AlertItem, locked bool) (ack, silence bool) {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	sev := [...]Severity{SeverityInfo, SeverityWarning, SeverityError}[a.Level]
	word := msg(c, [...]string{"alert_list.info", "alert_list.warning", "alert_list.critical"}[a.Level])
	silenced := a.Silence.After(c.Now())
	row := ui.Row(c).Key(a.ID).Gap(u*3).AlignItems(ui.Center).MinHeight(d.ControlHeight()+u*2).Padding(u*2, u*3).
		Label(word + ", " + a.Title)
	if tone, _ := severityLook(k, sev); a.Level == AlertCritical && !a.Acked {
		row.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 3, H: r.H}, tone, 0) })
	}
	row.Children(func() {
		ui.Box(c).Width(96).Shrink(0).Children(func() { devtoolsSeverity(c, sev, word) })
		ui.Column(c).Grow(1).Shrink(1).Gap(u / 2).Children(func() {
			t := ui.Text(c, a.Title).FontWeight(600)
			if a.Acked || silenced {
				t.TextColor(k.TextMuted)
			}
			line := a.Source
			if !a.At.IsZero() {
				line += " · " + devtoolsAgo(c, a.At)
			}
			agentCaption(c, line)
		})
		switch {
		case silenced:
			agentCaption(c, fmt.Sprintf(msg(c, "alert_list.silenced"), a.Silence.Format("15:04")))
		case a.Acked:
			agentCaption(c, msg(c, "alert_list.acked"))
		default:
			ack = Button(c, msg(c, "alert_list.ack"), ButtonOptions{Variant: Secondary, Disabled: locked}).Clicked()
		}
		if !silenced {
			silence = agentToggle(c, msg(c, "alert_list.silence"), "bell", locked)
		}
	})
	return ack, silence
}
