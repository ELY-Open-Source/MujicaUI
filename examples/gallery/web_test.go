package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// mirrorRig serves a one-field view through the web mirror.
type mirrorRig struct {
	t    *testing.T
	m    *mirror
	srv  *httptest.Server
	text string
}

func newMirrorRig(t *testing.T) *mirrorRig {
	r := &mirrorRig{t: t}
	var h http.Handler
	r.m, h = newMirror(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(40).Children(func() {
			input.TextInput(c, &r.text, input.TextInputOptions{Label: "Ledger"})
		})
	})
	r.srv = httptest.NewServer(h)
	t.Cleanup(r.srv.Close)
	return r
}

func (r *mirrorRig) post(in webEvent) int {
	r.t.Helper()
	body, err := json.Marshal(in)
	if err != nil {
		r.t.Fatal(err)
	}
	res, err := http.Post(r.srv.URL+"/input", "application/json", bytes.NewReader(body))
	if err != nil {
		r.t.Fatal(err)
	}
	res.Body.Close()
	return res.StatusCode
}

func (r *mirrorRig) get(path string) []byte {
	r.t.Helper()
	res, err := http.Get(r.srv.URL + path)
	if err != nil {
		r.t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		r.t.Fatalf("GET %s: %d", path, res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		r.t.Fatal(err)
	}
	return b
}

func (r *mirrorRig) value() string {
	r.m.mu.Lock()
	defer r.m.mu.Unlock()
	return r.text
}

func TestMirrorServesPageAndFirstNotice(t *testing.T) {
	r := newMirrorRig(t)
	if page := r.get("/"); !bytes.Contains(page, []byte("<canvas id=\"screen\"")) {
		t.Fatal("index has no canvas")
	}
	res, err := http.Get(r.srv.URL + "/events")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	line, err := bufio.NewReader(res.Body).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(line, `data: {"seq":`) {
		t.Fatalf("first event %q", line)
	}
}

func TestMirrorTypingAndIMEReachTheView(t *testing.T) {
	r := newMirrorRig(t)
	before := r.get("/frame.png")
	if code := r.post(webEvent{Type: "key", Code: "Tab"}); code != http.StatusNoContent {
		t.Fatalf("tab: %d", code)
	}
	for _, ch := range "Ave" {
		r.post(webEvent{Type: "key", Code: "Key" + strings.ToUpper(string(ch)), Text: string(ch), Mods: []bool{ch == 'A'}})
	}
	r.post(webEvent{Type: "compose", Text: "宫廷", Caret: 2})
	r.post(webEvent{Type: "compose"})
	r.post(webEvent{Type: "text", Text: "宫廷"})
	if got := r.value(); got != "Ave宫廷" {
		t.Fatalf("value %q, want %q", got, "Ave宫廷")
	}
	if bytes.Equal(before, r.get("/frame.png")) {
		t.Fatal("frame did not change after typing")
	}
}

func TestMirrorRejectsBadInput(t *testing.T) {
	r := newMirrorRig(t)
	for name, in := range map[string]webEvent{
		"unknown type":  {Type: "teleport"},
		"unmapped code": {Type: "key", Code: "Hyper"},
		"tiny window":   {Type: "resize", W: 10, H: 10, Scale: 1},
	} {
		if code := r.post(in); code != http.StatusBadRequest {
			t.Errorf("%s: %d, want 400", name, code)
		}
	}
	res, err := http.Post(r.srv.URL+"/input", "application/json", strings.NewReader("{"))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad json: %d", res.StatusCode)
	}
}

// mirrorWith serves view through a mirror for one test.
func mirrorWith(t *testing.T, view func(c *ui.Context)) (*mirror, *httptest.Server) {
	m, h := newMirror(view)
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return m, srv
}

func postInput(t *testing.T, srv *httptest.Server, in webEvent) {
	t.Helper()
	body, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.Post(srv.URL+"/input", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("%s: %d", in.Type, res.StatusCode)
	}
}

func TestMirrorTimersFireWithoutInput(t *testing.T) {
	start := time.Now()
	var late bool
	m, _ := mirrorWith(t, func(c *ui.Context) {
		core.Use(c, core.Settings{})
		late = c.Now().Sub(start) > 300*time.Millisecond
		if !late {
			c.After(350 * time.Millisecond)
		}
	})
	time.Sleep(900 * time.Millisecond)
	m.mu.Lock()
	defer m.mu.Unlock()
	if !late {
		t.Fatal("the view was not rebuilt after its timer with no input")
	}
}

