package mujica

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// ApiRequest is an HTTP request being built; the caller sends it.
type ApiRequest struct {
	// Method is GET when empty.
	Method  string
	URL     string
	Params  []KeyValue
	Headers []KeyValue
	Body    string
}

// ApiRequestBuilderOptions configure an ApiRequestBuilder.
type ApiRequestBuilderOptions struct {
	// Sending shows the request in flight and blocks another send.
	Sending  bool
	Disabled bool
}

// ApiRequestResult reports this frame's events.
type ApiRequestResult struct {
	Element *ui.Element
	// Send asks the caller to send the request.
	Send    bool
	Changed bool
}

type apiRequestState struct{ tab int }

var apiMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// ApiRequestBuilder edits method, URL, query parameters, headers and body, and reports a send.
func ApiRequestBuilder(c *ui.Context, req *ApiRequest, opts ApiRequestBuilderOptions) ApiRequestResult {
	if req == nil {
		panic("mujica: ApiRequestBuilder needs a request")
	}
	if req.Method == "" {
		req.Method = "GET"
	}
	if !slices.Contains(apiMethods, req.Method) {
		panic("mujica: ApiRequestBuilder method " + req.Method + " is not offered")
	}
	u := Density(c).Unit()
	var r ApiRequestResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3).Disabled(opts.Disabled)
	st := ui.Local(r.Element, "mujica.apirequest", func() apiRequestState { return apiRequestState{} })
	locked := r.Element.IsDisabled()
	full, urlErr := apiRequestURL(*req)
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
			methods := make([]SelectOption[string], len(apiMethods))
			for i, m := range apiMethods {
				methods[i] = SelectOption[string]{Value: m, Label: m}
			}
			r.Changed = Select(c, &req.Method, methods, SelectOptions{Label: msg(c, "api_request_builder.method"), Disabled: locked}).Changed()
			shownErr := ""
			if req.URL != "" && urlErr != nil {
				shownErr = urlErr.Error()
			}
			in := TextInput(c, &req.URL, TextInputOptions{Placeholder: "https://", Label: msg(c, "api_request_builder.url"), Disabled: locked, Error: shownErr})
			in.Element.Grow(1)
			r.Changed = r.Changed || in.Changed()
			ready := urlErr == nil && !opts.Sending && !locked
			send := Button(c, msg(c, "api_request_builder.send"), ButtonOptions{Icon: icons.Must("upload"), Loading: opts.Sending, Disabled: locked || urlErr != nil})
			r.Send = ready && (send.Clicked() || in.Submitted())
		})
		if urlErr == nil {
			codeMono(c, req.Method+" "+full).TextColor(Tokens(c).TextMuted)
		}
		tabs := []Tab{
			{Label: fmt.Sprintf(msg(c, "api_request_builder.params"), len(keyValueEnabled(req.Params))), Panel: func() {
				r.Changed = KeyValueInput(c, &req.Params, KeyValueInputOptions{Label: msg(c, "api_request_builder.params_label"), Disabled: locked}).Changed() || r.Changed
			}},
			{Label: fmt.Sprintf(msg(c, "api_request_builder.headers"), len(keyValueEnabled(req.Headers))), Panel: func() {
				r.Changed = KeyValueInput(c, &req.Headers, KeyValueInputOptions{Label: msg(c, "api_request_builder.headers_label"), Disabled: locked}).Changed() || r.Changed
			}},
			{Label: msg(c, "api_request_builder.body"), Panel: func() {
				if req.Method == "GET" || req.Method == "HEAD" {
					agentCaption(c, msg(c, "api_request_builder.no_body"))
				}
				body := TextArea(c, &req.Body, TextAreaOptions{Label: msg(c, "api_request_builder.body"), MinHeight: 120, MaxHeight: 320, Disabled: locked})
				body.Input.Font(theme.MonoFont)
				r.Changed = r.Changed || body.Changed()
			}},
		}
		Tabs(c, &st.tab, tabs, TabsOptions{Label: msg(c, "api_request_builder.parts")})
	})
	if r.Send {
		slog.Debug("mujica: api_request_builder send", "method", req.Method) // never the URL: it may hold a token
	}
	return r
}

// apiRequestURL joins the URL with its enabled query parameters, or says why it is invalid.
func apiRequestURL(req ApiRequest) (string, error) {
	if req.URL == "" {
		return "", fmt.Errorf("empty URL")
	}
	u, err := url.Parse(req.URL)
	if err != nil {
		return "", err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("the URL needs http:// or https:// and a host")
	}
	q := u.Query()
	for _, kv := range keyValueEnabled(req.Params) {
		q.Add(kv.Key, kv.Value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
