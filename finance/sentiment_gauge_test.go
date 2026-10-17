package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func TestSentimentGaugeBandsAndBoundary(t *testing.T) {
	for i, value := range []float64{0, 20, 40, 60, 80, 100} {
		tt := chatTester(400, 300, func(c *ui.Context) { SentimentGauge(c, value, SentimentGaugeOptions{Label: "Sentiment"}) })
		words := []string{"Extreme fear", "Fear", "Neutral", "Greed", "Extreme greed", "Extreme greed"}
		if !tt.HasText(words[i]) {
			t.Fatalf("value %v texts %q", value, tt.Texts())
		}
	}
	marketsPanic(t, func(c *ui.Context) { SentimentGauge(c, 101, SentimentGaugeOptions{}) })
}
func TestSentimentGaugeSnapshot(t *testing.T) {

	snapshot(t, "sentiment_gauge", 320, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Fear & Greed"})
		})
	})
}

func TestSentimentGaugeKeyboardReadout(t *testing.T) {
	tt := chatTester(400, 300, func(c *ui.Context) { SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Sentiment"}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if len(tt.Announcements()) == 0 {
		t.Fatal("missing keyboard readout")
	}
}

func TestSentimentGaugeStatesKeepSize(t *testing.T) {
	status := data.DataReady
	empty := false
	retried := 0
	var bounds ui.Rect
	tt := chatTester(400, 300, func(c *ui.Context) {
		bounds = SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Sentiment", Status: status, Empty: empty, Error: "Caller failed", Retry: func() { retried++ }}).Bounds()
	})
	settled := func() float32 { tt.Frame(); return bounds.H } // Bounds reads the previous frame
	height := settled()
	status = data.DataLoading
	tt.Frame()
	if h := settled(); h != height {
		t.Fatalf("loading height %v, ready %v", h, height)
	}
	status = data.DataFailed
	tt.Frame()
	_ = tt.Click("Retry")
	if retried != 1 {
		t.Fatalf("retries %d", retried)
	}
	status = data.DataReady
	empty = true
	tt.Frame()
	if h := settled(); !tt.HasText("No data") || h != height {
		t.Fatalf("empty height %v, ready %v", h, height)
	}
}
