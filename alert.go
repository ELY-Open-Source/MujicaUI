package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// Severity is the meaning of a message. The zero value is SeverityInfo.
type Severity int

const (
	SeverityInfo Severity = iota
	SeveritySuccess
	SeverityWarning
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeveritySuccess:
		return "success"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	}
	panic(fmt.Sprintf("mujica: unknown Severity %d", int(s)))
}

// severityLook returns the color and icon that, with text, mark s.
func severityLook(k theme.Tokens, s Severity) (ui.Color, *ui.SVG) {
	switch s {
	case SeveritySuccess:
		return k.Success, icons.Must("circle-check")
	case SeverityWarning:
		return k.Warning, icons.Must("triangle-alert")
	case SeverityError:
		return k.Danger, icons.Must("circle-x")
	case SeverityInfo:
		return k.Info, icons.Must("info")
	}
	panic(fmt.Sprintf("mujica: unknown Severity %d", int(s)))
}

// AlertOptions configure an Alert.
type AlertOptions struct {
	Severity Severity
	// Title goes above the message in bold; empty shows none.
	Title string
	// Action labels a button at the end; empty shows none.
	Action string
}

// AlertResult is what an Alert built and reported this frame.
type AlertResult struct {
	// Element is nil while the alert is closed.
	Element *ui.Element
	// Action reports a click on the action button.
	Action bool
	// Closed reports the user dismissing the alert.
	Closed bool
}

// Alert shows an inline message; a non-nil open adds a dismiss button that clears it.
func Alert(c *ui.Context, open *bool, message string, opts AlertOptions) AlertResult {
	if message == "" {
		panic("mujica: Alert needs a message")
	}
	if open != nil && !*open {
		return AlertResult{}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	tone, icon := severityLook(k, opts.Severity)
	name := msg(c, "alert."+opts.Severity.String())
	var res AlertResult
	box := ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(u*3).Padding(u*3, u*3, u*3, u*4).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Role(ui.RoleGroup).Label(name).Clip()
	// A bar in the severity color along the leading edge.
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 3, H: r.H}, tone, 0)
	})
	// Warnings and errors are read out when they appear or change.
	said := ui.Local(box, "alert.said", func() string { return "" })
	if news := name + ": " + message; opts.Severity >= SeverityWarning && *said != news {
		*said = news
		c.Announce(news)
	}
	box.Children(func() {
		ui.Icon(c, icon).Size(20, 20).TextColor(tone).Margin(1, 0, 0, 0)
		ui.Column(c).Grow(1).Shrink(1).Gap(u).Children(func() {
			if opts.Title != "" {
				ui.Text(c, opts.Title).Bold().TextColor(k.Text)
			}
			ui.Text(c, message).TextColor(k.Text).LineHeight(theme.LineHeight)
			if opts.Action != "" {
				ui.Row(c).Margin(u, 0, 0, 0).Children(func() {
					if Button(c, opts.Action, ButtonOptions{Variant: Outline}).Clicked() {
						slog.Debug("mujica: alert action", "severity", opts.Severity, "action", opts.Action)
						res.Action = true
					}
				})
			}
		})
		if open != nil {
			dismiss := Button(c, "", ButtonOptions{Variant: Ghost, Icon: icons.Must("x"), Label: msg(c, "alert.close")})
			dismiss.TextColor(k.TextMuted).Padding(0).MinHeight(theme.MinHitSize).Margin(-u, -u, -u, 0)
			if dismiss.Clicked() {
				slog.Debug("mujica: alert closed", "severity", opts.Severity)
				*open, res.Closed = false, true
			}
		}
	})
	res.Element = box
	return res
}
