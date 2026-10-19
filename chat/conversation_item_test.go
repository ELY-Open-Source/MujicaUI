package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func conversationItemTester(conv *ChatConversation, opts ConversationItemOptions, got *ConversationItemView) *ui.Tester {
	return chatTester(420, 360, func(c *ui.Context) {
		opts.Today = core.DateOf(chatNow)
		v := ConversationItem(c, conv, opts)
		if v.Opened || v.Deleted || v.Exported || v.PinToggled || v.Renamed {
			*got = v
		}
	})
}

func TestConversationItemShowsTitleAndTime(t *testing.T) {
	convs := chatConversations()
	var got ConversationItemView
	tt := conversationItemTester(&convs[0], ConversationItemOptions{}, &got)
	if !tt.HasText("Court etiquette notes") || !tt.HasText("13:30") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt = conversationItemTester(&convs[1], ConversationItemOptions{}, &got)
	if !tt.HasText("Yesterday") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Ballroom seating plan"); err != nil {
		t.Fatal(err)
	}
	if !got.Opened {
		t.Fatal("click did not open")
	}
}

func TestConversationItemMenuActions(t *testing.T) {
	for _, tc := range []struct {
		item string
		want func(v ConversationItemView) bool
	}{
		{"Pin", func(v ConversationItemView) bool { return v.PinToggled }},
		{"Export", func(v ConversationItemView) bool { return v.Exported }},
		{"Delete", func(v ConversationItemView) bool { return v.Deleted }},
	} {
		conv := chatConversations()[0]
		var got ConversationItemView
		tt := conversationItemTester(&conv, ConversationItemOptions{}, &got)
		if err := tt.Click("Conversation actions"); err != nil {
			t.Fatal(err)
		}
		if err := tt.Click(tc.item); err != nil {
			t.Fatal(err)
		}
		if !tc.want(got) {
			t.Errorf("%s: %+v", tc.item, got)
		}
	}
}

func TestConversationItemInlineRename(t *testing.T) {
	conv := chatConversations()[0]
	var got ConversationItemView
	tt := conversationItemTester(&conv, ConversationItemOptions{}, &got)
	_ = tt.Click("Conversation actions")
	_ = tt.Click("Rename")
	if !tt.Focused("Rename") {
		t.Fatal("rename field not focused")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Type(" II")
	tt.Key(0, ui.KeyEnter)
	if conv.Title != "Court etiquette notes II" || !got.Renamed {
		t.Fatalf("title %q renamed %v", conv.Title, got.Renamed)
	}
	// F2 starts a rename; Escape keeps the title; a blank title is refused.
	_ = tt.Click("Court etiquette notes II")
	tt.Key(0, ui.KeyF2)
	tt.Type("xyz")
	tt.Key(0, ui.KeyEscape)
	if conv.Title != "Court etiquette notes II" {
		t.Fatalf("Escape changed the title to %q", conv.Title)
	}
	tt.Key(0, ui.KeyF2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyEnter)
	if conv.Title != "Court etiquette notes II" {
		t.Fatalf("blank rename gave %q", conv.Title)
	}
}

func TestConversationItemDisabled(t *testing.T) {
	conv := chatConversations()[0]
	var got ConversationItemView
	tt := conversationItemTester(&conv, ConversationItemOptions{Disabled: true}, &got)
	_ = tt.Click("Court etiquette notes")
	_ = tt.Click("Conversation actions")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	tt.Key(0, ui.KeyDelete)
	if got.Opened || got.Deleted || tt.HasText("Rename") {
		t.Fatalf("disabled item acted: %+v", got)
	}
}

func TestConversationItemSnapshot(t *testing.T) {
	snapshot(t, "conversation_item", 420, 140, func(c *ui.Context) {
		convs := chatConversations()
		ui.Column(c).Padding(16).Gap(4).Children(func() {
			ConversationItem(c, &convs[0], ConversationItemOptions{Today: core.DateOf(chatNow), Selected: true})
			ConversationItem(c, &convs[4], ConversationItemOptions{Today: core.DateOf(chatNow)})
			ConversationItem(c, &convs[3], ConversationItemOptions{Today: core.DateOf(chatNow), Disabled: true})
		})
	})
}
