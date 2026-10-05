package mujica

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"math"
	"mujicaui/internal"
	"mujicaui/theme"
	"strings"
	"testing"
	"time"
)

func TestMarketsAncestorDisabled(t *testing.T) {
	s := NewsFeedState{}
	disabled := false
	var last NewsFeedView
	items := []NewsItem{{ID: "a", Title: "Alpha", Source: "Wire", Time: mktNow}, {ID: "b", Title: "Beta", Source: "Wire", Time: mktNow}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(360).Disabled(disabled).Children(func() { last = NewsFeed(c, &s, items, NewsFeedOptions{Today: DateOf(mktNow)}) })
	}, 640, 400)
	tt.Frame()
	_ = tt.Click("Alpha")
	disabled = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if s.Read["b"] || last.Opened != "" {
		t.Fatal("ancestor-disabled feed changed reads")
	}
}
func TestCurrencyConverterSwapInvalidatesRate(t *testing.T) {
	s := CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
	var last CurrencyConverterView
	tt := chatTester(760, 300, func(c *ui.Context) {
		last = CurrencyConverter(c, &s, []Currency{{"USD", "Dollar", 2}, {"JPY", "Yen", 0}}, CurrencyConverterOptions{Rate: 150})
	})
	_ = tt.Click("Swap currencies")
	if s.From != "JPY" || s.To != "USD" || !math.IsNaN(last.Result) {
		t.Fatalf("swapped pair %s/%s result %v", s.From, s.To, last.Result)
	}
}
func TestGreeksTableLarge(t *testing.T) {
	rows := make([]GreeksRow, 100000)
	for i := range rows {
		rows[i] = GreeksRow{Label: "Contract", Greeks: OptionGreeks{Delta: 0.5}}
	}
	tt := chatTester(640, 400, func(c *ui.Context) { GreeksTable(c, rows, GreeksTableOptions{}) })
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("built %d texts", n)
	}
}
func TestSymbolSearchLarge(t *testing.T) {
	rows := make([]MarketSymbol, 100000)
	for i := range rows {
		rows[i] = MarketSymbol{Code: statisticFormat(float64(i), 0), Name: "Company"}
	}
	s := SymbolSearchState{Query: "company"}
	tt := chatTester(640, 400, func(c *ui.Context) { SymbolSearch(c, &s, rows, SymbolSearchOptions{}) })
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("built %d texts", n)
	}
}

