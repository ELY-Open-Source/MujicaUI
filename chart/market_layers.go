package chart

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
)

// PriceAlertLine holds a draggable, caller-owned price threshold.
type PriceAlertLine struct {
	ID    string
	Price float64
	// Label shows beside the price; empty shows the price only.
	Label string
}

// CompareSymbol holds aligned closes for percent comparison.
type CompareSymbol struct {
	Name   string
	Closes []float64
}

// compareLines normalizes prices and compared symbols into percentages.
func (m *mktCore) compareLines() {
	if m.opts.Compare == nil || len(*m.opts.Compare) == 0 || m.to <= m.from {
		return
	}
	m.base = m.candles[m.from].Close
	if m.base == 0 {
		panic("mujica: percent comparison needs a nonzero base close")
	}
	pal := internal.Palette(m.k)
	for i := range m.overlays {
		line := &m.overlays[i]
		values := line.values
		base := m.base
		line.values = ChartMemo(m.ch, fmt.Sprint("compare.overlay.", i), [2]any{ChartID(values), base}, func() []float64 {
			out := make([]float64, len(values))
			for j, v := range values {
				out[j] = (v/base - 1) * 100
			}
			return out
		})
	}
	for i, s := range *m.opts.Compare {
		if s.Name == "" || len(s.Closes) != len(m.candles) {
			panic(fmt.Sprintf("mujica: CompareSymbol %q needs a Name and one close per candle", s.Name))
		}
		ChartMemo(m.ch, "compare.check."+s.Name, ChartID(s.Closes), func() bool {
			for _, v := range s.Closes {
				if math.IsInf(v, 0) || (!math.IsNaN(v) && v <= 0) {
					panic("mujica: compared closes must be positive or NaN")
				}
			}
			return true
		})
		base := math.NaN()
		for j := m.from; j < m.to && math.IsNaN(base); j++ {
			base = s.Closes[j]
		}
		closes := s.Closes
		vals := ChartMemo(m.ch, "compare."+s.Name, [2]any{ChartID(closes), math.Float64bits(base)}, func() []float64 {
			out := make([]float64, len(closes))
			for j, v := range closes {
				out[j] = (v/base - 1) * 100
			}
			return out
		})
		m.overlays = append(m.overlays, mktLine{name: s.Name, values: vals, color: pal[(i+1)%len(pal)], shape: i + 1})
	}
}

// CompareLegend shows series markers and removal actions.
func CompareLegend(c *ui.Context, symbols *[]CompareSymbol) core.ChoiceResult {
	if symbols == nil {
		panic("mujica: CompareLegend needs symbols")
	}
	k := core.Tokens(c)
	pal := internal.Palette(k)
	u := core.Density(c).Unit()
	removed := -1
	row := ui.Row(c).Wrap().Gap(u * 2).Role(ui.RoleList).Label(core.Msg(c, "compare_symbol.label"))
	row.Children(func() {
		for i, s := range *symbols {
			ui.Row(c.Key(s.Name)).AlignItems(ui.Center).Gap(u).Padding(0, 0, 0, u*2).Radius(6).Border(1, k.Border).Role(ui.RoleListItem).Label(s.Name).Children(func() {
				col := pal[(i+1)%len(pal)]
				ui.Box(c).Size(10, 10).Draw(func(p *ui.Painter, r ui.Rect) { internal.Marker(p, i+1, r.X+5, r.Y+5, 9, col) })
				ui.Text(c, s.Name).SingleLine()
				if core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "compare_symbol.remove"), s.Name)).Clicked() && !row.IsDisabled() {
					removed = i
				}
			})
		}
	})
	if removed >= 0 {
		slog.Debug("mujica: compare_symbol removed", "name", (*symbols)[removed].Name)
		*symbols = append((*symbols)[:removed:removed], (*symbols)[removed+1:]...)
	}
	return core.NewChoiceResult(row, removed >= 0)
}

// alertsCheck fires callbacks when the latest close crosses thresholds.
func (m *mktCore) alertsCheck() {
	if m.opts.Alerts != nil {
		seen := map[string]bool{}
		for _, a := range *m.opts.Alerts {
			if a.ID == "" || seen[a.ID] || math.IsNaN(a.Price) || math.IsInf(a.Price, 0) {
				panic("mujica: alerts need unique IDs and finite prices")
			}
			seen[a.ID] = true
		}
	}
	if m.opts.Alerts == nil || len(m.candles) == 0 {
		return
	}
	last := m.candles[len(m.candles)-1]
	if !m.ch.Res.Element.IsDisabled() && m.st.alertSeen && last.Close != m.st.lastClose && m.opts.OnAlert != nil {
		for _, a := range *m.opts.Alerts {
			if (m.st.lastClose-a.Price)*(last.Close-a.Price) <= 0 && m.st.lastClose != a.Price {
				slog.Debug("mujica: price_alert_line crossed", "id", a.ID)
				m.opts.OnAlert(a, last)
			}
		}
	}
	m.st.lastClose, m.st.alertSeen = last.Close, true
}

// alertsInput drags the threshold captured by the initial press.
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
			if g.plot.Contains(x, y-dy) && math.Abs(float64(g.ys.Map(a.Price)-(y-dy))) <= 14 {
				m.st.grab = i
			}
		}
	}
	if m.st.grab < 0 || m.st.grab >= len(alerts) {
		return false
	}
	spec := input.NumberInputResolve(input.NumberInputOptions{Step: math.Pow10(-m.dec), Precision: m.dec})
	alerts[m.st.grab].Price = spec.Round(g.ys.Invert(y))
	slog.Debug("mujica: price_alert_line moved", "id", alerts[m.st.grab].ID)
	return true
}

// paintAlerts draws dashed thresholds and price labels.
func (m *mktCore) paintAlerts(p *ui.Painter, g mktGeom, size float32) {
	if m.opts.Alerts == nil || m.base != 0 {
		return
	}
	for i, a := range *m.opts.Alerts {
		y := internal.Snap(p, g.ys.Map(a.Price))
		if i == m.st.alertFocus {
			p.FillPath(internal.Diamond(internal.Snap(p, g.plot.X+8), y, 9), m.k.Warning)
		}
		for x := g.plot.X; x < g.plot.X+g.plot.W; x += 8 {
			p.Line(x, y, min(x+4, g.plot.X+g.plot.W), y, 1, m.k.Warning)
		}
		text := core.PriceUngrouped(a.Price, m.dec)
		if a.Label != "" {
			text = a.Label + " " + text
		}
		mktTag(p, g.plot.X+g.plot.W, y, fmt.Sprintf(core.Msg(m.c, "price_alert_line.tag"), text), size, m.k.Warning, m.k.Surface)
	}
}
