package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chat"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// chatHistory is the gallery's sample conversations, dated from now.
func chatHistory(now time.Time) []chat.ChatConversation {
	day := 24 * time.Hour
	return []chat.ChatConversation{
		{ID: "a", Title: "Court etiquette notes", Updated: now.Add(-2 * time.Hour)},
		{ID: "b", Title: "Ballroom seating plan", Updated: now.Add(-day)},
		{ID: "c", Title: "Harpsichord repair", Updated: now.Add(-4 * day)},
		{ID: "d", Title: "Coronation menu", Updated: now.Add(-40 * day)},
		{ID: "e", Title: "Royal letters", Updated: now.Add(-90 * day), Pinned: true},
		{ID: "f", Title: "Garden maze", Updated: now.Add(-time.Hour)},
	}
}

func chatSampleMessages(now time.Time) []chat.SessionMessage {
	return []chat.SessionMessage{
		{Role: chat.SessionUser, Text: "How should guests be seated at the banquet?", Time: now.Add(-10 * time.Minute)},
		{Role: chat.SessionAssistant, Text: "Seat by rank: the host at the head, honored guests to the right.", Time: now.Add(-9 * time.Minute)},
		{Role: chat.SessionUser, Author: "Lady Ann", Text: "And the musicians?", Time: now.Add(-5 * time.Minute)},
		{Role: chat.SessionAssistant, Text: "Place the quartet behind a screen near the gallery.", Time: now.Add(-4 * time.Minute)},
	}
}

