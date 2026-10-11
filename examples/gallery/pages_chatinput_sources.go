package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func chatinputDemoTokenCounter(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "234", func() float64 { return 1234 })
	choiceDemo(c, func() {
		if input.NumberInput(c, s, input.NumberInputOptions{Min: 0, Max: 12000, Step: 100, Label: "Estimated tokens"}).Changed() {
			log(fmt.Sprint("estimate ", *s))
		}
	}, func() {
		chat.TokenCounter(c, int(*s), 8000)
		chat.TokenCounter(c, 7500, 8000)
		chat.TokenCounter(c, 9100, 8000)
	})
}

func chatinputDemoContextWindowMeter(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "235", func() float64 { return 58000 })
	choiceDemo(c, func() {
		if input.NumberInput(c, s, input.NumberInputOptions{Min: 0, Max: 150000, Step: 1000, Label: "History tokens"}).Changed() {
			log(fmt.Sprint("history ", *s))
		}
	}, func() {
		ui.Column(c).Width(520).Children(func() {
			chat.ContextWindowMeter(c, []chat.ContextSegment{{Label: "System", Tokens: 1200}, {Label: "History", Tokens: int(*s)}, {Label: "Attachments", Tokens: 24000}}, 128000, chat.ContextWindowMeterOptions{})
		})
	})
}

func chatinputDemoSendButton(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "236", func() chat.SendButtonOptions { return chat.SendButtonOptions{Shortcut: "Enter"} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.Streaming, "Streaming", log)
		choiceFlag(c, &s.Disabled, "Disabled send", log)
	}, func() {
		r := chat.SendButton(c, *s)
		if r.Sent() {
			s.Streaming = true
			log("send")
		}
		if r.Stopped() {
			s.Streaming = false
			log("stop")
		}
	})
}

func chatinputSources() []chat.Source {
	return []chat.Source{{Title: "The Palace of Versailles", URL: "https://example.com/versailles", Site: "example.com", Snippet: "Seat of the court from 1682.", Published: time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC)}, {Title: "Hall of Mirrors", URL: "https://example.org/hall", Site: "example.org", Snippet: "Seventeen arches face the garden."}, {Title: "Court etiquette", URL: "https://example.net/etiquette", Snippet: "The ranks of the court."}}
}

func chatinputDemoCitationBadge(c *ui.Context, log func(string)) {
	type state struct {
		cites    chat.CitationState
		disabled bool
	}
	s := ui.Local(c.Root(), "237", func() state { return state{} })
	choiceDemo(c, func() { choiceFlag(c, &s.disabled, "Disabled", log) }, func() {
		ui.Row(c).Disabled(s.disabled).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "The court moved in 1682")
			for i, src := range chatinputSources() {
				if chat.CitationBadge(c, i+1, src, &s.cites).Changed() {
					log(fmt.Sprint("locate source ", s.cites.Active))
				}
			}
		})
		ui.Column(c).Width(360).Height(260).Disabled(s.disabled).Children(func() {
			if n, ok := chat.SourcesPanel(c, chatinputSources(), &s.cites).Opened(); ok {
				log(fmt.Sprint("open source ", n))
			}
		})
	})
}

func chatinputDemoSourceCard(c *ui.Context, log func(string)) {
	type state struct {
		cites              chat.CitationState
		selected, disabled bool
	}
	s := ui.Local(c.Root(), "238", func() state { return state{} })
	choiceDemo(c, func() { choiceFlag(c, &s.selected, "Selected", log); choiceFlag(c, &s.disabled, "Disabled", log) }, func() {
		s.cites.Active = 0
		if s.selected {
			s.cites.Active = 1
		}
		ui.Column(c).Width(400).Disabled(s.disabled).Children(func() {
			if chat.SourceCard(c, 1, chatinputSources()[0], &s.cites).Clicked() {
				log("open source 1")
			}
		})
	})
}

func chatinputDemoSourcesPanel(c *ui.Context, log func(string)) {
	type state struct {
		cites                 chat.CitationState
		sources               []chat.Source
		disabled, empty, many bool
	}
	s := ui.Local(c.Root(), "239", func() state { return state{sources: chatinputSources()} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.empty, "Empty", log)
		if input.Switch(c, &s.many, "100,000 sources", input.SwitchOptions{}).Changed() {
			s.sources = chatinputSources()
			s.cites.Active = 0
			if s.many {
				s.sources = make([]chat.Source, 100000)
				for i := range s.sources {
					s.sources[i] = chat.Source{Title: fmt.Sprintf("Archive %d", i+1), URL: "https://example.com/archive"}
				}
			}
			log(fmt.Sprint("sources ", len(s.sources)))
		}
		if input.Button(c, "Locate last source", input.ButtonOptions{Disabled: s.empty}).Clicked() {
			s.cites.Active = len(s.sources)
			log("locate last source")
		}
	}, func() {
		sources := s.sources
		if s.empty {
			sources = nil
			s.cites.Active = 0
		}
		ui.Column(c).Width(400).Height(320).Disabled(s.disabled).Children(func() {
			if n, ok := chat.SourcesPanel(c, sources, &s.cites).Opened(); ok {
				log(fmt.Sprint("open source ", n))
			}
		})
	})
}

func chatinputDemoSearchProgress(c *ui.Context, log func(string)) {
	type state struct {
		done, disabled bool
		read           float64
	}
	s := ui.Local(c.Root(), "240", func() state { return state{read: 3} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.done, "Done", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
		if input.NumberInput(c, &s.read, input.NumberInputOptions{Min: 0, Max: 8, Step: 1, Label: "Read sources"}).Changed() {
			log(fmt.Sprint("read ", s.read))
		}
	}, func() {
		ui.Column(c).Width(440).Disabled(s.disabled).Children(func() {
			if chat.SearchProgress(c, chat.SearchProgressOptions{Query: "Versailles 1682", Read: int(s.read), Total: 8, Done: s.done, Stages: []chat.SearchStage{{Label: "Planning queries", Done: true}, {Label: "Searching the web", Done: true}, {Label: "Reading sources", Done: s.done}}}).Changed() {
				log("toggle search details")
			}
		})
	})
}