func TestMirrorModifiedClickNeedsReleaseInPlace(t *testing.T) {
	clicks := 0
	_, srv := mirrorWith(t, func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Children(func() {
			if input.Button(c, "Seal", input.ButtonOptions{}).Clicked() {
				clicks++
			}
		})
	})
	shift := []bool{true}
	postInput(t, srv, webEvent{Type: "down", X: 40, Y: 35, Mods: shift})
	postInput(t, srv, webEvent{Type: "up", X: 400, Y: 400, Mods: shift})
	postInput(t, srv, webEvent{Type: "down", X: 40, Y: 35, Mods: shift})
	postInput(t, srv, webEvent{Type: "up", X: 41, Y: 35, Mods: shift})
	if clicks != 1 {
		t.Fatalf("clicks %d, want 1 (released elsewhere must not click)", clicks)
	}
}

func TestMirrorForwardsRepeatedCopy(t *testing.T) {
	m, _ := mirrorWith(t, func(c *ui.Context) { core.Use(c, core.Settings{}) })
	ch := make(chan []byte, 8)
	m.mu.Lock()
	m.subs[ch] = struct{}{}
	notices := 0
	for range 2 {
		m.tt.SetClipboard("same")
		m.copied = true
		if m.capture() {
			notices++
		}
	}
	m.mu.Unlock()
	if notices != 2 {
		t.Fatalf("repeated copy of the same text notified %d times, want 2", notices)
	}
	for range 2 {
		msg := <-ch
		if !bytes.Contains(msg, []byte(`"clipboard":"same"`)) {
			t.Fatalf("notice %s", msg)
		}
	}
}

func TestMirrorMapsNumpad(t *testing.T) {
	r := newMirrorRig(t)
	r.post(webEvent{Type: "key", Code: "Tab"})
	if code := r.post(webEvent{Type: "key", Code: "Numpad7", Text: "7"}); code != http.StatusNoContent {
		t.Fatalf("numpad: %d", code)
	}
	if got := r.value(); got != "7" {
		t.Fatalf("value %q", got)
	}
}

func TestMirrorForwardsEveryAnnouncement(t *testing.T) {
	m, srv := mirrorWith(t, func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Children(func() {
			if input.Button(c, "Ring", input.ButtonOptions{}).Clicked() {
				c.Announce("bell rung")
			}
		})
	})
	ch := make(chan []byte, 64)
	m.mu.Lock()
	m.subs[ch] = struct{}{}
	m.mu.Unlock()
	for range 2 {
		postInput(t, srv, webEvent{Type: "down", X: 40, Y: 35})
		postInput(t, srv, webEvent{Type: "up", X: 40, Y: 35})
	}
	heard := 0
	for len(ch) > 0 {
		if bytes.Contains(<-ch, []byte("bell rung")) {
			heard++
		}
	}
	if heard != 2 {
		t.Fatalf("announcements forwarded %d times, want 2", heard)
	}
}

func TestMirrorModifiedClickReleasedJustOutsideDoesNotClick(t *testing.T) {
	clicks := 0
	var box ui.Rect
	m, srv := mirrorWith(t, func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			b := input.Button(c, "Seal", input.ButtonOptions{})
			box = b.Bounds()
			if b.Clicked() {
				clicks++
			}
		})
	})
	m.mu.Lock()
	m.tt.Frame()
	r := box
	m.mu.Unlock()
	shift := []bool{true}
	y := r.Y + r.H/2
	postInput(t, srv, webEvent{Type: "down", X: r.X + r.W - 1, Y: y, Mods: shift})
	postInput(t, srv, webEvent{Type: "up", X: r.X + r.W + 2, Y: y, Mods: shift})
	if clicks != 0 {
		t.Fatalf("clicks %d after releasing 2 DIP outside", clicks)
	}
}

func TestMirrorNotifiesEveryCopyOfEqualText(t *testing.T) {
	m, srv := mirrorWith(t, func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Children(func() {
			if input.Button(c, "Copy", input.ButtonOptions{}).Clicked() {
				c.WriteClipboard("snippet")
			}
		})
	})
	ch := make(chan []byte, 64)
	m.mu.Lock()
	m.subs[ch] = struct{}{}
	m.mu.Unlock()
	for range 2 {
		postInput(t, srv, webEvent{Type: "down", X: 40, Y: 35})
		postInput(t, srv, webEvent{Type: "up", X: 40, Y: 35})
	}
	copies := 0
	for len(ch) > 0 {
		if bytes.Contains(<-ch, []byte(`"clipboard":"snippet"`)) {
			copies++
		}
	}
	if copies != 2 {
		t.Fatalf("copies forwarded %d, want 2", copies)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if got := m.tt.Clipboard(); got != "snippet" {
		t.Fatalf("clipboard %q", got)
	}
}
