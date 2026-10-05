package mujica

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// StreamingTextOptions configure a StreamingText.
type StreamingTextOptions struct {
	// Streaming shows the cursor after the text while more is coming.
	Streaming bool
	// Smooth reveals new text gradually instead of at once.
	Smooth bool
	// Rate is the smooth reveal speed in characters per second; zero means 80.
	Rate float64
}

// streamingState is how much of the text shows, and since when.
type streamingState struct {
	text  string
	shown float64
	at    time.Time
}

// StreamingText shows text that grows as it streams; reduced motion shows it at once.
func StreamingText(c *ui.Context, text string, opts StreamingTextOptions) *ui.Element {
	if opts.Rate < 0 {
		panic("mujica: StreamingText Rate must not be negative")
	}
	rate := opts.Rate
	if rate == 0 {
		rate = 80
	}
	k := Tokens(c)
	e := ui.Box(c).Label(text)
	st := ui.Local(e, "stream", func() streamingState {
		return streamingState{text: text, shown: float64(utf8.RuneCountInString(text)), at: c.Now()}
	})
	total := float64(utf8.RuneCountInString(text))
	if !strings.HasPrefix(text, st.text) {
		// Replaced text starts over from what it shares.
		st.shown = min(st.shown, float64(utf8.RuneCountInString(streamingCommon(st.text, text))))
	}
	st.text = text
	now := c.Now()
	if !opts.Smooth || Motion(c, time.Second) == 0 {
		st.shown = total
	} else if st.shown < total {
		st.shown = min(total, st.shown+now.Sub(st.at).Seconds()*rate)
		c.AnimationFrame()
	}
	st.at = now
	shown := text
	if n := int(st.shown); n < int(total) {
		shown = string([]rune(text)[:n])
	}
	spans := []ui.Span{{Text: shown}}
	if opts.Streaming {
		spans = append(spans, ui.Span{Text: "▍", Color: k.AccentText})
	}
	return e.Children(func() { ui.RichText(c, spans...) })
}

// streamingCommon returns the longest common prefix of a and b, on rune boundaries.
func streamingCommon(a, b string) string {
	n := 0
	for n < len(a) && n < len(b) && a[n] == b[n] {
		n++
	}
	for n > 0 && n < len(a) && !utf8.RuneStart(a[n]) {
		n--
	}
	return a[:n]
}
