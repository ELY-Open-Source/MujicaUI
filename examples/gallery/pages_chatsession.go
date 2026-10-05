package main

const chatsessionCategory = "Chat Sessions"

var chatsessionPages = []page{
	{ID: "241", Name: "ConversationList", ZhName: "会话列表", Category: chatsessionCategory, Demo: chatDemoConversationList,
		Keys: "Up/Down choose a conversation, skipping group headers; typing jumps; Enter opens.", Code: `var state chat.ConversationListState
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
ui.Render(view, 400, 300, 1)`},
	{ID: "242", Name: "ConversationItem", ZhName: "会话条目", Category: chatsessionCategory, Demo: chatDemoConversationItem,
		Keys: "Enter opens; F2 renames in place (Enter keeps, Escape cancels); Delete asks to delete; the ellipsis opens the menu.", Code: `conv := chat.ChatConversation{ID: "1", Title: "Court etiquette", Updated: time.Now()}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	v := chat.ConversationItem(c, &conv, chat.ConversationItemOptions{})
	if v.Deleted {
		log.Println("delete", conv.ID)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "243", Name: "ConversationSearch", ZhName: "会话搜索", Category: chatsessionCategory, Demo: chatDemoConversationSearch,
		Keys: "Type to search; Down moves into the results; Up/Down move; Enter opens; Escape clears.", Code: `query := "seat"
hits := []chat.ConversationSearchHit{{ID: "1", Title: "Ballroom seating plan"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	v := chat.ConversationSearch(c, &query, hits, chat.ConversationSearchOptions{})
	if v.QueryChanged {
		log.Println("search", query)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "244", Name: "ProjectList", ZhName: "项目列表", Category: chatsessionCategory, Demo: chatDemoProjectList,
		Keys: "Up/Down move through folders, projects and conversations; Right/Left show or hide a project's conversations; the ellipsis menu moves a conversation.", Code: `var state chat.ProjectListState
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
ui.Render(view, 400, 300, 1)`},
	{ID: "245", Name: "ProjectKnowledgePanel", ZhName: "项目知识面板", Category: chatsessionCategory, Demo: chatDemoProjectKnowledgePanel,
		Keys: "Tab reaches Add files and the file list; Up/Down move between remove buttons; Enter removes.", Code: `files := []chat.KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 12000}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	v := chat.ProjectKnowledgePanel(c, files, chat.ProjectKnowledgePanelOptions{Capacity: 200000})
	if v.Add {
		log.Println("pick files")
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "246", Name: "WelcomeScreen", ZhName: "欢迎页", Category: chatsessionCategory, Demo: chatDemoWelcomeScreen,
		Keys: "Tab moves through the cards' examples and the prompt slot; Enter picks.", Code: `cards := []chat.CapabilityCard{{Title: "Write", Examples: []string{"Draft a decree"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.WelcomeScreen(c, chat.WelcomeScreenOptions{Greeting: "Good evening", Cards: cards}); v.Example != "" {
		log.Println("ask", v.Example)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "247", Name: "SuggestionChips", ZhName: "建议标签", Category: chatsessionCategory, Demo: chatDemoSuggestionChips,
		Keys: "Tab enters the chips; Left/Right move; Enter or Space picks.", Code: `input := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.SuggestionChips(c, &input, []string{"Plan a ball", "Write a sonnet"}, chat.SuggestionChipsOptions{})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "248", Name: "PromptLibrary", ZhName: "提示模板库", Category: chatsessionCategory, Demo: chatDemoPromptLibrary,
		Keys: "Tab moves between categories, search, templates and fields; Up/Down move in a column; Enter picks.", Code: `var state chat.PromptLibraryState
templates := []chat.PromptTemplate{{ID: "1", Category: "Letters", Title: "Invitation", Body: "Invite {{guest}}."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.PromptLibrary(c, &state, templates, chat.PromptLibraryOptions{}); v.Inserted != "" {
		log.Println("insert", v.Inserted)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "249", Name: "CapabilityCards", ZhName: "能力卡片", Category: chatsessionCategory, Demo: chatDemoCapabilityCards,
		Keys: "Tab moves between cards; Up/Down move among a card's examples; Enter picks.", Code: `cards := []chat.CapabilityCard{{Title: "Plan", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.CapabilityCards(c, cards, chat.CapabilityCardsOptions{}); v.Example != "" {
		log.Println("ask", v.Example)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "250", Name: "FeedbackForm", ZhName: "反馈表单", Category: chatsessionCategory, Demo: chatDemoFeedbackForm,
		Keys: "Left/Right move among reasons, Space toggles; Escape cancels.", Code: `var value chat.FeedbackValue
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.FeedbackForm(c, &value, chat.FeedbackFormOptions{Rating: chat.FeedbackDown}).Submitted {
		log.Println("reasons", value.Reasons)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "251", Name: "DateSeparator", ZhName: "日期分隔", Category: chatsessionCategory, Demo: chatDemoDateSeparator,
		Keys: "Not focusable; read as the date.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.DateSeparator(c, core.DateOf(time.Now()), chat.DateSeparatorOptions{})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "252", Name: "ThinkingIndicator", ZhName: "思考指示", Category: chatsessionCategory, Demo: chatDemoThinkingIndicator,
		Keys: "Not focusable; read as a status.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.ThinkingIndicator(c, chat.ThinkingIndicatorOptions{Phase: "Searching the archives"})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "253", Name: "StopGeneratingButton", ZhName: "停止生成按钮", Category: chatsessionCategory, Demo: chatDemoStopGeneratingButton,
		Keys: "Enter or Space stops; Cmd+. (Ctrl+. elsewhere) stops from anywhere.", Code: `generating := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.StopGeneratingButton(c, chat.StopGeneratingButtonOptions{Generating: generating}).Changed() {
		generating = false
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "254", Name: "SystemPromptEditor", ZhName: "系统提示编辑器", Category: chatsessionCategory, Demo: chatDemoSystemPromptEditor,
		Keys: "Cmd+S (Ctrl+S elsewhere) saves.", Code: `prompt, saved := "You are helpful.", "You are helpful."
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.SystemPromptEditor(c, &prompt, chat.SystemPromptEditorOptions{Default: "You are helpful.", Saved: saved}).Save {
		saved = prompt
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "255", Name: "CostEstimator", ZhName: "费用估算", Category: chatsessionCategory, Demo: chatDemoCostEstimator,
		Keys: "Not focusable; read as the estimated cost.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.CostEstimator(c, chat.CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "256", Name: "SharedConversationView", ZhName: "分享会话视图", Category: chatsessionCategory, Demo: chatDemoSharedConversationView,
		Keys: "Tab reaches Copy link and the messages; message text is selectable.", Code: `msgs := []chat.SessionMessage{{Role: chat.SessionUser, Text: "Hello"}, {Role: chat.SessionAssistant, Text: "Good day."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.SharedConversationView(c, msgs, chat.SharedConversationViewOptions{Title: "Greeting", SharedBy: "Ann"}).Changed() {
		c.WriteClipboard("https://example.com/share/1")
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "257", Name: "ConversationExport", ZhName: "会话导出", Category: chatsessionCategory, Demo: chatDemoConversationExport,
		Keys: "Arrows choose a format and range; Tab reaches Export.", Code: `var state chat.ConversationExportState
msgs := []chat.SessionMessage{{Role: chat.SessionUser, Text: "Hello"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.ConversationExport(c, &state, msgs, chat.ConversationExportOptions{Title: "Greeting"}); v.Exported {
		log.Println(len(v.Content), "bytes to save")
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "258", Name: "RegenerateMenu", ZhName: "重新生成菜单", Category: chatsessionCategory, Demo: chatDemoRegenerateMenu,
		Keys: "Enter, Space or Down open; Up/Down move; Right opens the model list; Escape closes.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.RegenerateMenu(c, chat.RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"}); v.Chosen {
		log.Println("retry", v.Choice.Kind, v.Choice.Model, v.Choice.Variant)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "259", Name: "ArtifactCard", ZhName: "产物卡片", Category: chatsessionCategory, Demo: chatDemoArtifactCard,
		Keys: "Enter or Space opens the artifact panel.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.ArtifactCard(c, chat.ArtifactCardOptions{Title: "seating.go", Detail: "Go · 120 lines"}).Changed() {
		log.Println("open artifact panel")
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "260", Name: "ChatEmptyState", ZhName: "会话空状态", Category: chatsessionCategory, Demo: chatDemoChatEmptyState,
		Keys: "Tab reaches the action; Enter runs it.", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.ChatEmptyState(c, chat.ChatEmptyStateOptions{Kind: chat.ChatEmptyOffline}).Action {
		log.Println("reconnect")
	}
}
ui.Render(view, 400, 300, 1)`},
}
