package mujica

import (
	"cmp"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// OptionGreeks are an option's sensitivities.
type OptionGreeks struct {
	Delta, Gamma, Theta, Vega, Rho float64
}

// GreeksRow is a labeled contract's greeks; NaN shows a dash.
type GreeksRow struct {
	Label  string
	Greeks OptionGreeks
}

// GreeksTableOptions configure a GreeksTable.
type GreeksTableOptions struct {
	// Decimals per value; zero means 4.
	Decimals int
	// Notes shows what each greek means under the table.
	Notes bool
}

var greekNames = [5]string{"Delta", "Gamma", "Theta", "Vega", "Rho"}

// GreeksTable lists Delta, Gamma, Theta, Vega and Rho per contract; headers explain each.
func GreeksTable(c *ui.Context, rows []GreeksRow, opts GreeksTableOptions) *ui.Element {
	if opts.Decimals < 0 {
		panic("mujica: GreeksTable Decimals is negative")
	}
	dec := cmp.Or(opts.Decimals, 4)
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	col := ui.Column(c).Gap(u * 2)
	col.Children(func() {
		g := ui.Grid(c).ColumnTracks(ui.Fr(1), ui.FitContent(), ui.FitContent(), ui.FitContent(), ui.FitContent(), ui.FitContent()).
			GapX(u*5).Padding(u*2, u*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
			Role(ui.RoleTable).Label(msg(c, "greeks_table.title"))
		g.Children(func() {
			ui.Text(c, msg(c, "greeks_table.contract")).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).
				MinHeight(d.ControlHeight()).Role(ui.RoleColumnHeader)
			for i, name := range greekNames {
				h := ui.Row(c).Key(name).Gap(2).AlignItems(ui.Center).Justify(ui.End).MinHeight(d.ControlHeight()).
					Role(ui.RoleColumnHeader).Label(name).Description(msg(c, "greeks_table.note_"+greekKey(i)))
				h.Tooltip(msg(c, "greeks_table.note_"+greekKey(i)))
				h.Children(func() {
					ui.Text(c, name).Font(theme.SerifFont).TextColor(k.Text)
					ui.Icon(c, icons.Must("info")).Size(12, 12).TextColor(k.TextMuted)
				})
			}
			for ri, r := range rows {
				vals := [5]float64{r.Greeks.Delta, r.Greeks.Gamma, r.Greeks.Theta, r.Greeks.Vega, r.Greeks.Rho}
				ui.Text(c, r.Label).Key(ri).SingleLine().MinHeight(d.ControlHeight()).Role(ui.RoleCell)
				for i, v := range vals {
					s := "—"
					if !math.IsNaN(v) {
						s = statisticSigned(v, dec)
					}
					ui.Text(c, s).Key([2]int{ri, i}).JustifySelf(ui.End).FontFeatures("tnum").TextColor(k.Text).
						Role(ui.RoleCell).Label(greekNames[i] + " " + s)
				}
			}
		})
		if opts.Notes {
			ui.Column(c).Gap(u).Children(func() {
				for i, name := range greekNames {
					ui.RichText(c, ui.Span{Text: name + "  ", Weight: 700}, ui.Span{Text: msg(c, "greeks_table.note_"+greekKey(i))}).
						FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				}
			})
		}
	})
	return col
}

func greekKey(i int) string { return [5]string{"delta", "gamma", "theta", "vega", "rho"}[i] }
