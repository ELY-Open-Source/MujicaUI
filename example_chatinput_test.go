package mujica_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExamplePromptComposer() {
	var prompt string
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.PromptComposer(c, &prompt, mujica.PromptComposerOptions{Placeholder: "Ask the herald"})
		if r.Submitted() {
			fmt.Println("send", prompt)
			prompt = ""
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAttachmentChip() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.AttachmentChip(c, mujica.Attachment{Name: "guests.csv", Size: 4096, Uploading: true, Progress: 0.4})
		if r.Removed() {
			fmt.Println("cancel the upload")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePasteImagePreview() {
	pic := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="3"/>`))
	images := []mujica.PastedImage{{ID: "p1", Source: pic, Name: "screenshot"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if id, ok := mujica.PasteImagePreview(c, images).Removed(); ok {
			fmt.Println("remove", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDragDropOverlay() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		images := func(path string) bool { return strings.HasSuffix(path, ".png") }
		r := mujica.DragDropOverlay(c, mujica.DragDropOverlayOptions{Accept: images}, func() {
			ui.Text(c, "Drag images into the chat")
		})
		for _, path := range r.Accepted() {
			fmt.Println("attach", path)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVoiceInputButton() {
	var rec mujica.VoiceInputOptions
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.VoiceInputButton(c, rec).Changed() {
			rec = mujica.VoiceInputOptions{Recording: !rec.Recording, Started: time.Now()}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVoiceWaveform() {
	levels := []float32{0.2, 0.6, 0.9, 0.4, 0.1}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VoiceWaveform(c, levels, mujica.VoiceWaveformOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSlashCommandMenu() {
	var prompt string
	commands := []mujica.SlashCommand{{Name: "summarize", Description: "Summarize the thread"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		composer := mujica.PromptComposer(c, &prompt, mujica.PromptComposerOptions{})
		if cmd, ok := mujica.SlashCommandMenu(c, composer, &prompt, commands).Chosen(); ok {
			fmt.Println("run", cmd.Name)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextMentionMenu() {
	var prompt string
	var found []mujica.ContextItem
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		composer := mujica.PromptComposer(c, &prompt, mujica.PromptComposerOptions{})
		r := mujica.ContextMentionMenu(c, composer, &prompt, found, mujica.ContextMentionOptions{})
		if r.QueryChanged() {
			found = []mujica.ContextItem{{ID: "f1", Label: r.Query() + ".md", Kind: mujica.ContextFile}}
		}
		if item, ok := r.Chosen(); ok {
			fmt.Println("attach", item.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextChips() {
	items := []mujica.ContextItem{{ID: "f1", Label: "guests.csv", Kind: mujica.ContextFile}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if id, ok := mujica.ContextChips(c, items, mujica.ContextChipsOptions{}).Removed(); ok {
			fmt.Println("detach", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleModelSelector() {
	model := "herald"
	models := []mujica.ModelOption{{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision"}, ContextTokens: 200000}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ModelSelector(c, &model, models, mujica.ModelSelectorOptions{}).Changed() {
			fmt.Println("model", model)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleModeSelector() {
	mode := mujica.ModeChat
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ModeSelector(c, &mode, mujica.ModeSelectorOptions{}).Changed() {
			fmt.Println("mode", mode)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolToggleMenu() {
	web := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		tools := []mujica.ToolItem{{ID: "web", Label: "Web search", On: &web}}
		if id, ok := mujica.ToolToggleMenu(c, tools).Changed(); ok {
			fmt.Println("toggled", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleParameterPanel() {
	temperature := 1.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		params := []mujica.ModelParam{{Label: "Temperature", Value: &temperature, Min: 0, Max: 2, Step: 0.1, Default: 1}}
		mujica.ParameterPanel(c, params)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTokenCounter() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TokenCounter(c, 1234, 8000)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextWindowMeter() {
	segments := []mujica.ContextSegment{{Label: "System", Tokens: 1200}, {Label: "History", Tokens: 58000}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ContextWindowMeter(c, segments, 128000, mujica.ContextWindowMeterOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSendButton() {
	streaming := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SendButton(c, mujica.SendButtonOptions{Streaming: streaming, Shortcut: "Enter"})
		if r.Sent() || r.Stopped() {
			streaming = r.Sent()
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCitationBadge() {
	var cites mujica.CitationState
	src := mujica.Source{Title: "The Palace of Versailles", URL: "https://example.com/versailles"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.CitationBadge(c, 1, src, &cites).Changed() {
			fmt.Println("show source", cites.Active)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSourceCard() {
	var cites mujica.CitationState
	src := mujica.Source{Title: "Hall of Mirrors", URL: "https://example.org/hall", Snippet: "Seventeen arches."}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SourceCard(c, 1, src, &cites).Clicked() {
			fmt.Println("open", src.URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSourcesPanel() {
	var cites mujica.CitationState
	sources := []mujica.Source{{Title: "Hall of Mirrors", URL: "https://example.org/hall"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if n, ok := mujica.SourcesPanel(c, sources, &cites).Opened(); ok {
			fmt.Println("open", sources[n-1].URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSearchProgress() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		stages := []mujica.SearchStage{{Label: "Searching the web", Done: true}, {Label: "Reading sources"}}
		mujica.SearchProgress(c, mujica.SearchProgressOptions{Query: "Versailles", Read: 3, Total: 8, Stages: stages})
	}
	ui.Render(view, 400, 300, 1)
}
