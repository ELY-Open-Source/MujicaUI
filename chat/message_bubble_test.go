package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
)

func TestMessageBubbleRolesAlignAndLimit(t *testing.T) {
	var user, bot, sys ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(600).Children(func() {
			MessageBubble(c, MessageUser, MessageBubbleOptions{Name: "Ada", MaxWidth: 200}, func() {
				user = ui.Text(c, "A long user message that must wrap inside the bubble's width limit.").Bounds()
			})
			MessageBubble(c, MessageAssistant, MessageBubbleOptions{Avatar: func() { display.Avatar(c, "Bot", display.AvatarOptions{}) }}, func() {
				bot = ui.Text(c, "Hello").Bounds()
			})
			MessageBubble(c, MessageSystem, MessageBubbleOptions{}, func() { sys = ui.Text(c, "Joined").Bounds() })
			MessageBubble(c, MessageError, MessageBubbleOptions{}, func() { ui.Text(c, "Timeout") })
		})
	}, 600, 400)
	tt.Frame()
	if user.W > 200 || user.X < 300 {
		t.Errorf("user bubble %+v not right-aligned within 200", user)
	}
	if bot.X > 100 || sys.X < 200 || sys.X > 300 {
		t.Errorf("assistant at %v, system at %v", bot.X, sys.X)
	}
	for _, s := range []string{"Ada", "Bot", "Error"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q", s)
		}
	}
}

func TestMessageBubbleBadOptions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a negative MaxWidth")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		MessageBubble(c, MessageUser, MessageBubbleOptions{MaxWidth: -1}, func() {})
	}, 100, 100)
}

func TestMessageBubbleSnapshot(t *testing.T) {
	snapshot(t, "message_bubble", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			MessageBubble(c, MessageUser, MessageBubbleOptions{Name: "Marie"}, func() { ui.Text(c, "Is the ball tonight?") })
			MessageBubble(c, MessageAssistant, MessageBubbleOptions{Name: "Herald", Avatar: func() { display.Avatar(c, "Herald", display.AvatarOptions{}) }},
				func() { ui.Text(c, "Yes, at nine, in the Hall of Mirrors.") })
			MessageBubble(c, MessageSystem, MessageBubbleOptions{}, func() { ui.Text(c, "Herald joined the chat") })
			MessageBubble(c, MessageError, MessageBubbleOptions{}, func() { ui.Text(c, "The messenger did not return.") })
		})
	})
}
