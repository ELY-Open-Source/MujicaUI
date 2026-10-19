package messaging

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// QuotedText folds or expands one level of quoted text.
func QuotedText(c *ui.Context, text string, open *bool) core.ChoiceResult {
	body := QuoteStrip(text)
	n := strings.Count(body, "\n") + 1
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	title := fmt.Sprintf(core.Msg(c, "quoted_text.show"), n)
	if n == 1 {
		title = core.Msg(c, "quoted_text.showOne")
	}
	if *open {
		title = core.Msg(c, "quoted_text.hide")
	}
	return data.Collapsible(c, open, data.CollapsibleOptions{Title: title}, func() {
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
