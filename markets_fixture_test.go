package mujica

import (
	"github.com/egoist/mygo/ui"
	"math"
	"testing"
	"time"
)

// mktNow is the fixed clock of the markets tests: 10:30 in New York.
var mktNow = time.Date(2026, time.October, 5, 14, 30, 0, 0, time.UTC)

func mktNewYork() *time.Location { return time.FixedZone("EDT", -4*3600) }

func mktStrikes() []OptionStrike {
	var out []OptionStrike
	for i := range 21 {
		k := 90 + float64(i)
		call := math.Max(100-k, 0) + 1.2 + float64(20-i)*0.05
		put := math.Max(k-100, 0) + 1.2 + float64(i)*0.05
		out = append(out, OptionStrike{Strike: k,
			Call: OptionQuote{Bid: call - 0.05, Ask: call + 0.05, Last: call},
			Put:  OptionQuote{Bid: put - 0.05, Ask: put + 0.05, Last: put}})
	}
	return out
}

func mktSeries(n int, seed float64) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 100 + 10*math.Sin(float64(i)/6+seed) + float64(i)*0.3
	}
	return out
}

func mktTimes(n int) []time.Time {
	out := make([]time.Time, n)
	for i := range out {
		out[i] = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i)
	}
	return out
}

func mktHeat() HeatNode {
	return HeatNode{Name: "Market", Children: []HeatNode{
		{Name: "Technology", Children: []HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}, {Name: "NVDA", Cap: 2200, Change: 3.4}}},
		{Name: "Energy", Children: []HeatNode{{Name: "XOM", Cap: 450, Change: -2.1}, {Name: "CVX", Cap: 300, Change: -1.4}}},
		{Name: "Health", Children: []HeatNode{{Name: "LLY", Cap: 700, Change: 0.4}, {Name: "JNJ", Cap: 380, Change: 0}}},
	}}
}

func mktWatch() []WatchGroup {
	return []WatchGroup{
		{Name: "Tech", Items: []WatchItem{
			{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2, Volume: 51e6},
			{Symbol: "MSFT", Name: "Microsoft", Price: 418.2, Change: -2.5, ChangePercent: -0.6, Volume: 21e6},
			{Symbol: "NVDA", Name: "Nvidia", Price: 124.9, Change: 4.1, ChangePercent: 3.4, Volume: 310e6},
		}},
		{Name: "Energy", Items: []WatchItem{{Symbol: "XOM", Name: "Exxon", Price: 117.4, Change: -2.5, ChangePercent: -2.1, Volume: 18e6}}},
	}
}

func marketsPanic(t *testing.T, view func(*ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Use(c, Settings{}); view(c) }, 960, 500, 1)
}
