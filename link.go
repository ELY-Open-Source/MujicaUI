package mujica

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// LinkOptions configure a Link. An empty URL makes an in-app action link.
type LinkOptions struct {
	// URL opens in the browser, or goes there in a MyGo router for a path.
	URL string
	// Error is the caller's open failure, shown below the link.
	Error string
}

// Link is underlined text activated by a click or Enter. Clicked reports
// every activation. A malformed URL is not opened and shows an error.
func Link(c *ui.Context, label string, opts LinkOptions) *ui.Element {
	if label == "" {
		panic("mujica: Link needs a label")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	external := false
	problem := opts.Error
	target := opts.URL
	if opts.URL != "" {
		ok := false
		external, ok = linkCheck(opts.URL)
		if !ok {
			external, target = false, ""
			if problem == "" {
				problem = fmt.Sprintf(msg(c, "link.invalid"), opts.URL)
			}
		}
	}
	var l *ui.Element
	ui.Column(c).Gap(u).AlignItems(ui.Start).Children(func() {
		ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
			l = ui.Link(c, label, target).Underline().DecorationThickness(1)
			internal.Ring(l, k.Focus, theme.ControlRadius)
			color := k.AccentText
			if l.Pressed() || l.Hovered() {
				color = k.Text
			}
			l.TextColor(color).DecorationColor(color)
			if external {
				Icon(c, "external-link", IconOptions{Size: 12, Tone: IconAccent})
			}
		})
		if problem != "" {
			l.Error(problem)
			ui.Text(c, problem).TextColor(k.Danger).FontSize(theme.CaptionSize)
		}
	})
	if l.Clicked() {
		slog.Debug("mujica: link activate", "url", opts.URL, "opened", target != "")
	}
	return l
}

// linkCheck reports whether raw has a scheme, and whether it can be opened.
func linkCheck(raw string) (external, ok bool) {
	p, err := url.Parse(raw)
	if err != nil {
		return false, false
	}
	switch p.Scheme {
	case "":
		return false, p.Path != "" || p.RawQuery != ""
	case "http", "https":
		return true, p.Host != ""
	}
	return true, p.Opaque != "" || p.Host != "" || p.Path != ""
}