func TestMarketsLargeViews(t *testing.T) {
	const n = 100000
	news := make([]NewsItem, n)
	econ := make([]EconomicEvent, n)
	earn := make([]EarningsEvent, n)
	rows := make([]ScreenerRow, n)
	strikes := make([]OptionStrike, n)
	for i := range n {
		id := fmt.Sprint(i)
		news[i] = NewsItem{ID: id, Title: "Headline " + id, Time: mktNow}
		econ[i] = EconomicEvent{ID: id, Country: "US", Title: "Event " + id, Importance: 3}
		earn[i] = EarningsEvent{ID: id, Symbol: "S" + id}
		rows[i] = ScreenerRow{Symbol: "S" + id, Values: map[string]float64{"pe": 10}}
		strikes[i] = OptionStrike{Strike: float64(i)}
	}
	for _, tc := range []struct {
		name string
		view func(*ui.Context)
	}{
		{"news", func(c *ui.Context) {
			s := ui.Local(c.Root(), "large.news", func() NewsFeedState { return NewsFeedState{} })
			NewsFeed(c, s, news, NewsFeedOptions{}).Element.Height(350)
		}},
		{"economic", func(c *ui.Context) {
			s := ui.Local(c.Root(), "large.econ", func() EconomicCalendarState { return EconomicCalendarState{} })
			EconomicCalendar(c, s, econ, EconomicCalendarOptions{}).Element.Height(350)
		}},
		{"earnings", func(c *ui.Context) {
			s := ui.Local(c.Root(), "large.earn", func() DataTableState[string] { return DataTableState[string]{} })
			EarningsCalendar(c, s, earn, EarningsCalendarOptions{}).Element.Height(350)
		}},
		{"screener", func(c *ui.Context) {
			s := ui.Local(c.Root(), "large.screen", func() ScreenerState { return ScreenerState{} })
			Screener(c, s, []ScreenerField{{ID: "pe", Label: "P/E"}}, rows, ScreenerOptions{}).Element.Height(350)
		}},
		{"chain", func(c *ui.Context) {
			s := ui.Local(c.Root(), "large.chain", func() OptionChainState { return OptionChainState{} })
			OptionChain(c, s, []string{"Oct"}, strikes, OptionChainOptions{Spot: 50000}).Element.Height(350)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := chatTester(960, 400, tc.view)
			tt.Frame()
			if count := len(tt.Texts()); count > 300 {
				t.Fatalf("built %d texts", count)
			}
		})
	}
}
func TestMarketHeatmapLargeGeometry(t *testing.T) {
	m := marketHeatMemo{}
	for i := range 100000 {
		m.tiles = append(m.tiles, marketHeatTile{name: fmt.Sprint(i), cap: 1, change: 1})
		m.caps = append(m.caps, 1)
	}
	m.layout(m.caps, ui.Rect{W: 640, H: 240})
	if count := len(m.visible) + len(m.small); count > 640*240/16+640/4+240/4 {
		t.Fatalf("built %d marks", count)
	}
	if len(m.rects) != 100000 {
		t.Fatalf("keyboard points %d", len(m.rects))
	}
	first := &m.rects[0]
	m.layout(m.caps, ui.Rect{W: 640, H: 240})
	if first != &m.rects[0] {
		t.Fatal("idle layout recomputed")
	}
}
func TestMarketsPaletteContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for i, col := range internal.Palette(k) {
			if ratio := theme.Contrast(col, k.Surface); ratio < 3 {
				t.Fatalf("color %d ratio %v", i, ratio)
			}
		}
	}
}
func TestMarketsChartStatesAndDisabled(t *testing.T) {
	for _, name := range []string{"payoff", "heatmap", "allocation", "performance"} {
		t.Run(name, func(t *testing.T) {
			for _, status := range []DataStatus{DataLoading, DataFailed} {
				frame := ChartFrame{Label: "Chart", Height: 180, Status: status, Error: "Caller failed"}
				retry := false
				tt := chatTester(960, 420, func(c *ui.Context) {
					var r ChartResult
					switch name {
					case "payoff":
						r = PayoffDiagram(c, bullSpread(), PayoffDiagramOptions{Frame: frame})
					case "heatmap":
						path := ui.Local(c.Root(), "path", func() []string { return nil })
						r = MarketHeatmap(c, path, mktHeat(), MarketHeatmapOptions{Frame: frame})
					case "allocation":
						r = AssetAllocationChart(c, []AllocationItem{{"Stocks", 1, 100}}, AssetAllocationChartOptions{Frame: frame})
					case "performance":
						r = PerformanceChart(c, []float64{100, 110}, PerformanceChartOptions{Frame: frame, Times: mktTimes(2)})
					}
					retry = retry || r.Retried()
				})
				if status == DataFailed {
					_ = tt.Click("Retry")
					if !retry {
						t.Fatal("retry missing")
					}
				} else if !strings.Contains(strings.Join(tt.Texts(), "|"), "Loading") {
					t.Fatalf("loading %q", tt.Texts())
				}
			}
		})
	}
}

func TestMarketsLargeChartGeometry(t *testing.T) {
	values := mktSeries(100000, 0)
	times := make([]time.Time, 100000)
	items := make([]AllocationItem, 100000)
	for i := range times {
		times[i] = mktNow.Add(time.Duration(i) * time.Minute)
		items[i] = AllocationItem{Name: fmt.Sprint("Asset ", i), Value: 1}
	}
	for _, name := range []string{"payoff", "performance", "allocation"} {
		t.Run(name, func(t *testing.T) {
			var geometry *chartMemoEntry
			tt := chatTester(960, 600, func(c *ui.Context) {
				var r ChartResult
				frame := ChartFrame{Label: "Large"}
				switch name {
				case "payoff":
					r = PayoffDiagram(c, bullSpread(), PayoffDiagramOptions{Frame: frame, Points: 100000})
					geometry = ui.Local(r.Element, "mujica.memo.payoff.geometry", func() chartMemoEntry { return chartMemoEntry{} })
				case "performance":
					r = PerformanceChart(c, values, PerformanceChartOptions{Frame: frame, Times: times})
					geometry = ui.Local(r.Element, "mujica.memo.performance.geometry.Portfolio", func() chartMemoEntry { return chartMemoEntry{} })
				case "allocation":
					r = AssetAllocationChart(c, items, AssetAllocationChartOptions{Frame: frame})
				}
			})
			tt.Frame()
			if len(tt.Texts()) > 300 {
				t.Fatalf("built %d texts", len(tt.Texts()))
			}
			if geometry != nil {
				points := geometry.val.([]int)
				if len(points) > 4*960 {
					t.Fatalf("built %d points", len(points))
				}
				first := &points[0]
				tt.Frame()
				if first != &geometry.val.([]int)[0] {
					t.Fatal("idle geometry recomputed")
				}
			}
		})
	}
}
