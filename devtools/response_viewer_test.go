package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
)

func responseSample() *ApiResponse {
	return &ApiResponse{Status: 201, Duration: 182 * time.Millisecond, Size: 2048,
		Headers: []KeyValue{{Key: "Content-Type", Value: "application/json"}, {Key: "X-Trace", Value: "ab12"}},
		Body:    `{"id":7,"hall":"east wing"}`}
}

func responseTester(resp **ApiResponse, opts *ResponseViewerOptions, retries *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ResponseViewer(c, *resp, *opts).Changed() {
			*retries++
		}
	}, 640, 520)
}

func TestResponseViewerJsonTextHeaders(t *testing.T) {
	resp := responseSample()
	opts := ResponseViewerOptions{}
	n := 0
	tt := responseTester(&resp, &opts, &n)
	for _, s := range []string{"201 Created", "182 ms", "2 KB", `id: 7`} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Headers (2)")
	if !tt.HasText("X-Trace") || !tt.HasText("ab12") {
		t.Fatalf("headers: %q", tt.Texts())
	}
	resp = &ApiResponse{Status: 503, Body: "upstream down"}
	tt.Frame()
	_ = tt.Click("Body")
	if !tt.HasText("503 Service Unavailable") || !tt.HasText("upstream down") {
		t.Fatalf("text body: %q", tt.Texts())
	}
}

func TestResponseViewerStates(t *testing.T) {
	var resp *ApiResponse
	opts := ResponseViewerOptions{}
	n := 0
	tt := responseTester(&resp, &opts, &n)
	if !tt.HasText("Send a request to see its response") {
		t.Fatalf("none: %q", tt.Texts())
	}
	opts = ResponseViewerOptions{Status: data.DataFailed, Error: "connection refused"}
	tt.Frame()
	_ = tt.Click("Retry")
	if n != 1 || !tt.HasText("connection refused") {
		t.Fatalf("failed: %q retries %d", tt.Texts(), n)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("status 42 did not panic")
		}
	}()
	bad := &ApiResponse{Status: 42}
	responseTester(&bad, &ResponseViewerOptions{}, &n)
}

func TestResponseViewerSnapshot(t *testing.T) {
	snapshot(t, "response_viewer", 620, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ResponseViewer(c, responseSample(), ResponseViewerOptions{Height: 240})
		})
	})
}

func TestResponseViewerKeyboardDisabledAndLoading(t *testing.T) {
	resp := responseSample()
	disabled := false
	retries := 0
	opts := ResponseViewerOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if ResponseViewer(c, resp, opts).Changed() {
				retries++
			}
		})
	}, 640, 520)
	if err := tt.Click("Body"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("X-Trace") {
		t.Fatal("keyboard header selection")
	}
	opts.Status = data.DataLoading
	tt.Frame()
	if !tt.HasText("Loading") || tt.HasText("X-Trace") {
		t.Fatal("loading state")
	}
	opts.Status = data.DataFailed
	opts.Error = "offline"
	disabled = true
	tt.Frame()
	_ = tt.Click("Retry")
	tt.Key(0, ui.KeyEnter)
	if retries != 0 {
		t.Fatal("disabled retry")
	}
}

func TestResponseViewerRejectsUnknownDataStatus(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unknown data status accepted")
		}
	}()
	resp := responseSample()
	n := 0
	responseTester(&resp, &ResponseViewerOptions{Status: data.DataStatus(99)}, &n)
}

func TestResponseViewerEmptyBodyAndHeaders(t *testing.T) {
	resp := &ApiResponse{Status: 204}
	n := 0
	tt := responseTester(&resp, &ResponseViewerOptions{}, &n)
	if !tt.HasText("204 No Content") || !tt.HasText("Empty body") {
		t.Fatal("empty response")
	}
	_ = tt.Click("Headers (0)")
	if !tt.HasText("No headers") {
		t.Fatal("empty headers")
	}
}

func TestResponseViewerAncestorDisabledTabs(t *testing.T) {
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			ResponseViewer(c, responseSample(), ResponseViewerOptions{})
		})
	}, 640, 520)
	if err := tt.Click("Body"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyRight)
	if tt.HasText("X-Trace") || !tt.HasText("id: 7") {
		t.Fatal("disabled response switched tabs")
	}
}
