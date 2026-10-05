package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
	"strings"
)

func chatinputDemoSlashCommandMenu(c *ui.Context, log func(string)) {
	type state struct {
		text            string
		disabled, empty bool
	}
	s := ui.Local(c.Root(), "227", func() state { return state{text: "/"} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.empty, "Empty commands", log)
	}, func() {
		ui.Column(c).Width(520).Children(func() {
			p := mujica.PromptComposer(c, &s.text, mujica.PromptComposerOptions{Label: "Command prompt", Disabled: s.disabled})
			commands := []mujica.SlashCommand{{Name: "summarize", Description: "Summarize the thread"}, {Name: "search", Description: "Search the web"}, {Name: "translate", Description: "Translate the message"}}
			if s.empty {
				commands = nil
			}
			before := s.text
			if cmd, ok := mujica.SlashCommandMenu(c, p, &s.text, commands).Chosen(); ok {
				log("command " + cmd.Name)
			} else if before != s.text {
				log("complete command")
			}
			if p.Changed() {
				log("edit command query")
			}
			if p.Submitted() {
				log("submit prompt")
			}
		})
	})
}

func chatinputContextItems() []mujica.ContextItem {
	return []mujica.ContextItem{{ID: "file", Label: "guests.csv", Kind: mujica.ContextFile}, {ID: "doc", Label: "Court protocol", Kind: mujica.ContextDoc}, {ID: "web", Label: "Versailles", Detail: "example.com", Kind: mujica.ContextWeb}}
}

func chatinputDemoContextMentionMenu(c *ui.Context, log func(string)) {
	type state struct {
		text, query       string
		items             []mujica.ContextItem
		loading, disabled bool
	}
	s := ui.Local(c.Root(), "228", func() state { return state{text: "@", items: chatinputContextItems()} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.loading, "Loading results", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Column(c).Width(520).Children(func() {
			p := mujica.PromptComposer(c, &s.text, mujica.PromptComposerOptions{Label: "Context prompt", Disabled: s.disabled})
			items := s.items
			if s.loading {
				items = nil
			}
			r := mujica.ContextMentionMenu(c, p, &s.text, items, mujica.ContextMentionOptions{Loading: s.loading})
			if r.QueryChanged() {
				s.query = r.Query()
				s.items = nil
				for _, it := range chatinputContextItems() {
					if strings.Contains(strings.ToLower(it.Label), strings.ToLower(s.query)) {
						s.items = append(s.items, it)
					}
				}
				c.Invalidate()
				log("context query changed")
			}
			if it, ok := r.Chosen(); ok {
				log("attach context " + it.ID)
			}
			if p.Changed() {
				log("edit context prompt")
			}
			if p.Submitted() {
				log("submit prompt")
			}
		})
	})
}

func chatinputDemoContextChips(c *ui.Context, log func(string)) {
	type state struct {
		items    []mujica.ContextItem
		max      float64
		disabled bool
	}
	s := ui.Local(c.Root(), "229", func() state { return state{items: chatinputContextItems(), max: 2} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		if mujica.NumberInput(c, &s.max, mujica.NumberInputOptions{Min: 1, Max: 5, Step: 1, Label: "Visible chips"}).Changed() {
			log(fmt.Sprint("max ", s.max))
		}
		if mujica.Button(c, "Restore context", mujica.ButtonOptions{}).Clicked() {
			s.items = chatinputContextItems()
			log("restore context")
		}
	}, func() {
		ui.Column(c).Width(520).Disabled(s.disabled).Children(func() {
			if id, ok := mujica.ContextChips(c, s.items, mujica.ContextChipsOptions{Max: int(s.max)}).Removed(); ok {
				for i, it := range s.items {
					if it.ID == id {
						s.items = append(s.items[:i], s.items[i+1:]...)
						break
					}
				}
				log("remove context " + id)
			}
		})
	})
}

func chatinputDemoModelSelector(c *ui.Context, log func(string)) {
	type state struct {
		model                 string
		disabled, unavailable bool
	}
	s := ui.Local(c.Root(), "230", func() state { return state{model: "herald", unavailable: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.unavailable, "Restricted model", log)
	}, func() {
		models := []mujica.ModelOption{{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision", "Tools"}, ContextTokens: 200000}, {ID: "page", Name: "Page Mini", Capabilities: []string{"Fast"}, ContextTokens: 32000}, {ID: "oracle", Name: "Oracle", ContextTokens: 1000000}}
		if s.unavailable {
			models[2].Unavailable = "Needs the Court plan"
		}
		if mujica.ModelSelector(c, &s.model, models, mujica.ModelSelectorOptions{Disabled: s.disabled}).Changed() {
			log("model " + s.model)
		}
	})
}

func chatinputDemoModeSelector(c *ui.Context, log func(string)) {
	type state struct {
		mode            mujica.ChatMode
		icons, disabled bool
	}
	s := ui.Local(c.Root(), "231", func() state { return state{} })
	choiceDemo(c, func() { choiceFlag(c, &s.icons, "Icons only", log); choiceFlag(c, &s.disabled, "Disabled", log) }, func() {
		if mujica.ModeSelector(c, &s.mode, mujica.ModeSelectorOptions{IconsOnly: s.icons, Disabled: s.disabled}).Changed() {
			log(fmt.Sprint("mode ", s.mode))
		}
	})
}

func chatinputDemoToolToggleMenu(c *ui.Context, log func(string)) {
	type state struct{ web, code, mail, authorized, disabled bool }
	s := ui.Local(c.Root(), "232", func() state { return state{web: true} })
	choiceDemo(c, func() { choiceFlag(c, &s.authorized, "Authorized", log); choiceFlag(c, &s.disabled, "Disabled", log) }, func() {
		ui.Column(c).Disabled(s.disabled).Children(func() {
			r := mujica.ToolToggleMenu(c, []mujica.ToolItem{{ID: "web", Label: "Web search", On: &s.web}, {ID: "code", Label: "Code execution", On: &s.code}, {ID: "mail", Label: "Mail", On: &s.mail, NeedsAuth: true, Authorized: s.authorized}})
			if id, ok := r.Changed(); ok {
				log("toggle " + id)
			}
			if id, ok := r.Authorize(); ok {
				s.authorized = true
				log("authorize " + id)
			}
		})
	})
}

func chatinputDemoParameterPanel(c *ui.Context, log func(string)) {
	type state struct {
		temp, top, length float64
		disabled          bool
	}
	s := ui.Local(c.Root(), "233", func() state { return state{temp: 0.7, top: 1, length: 2048} })
	choiceDemo(c, func() { choiceFlag(c, &s.disabled, "Disabled", log) }, func() {
		ui.Column(c).Width(420).Disabled(s.disabled).Children(func() {
			if mujica.ParameterPanel(c, []mujica.ModelParam{{Label: "Temperature", Value: &s.temp, Min: 0, Max: 2, Step: 0.1, Default: 1}, {Label: "Top-p", Value: &s.top, Min: 0, Max: 1, Step: 0.05, Default: 1}, {Label: "Max tokens", Value: &s.length, Min: 256, Max: 8192, Step: 256, Default: 2048}}).Changed() {
				log(fmt.Sprintf("parameters %.1f / %.2f / %.0f", s.temp, s.top, s.length))
			}
		})
	})
}
