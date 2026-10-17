package finance

import (
	"cmp"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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

func GreeksTable(c *ui.Context, rows []GreeksRow, opts GreeksTableOptions) ui.Element {
	if opts.Decimals < 0 || opts.Decimals > 12 {
		panic("mujica: GreeksTable Decimals is negative")
	}
	dec := cmp.Or(opts.Decimals, 4)
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	col := ui.Column(c).Gap(u * 2)
	col.Children(func() {
		g := ui.Grid(c).ColumnTracks(ui.Fr(2), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).
			GapX(u*2).Padding(u*2, u*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
			Role(ui.RoleTable).Label(core.Msg(c, "greeks_table.title"))
		g.Children(func() {
			ui.Text(c, core.Msg(c, "greeks_table.contract")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).
				MinHeight(d.ControlHeight()).Role(ui.RoleColumnHeader)
			for i, name := range greekNames {
				h := ui.Row(c.Key(name)).Gap(2).AlignItems(ui.Center).Justify(ui.End).MinHeight(d.ControlHeight()).
					Role(ui.RoleColumnHeader).Label(name).Description(core.Msg(c, "greeks_table.note_"+greekKey(i)))
				h.Tooltip(core.Msg(c, "greeks_table.note_"+greekKey(i)))
				h.Children(func() {
					ui.Text(c, name).SingleLine().Font(theme.SerifFont).TextColor(k.Text)
					ui.Icon(c, icons.Must("info")).Size(12, 12).TextColor(k.TextMuted)
				})
			}

		})
		state := ui.Local(col, "mujica.greeks.rows", func() ui.ListState { return ui.ListState{} })
		ui.List(c, state, len(rows), func(ri int) {
			ui.Grid(c).ColumnTracks(ui.Fr(2), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).GapX(u*2).Padding(0, u*3).MinHeight(d.ControlHeight()).Children(func() {
				r := rows[ri]
				vals := [5]float64{r.Greeks.Delta, r.Greeks.Gamma, r.Greeks.Theta, r.Greeks.Vega, r.Greeks.Rho}
				ui.Text(c.Key(ri), r.Label).SingleLine().MinHeight(d.ControlHeight()).Role(ui.RoleCell)
				for i, v := range vals {
					s := "—"
					if !math.IsNaN(v) {
						s = data.StatisticSigned(v, dec)
					}
					ui.Text(c.Key([2]int{ri, i}), s).JustifySelf(ui.End).FontFeatures("tnum").TextColor(k.Text).
						Role(ui.RoleCell).Label(greekNames[i] + " " + s)
				}
			})
		}).Height(min(float32(len(rows))*d.ControlHeight(), 240)).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius)

		if opts.Notes {
			ui.Column(c).Gap(u).Children(func() {
				for i, name := range greekNames {
					ui.RichText(c, ui.Span{Text: name + "  ", Weight: 700}, ui.Span{Text: core.Msg(c, "greeks_table.note_"+greekKey(i))}).
						FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				}
			})
		}
	})
	return col
}

func greekKey(i int) string { return [5]string{"delta", "gamma", "theta", "vega", "rho"}[i] }
