package mujica

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
)

// PriceAlertLine is a caller-held price alert drawn as a dashed line that the pointer drags.
type PriceAlertLine struct {
	ID    string
	Price float64
	// Label shows beside the price; empty shows the price only.
	Label string
}

// CompareSymbol is another symbol's closes, aligned to the chart's candles, shown in percent.
type CompareSymbol struct {
	Name   string
	Closes []float64
}

// compareLines adds the compare symbols' percent lines and sets percent mode.
func (m *mktCore) compareLines() {
	if m.opts.Compare == nil || len(*m.opts.Compare) == 0 || m.to <= m.from {
		return
	}
	m.base = m.candles[m.from].Close
	pal := internal.Palette(m.k)
	for i, s := range *m.opts.Compare {
		if s.Name == "" || len(s.Closes) != len(m.candles) {
			panic(fmt.Sprintf("mujica: CompareSymbol %q needs a Name and one close per candle", s.Name))
		}
		base := math.NaN()
		for j := m.from; j < m.to && math.IsNaN(base); j++ {
			base = s.Closes[j]
		}
		closes := s.Closes
		vals := chartMemo(m.ch, "compare."+s.Name, [2]any{chartID(closes), base}, func() []float64 {
			out := make([]float64, len(closes))
			for j, v := range closes {
				out[j] = (v/base - 1) * 100
			}
			return out
		})
		m.overlays = append(m.overlays, mktLine{name: s.Name, values: vals, color: pal[(i+1)%len(pal)], shape: len(m.overlays)})
	}
}

// CompareLegend lists the compare symbols with their colors and a remove button each; Changed reports a removal.
func CompareLegend(c *ui.Context, symbols *[]CompareSymbol) ChoiceResult {
	if symbols == nil {
		panic("mujica: CompareLegend needs symbols")
	}
	k := Tokens(c)
	pal := internal.Palette(k)
	u := Density(c).Unit()
	removed := -1
	row := ui.Row(c).Wrap().Gap(u * 2).Role(ui.RoleList).Label(msg(c, "compare_symbol.label"))
	row.Children(func() {
		for i, s := range *symbols {
			ui.Row(c).Key(s.Name).AlignItems(ui.Center).Gap(u).Padding(0, 0, 0, u*2).Radius(6).Border(1, k.Border).Role(ui.RoleListItem).Label(s.Name).Children(func() {
				col := pal[(i+1)%len(pal)]
				ui.Box(c).Size(10, 10).Draw(func(p *ui.Painter, r ui.Rect) { internal.Marker(p, i+1, r.X+5, r.Y+5, 9, col) })
				ui.Text(c, s.Name).SingleLine()
				if iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "compare_symbol.remove"), s.Name)).Clicked() && !row.IsDisabled() {
					removed = i
				}
			})
		}
	})
	if removed >= 0 {
		slog.Debug("mujica: compare_symbol removed", "name", (*symbols)[removed].Name)
		*symbols = append((*symbols)[:removed:removed], (*symbols)[removed+1:]...)
	}
	return ChoiceResult{row, removed >= 0}
}

// alertsCheck runs OnAlert for each alert the newest close crossed since the last frame.
func (m *mktCore) alertsCheck() {
	if m.opts.Alerts == nil || len(m.candles) == 0 {
		return
	}
	last := m.candles[len(m.candles)-1]
	if m.st.alertSeen && last.Close != m.st.lastClose && m.opts.OnAlert != nil {
		for _, a := range *m.opts.Alerts {
			if (m.st.lastClose-a.Price)*(last.Close-a.Price) <= 0 && m.st.lastClose != a.Price {
				slog.Debug("mujica: price_alert_line crossed", "id", a.ID)
				m.opts.OnAlert(a, last)
			}
		}
	}
	m.st.lastClose, m.st.alertSeen = last.Close, true
}

// alertsInput drags an alert line that the press started on; it reports taking the drag.
func (m *mktCore) alertsInput(g mktGeom, x, y float32, held bool, dy float32) bool {
	if m.opts.Alerts == nil || m.base != 0 {
		return false
	}
	alerts := *m.opts.Alerts
	if !held {
		m.st.grab = -1
		return false
	}
	if m.st.grab < 0 {
		for i, a := range alerts {
			if math.Abs(float64(g.ys.Map(a.Price)-(y-dy))) <= 5 {
				m.st.grab = i
			}
		}
	}
	if m.st.grab < 0 || m.st.grab >= len(alerts) {
		return false
	}
	p := math.Pow(10, float64(m.dec))
	alerts[m.st.grab].Price = math.Round(g.ys.Invert(y)*p) / p
	slog.Debug("mujica: price_alert_line moved", "id", alerts[m.st.grab].ID)
	return true
}

// paintAlerts draws each alert as a dashed line with its price tag.
func (m *mktCore) paintAlerts(p *ui.Painter, g mktGeom, size float32) {
	if m.opts.Alerts == nil || m.base != 0 {
		return
	}
	for _, a := range *m.opts.Alerts {
		y := internal.Snap(p, g.ys.Map(a.Price))
		for x := g.plot.X; x < g.plot.X+g.plot.W; x += 8 {
			p.Line(x, y, min(x+4, g.plot.X+g.plot.W), y, 1, m.k.Warning)
		}
		text := priceUngrouped(a.Price, m.dec)
		if a.Label != "" {
			text = a.Label + " " + text
		}
		mktTag(p, g.plot.X+g.plot.W, y, fmt.Sprintf(msg(m.c, "price_alert_line.tag"), text), size, m.k.Warning, m.k.Surface)
	}
}
