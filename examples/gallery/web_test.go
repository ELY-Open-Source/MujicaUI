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

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
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
		mujica.Use(c, mujica.Settings{})
		ui.Column(c).Padding(40).Children(func() {
			mujica.TextInput(c, &r.text, mujica.TextInputOptions{Label: "Ledger"})
		})
	})
	r.srv = httptest.NewServer(h)
	t.Cleanup(r.srv.Close)
	return r
}

func (r *mirrorRig) post(in input) int {
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
	if code := r.post(input{Type: "key", Code: "Tab"}); code != http.StatusNoContent {
		t.Fatalf("tab: %d", code)
	}
	for _, ch := range "Ave" {
		r.post(input{Type: "key", Code: "Key" + strings.ToUpper(string(ch)), Text: string(ch), Mods: []bool{ch == 'A'}})
	}
	r.post(input{Type: "compose", Text: "宫廷", Caret: 2})
	r.post(input{Type: "compose"})
	r.post(input{Type: "text", Text: "宫廷"})
	if got := r.value(); got != "Ave宫廷" {
		t.Fatalf("value %q, want %q", got, "Ave宫廷")
	}
	if bytes.Equal(before, r.get("/frame.png")) {
		t.Fatal("frame did not change after typing")
	}
}

func TestMirrorRejectsBadInput(t *testing.T) {
	r := newMirrorRig(t)
	for name, in := range map[string]input{
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
