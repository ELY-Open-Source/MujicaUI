package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

// chatContainerApp is a chat with a sidebar and an input.
type chatContainerApp struct {
	msgs  []string
	state MessageListState
	draft string
}

func (a *chatContainerApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ChatContainer(c, &a.state, len(a.msgs), ChatContainerOptions{
		List:    MessageListOptions{ID: func(i int) string { return a.msgs[i] }},
		Sidebar: func() { ui.Text(c, "Sessions").Padding(12) },
		Input:   func() { input.TextInput(c, &a.draft, input.TextInputOptions{Label: "Prompt"}) },
	}, func(i int) {
		MessageBubble(c, MessageAssistant, MessageBubbleOptions{}, func() { ui.Text(c, a.msgs[i]) })
	})
}

func TestChatContainerSlotsAndFollow(t *testing.T) {
	a := &chatContainerApp{}
	for i := range 40 {
		a.msgs = append(a.msgs, "reply "+string(rune('A'+i%26))+string(rune('a'+i/26)))
	}
	tt := ui.NewTester(a.view, 600, 400)
	tt.Frame()
	if !tt.HasText("Sessions") || !tt.HasText("reply Nb") {
		t.Fatalf("slots or newest message missing: %q", tt.Texts())
	}
	_ = tt.Click("Prompt")
	tt.Type("hi")
	if a.draft != "hi" {
		t.Fatalf("input slot draft %q", a.draft)
	}
	a.msgs = append(a.msgs, "fresh")
	tt.Frame()
	tt.Frame()
	if !tt.HasText("fresh") {
		t.Fatal("new message not followed")
	}
	tt.Scroll(400, 150, 0, -3000)
	tt.Frame()
	a.msgs = append(a.msgs, "later")
	tt.Frame()
	if tt.HasText("later") {
		t.Fatal("followed after the user scrolled up")
	}
}

func TestChatContainerSnapshot(t *testing.T) {
	snapshot(t, "chat_container", 560, 300, func(c *ui.Context) {
		a := ui.Local(c.Root(), "app", func() *chatContainerApp {
			return &chatContainerApp{msgs: []string{"Bonjour, Madame.", "How may I serve the court today?"}}
		})
		(*a).view(c)
	})
}
