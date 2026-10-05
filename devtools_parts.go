package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// devtoolsAgo formats elapsed time using localized copy.
func devtoolsAgo(c *ui.Context, t time.Time) string {
	d := max(c.Now().Sub(t), 0).Truncate(time.Second)
	return fmt.Sprintf(msg(c, "devtools.ago"), agentDuration(d))
}

// devtoolsSeverity shows severity with an icon and word.
func devtoolsSeverity(c *ui.Context, s Severity, word string) {
	tone, icon := severityLook(Tokens(c), s)
	ui.Row(c).Gap(Density(c).Unit()).AlignItems(ui.Center).Shrink(0).Children(func() {
		ui.Icon(c, icon).Size(14, 14).TextColor(tone)
		ui.Text(c, word).TextColor(tone).FontSize(FontSize(c, theme.CaptionSize+1)).FontWeight(600).SingleLine()
	})
}

// devtoolsTitle is a serif panel title.
func devtoolsTitle(c *ui.Context, s string) *ui.Element {
	return ui.Text(c, s).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600)
}
