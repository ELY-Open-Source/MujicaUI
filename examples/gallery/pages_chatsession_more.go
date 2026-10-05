package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func chatDemoWelcomeScreen(c *ui.Context, log func(string)) {
	type state struct {
		input           string
		cards, disabled bool
	}
	s := *ui.Local(c.Root(), "246", func() *state { return &state{cards: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.cards, "Capability cards", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		var cards []mujica.CapabilityCard
		if s.cards {
			cards = chatSampleCards()
		}
		ui.Column(c).Width(860).Children(func() {
			v := mujica.WelcomeScreen(c, mujica.WelcomeScreenOptions{Greeting: "Good afternoon", Subtitle: "How may I serve the court?", Cards: cards, Disabled: s.disabled,
				Prompts: func() {
					if mujica.SuggestionChips(c, &s.input, []string{"Plan a ball", "Write a sonnet"}, mujica.SuggestionChipsOptions{Disabled: s.disabled}).Chosen != "" {
						log(fmt.Sprintf("fill %q", s.input))
					}
				}})
			if v.Example != "" {
				log(fmt.Sprintf("ask %q", v.Example))
			}
		})
	})
}

func chatDemoSuggestionChips(c *ui.Context, log func(string)) {
	type state struct {
		input          string
		send, disabled bool
	}
	s := *ui.Local(c.Root(), "247", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.send, "Send directly", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		v := mujica.SuggestionChips(c, &s.input, []string{"Plan a ball", "Write a sonnet", "Explain the minuet"},
			mujica.SuggestionChipsOptions{Send: s.send, Label: "Suggestions", Disabled: s.disabled})
		switch {
		case v.Sent:
			log(fmt.Sprintf("send %q", v.Chosen))
		case v.Chosen != "":
			log(fmt.Sprintf("fill %q", v.Chosen))
		}
		ui.Column(c).Width(420).Children(func() {
			mujica.TextInput(c, &s.input, mujica.TextInputOptions{Label: "Message", Placeholder: "Message"})
		})
	})
}

func chatDemoPromptLibrary(c *ui.Context, log func(string)) {
	type state struct {
		lib      mujica.PromptLibraryState
		disabled bool
	}
	s := *ui.Local(c.Root(), "248", func() *state { return &state{} })
	templates := []mujica.PromptTemplate{
		{ID: "t1", Category: "Letters", Title: "Invitation", Body: "Write an invitation to {{guest}} for the {{event}}."},
		{ID: "t2", Category: "Letters", Title: "Thank-you note", Body: "Thank {{host}} for the evening."},
		{ID: "t3", Category: "Music", Title: "Recital program", Body: "List pieces for a {{instrument}} recital."},
		{ID: "t4", Category: "Planning", Title: "Seating chart", Body: "Seat {{count}} guests by rank."},
	}
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		before := s.lib.Selected
		v := mujica.PromptLibrary(c, &s.lib, templates, mujica.PromptLibraryOptions{Label: "Prompt library", Disabled: s.disabled})
		if s.lib.Selected != before {
			log("preview " + s.lib.Selected)
		}
		if v.Inserted != "" {
			log(fmt.Sprintf("insert %q", v.Inserted))
		}
	})
}

func chatDemoCapabilityCards(c *ui.Context, log func(string)) {
	type state struct {
		one, disabled bool
	}
	s := *ui.Local(c.Root(), "249", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.one, "One column", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		cols := 0
		if s.one {
			cols = 1
		}
		ui.Column(c).Width(860).Children(func() {
			if v := mujica.CapabilityCards(c, chatSampleCards(), mujica.CapabilityCardsOptions{Columns: cols, Disabled: s.disabled}); v.Example != "" {
				log(fmt.Sprintf("ask %q", v.Example))
			}
		})
	})
}

