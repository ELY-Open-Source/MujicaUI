package mujica_test

import (
	"log"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

func ExampleConversationList() {
	var state mujica.ConversationListState
	convs := []mujica.ChatConversation{
		{ID: "1", Title: "Court etiquette", Updated: time.Now()},
		{ID: "2", Title: "Royal letters", Updated: time.Now().AddDate(0, 0, -30), Pinned: true},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ConversationList(c, &state, convs, mujica.ConversationListOptions{Label: "History"}, nil).Changed() {
			log.Println("open", state.Selected)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationItem() {
	conv := mujica.ChatConversation{ID: "1", Title: "Court etiquette", Updated: time.Now()}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		v := mujica.ConversationItem(c, &conv, mujica.ConversationItemOptions{})
		if v.Deleted {
			log.Println("delete", conv.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationSearch() {
	query := "seat"
	hits := []mujica.ConversationSearchHit{{ID: "1", Title: "Ballroom seating plan"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		v := mujica.ConversationSearch(c, &query, hits, mujica.ConversationSearchOptions{})
		if v.QueryChanged {
			log.Println("search", query)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProjectList() {
	var state mujica.ProjectListState
	projects := []mujica.ChatProject{{ID: "p1", Name: "Coronation", Folder: "Court"}}
	convs := func(id string) []mujica.ChatConversation {
		return []mujica.ChatConversation{{ID: "1", Title: "Menu"}}
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.ProjectList(c, &state, projects, convs, mujica.ProjectListOptions{}); v.Moved {
			log.Println("move", v.MoveID, "to", v.MoveTo)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProjectKnowledgePanel() {
	files := []mujica.KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 12000}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		v := mujica.ProjectKnowledgePanel(c, files, mujica.ProjectKnowledgePanelOptions{Capacity: 200000})
		if v.Add {
			log.Println("pick files")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWelcomeScreen() {
	cards := []mujica.CapabilityCard{{Title: "Write", Examples: []string{"Draft a decree"}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.WelcomeScreen(c, mujica.WelcomeScreenOptions{Greeting: "Good evening", Cards: cards}); v.Example != "" {
			log.Println("ask", v.Example)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSuggestionChips() {
	input := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.SuggestionChips(c, &input, []string{"Plan a ball", "Write a sonnet"}, mujica.SuggestionChipsOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePromptLibrary() {
	var state mujica.PromptLibraryState
	templates := []mujica.PromptTemplate{{ID: "1", Category: "Letters", Title: "Invitation", Body: "Invite {{guest}}."}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.PromptLibrary(c, &state, templates, mujica.PromptLibraryOptions{}); v.Inserted != "" {
			log.Println("insert", v.Inserted)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCapabilityCards() {
	cards := []mujica.CapabilityCard{{Title: "Plan", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball"}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.CapabilityCards(c, cards, mujica.CapabilityCardsOptions{}); v.Example != "" {
			log.Println("ask", v.Example)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFeedbackForm() {
	var value mujica.FeedbackValue
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.FeedbackForm(c, &value, mujica.FeedbackFormOptions{Rating: mujica.FeedbackDown}).Submitted {
			log.Println("reasons", value.Reasons)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDateSeparator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DateSeparator(c, mujica.DateOf(time.Now()), mujica.DateSeparatorOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThinkingIndicator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ThinkingIndicator(c, mujica.ThinkingIndicatorOptions{Phase: "Searching the archives"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStopGeneratingButton() {
	generating := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.StopGeneratingButton(c, mujica.StopGeneratingButtonOptions{Generating: generating}).Changed() {
			generating = false
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSystemPromptEditor() {
	prompt, saved := "You are helpful.", "You are helpful."
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SystemPromptEditor(c, &prompt, mujica.SystemPromptEditorOptions{Default: "You are helpful.", Saved: saved}).Save {
			saved = prompt
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCostEstimator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CostEstimator(c, mujica.CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSharedConversationView() {
	msgs := []mujica.SessionMessage{{Role: mujica.SessionUser, Text: "Hello"}, {Role: mujica.SessionAssistant, Text: "Good day."}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SharedConversationView(c, msgs, mujica.SharedConversationViewOptions{Title: "Greeting", SharedBy: "Ann"}).Changed() {
			c.WriteClipboard("https://example.com/share/1")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationExport() {
	var state mujica.ConversationExportState
	msgs := []mujica.SessionMessage{{Role: mujica.SessionUser, Text: "Hello"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.ConversationExport(c, &state, msgs, mujica.ConversationExportOptions{Title: "Greeting"}); v.Exported {
			log.Println(len(v.Content), "bytes to save")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRegenerateMenu() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if v := mujica.RegenerateMenu(c, mujica.RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"}); v.Chosen {
			log.Println("retry", v.Choice.Kind, v.Choice.Model, v.Choice.Variant)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleArtifactCard() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ArtifactCard(c, mujica.ArtifactCardOptions{Title: "seating.go", Detail: "Go · 120 lines"}).Changed() {
			log.Println("open artifact panel")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChatEmptyState() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ChatEmptyState(c, mujica.ChatEmptyStateOptions{Kind: mujica.ChatEmptyOffline}).Action {
			log.Println("reconnect")
		}
	}
	ui.Render(view, 400, 300, 1)
}
