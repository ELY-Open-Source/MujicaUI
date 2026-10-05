package chat_test

import (
	"log"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chat"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

func ExampleConversationList() {
	var state chat.ConversationListState
	convs := []chat.ChatConversation{
		{ID: "1", Title: "Court etiquette", Updated: time.Now()},
		{ID: "2", Title: "Royal letters", Updated: time.Now().AddDate(0, 0, -30), Pinned: true},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.ConversationList(c, &state, convs, chat.ConversationListOptions{Label: "History"}, nil).Changed() {
			log.Println("open", state.Selected)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationItem() {
	conv := chat.ChatConversation{ID: "1", Title: "Court etiquette", Updated: time.Now()}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		v := chat.ConversationItem(c, &conv, chat.ConversationItemOptions{})
		if v.Deleted {
			log.Println("delete", conv.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationSearch() {
	query := "seat"
	hits := []chat.ConversationSearchHit{{ID: "1", Title: "Ballroom seating plan"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		v := chat.ConversationSearch(c, &query, hits, chat.ConversationSearchOptions{})
		if v.QueryChanged {
			log.Println("search", query)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProjectList() {
	var state chat.ProjectListState
	projects := []chat.ChatProject{{ID: "p1", Name: "Coronation", Folder: "Court"}}
	convs := func(id string) []chat.ChatConversation {
		return []chat.ChatConversation{{ID: "1", Title: "Menu"}}
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.ProjectList(c, &state, projects, convs, chat.ProjectListOptions{}); v.Moved {
			log.Println("move", v.MoveID, "to", v.MoveTo)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProjectKnowledgePanel() {
	files := []chat.KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 12000}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		v := chat.ProjectKnowledgePanel(c, files, chat.ProjectKnowledgePanelOptions{Capacity: 200000})
		if v.Add {
			log.Println("pick files")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWelcomeScreen() {
	cards := []chat.CapabilityCard{{Title: "Write", Examples: []string{"Draft a decree"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.WelcomeScreen(c, chat.WelcomeScreenOptions{Greeting: "Good evening", Cards: cards}); v.Example != "" {
			log.Println("ask", v.Example)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSuggestionChips() {
	input := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.SuggestionChips(c, &input, []string{"Plan a ball", "Write a sonnet"}, chat.SuggestionChipsOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePromptLibrary() {
	var state chat.PromptLibraryState
	templates := []chat.PromptTemplate{{ID: "1", Category: "Letters", Title: "Invitation", Body: "Invite {{guest}}."}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.PromptLibrary(c, &state, templates, chat.PromptLibraryOptions{}); v.Inserted != "" {
			log.Println("insert", v.Inserted)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCapabilityCards() {
	cards := []chat.CapabilityCard{{Title: "Plan", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.CapabilityCards(c, cards, chat.CapabilityCardsOptions{}); v.Example != "" {
			log.Println("ask", v.Example)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFeedbackForm() {
	var value chat.FeedbackValue
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.FeedbackForm(c, &value, chat.FeedbackFormOptions{Rating: chat.FeedbackDown}).Submitted {
			log.Println("reasons", value.Reasons)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDateSeparator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.DateSeparator(c, core.DateOf(time.Now()), chat.DateSeparatorOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThinkingIndicator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.ThinkingIndicator(c, chat.ThinkingIndicatorOptions{Phase: "Searching the archives"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStopGeneratingButton() {
	generating := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.StopGeneratingButton(c, chat.StopGeneratingButtonOptions{Generating: generating}).Changed() {
			generating = false
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSystemPromptEditor() {
	prompt, saved := "You are helpful.", "You are helpful."
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.SystemPromptEditor(c, &prompt, chat.SystemPromptEditorOptions{Default: "You are helpful.", Saved: saved}).Save {
			saved = prompt
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCostEstimator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.CostEstimator(c, chat.CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSharedConversationView() {
	msgs := []chat.SessionMessage{{Role: chat.SessionUser, Text: "Hello"}, {Role: chat.SessionAssistant, Text: "Good day."}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.SharedConversationView(c, msgs, chat.SharedConversationViewOptions{Title: "Greeting", SharedBy: "Ann"}).Changed() {
			c.WriteClipboard("https://example.com/share/1")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConversationExport() {
	var state chat.ConversationExportState
	msgs := []chat.SessionMessage{{Role: chat.SessionUser, Text: "Hello"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.ConversationExport(c, &state, msgs, chat.ConversationExportOptions{Title: "Greeting"}); v.Exported {
			log.Println(len(v.Content), "bytes to save")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRegenerateMenu() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := chat.RegenerateMenu(c, chat.RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"}); v.Chosen {
			log.Println("retry", v.Choice.Kind, v.Choice.Model, v.Choice.Variant)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleArtifactCard() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.ArtifactCard(c, chat.ArtifactCardOptions{Title: "seating.go", Detail: "Go · 120 lines"}).Changed() {
			log.Println("open artifact panel")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChatEmptyState() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.ChatEmptyState(c, chat.ChatEmptyStateOptions{Kind: chat.ChatEmptyOffline}).Action {
			log.Println("reconnect")
		}
	}
	ui.Render(view, 400, 300, 1)
}
