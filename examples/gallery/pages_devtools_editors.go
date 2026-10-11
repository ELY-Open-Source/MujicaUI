package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/devtools"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func devJsonViewerDemo(c *ui.Context, log func(string)) {
	dtSection(c, "JsonViewer", func() {
		s := devControls(c, "453", log, false)
		type state struct {
			tree    devtools.JsonViewerState
			invalid bool
		}
		st := ui.Local(c.Root(), "453", func() state { return state{} })
		if input.Checkbox(c, &st.invalid, "Invalid JSON", input.CheckboxOptions{}).Changed() {
			log("JSON validity changed")
		}
		devPreview(c, s, func() {
			data := `{"server":{"port":8080,"hosts":["east","west"]},"ready":true}`
			if st.invalid {
				data = `{"port":}`
			}
			if s.empty {
				data = `{}`
			}
			if devtools.JsonViewer(c, data, &st.tree, devtools.JsonViewerOptions{}).Copied != "" {
				log("path copied")
			}
		})
	})
}

func devJsonEditorDemo(c *ui.Context, log func(string)) {
	dtSection(c, "JsonEditor", func() {
		s := devControls(c, "454", log, false)
		text := ui.Local(c.Root(), "454", func() string { return `{"port":8080,"ready":true}` })
		if input.Button(c, "Invalid sample", input.ButtonOptions{}).Clicked() {
			*text = "{\n  \"port\": [1,2,]\n}"
			log("invalid sample loaded")
		}
		devPreview(c, s, func() {
			r := devtools.JsonEditor(c, text, devtools.JsonEditorOptions{})
			if r.Changed {
				log(fmt.Sprintf("JSON edited; valid=%t", r.Valid))
			}
		})
	})
}

func devSchemaTreeDemo(c *ui.Context, log func(string)) {
	dtSection(c, "SchemaTree", func() {
		s := devControls(c, "455", log, false)
		type state struct {
			tree  devtools.SchemaTreeState
			roots []devtools.SchemaNode
		}
		st := ui.Local(c.Root(), "455", func() state {
			return state{roots: []devtools.SchemaNode{{Name: "app", Children: []devtools.SchemaNode{{Name: "users", Kind: devtools.SchemaTable, Children: []devtools.SchemaNode{{Name: "id", Kind: devtools.SchemaColumn, Detail: "bigint"}, {Name: "name", Kind: devtools.SchemaColumn, Detail: "text"}, {Name: "users_pk", Kind: devtools.SchemaIndex, Detail: "(id)"}}}, {Name: "active_users", Kind: devtools.SchemaView}}}}}
		})
		devPreview(c, s, func() {
			roots := st.roots
			if s.empty {
				roots = nil
			}
			if devtools.SchemaTree(c, &st.tree, roots, devtools.SchemaTreeOptions{}).Changed() {
				log("schema node selected")
			}
		})
	})
}

func devApiRequestBuilderDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ApiRequestBuilder", func() {
		s := devControls(c, "456", log, false)
		request := ui.Local(c.Root(), "456", func() devtools.ApiRequest {
			return devtools.ApiRequest{URL: "https://example.com/api", Headers: []devtools.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
		})
		sending := ui.Local(c.Root(), "456-sending", func() bool { return false })
		if input.Checkbox(c, sending, "Sending", input.CheckboxOptions{}).Changed() {
			log("send state changed")
		}
		devPreview(c, s, func() {
			r := devtools.ApiRequestBuilder(c, request, devtools.ApiRequestBuilderOptions{Sending: *sending})
			if r.Send {
				log("send requested")
			}
			if r.Changed {
				log("request edited")
			}
		})
	})
}

func devResponseViewerDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ResponseViewer", func() {
		s := devControls(c, "457", log, true)
		text := ui.Local(c.Root(), "457", func() bool { return false })
		if input.Checkbox(c, text, "Text body", input.CheckboxOptions{}).Changed() {
			log("response format changed")
		}
		devPreview(c, s, func() {
			resp := &devtools.ApiResponse{Status: 200, Duration: 182 * time.Millisecond, Size: 2048, Body: `{"ready":true,"requests":1248}`, Headers: []devtools.KeyValue{{Key: "Content-Type", Value: "application/json"}, {Key: "X-Trace", Value: "demo"}}}
			if *text {
				resp.Body = "Response received successfully."
			}
			if s.empty {
				resp = nil
			}
			if devtools.ResponseViewer(c, resp, devtools.ResponseViewerOptions{Height: 260, Status: s.status, Error: "Connection refused"}).Changed() {
				log("retry requested")
			}
		})
	})
}

func devKeyValueInputDemo(c *ui.Context, log func(string)) {
	dtSection(c, "KeyValueInput", func() {
		s := devControls(c, "458", log, false)
		rows := ui.Local(c.Root(), "458", func() []devtools.KeyValue {
			return []devtools.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}, {Key: "debug", Value: "true"}}
		})
		devPreview(c, s, func() {
			if devtools.KeyValueInput(c, rows, devtools.KeyValueInputOptions{Label: "Headers"}).Changed() {
				log("pairs edited")
			}
		})
	})
}
