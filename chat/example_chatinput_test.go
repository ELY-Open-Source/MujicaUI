package chat_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func ExamplePromptComposer() {
	var prompt string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.PromptComposer(c, &prompt, chat.PromptComposerOptions{Placeholder: "Ask the herald"})
		if r.Submitted() {
			fmt.Println("send", prompt)
			prompt = ""
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAttachmentChip() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.AttachmentChip(c, chat.Attachment{Name: "guests.csv", Size: 4096, Uploading: true, Progress: 0.4})
		if r.Removed() {
			fmt.Println("cancel the upload")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePasteImagePreview() {
	pic := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="3"/>`))
	images := []chat.PastedImage{{ID: "p1", Source: pic, Name: "screenshot"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if id, ok := chat.PasteImagePreview(c, images).Removed(); ok {
			fmt.Println("remove", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDragDropOverlay() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		images := func(path string) bool { return strings.HasSuffix(path, ".png") }
		r := chat.DragDropOverlay(c, chat.DragDropOverlayOptions{Accept: images}, func() {
			ui.Text(c, "Drag images into the chat")
		})
		for _, path := range r.Accepted() {
			fmt.Println("attach", path)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVoiceInputButton() {
	var rec chat.VoiceInputOptions
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.VoiceInputButton(c, rec).Changed() {
			rec = chat.VoiceInputOptions{Recording: !rec.Recording, Started: time.Now()}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVoiceWaveform() {
	levels := []float32{0.2, 0.6, 0.9, 0.4, 0.1}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.VoiceWaveform(c, levels, chat.VoiceWaveformOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSlashCommandMenu() {
	var prompt string
	commands := []chat.SlashCommand{{Name: "summarize", Description: "Summarize the thread"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		composer := chat.PromptComposer(c, &prompt, chat.PromptComposerOptions{})
		if cmd, ok := chat.SlashCommandMenu(c, composer, &prompt, commands).Chosen(); ok {
			fmt.Println("run", cmd.Name)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextMentionMenu() {
	var prompt string
	var found []chat.ContextItem
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		composer := chat.PromptComposer(c, &prompt, chat.PromptComposerOptions{})
		r := chat.ContextMentionMenu(c, composer, &prompt, found, chat.ContextMentionOptions{})
		if r.QueryChanged() {
			found = []chat.ContextItem{{ID: "f1", Label: r.Query() + ".md", Kind: chat.ContextFile}}
		}
		if item, ok := r.Chosen(); ok {
			fmt.Println("attach", item.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextChips() {
	items := []chat.ContextItem{{ID: "f1", Label: "guests.csv", Kind: chat.ContextFile}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if id, ok := chat.ContextChips(c, items, chat.ContextChipsOptions{}).Removed(); ok {
			fmt.Println("detach", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleModelSelector() {
	model := "herald"
	models := []chat.ModelOption{{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision"}, ContextTokens: 200000}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.ModelSelector(c, &model, models, chat.ModelSelectorOptions{}).Changed() {
			fmt.Println("model", model)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleModeSelector() {
	mode := chat.ModeChat
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.ModeSelector(c, &mode, chat.ModeSelectorOptions{}).Changed() {
			fmt.Println("mode", mode)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolToggleMenu() {
	web := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		tools := []chat.ToolItem{{ID: "web", Label: "Web search", On: &web}}
		if id, ok := chat.ToolToggleMenu(c, tools).Changed(); ok {
			fmt.Println("toggled", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleParameterPanel() {
	temperature := 1.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		params := []chat.ModelParam{{Label: "Temperature", Value: &temperature, Min: 0, Max: 2, Step: 0.1, Default: 1}}
		chat.ParameterPanel(c, params)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTokenCounter() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.TokenCounter(c, 1234, 8000)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextWindowMeter() {
	segments := []chat.ContextSegment{{Label: "System", Tokens: 1200}, {Label: "History", Tokens: 58000}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.ContextWindowMeter(c, segments, 128000, chat.ContextWindowMeterOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSendButton() {
	streaming := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.SendButton(c, chat.SendButtonOptions{Streaming: streaming, Shortcut: "Enter"})
		if r.Sent() || r.Stopped() {
			streaming = r.Sent()
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCitationBadge() {
	var cites chat.CitationState
	src := chat.Source{Title: "The Palace of Versailles", URL: "https://example.com/versailles"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.CitationBadge(c, 1, src, &cites).Changed() {
			fmt.Println("show source", cites.Active)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSourceCard() {
	var cites chat.CitationState
	src := chat.Source{Title: "Hall of Mirrors", URL: "https://example.org/hall", Snippet: "Seventeen arches."}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.SourceCard(c, 1, src, &cites).Clicked() {
			fmt.Println("open", src.URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSourcesPanel() {
	var cites chat.CitationState
	sources := []chat.Source{{Title: "Hall of Mirrors", URL: "https://example.org/hall"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if n, ok := chat.SourcesPanel(c, sources, &cites).Opened(); ok {
			fmt.Println("open", sources[n-1].URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSearchProgress() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		stages := []chat.SearchStage{{Label: "Searching the web", Done: true}, {Label: "Reading sources"}}
		chat.SearchProgress(c, chat.SearchProgressOptions{Query: "Versailles", Read: 3, Total: 8, Stages: stages})
	}
	ui.Render(view, 400, 300, 1)
}
