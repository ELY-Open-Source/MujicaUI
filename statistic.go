package mujica

import (
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// StatisticOptions configure a Statistic. It formats values; it computes none.
type StatisticOptions struct {
	Title string
	// Precision is the number of decimals, 0 or more.
	Precision int
	// Prefix and Suffix wrap the number, as "$" or "×"; Unit follows it, smaller.
	Prefix, Suffix, Unit string
	// Format replaces the built-in grouped decimal formatting.
	Format func(v float64) string
	// Change, when set, shows the caller's change with a sign and an arrow.
	Change     *float64
	ChangeUnit string
	// ChangeNote explains the change, as "vs last week".
	ChangeNote string
}

// Statistic shows a headline number with its title, unit and change.
func Statistic(c *ui.Context, value float64, opts StatisticOptions) *ui.Element {
	if opts.Precision < 0 {
		panic("mujica: Statistic with negative Precision")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return statisticFormat(v, opts.Precision) }
	}
	num := opts.Prefix + format(value) + opts.Suffix
	spoken := strings.TrimSpace(opts.Title + " " + num + " " + opts.Unit)
	e := ui.Column(c).Gap(u).Role(ui.RoleGroup)
	e.Children(func() {
		if opts.Title != "" {
			ui.Text(c, opts.Title).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
		}
		ui.Row(c).Gap(u).AlignItems(ui.End).Children(func() {
			ui.Text(c, num).Font(theme.SerifFont).FontSize(theme.H1Size).FontWeight(600).LineHeight(1.2).SingleLine()
			if opts.Unit != "" {
				ui.Text(c, opts.Unit).TextColor(k.TextMuted).Padding(0, 0, u, 0)
			}
		})
		if opts.Change == nil {
			return
		}
		ch := *opts.Change
		col, icon, word := k.TextMuted, "minus", msg(c, "statistic.flat")
		switch {
		case ch > 0:
			col, icon, word = k.Success, "arrow-up", msg(c, "statistic.up")
		case ch < 0:
			col, icon, word = k.Danger, "arrow-down", msg(c, "statistic.down")
		}
		text := statisticSigned(ch, opts.Precision) + opts.ChangeUnit
		spoken += ", " + word + " " + text
		ui.Row(c).Gap(u).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Row(c).Gap(u / 2).AlignItems(ui.Center).TextColor(col).Children(func() {
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, text).FontWeight(600).SingleLine()
			})
			if opts.ChangeNote != "" {
				ui.Text(c, opts.ChangeNote).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
			}
		})
	})
	return e.Label(spoken)
}

// statisticFormat writes v with p decimals and comma-grouped thousands.
func statisticFormat(v float64, p int) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	s := strconv.FormatFloat(math.Abs(v), 'f', p, 64)
	whole, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	if v < 0 && strings.Trim(s, "0.") != "" {
		b.WriteString("−")
	}
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if frac != "" {
		b.WriteString("." + frac)
	}
	return b.String()
}

// statisticSigned writes a change with an explicit sign.
func statisticSigned(v float64, p int) string {
	s := statisticFormat(math.Abs(v), p)
	switch {
	case v > 0:
		return "+" + s
	case v < 0:
		return "−" + s
	}
	return s
}
