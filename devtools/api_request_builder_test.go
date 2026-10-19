package devtools

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func apiTester(req *ApiRequest, opts *ApiRequestBuilderOptions, sends *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(8).Children(func() {
			if ApiRequestBuilder(c, req, *opts).Send {
				*sends++
			}
		})
	}, 760, 520)
}

func TestApiRequestBuilderBuildAndSend(t *testing.T) {
	var req ApiRequest
	opts := ApiRequestBuilderOptions{}
	sends := 0
	tt := apiTester(&req, &opts, &sends)
	if req.Method != "GET" {
		t.Fatalf("default method %q", req.Method)
	}
	_ = tt.Click("Send") // disabled without a URL
	_ = tt.Click("Request URL")
	tt.Type("ftp://court")
	if !tt.HasText("Enter a valid HTTP or HTTPS URL") || sends != 0 {
		t.Fatalf("invalid URL: %q", tt.Texts())
	}
	req.URL = ""
	tt.Frame()
	_ = tt.Click("Request URL")
	tt.Type("https://court.example/guests")
	_ = tt.Click("Add row")
	_ = tt.Click("Row 1 Key")
	tt.Type("page")
	tt.Key(0, ui.KeyTab)
	tt.Type("2")
	if !tt.HasText("GET https://court.example/guests?page=2") || !tt.HasText("Params (1)") {
		t.Fatalf("preview: %q", tt.Texts())
	}
	_ = tt.Click("Request URL")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Send")
	if sends != 2 {
		t.Fatalf("sends %d", sends)
	}
	_ = tt.Click("Method")
	_ = tt.Click("POST")
	_ = tt.Click("Body")
	if req.Method != "POST" || tt.HasText("GET and HEAD requests usually carry no body.") {
		t.Fatalf("method %q", req.Method)
	}
}

func TestApiRequestBuilderSendingAndDisabled(t *testing.T) {
	req := ApiRequest{URL: "https://court.example"}
	opts := ApiRequestBuilderOptions{Sending: true}
	sends := 0
	tt := apiTester(&req, &opts, &sends)
	_ = tt.Click("Request URL")
	tt.Key(0, ui.KeyEnter)
	if sends != 0 {
		t.Fatal("sent while sending")
	}
	opts = ApiRequestBuilderOptions{Disabled: true}
	tt.Frame()
	_ = tt.Click("Send")
	_ = tt.Click("Request URL")
	tt.Type("x")
	if sends != 0 || req.URL != "https://court.example" {
		t.Fatalf("disabled: %q sends %d", req.URL, sends)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown method did not panic")
		}
	}()
	apiTester(&ApiRequest{Method: "BREW"}, &ApiRequestBuilderOptions{}, &sends)
}

func TestApiRequestBuilderSnapshot(t *testing.T) {
	snapshot(t, "api_request_builder", 720, 330, func(c *ui.Context) {
		req := ApiRequest{Method: "POST", URL: "https://court.example/guests",
			Params:  []KeyValue{{Key: "notify", Value: "true", Enabled: true}},
			Headers: []KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
		ui.Column(c).Padding(16).Children(func() { ApiRequestBuilder(c, &req, ApiRequestBuilderOptions{}) })
	})
}

func TestApiRequestBuilderDisabledPreservesDefault(t *testing.T) {
	var req ApiRequest
	sends := 0
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() { ApiRequestBuilder(c, &req, ApiRequestBuilderOptions{}) })
	}, 600, 400)
	if req.Method != "" || sends != 0 {
		t.Fatal("disabled request changed")
	}
}

func TestApiRequestBuilderAncestorDisabledTabs(t *testing.T) {
	disabled := false
	request := ApiRequest{URL: "https://court.example", Params: []KeyValue{{Key: "page", Value: "1", Enabled: true}}, Headers: []KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			ApiRequestBuilder(c, &request, ApiRequestBuilderOptions{})
		})
	}, 760, 520)
	if err := tt.Click("Params (1)"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("page") || tt.HasText("Accept") {
		t.Fatal("disabled request switched tabs")
	}
}
