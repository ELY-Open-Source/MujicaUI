package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func devJsonViewerDemo(c *ui.Context, log func(string)) {
	dtSection(c, "JsonViewer", func() {
		s := devControls(c, "453", log, false)
		type state struct {
			tree    mujica.JsonViewerState
			invalid bool
		}
		st := ui.Local(c.Root(), "453", func() state { return state{} })
		if mujica.Checkbox(c, &st.invalid, "Invalid JSON", mujica.CheckboxOptions{}).Changed() {
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
			if mujica.JsonViewer(c, data, &st.tree, mujica.JsonViewerOptions{}).Copied != "" {
				log("path copied")
			}
		})
	})
}

func devJsonEditorDemo(c *ui.Context, log func(string)) {
	dtSection(c, "JsonEditor", func() {
		s := devControls(c, "454", log, false)
		text := ui.Local(c.Root(), "454", func() string { return `{"port":8080,"ready":true}` })
		if mujica.Button(c, "Invalid sample", mujica.ButtonOptions{}).Clicked() {
			*text = "{\n  \"port\": [1,2,]\n}"
			log("invalid sample loaded")
		}
		devPreview(c, s, func() {
			r := mujica.JsonEditor(c, text, mujica.JsonEditorOptions{})
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
			tree  mujica.SchemaTreeState
			roots []mujica.SchemaNode
		}
		st := ui.Local(c.Root(), "455", func() state {
			return state{roots: []mujica.SchemaNode{{Name: "app", Children: []mujica.SchemaNode{{Name: "users", Kind: mujica.SchemaTable, Children: []mujica.SchemaNode{{Name: "id", Kind: mujica.SchemaColumn, Detail: "bigint"}, {Name: "name", Kind: mujica.SchemaColumn, Detail: "text"}, {Name: "users_pk", Kind: mujica.SchemaIndex, Detail: "(id)"}}}, {Name: "active_users", Kind: mujica.SchemaView}}}}}
		})
		devPreview(c, s, func() {
			roots := st.roots
			if s.empty {
				roots = nil
			}
			if mujica.SchemaTree(c, &st.tree, roots, mujica.SchemaTreeOptions{}).Changed() {
				log("schema node selected")
			}
		})
	})
}

func devApiRequestBuilderDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ApiRequestBuilder", func() {
		s := devControls(c, "456", log, false)
		request := ui.Local(c.Root(), "456", func() mujica.ApiRequest {
			return mujica.ApiRequest{URL: "https://example.com/api", Headers: []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
		})
		sending := ui.Local(c.Root(), "456-sending", func() bool { return false })
		if mujica.Checkbox(c, sending, "Sending", mujica.CheckboxOptions{}).Changed() {
			log("send state changed")
		}
		devPreview(c, s, func() {
			r := mujica.ApiRequestBuilder(c, request, mujica.ApiRequestBuilderOptions{Sending: *sending})
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
		if mujica.Checkbox(c, text, "Text body", mujica.CheckboxOptions{}).Changed() {
			log("response format changed")
		}
		devPreview(c, s, func() {
			resp := &mujica.ApiResponse{Status: 200, Duration: 182 * time.Millisecond, Size: 2048, Body: `{"ready":true,"requests":1248}`, Headers: []mujica.KeyValue{{Key: "Content-Type", Value: "application/json"}, {Key: "X-Trace", Value: "demo"}}}
			if *text {
				resp.Body = "Response received successfully."
			}
			if s.empty {
				resp = nil
			}
			if mujica.ResponseViewer(c, resp, mujica.ResponseViewerOptions{Height: 260, Status: s.status, Error: "Connection refused"}).Changed() {
				log("retry requested")
			}
		})
	})
}

func devKeyValueInputDemo(c *ui.Context, log func(string)) {
	dtSection(c, "KeyValueInput", func() {
		s := devControls(c, "458", log, false)
		rows := ui.Local(c.Root(), "458", func() []mujica.KeyValue {
			return []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}, {Key: "debug", Value: "true"}}
		})
		devPreview(c, s, func() {
			if mujica.KeyValueInput(c, rows, mujica.KeyValueInputOptions{Label: "Headers"}).Changed() {
				log("pairs edited")
			}
		})
	})
}
