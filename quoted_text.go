package mujica

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// QuotedText folds or expands one level of quoted text.
func QuotedText(c *ui.Context, text string, open *bool) ChoiceResult {
	body := QuoteStrip(text)
	n := strings.Count(body, "\n") + 1
	k := Tokens(c)
	u := Density(c).Unit()
	title := fmt.Sprintf(msg(c, "quoted_text.show"), n)
	if n == 1 {
		title = msg(c, "quoted_text.showOne")
	}
	if *open {
		title = msg(c, "quoted_text.hide")
	}
	return Collapsible(c, open, CollapsibleOptions{Title: title}, func() {
		ui.Text(c, body).Selectable().FocusRing(false).TextColor(k.TextMuted).LineHeight(theme.LineHeight).
			Padding(0, 0, 0, u*3).BorderWidth(0, 0, 0, 2).BorderColor(k.Border)
	})
}

// QuoteStrip removes one quote prefix from each line.
func QuoteStrip(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i, l := range lines {
		l = strings.TrimPrefix(l, ">")
		lines[i] = strings.TrimPrefix(l, " ")
	}
	return strings.Join(lines, "\n")
}