func chatDemoFeedbackForm(c *ui.Context, log func(string)) {
	type state struct {
		value            mujica.FeedbackValue
		down, submitting bool
		disabled         bool
		at               time.Time
	}
	s := *ui.Local(c.Root(), "250", func() *state { return &state{down: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.down, "Thumbs down", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		rating := mujica.FeedbackUp
		if s.down {
			rating = mujica.FeedbackDown
		}
		if s.submitting && choiceDue(c, s.at) {
			s.submitting, s.value = false, mujica.FeedbackValue{}
			log("feedback sent")
		}
		ui.Column(c).Width(480).Children(func() {
			v := mujica.FeedbackForm(c, &s.value, mujica.FeedbackFormOptions{Rating: rating, Submitting: s.submitting, Disabled: s.disabled})
			if v.Submitted {
				s.submitting, s.at = true, c.Now()
				log(fmt.Sprintf("submit reasons %q, %d characters", s.value.Reasons, len(s.value.Comment)))
			}
			if v.Cancelled {
				s.value = mujica.FeedbackValue{}
				log("cancel")
			}
		})
	})
}

func chatDemoDateSeparator(c *ui.Context, log func(string)) {
	today := mujica.DateOf(time.Now())
	choiceDemo(c, func() {
		ui.Text(c, "No parameters; the language follows the window.")
	}, func() {
		ui.Column(c).Width(520).Children(func() {
			for _, d := range []int{0, -1, -12} {
				mujica.DateSeparator(c, today.AddDays(d), mujica.DateSeparatorOptions{})
			}
		})
	})
}

func chatDemoThinkingIndicator(c *ui.Context, log func(string)) {
	type state struct{ phase int }
	s := *ui.Local(c.Root(), "252", func() *state { return &state{} })
	phases := []string{"", "Searching the archives", "Reading 3 letters", "Composing a reply"}
	choiceDemo(c, func() {
		if mujica.Button(c, "Next phase", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.phase = (s.phase + 1) % len(phases)
			log(fmt.Sprintf("phase %q", phases[s.phase]))
		}
	}, func() {
		mujica.ThinkingIndicator(c, mujica.ThinkingIndicatorOptions{Phase: phases[s.phase]})
	})
}

func chatDemoStopGeneratingButton(c *ui.Context, log func(string)) {
	type state struct {
		generating, disabled bool
		at                   time.Time
	}
	s := *ui.Local(c.Root(), "253", func() *state { return &state{} })
	choiceDemo(c, func() {
		if mujica.Button(c, "Generate", mujica.ButtonOptions{Disabled: s.generating}).Clicked() {
			s.generating, s.at = true, c.Now()
			log("generate")
		}
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		if s.generating && c.Now().Sub(s.at) > 8*time.Second {
			s.generating = false
			log("finished")
		} else if s.generating {
			c.After(time.Second)
		}
		if mujica.StopGeneratingButton(c, mujica.StopGeneratingButtonOptions{Generating: s.generating, Disabled: s.disabled}).Changed() {
			s.generating = false
			log("stop")
		}
		if !s.generating {
			ui.Text(c, "Idle: the button is hidden.").TextColor(mujica.Tokens(c).TextMuted)
		}
	})
}

func chatDemoSystemPromptEditor(c *ui.Context, log func(string)) {
	type state struct {
		value, saved    string
		limit, disabled bool
		saving          bool
		at              time.Time
	}
	s := *ui.Local(c.Root(), "254", func() *state {
		return &state{value: "You are a courtly assistant.", saved: "You are a courtly assistant."}
	})
	choiceDemo(c, func() {
		choiceFlag(c, &s.limit, "Limit 60", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		if s.saving && choiceDue(c, s.at) {
			s.saving, s.saved = false, s.value
			log("saved")
		}
		limit := 0
		if s.limit {
			limit = 60
		}
		ui.Column(c).Width(520).Children(func() {
			v := mujica.SystemPromptEditor(c, &s.value, mujica.SystemPromptEditorOptions{Default: "You are helpful.", Saved: s.saved,
				MaxLength: limit, Saving: s.saving, Disabled: s.disabled})
			if v.Save {
				s.saving, s.at = true, c.Now()
				log(fmt.Sprintf("save %d characters", len(s.value)))
			}
			if v.Restored {
				log("restore default")
			}
		})
	})
}

func chatDemoCostEstimator(c *ui.Context, log func(string)) {
	type state struct{ in, out, inPrice, outPrice float64 }
	s := *ui.Local(c.Root(), "255", func() *state { return &state{in: 12000, out: 3400, inPrice: 3, outPrice: 15} })
	choiceDemo(c, func() {
		for _, f := range []struct {
			v     *float64
			label string
		}{{&s.in, "Input tokens"}, {&s.out, "Output tokens"}, {&s.inPrice, "Input $/M"}, {&s.outPrice, "Output $/M"}} {
			choiceLabeled(c, f.label, func() {
				r := mujica.NumberInput(c, f.v, mujica.NumberInputOptions{Min: 0, Max: 1e7, Step: 100, Label: f.label})
				r.Element.Width(150)
				if r.Changed() {
					log(fmt.Sprintf("%s = %v", f.label, *f.v))
				}
			})
		}
	}, func() {
		ui.Column(c).Width(420).Children(func() {
			mujica.CostEstimator(c, mujica.CostEstimatorOptions{InputTokens: int(s.in), OutputTokens: int(s.out), InputPrice: s.inPrice, OutputPrice: s.outPrice})
		})
	})
}

func chatDemoSharedConversationView(c *ui.Context, log func(string)) {
	type state struct{ disabled bool }
	s := *ui.Local(c.Root(), "256", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Column(c).Width(600).Height(420).Children(func() {
			if mujica.SharedConversationView(c, chatSampleMessages(time.Now()), mujica.SharedConversationViewOptions{Title: "Banquet seating",
				SharedBy: "Lady Ann", SharedAt: time.Now(), Disabled: s.disabled}).Changed() {
				c.WriteClipboard("https://example.com/share/banquet")
				log("copy link")
			}
		})
	})
}

func chatDemoConversationExport(c *ui.Context, log func(string)) {
	type state struct {
		export   mujica.ConversationExportState
		disabled bool
	}
	s := *ui.Local(c.Root(), "257", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Column(c).Width(580).Children(func() {
			v := mujica.ConversationExport(c, &s.export, chatSampleMessages(time.Now()), mujica.ConversationExportOptions{Title: "Banquet seating",
				Selected: []int{1, 3}, Disabled: s.disabled})
			if v.Exported {
				log(fmt.Sprintf("export %d bytes as format %d (the caller writes the file)", len(v.Content), v.Format))
			}
		})
	})
}

