package devtools

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ApiRequest holds request data sent by the caller.
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

// ApiRequestBuilder edits requests and reports send actions.
func ApiRequestBuilder(c *ui.Context, req *ApiRequest, opts ApiRequestBuilderOptions) ApiRequestResult {
	if req == nil {
		panic("mujica: ApiRequestBuilder needs a request")
	}
	original := req
	if req.Method == "" {
		copy := *req
		copy.Method = "GET"
		req = &copy
	}
	if !slices.Contains(apiMethods, req.Method) {
		panic("mujica: ApiRequestBuilder method " + req.Method + " is not offered")
	}
	u := core.Density(c).Unit()
	var r ApiRequestResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3).Disabled(opts.Disabled)
	st := ui.Local(r.Element, "mujica.apirequest", func() apiRequestState { return apiRequestState{} })
	locked := r.Element.IsDisabled()
	if locked {
		copy := *req
		copy.Params = slices.Clone(req.Params)
		copy.Headers = slices.Clone(req.Headers)
		req = &copy
	} else {
		*original = *req
		req = original
	}
	full, urlErr := apiRequestURL(*req)
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
			methods := make([]input.SelectOption[string], len(apiMethods))
			for i, m := range apiMethods {
				methods[i] = input.SelectOption[string]{Value: m, Label: m}
			}
			r.Changed = input.Select(c, &req.Method, methods, input.SelectOptions{Label: core.Msg(c, "api_request_builder.method"), Disabled: locked}).Changed()
			shownErr := ""
			if req.URL != "" && urlErr != nil {
				shownErr = core.Msg(c, "api_request_builder.invalid_url")
			}
			in := input.TextInput(c, &req.URL, input.TextInputOptions{Placeholder: "https://", Label: core.Msg(c, "api_request_builder.url"), Disabled: locked, Error: shownErr})
			in.Element.Grow(1)
			r.Changed = r.Changed || in.Changed()
			ready := urlErr == nil && !opts.Sending && !locked
			send := input.Button(c, core.Msg(c, "api_request_builder.send"), input.ButtonOptions{Icon: icons.Must("upload"), Loading: opts.Sending, Disabled: locked || urlErr != nil})
			r.Send = ready && (send.Clicked() || in.Submitted())
		})
		if urlErr == nil {
			code.CodeMono(c, req.Method+" "+full).TextColor(core.Tokens(c).TextMuted)
		}
		tabs := []navigation.Tab{
			{Label: fmt.Sprintf(core.Msg(c, "api_request_builder.params"), len(keyValueEnabled(req.Params))), Panel: func() {
				r.Changed = KeyValueInput(c, &req.Params, KeyValueInputOptions{Label: core.Msg(c, "api_request_builder.params_label"), Disabled: locked}).Changed() || r.Changed
			}},
			{Label: fmt.Sprintf(core.Msg(c, "api_request_builder.headers"), len(keyValueEnabled(req.Headers))), Panel: func() {
				r.Changed = KeyValueInput(c, &req.Headers, KeyValueInputOptions{Label: core.Msg(c, "api_request_builder.headers_label"), Disabled: locked}).Changed() || r.Changed
			}},
			{Label: core.Msg(c, "api_request_builder.body"), Panel: func() {
				if req.Method == "GET" || req.Method == "HEAD" {
					agent.AgentCaption(c, core.Msg(c, "api_request_builder.no_body"))
				}
				body := input.TextArea(c, &req.Body, input.TextAreaOptions{Label: core.Msg(c, "api_request_builder.body"), MinHeight: 120, MaxHeight: 320, Disabled: locked})
				body.Input.Font(theme.MonoFont)
				r.Changed = r.Changed || body.Changed()
			}},
		}
		navigation.Tabs(c, &st.tab, tabs, navigation.TabsOptions{Label: core.Msg(c, "api_request_builder.parts")})
	})
	if r.Send {
		slog.Debug("mujica: api_request_builder send", "method", req.Method) // Request URLs may contain secrets.
	}
	return r
}

// apiRequestURL joins enabled query parameters with a validated URL.
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
