package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"image/png"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

//go:embed web/index.html
var indexHTML []byte

// mirror runs the gallery view headless and streams its frames.
type mirror struct {
	mu    sync.Mutex
	tt    *ui.Tester
	frame []byte
	seq   int
	hash  uint64
	seed  maphash.Seed
	subs  map[chan []byte]struct{}
	wake  chan struct{}
	urls  int
	told  int
	clip  string
	enc   png.Encoder
}

func serveWeb(addr string, view func(c *ui.Context)) error {
	_, h := newMirror(view)
	slog.Info("gallery: web mirror listening", "addr", "http://"+addr)
	return http.ListenAndServe(addr, h)
}

// newMirror starts the frame loop of view and returns its HTTP handler.
func newMirror(view func(c *ui.Context)) (*mirror, http.Handler) {
	m := &mirror{
		tt:   ui.NewTester(view, 1200, 800),
		seed: maphash.MakeSeed(),
		subs: map[chan []byte]struct{}{},
		wake: make(chan struct{}, 1),
		enc:  png.Encoder{CompressionLevel: png.BestSpeed},
	}
	m.capture()
	go m.loop()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	mux.HandleFunc("GET /frame.png", m.serveFrame)
	mux.HandleFunc("GET /events", m.serveEvents)
	mux.HandleFunc("POST /input", m.serveInput)
	return m, mux
}

// loop draws frames while the picture keeps changing, then sleeps.
func (m *mirror) loop() {
	for range m.wake {
		for still := 0; still < 6; {
			time.Sleep(33 * time.Millisecond)
			m.mu.Lock()
			m.tt.Frame()
			changed := m.capture()
			m.mu.Unlock()
			if changed {
				still = 0
			} else {
				still++
			}
		}
	}
}

func (m *mirror) kick() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}

// capture encodes the frame and notifies browsers; the caller holds mu.
func (m *mirror) capture() bool {
	img := m.tt.Image()
	h := maphash.Bytes(m.seed, img.Pix)
	notice := m.notices()
	if h == m.hash && notice == nil {
		return false
	}
	if h != m.hash {
		var b bytes.Buffer
		if err := m.enc.Encode(&b, img); err != nil {
			panic(err)
		}
		m.frame, m.hash = b.Bytes(), h
		m.seq++
	}
	if notice == nil {
		notice = map[string]any{}
	}
	notice["seq"] = m.seq
	notice["cursor"] = cursorCSS(m.tt.Cursor())
	msg, err := json.Marshal(notice)
	if err != nil {
		panic(err)
	}
	for ch := range m.subs {
		select {
		case ch <- msg:
		default:
			slog.Warn("gallery: web client lagging, frame notice dropped", "seq", m.seq)
		}
	}
	return true
}

// notices gathers opened URLs, announcements and clipboard changes.
func (m *mirror) notices() map[string]any {
	n := map[string]any{}
	if urls := m.tt.OpenedURLs(); len(urls) > m.urls {
		n["open"], m.urls = urls[m.urls:], len(urls)
	}
	if told := m.tt.Announcements(); len(told) > m.told {
		n["announce"], m.told = told[m.told:], len(told)
	}
	if clip := m.tt.Clipboard(); clip != m.clip {
		n["clipboard"], m.clip = clip, clip
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

func (m *mirror) serveFrame(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	frame := m.frame
	m.mu.Unlock()
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(frame)
}

func (m *mirror) serveEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	ch := make(chan []byte, 8)
	m.mu.Lock()
	m.subs[ch] = struct{}{}
	first := fmt.Sprintf(`{"seq":%d}`, m.seq)
	m.mu.Unlock()
	slog.Info("gallery: web client connected", "remote", r.RemoteAddr)
	defer func() {
		m.mu.Lock()
		delete(m.subs, ch)
		m.mu.Unlock()
		slog.Info("gallery: web client left", "remote", r.RemoteAddr)
	}()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprintf(w, "data: %s\n\n", first)
	flusher.Flush()
	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// input is one browser event.
type input struct {
	Type   string  `json:"type"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	DX     float32 `json:"dx"`
	DY     float32 `json:"dy"`
	Button int     `json:"button"`
	Code   string  `json:"code"`
	Text   string  `json:"text"`
	Caret  int     `json:"caret"`
	Mods   []bool  `json:"mods"`
	W      int     `json:"w"`
	H      int     `json:"h"`
	Scale  float32 `json:"scale"`
	Dark   bool    `json:"dark"`
	Motion bool    `json:"reduceMotion"`
	More   bool    `json:"moreContrast"`
}

func (m *mirror) serveInput(w http.ResponseWriter, r *http.Request) {
	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.mu.Lock()
	err := m.apply(in)
	if err == nil {
		m.capture()
	}
	m.mu.Unlock()
	if err != nil {
		slog.Warn("gallery: web input rejected", "type", in.Type, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.kick()
	w.WriteHeader(http.StatusNoContent)
}

// apply sends one event to the tester; the caller holds mu.
func (m *mirror) apply(in input) error {
	mods := modifiers(in.Mods)
	switch in.Type {
	case "down":
		if in.Button == 2 {
			m.tt.RightClickAt(in.X, in.Y)
		} else if mods != 0 {
			m.tt.ClickAtWith(mods, in.X, in.Y)
		} else {
			m.tt.Press(in.X, in.Y)
		}
	case "up":
		if in.Button == 0 && mods == 0 {
			m.tt.Release(in.X, in.Y)
		}
	case "move":
		m.tt.Move(in.X, in.Y)
	case "wheel":
		m.tt.Scroll(in.X, in.Y, in.DX, in.DY)
	case "key":
		k, ok := keyCodes[in.Code]
		if !ok {
			return fmt.Errorf("unmapped key code %q", in.Code)
		}
		if in.Text != "" && mods&^ui.Shift == 0 {
			m.tt.TypeKey(mods, k, in.Text)
		} else {
			m.tt.Key(mods, k)
		}
	case "text":
		m.tt.Type(in.Text)
	case "compose":
		m.tt.Compose(in.Text, in.Caret)
	case "paste":
		m.tt.SetClipboard(in.Text)
		m.clip = in.Text
		m.tt.Command("paste")
	case "command":
		m.tt.Command(in.Text)
	case "resize":
		if in.W < 320 || in.H < 240 || in.Scale <= 0 {
			return fmt.Errorf("bad size %dx%d@%v", in.W, in.H, in.Scale)
		}
		m.tt.SetSize(in.W, in.H)
		m.tt.SetScale(in.Scale)
		m.tt.SetDark(in.Dark)
		m.tt.SetPreferences(ui.Preferences{ReduceMotion: in.Motion, HighContrast: in.More, TextScale: 1})
	case "focus":
		m.tt.SetFocused(in.Text == "in")
	default:
		return fmt.Errorf("unknown input type %q", in.Type)
	}
	return nil
}

func modifiers(m []bool) ui.Modifiers {
	var out ui.Modifiers
	for i, mod := range []ui.Modifiers{ui.Shift, ui.Ctrl, ui.Alt, ui.Super} {
		if i < len(m) && m[i] {
			out |= mod
		}
	}
	return out
}