func chatDemoRegenerateMenu(c *ui.Context, log func(string)) {
	type state struct{ disabled bool }
	s := *ui.Local(c.Root(), "258", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		v := mujica.RegenerateMenu(c, mujica.RegenerateMenuOptions{Models: []string{"Opus", "Sonnet", "Haiku"}, Current: "Opus", Disabled: s.disabled})
		if v.Chosen {
			log(fmt.Sprintf("retry kind %d model %q variant %q", v.Choice.Kind, v.Choice.Model, v.Choice.Variant))
		}
	})
}

func chatDemoArtifactCard(c *ui.Context, log func(string)) {
	type state struct{ streaming, disabled bool }
	s := *ui.Local(c.Root(), "259", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.streaming, "Streaming", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		for _, a := range []mujica.ArtifactCardOptions{
			{Kind: mujica.ArtifactCode, Title: "seating.go", Detail: "Go · 120 lines"},
			{Kind: mujica.ArtifactDocument, Title: "Royal decree", Detail: "2 pages"},
			{Kind: mujica.ArtifactWeb, Title: "Ball invitation"},
		} {
			a.Streaming, a.Disabled = s.streaming, s.disabled
			if mujica.ArtifactCard(c, a).Changed() {
				log("open " + a.Title)
			}
		}
	})
}

func chatDemoChatEmptyState(c *ui.Context, log func(string)) {
	type state struct{ kind int }
	s := *ui.Local(c.Root(), "260", func() *state { return &state{} })
	choiceDemo(c, func() {
		mujica.RadioGroup(c, &s.kind, []mujica.RadioItem[int]{{Value: 0, Label: "New"}, {Value: 1, Label: "No results"}, {Value: 2, Label: "Offline"}},
			mujica.RadioGroupOptions{Horizontal: true, Label: "Kind"})
	}, func() {
		ui.Column(c).Width(420).Children(func() {
			if mujica.ChatEmptyState(c, mujica.ChatEmptyStateOptions{Kind: mujica.ChatEmptyKind(s.kind), Query: "minuet"}).Action {
				log(fmt.Sprintf("action for kind %d", s.kind))
			}
		})
	})
}
