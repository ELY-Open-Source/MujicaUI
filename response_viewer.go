package mujica

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ApiResponse is an HTTP response the caller received.
type ApiResponse struct {
	Status   int
	Duration time.Duration
	// Size is the body size in bytes.
	Size    int64
	Headers []KeyValue
	Body    string
}

// ResponseViewerOptions configure a ResponseViewer.
type ResponseViewerOptions struct {
	// Status replaces the response while loading or failed.
	Status DataStatus
	Error  string
	// Height of the body; zero means 320.
	Height float32
}

type responseViewerState struct {
	tab    int
	json   JsonViewerState
	body   string
	isJSON bool
}

// ResponseViewer shows response metadata, body and headers.
func ResponseViewer(c *ui.Context, resp *ApiResponse, opts ResponseViewerOptions) ChoiceResult {
	if opts.Status != DataReady && opts.Status != DataLoading && opts.Status != DataFailed {
		panic("mujica: ResponseViewer needs a known status")
	}
	if opts.Height < 0 || resp != nil && (resp.Status < 100 || resp.Status > 599 || resp.Size < 0 || resp.Duration < 0) {
		panic("mujica: ResponseViewer needs Height >= 0, a 1xx–5xx Status and Size >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3).Label(msg(c, "response_viewer.label"))
	st := ui.Local(r.Element, "mujica.response", func() responseViewerState { return responseViewerState{} })
	r.Element.Children(func() {
		if opts.Status != DataReady || resp == nil {
			ui.Column(c).Height(h).Children(func() {
				r.changed = dataPlaceholder(c, opts.Status, opts.Error, msg(c, "response_viewer.none"), "")
			})
			return
		}
		ui.Row(c).Gap(u * 4).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
			responseStatus(c, resp.Status)
			agentCaption(c, agentDuration(resp.Duration))
			agentCaption(c, processListBytes(uint64(resp.Size)))
		})
		if st.body != resp.Body {
			st.body, st.isJSON = resp.Body, json.Valid([]byte(resp.Body))
		}
		tabs := []Tab{
			{Label: msg(c, "response_viewer.body"), Panel: func() {
				ui.Column(c).Height(h).Children(func() {
					switch {
					case resp.Body == "":
						agentCaption(c, msg(c, "response_viewer.empty"))
					case st.isJSON:
						JsonViewer(c, resp.Body, &st.json, JsonViewerOptions{Height: h})
					default:
						agentCode(c, "", resp.Body, h)
					}
				})
			}},
			{Label: fmt.Sprintf(msg(c, "response_viewer.headers"), len(resp.Headers)), Panel: func() {
				ui.Column(c).Height(h).Children(func() { responseHeaders(c, resp.Headers, h) })
			}},
		}
		Tabs(c, &st.tab, tabs, TabsOptions{Label: msg(c, "response_viewer.label")})
	})
	if r.changed {
		slog.Debug("mujica: response_viewer retry")
	}
	return r
}

// responseStatus shows "200 OK" with a severity icon and tone.
func responseStatus(c *ui.Context, code int) {
	sev := SeveritySuccess
	switch {
	case code >= 500:
		sev = SeverityError
	case code >= 400:
		sev = SeverityWarning
	case code >= 300 || code < 200:
		sev = SeverityInfo
	}
	text := strings.TrimSpace(fmt.Sprintf("%d %s", code, http.StatusText(code)))
	devtoolsSeverity(c, sev, text)
}

// responseHeaders lists headers as name and value rows.
func responseHeaders(c *ui.Context, headers []KeyValue, h float32) {
	k := Tokens(c)
	if len(headers) == 0 {
		agentCaption(c, msg(c, "response_viewer.no_headers"))
		return
	}
	agentWellRows(c, len(headers), func(i int) {
		kv := headers[i]
		ui.Row(c).MinHeight(Density(c).ControlHeight()).Gap(Density(c).Unit() * 3).Children(func() {
			ui.Text(c, kv.Key).TextColor(k.AccentText).FontWeight(600).Width(200).Shrink(0)
			ui.Text(c, kv.Value).Selectable().Grow(1).Shrink(1)
		})
	}, h).Font(theme.MonoFont)
}
