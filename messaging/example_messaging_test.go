package messaging_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/messaging"
)

func ExampleChannelList() {
	groups := []messaging.ChannelGroup{{ID: "ch", Title: "Channels", Channels: []messaging.Channel{
		{ID: "general", Name: "general", Unread: 3},
		{ID: "rehearsal", Name: "rehearsal", Pinned: true},
	}}}
	selected := "general"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if messaging.ChannelList(c, &selected, groups, messaging.ChannelListOptions{}).Changed() {
			_ = selected // show that channel
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChatMessage() {
	prev := messaging.ChatEntry{ID: "1", Author: "Uika", Time: time.Now(), Text: "Rehearsal at nine."}
	m := messaging.ChatEntry{ID: "2", Author: "Uika", Time: time.Now(), Text: "Bring the masks.", Edited: true}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.ChatMessage(c, prev, messaging.ChatMessageOptions{})
		messaging.ChatMessage(c, m, messaging.ChatMessageOptions{
			Continued: messaging.ChatContinues(prev, m),
			Reactions: func() { display.Tag(c, "👍 2", display.TagOptions{}) },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUnreadDivider() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.UnreadDivider(c, 3)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReadReceipt() {
	readers := []display.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.ReadReceipt(c, messaging.ReceiptRead, messaging.ReadReceiptOptions{Readers: readers})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMemberList() {
	var s messaging.MemberListState
	members := []messaging.Member{
		{ID: "1", Name: "Sakiko", Presence: messaging.PresenceOnline, Role: "Leader"},
		{ID: "2", Name: "Umiri", Presence: messaging.PresenceOffline},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := messaging.MemberList(c, &s, members, messaging.MemberListOptions{}); r.Opened != "" {
			_ = r.Opened // show that member's profile
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUserProfileCard() {
	profile := messaging.UserProfile{Name: "Sakiko Togawa", Presence: messaging.PresenceBusy, Role: "Keyboard", Status: "In rehearsal"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.UserProfileCard(c, profile, messaging.UserProfileCardOptions{Actions: func() {
			input.Button(c, "Message", input.ButtonOptions{})
		}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOnlineStatus() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.OnlineStatus(c, messaging.PresenceAway, messaging.OnlineStatusOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatusSetter() {
	status := messaging.CustomStatus{Emoji: "🎻", Text: "Rehearsing"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if messaging.StatusSetter(c, &status, messaging.StatusSetterOptions{}).Changed() {
			_ = status // publish it; clear it once status.Until passes
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCallControls() {
	st := messaging.CallState{}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.CallControls(c, st, messaging.CallControlsOptions{
			OnMute:   func(on bool) { st.Muted = on },
			OnCamera: func(off bool) { st.CameraOff = off },
			OnShare:  func(on bool) { st.Sharing = on },
			OnHangUp: func() {}, // leave the call
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVideoCallGrid() {
	people := []messaging.CallParticipant{
		{ID: "1", Name: "Sakiko", Speaking: true}, // Frame: the latest decoded frame
		{ID: "2", Name: "Uika", Muted: true},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.VideoCallGrid(c, people, messaging.VideoCallGridOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIncomingCallDialog() {
	ringing := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.IncomingCallDialog(c, &ringing, messaging.IncomingCall{Caller: "Sakiko", Video: true}, messaging.IncomingCallOptions{
			OnAccept:  func(video bool) {}, // join, with the camera on when video
			OnDecline: func() {},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEmojiPicker() {
	var recent []string
	tone := 0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := messaging.EmojiPicker(c, messaging.EmojiPickerOptions{Recent: &recent, Tone: &tone}); r.Picked != "" {
			_ = r.Picked // insert it
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePinnedMessages() {
	pins := []messaging.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := messaging.PinnedMessages(c, pins, messaging.PinnedMessagesOptions{})
		if r.Unpinned != "" {
			pins = pins[:0] // unpin it
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThreadPanel() {
	open := true
	root := messaging.ChatEntry{ID: "r", Author: "Sakiko", Text: "Set list for Friday?"}
	var replies []messaging.ChatEntry
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := messaging.ThreadPanel(c, &open, root, replies, messaging.ThreadPanelOptions{}); r.Sent != "" {
			replies = append(replies, messaging.ChatEntry{ID: r.Sent, Author: "Me", Text: r.Sent, Own: true})
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailList() {
	var s data.ListState[string]
	mails := []messaging.MailItem{{ID: "1", From: "Sakiko", Subject: "Second act", Unread: true, Attachments: 1}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := messaging.MailList(c, &s, mails, messaging.MailListOptions{}); r.Starred != "" {
			mails[0].Starred = !mails[0].Starred
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailReader() {
	m := messaging.Mail{From: "Sakiko", Subject: "Second act", Body: "The score is attached.\n> Can you send it?",
		Attachments: []messaging.MailAttachment{{Name: "score.pdf", Size: 2048}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := messaging.MailReader(c, m); r.Link != "" {
			c.OpenURL(r.Link)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailComposer() {
	draft := messaging.MailDraft{To: []string{"uika@example.com"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := messaging.MailComposer(c, &draft, messaging.MailComposerOptions{
			OnSaveDraft: func(d messaging.MailDraft) {}, // store d
		})
		if r.Sent {
			draft = messaging.MailDraft{}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRecipientInput() {
	var to []string
	contacts := []messaging.Contact{{Name: "Uika Misumi", Email: "uika@example.com"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.RecipientInput(c, &to, messaging.RecipientInputOptions{Contacts: contacts, Label: "To"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuotedText() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		messaging.QuotedText(c, "> Can you send the score?", &open)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSnoozePicker() {
	var until time.Time
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if messaging.SnoozePicker(c, &until, messaging.SnoozePickerOptions{}).Changed() {
			_ = until // hide the mail until then
		}
	}
	ui.Render(view, 400, 300, 1)
}