func chatSampleCards() []chat.CapabilityCard {
	return []chat.CapabilityCard{
		{Title: "Write", Description: "Letters, speeches and verse.", Icon: icons.Must("file-text"), Examples: []string{"Draft a royal decree", "Write a sonnet"}},
		{Title: "Plan", Description: "Banquets, seating and schedules.", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball for 200 guests"}},
		{Title: "Research", Description: "History, heraldry and music.", Icon: icons.Must("search"), Examples: []string{"Explain the minuet"}},
	}
}

func chatDemoConversationList(c *ui.Context, log func(string)) {
	type state struct {
		list       chat.ConversationListState
		items      []chat.ChatConversation
		many, rich bool
		disabled   bool
	}
	s := *ui.Local(c.Root(), "241", func() *state { return &state{items: chatHistory(time.Now())} })
	choiceDemo(c, func() {
		if input.Switch(c, &s.many, "100,000 conversations", input.SwitchOptions{}).Changed() {
			s.items = chatHistory(time.Now())
			if s.many {
				now := time.Now()
				for i := range 100000 {
					s.items = append(s.items, chat.ChatConversation{ID: fmt.Sprint("n", i), Title: fmt.Sprintf("Archive %d", i), Updated: now.Add(-time.Duration(i) * 7 * time.Minute)})
				}
			}
			log(fmt.Sprintf("conversations = %d", len(s.items)))
		}
		choiceFlag(c, &s.rich, "ConversationItem rows", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		var row func(conv chat.ChatConversation)
		if s.rich {
			row = func(conv chat.ChatConversation) {
				v := chat.ConversationItem(c, &conv, chat.ConversationItemOptions{Selected: conv.ID == s.list.Selected})
				if v.Opened {
					s.list.Selected = conv.ID
					log("open " + conv.Title)
				}
			}
		}
		v := chat.ConversationList(c, &s.list, s.items, chat.ConversationListOptions{Label: "History", Disabled: s.disabled}, row)
		v.Element.Width(380).Height(380)
		if v.Changed() {
			log("select " + s.list.Selected)
		}
		if v.Submitted() {
			log("open " + s.list.Selected)
		}
	})
}

func chatDemoConversationItem(c *ui.Context, log func(string)) {
	type state struct {
		convs              []chat.ChatConversation
		selected, disabled bool
	}
	s := *ui.Local(c.Root(), "242", func() *state { return &state{convs: chatHistory(time.Now())[:3]} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.selected, "Selected", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Column(c).Width(420).Gap(4).Children(func() {
			for i := range s.convs {
				conv := &s.convs[i]
				ui.Row(c.Key(conv.ID)).Children(func() {
					v := chat.ConversationItem(c, conv, chat.ConversationItemOptions{Selected: s.selected && i == 0, Disabled: s.disabled})
					switch {
					case v.Opened:
						log("open " + conv.Title)
					case v.Renamed:
						log("rename to " + conv.Title)
					case v.PinToggled:
						conv.Pinned = !conv.Pinned
						log(fmt.Sprintf("pinned %s = %v", conv.Title, conv.Pinned))
					case v.Exported:
						log("export " + conv.Title)
					case v.Deleted:
						log("delete " + conv.Title)
					}
				})
			}
		})
	})
}

func chatDemoConversationSearch(c *ui.Context, log func(string)) {
	type state struct {
		query, asked string
		at           time.Time
		hits         []chat.ConversationSearchHit
		fail, slow   bool
		disabled     bool
	}
	s := *ui.Local(c.Root(), "243", func() *state { return &state{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.slow, "Slow server", log)
		choiceFlag(c, &s.fail, "Fail searches", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		loading := s.asked != s.query
		if loading && (!s.slow || choiceDue(c, s.at)) {
			s.asked, s.hits, loading = s.query, nil, false
			for _, conv := range chatHistory(time.Now()) {
				if q := strings.ToLower(s.query); q != "" && strings.Contains(strings.ToLower(conv.Title), q) {
					s.hits = append(s.hits, chat.ConversationSearchHit{ID: conv.ID, Title: conv.Title, Snippet: "…mentions " + strings.ToLower(conv.Title) + "…"})
				}
			}
			log(fmt.Sprintf("results for %q: %d", s.query, len(s.hits)))
		}
		errText := ""
		if s.fail && s.query != "" {
			errText = "The archive did not answer."
		}
		ui.Column(c).Width(420).Children(func() {
			v := chat.ConversationSearch(c, &s.query, s.hits, chat.ConversationSearchOptions{Loading: loading, Error: errText, Disabled: s.disabled, Label: "Search conversations"})
			if v.QueryChanged {
				s.at = c.Now()
				log(fmt.Sprintf("query %q", s.query))
			}
			if v.Chosen != "" {
				log("open " + v.Chosen)
			}
			if v.Retry {
				s.fail, s.asked = false, ""
				log("retry")
			}
		})
	})
}

func chatDemoProjectList(c *ui.Context, log func(string)) {
	type state struct {
		list     chat.ProjectListState
		owner    map[string]string
		disabled bool
	}
	s := *ui.Local(c.Root(), "244", func() *state {
		return &state{owner: map[string]string{"a": "p1", "b": "p1", "c": "p3", "f": "p2"}}
	})
	projects := []chat.ChatProject{{ID: "p1", Name: "Coronation", Folder: "Court"}, {ID: "p2", Name: "Gardens", Folder: "Court"}, {ID: "p3", Name: "Music"}}
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Row(c).Gap(24).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Width(380).Children(func() {
				convs := func(id string) []chat.ChatConversation {
					var out []chat.ChatConversation
					for _, conv := range chatHistory(time.Now()) {
						if s.owner[conv.ID] == id {
							out = append(out, conv)
						}
					}
					return out
				}
				v := chat.ProjectList(c, &s.list, projects, convs, chat.ProjectListOptions{Label: "Projects", Disabled: s.disabled})
				switch {
				case v.Moved:
					s.owner[v.MoveID] = v.MoveTo
					log(fmt.Sprintf("move %s to %s", v.MoveID, v.MoveTo))
				case v.Changed:
					log("project " + s.list.Selected)
				case v.Opened != "":
					log("open " + v.Opened)
				}
			})
			choiceLabeled(c, "Drag onto a project", func() {
				conv := chatHistory(time.Now())[3]
				v := chat.ConversationItem(c, &conv, chat.ConversationItemOptions{Draggable: true})
				v.Element.Width(300)
				if v.Opened {
					log("open " + conv.Title)
				}
			})
		})
	})
}

func chatDemoProjectKnowledgePanel(c *ui.Context, log func(string)) {
	type state struct {
		files    []chat.KnowledgeFile
		next     int
		small    bool
		disabled bool
	}
	s := *ui.Local(c.Root(), "245", func() *state {
		return &state{files: []chat.KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 120000}, {ID: "2", Name: "menu.md", Tokens: 30000}}}
	})
	choiceDemo(c, func() {
		choiceFlag(c, &s.small, "Capacity 150,000", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		capacity := 400000
		if s.small {
			capacity = 150000
		}
		ui.Column(c).Width(440).Children(func() {
			v := chat.ProjectKnowledgePanel(c, s.files, chat.ProjectKnowledgePanelOptions{Capacity: capacity, Disabled: s.disabled})
			if v.Add {
				s.next++
				s.files = append(s.files, chat.KnowledgeFile{ID: fmt.Sprint("n", s.next), Name: fmt.Sprintf("notes-%d.txt", s.next), Tokens: 15000})
				log("add (the demo invents a file)")
			}
			if v.Removed != "" {
				for i, f := range s.files {
					if f.ID == v.Removed {
						s.files = append(s.files[:i:i], s.files[i+1:]...)
						log("remove " + f.Name)
						break
					}
				}
			}
		})
	})
}
