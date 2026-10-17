package devtools

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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
	Element ui.Element
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
	k := core.Tokens(c)
	var r AlertListResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled)
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		if len(alerts) == 0 {
			agent.AgentCaption(c, core.Msg(c, "alert_list.empty"))
			return
		}
		agent.AgentRows(c, len(alerts), func(i int) {
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
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	sev := [...]core.Severity{core.SeverityInfo, core.SeverityWarning, core.SeverityError}[a.Level]
	word := core.Msg(c, [...]string{"alert_list.info", "alert_list.warning", "alert_list.critical"}[a.Level])
	silenced := a.Silence.After(c.Now())
	row := ui.Row(c.Key(a.ID)).Gap(u*3).AlignItems(ui.Center).MinHeight(d.ControlHeight()+u*2).Padding(u*2, u*3).
		Label(word + ", " + a.Title)
	if tone, _ := core.SeverityLook(k, sev); a.Level == AlertCritical && !a.Acked {
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
			agent.AgentCaption(c, line)
		})
		switch {
		case silenced:
			agent.AgentCaption(c, fmt.Sprintf(core.Msg(c, "alert_list.silenced"), a.Silence.Format("15:04")))
		case a.Acked:
			agent.AgentCaption(c, core.Msg(c, "alert_list.acked"))
		default:
			ack = input.Button(c, core.Msg(c, "alert_list.ack"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked()
		}
		if !silenced {
			silence = agent.AgentToggle(c, core.Msg(c, "alert_list.silence"), "bell", locked)
		}
	})
	return ack, silence
}
