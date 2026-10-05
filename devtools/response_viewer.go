package devtools

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/agent"
	"github.com/ZacharyZhang-NY/MujicaUI/code"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
	Status data.DataStatus
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
func ResponseViewer(c *ui.Context, resp *ApiResponse, opts ResponseViewerOptions) core.ChoiceResult {
	if opts.Status != data.DataReady && opts.Status != data.DataLoading && opts.Status != data.DataFailed {
		panic("mujica: ResponseViewer needs a known status")
	}
	if opts.Height < 0 || resp != nil && (resp.Status < 100 || resp.Status > 599 || resp.Size < 0 || resp.Duration < 0) {
		panic("mujica: ResponseViewer needs Height >= 0, a 1xx–5xx Status and Size >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3).Label(core.Msg(c, "response_viewer.label"))
	st := ui.Local(r.Element, "mujica.response", func() responseViewerState { return responseViewerState{} })
	r.Element.Children(func() {
		if opts.Status != data.DataReady || resp == nil {
			ui.Column(c).Height(h).Children(func() {
				r = core.NewChoiceResult(r.Element, data.DataPlaceholder(c, opts.Status, opts.Error, core.Msg(c, "response_viewer.none"), ""))
			})
			return
		}
		ui.Row(c).Gap(u * 4).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
			responseStatus(c, resp.Status)
			agent.AgentCaption(c, agent.AgentDuration(resp.Duration))
			agent.AgentCaption(c, code.ProcessListBytes(uint64(resp.Size)))
		})
		if st.body != resp.Body {
			st.body, st.isJSON = resp.Body, json.Valid([]byte(resp.Body))
		}
		tabs := []navigation.Tab{
			{Label: core.Msg(c, "response_viewer.body"), Panel: func() {
				ui.Column(c).Height(h).Children(func() {
					switch {
					case resp.Body == "":
						agent.AgentCaption(c, core.Msg(c, "response_viewer.empty"))
					case st.isJSON:
						JsonViewer(c, resp.Body, &st.json, JsonViewerOptions{Height: h})
					default:
						agent.AgentCode(c, "", resp.Body, h)
					}
				})
			}},
			{Label: fmt.Sprintf(core.Msg(c, "response_viewer.headers"), len(resp.Headers)), Panel: func() {
				ui.Column(c).Height(h).Children(func() { responseHeaders(c, resp.Headers, h) })
			}},
		}
		navigation.Tabs(c, &st.tab, tabs, navigation.TabsOptions{Label: core.Msg(c, "response_viewer.label")})
	})
	if r.Changed() {
		slog.Debug("mujica: response_viewer retry")
	}
	return r
}

// responseStatus shows "200 OK" with a severity icon and tone.
func responseStatus(c *ui.Context, code int) {
	sev := core.SeveritySuccess
	switch {
	case code >= 500:
		sev = core.SeverityError
	case code >= 400:
		sev = core.SeverityWarning
	case code >= 300 || code < 200:
		sev = core.SeverityInfo
	}
	text := strings.TrimSpace(fmt.Sprintf("%d %s", code, http.StatusText(code)))
	devtoolsSeverity(c, sev, text)
}

// responseHeaders lists headers as name and value rows.
func responseHeaders(c *ui.Context, headers []KeyValue, h float32) {
	k := core.Tokens(c)
	if len(headers) == 0 {
		agent.AgentCaption(c, core.Msg(c, "response_viewer.no_headers"))
		return
	}
	agent.AgentWellRows(c, len(headers), func(i int) {
		kv := headers[i]
		ui.Row(c).MinHeight(core.Density(c).ControlHeight()).Gap(core.Density(c).Unit() * 3).Children(func() {
			ui.Text(c, kv.Key).TextColor(k.AccentText).FontWeight(600).Width(200).Shrink(0)
			ui.Text(c, kv.Value).Selectable().Grow(1).Shrink(1)
		})
	}, h).Font(theme.MonoFont)
}
