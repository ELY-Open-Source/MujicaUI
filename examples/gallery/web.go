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
	mu     sync.Mutex
	tt     *ui.Tester
	frame  []byte
	seq    int
	hash   uint64
	seed   maphash.Seed
	subs   map[chan []byte]struct{}
	wake   chan struct{}
	urls   int
	clip   string
	copied bool
	probes int
	cursor ui.Cursor
	down   *webEvent
	enc    png.Encoder
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
		if _, err := w.Write(indexHTML); err != nil {
			slog.Warn("gallery: web page write failed", "err", err)
		}
	})
	mux.HandleFunc("GET /frame.png", m.serveFrame)
	mux.HandleFunc("GET /events", m.serveEvents)
	mux.HandleFunc("POST /input", m.serveInput)
	return m, mux
}

// loop frames at 33 ms while changing, else polls at 200 ms for timers.
func (m *mirror) loop() {
	still := 0
	for {
		wait := 33 * time.Millisecond
		if still > 30 {
			wait = 200 * time.Millisecond
		}
		select {
		case <-m.wake:
			still = 0
		case <-time.After(wait):
		}
		m.mu.Lock()
		m.sync()
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

// sync moves the Tester's virtual clock up to the wall clock; the caller holds mu.
func (m *mirror) sync() {
	if d := time.Since(m.tt.Now()); d > 0 {
		m.tt.Advance(d)
	}
}

func (m *mirror) kick() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}

// capture encodes a changed frame and notifies browsers; the caller holds mu.
func (m *mirror) capture() bool {
	img := m.tt.Image()
	h := maphash.Bytes(m.seed, img.Pix)
	notice := m.notices()
	if h == m.hash && len(notice) == 0 {
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
	notice["seq"] = m.seq
	notice["cursor"] = cursorCSS(m.cursor)
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

// notices gathers opened URLs, announcements, clipboard writes and cursor changes.
func (m *mirror) notices() map[string]any {
	n := map[string]any{}
	if urls := m.tt.OpenedURLs(); len(urls) > m.urls {
		n["open"], m.urls = urls[m.urls:], len(urls)
	}
	if told := m.tt.Announcements(); len(told) > 0 {
		n["announce"] = told
	}
	if clip := m.tt.Clipboard(); clip != m.clip || m.copied {
		n["clipboard"], m.clip, m.copied = clip, clip, false
	}
	if c := m.tt.Cursor(); c != m.cursor {
		m.cursor = c
		n["cursorChanged"] = true
	}
	return n
}

func (m *mirror) serveFrame(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	frame := m.frame
	m.mu.Unlock()
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	if _, err := w.Write(frame); err != nil {
		slog.Warn("gallery: web frame write failed", "err", err)
	}
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
	if _, err := fmt.Fprintf(w, "data: %s\n\n", first); err != nil {
		slog.Warn("gallery: web events write failed", "remote", r.RemoteAddr, "err", err)
		return
	}
	flusher.Flush()
	for {
		select {
		case msg := <-ch:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg); err != nil {
				slog.Warn("gallery: web events write failed", "remote", r.RemoteAddr, "err", err)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
