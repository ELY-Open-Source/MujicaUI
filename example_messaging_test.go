package mujica_test

import (
	mujica "mujicaui"
	"time"

	"github.com/egoist/mygo/ui"
)

func ExampleChannelList() {
	groups := []mujica.ChannelGroup{{ID: "ch", Title: "Channels", Channels: []mujica.Channel{
		{ID: "general", Name: "general", Unread: 3},
		{ID: "rehearsal", Name: "rehearsal", Pinned: true},
	}}}
	selected := "general"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ChannelList(c, &selected, groups, mujica.ChannelListOptions{}).Changed() {
			_ = selected // show that channel
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChatMessage() {
	prev := mujica.ChatEntry{ID: "1", Author: "Uika", Time: time.Now(), Text: "Rehearsal at nine."}
	m := mujica.ChatEntry{ID: "2", Author: "Uika", Time: time.Now(), Text: "Bring the masks.", Edited: true}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ChatMessage(c, prev, mujica.ChatMessageOptions{})
		mujica.ChatMessage(c, m, mujica.ChatMessageOptions{
			Continued: mujica.ChatContinues(prev, m),
			Reactions: func() { mujica.Tag(c, "👍 2", mujica.TagOptions{}) },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUnreadDivider() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.UnreadDivider(c, 3)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReadReceipt() {
	readers := []mujica.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ReadReceipt(c, mujica.ReceiptRead, mujica.ReadReceiptOptions{Readers: readers})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMemberList() {
	var s mujica.MemberListState
	members := []mujica.Member{
		{ID: "1", Name: "Sakiko", Presence: mujica.PresenceOnline, Role: "Leader"},
		{ID: "2", Name: "Umiri", Presence: mujica.PresenceOffline},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if r := mujica.MemberList(c, &s, members, mujica.MemberListOptions{}); r.Opened != "" {
			_ = r.Opened // show that member's profile
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUserProfileCard() {
	profile := mujica.UserProfile{Name: "Sakiko Togawa", Presence: mujica.PresenceBusy, Role: "Keyboard", Status: "In rehearsal"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.UserProfileCard(c, profile, mujica.UserProfileCardOptions{Actions: func() {
			mujica.Button(c, "Message", mujica.ButtonOptions{})
		}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOnlineStatus() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.OnlineStatus(c, mujica.PresenceAway, mujica.OnlineStatusOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatusSetter() {
	status := mujica.CustomStatus{Emoji: "🎻", Text: "Rehearsing"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.StatusSetter(c, &status, mujica.StatusSetterOptions{}).Changed() {
			_ = status // publish it; clear it once status.Until passes
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCallControls() {
	st := mujica.CallState{}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CallControls(c, st, mujica.CallControlsOptions{
			OnMute:   func(on bool) { st.Muted = on },
			OnCamera: func(off bool) { st.CameraOff = off },
			OnShare:  func(on bool) { st.Sharing = on },
			OnHangUp: func() {}, // leave the call
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVideoCallGrid() {
	people := []mujica.CallParticipant{
		{ID: "1", Name: "Sakiko", Speaking: true}, // Frame: the latest decoded frame
		{ID: "2", Name: "Uika", Muted: true},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VideoCallGrid(c, people, mujica.VideoCallGridOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIncomingCallDialog() {
	ringing := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.IncomingCallDialog(c, &ringing, mujica.IncomingCall{Caller: "Sakiko", Video: true}, mujica.IncomingCallOptions{
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
		mujica.Use(c, mujica.Settings{})
		if r := mujica.EmojiPicker(c, mujica.EmojiPickerOptions{Recent: &recent, Tone: &tone}); r.Picked != "" {
			_ = r.Picked // insert it
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePinnedMessages() {
	pins := []mujica.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.PinnedMessages(c, pins, mujica.PinnedMessagesOptions{})
		if r.Unpinned != "" {
			pins = pins[:0] // unpin it
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThreadPanel() {
	open := true
	root := mujica.ChatEntry{ID: "r", Author: "Sakiko", Text: "Set list for Friday?"}
	var replies []mujica.ChatEntry
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if r := mujica.ThreadPanel(c, &open, root, replies, mujica.ThreadPanelOptions{}); r.Sent != "" {
			replies = append(replies, mujica.ChatEntry{ID: r.Sent, Author: "Me", Text: r.Sent, Own: true})
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailList() {
	var s mujica.ListState[string]
	mails := []mujica.MailItem{{ID: "1", From: "Sakiko", Subject: "Second act", Unread: true, Attachments: 1}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if r := mujica.MailList(c, &s, mails, mujica.MailListOptions{}); r.Starred != "" {
			mails[0].Starred = !mails[0].Starred
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailReader() {
	m := mujica.Mail{From: "Sakiko", Subject: "Second act", Body: "The score is attached.\n> Can you send it?",
		Attachments: []mujica.MailAttachment{{Name: "score.pdf", Size: 2048}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if r := mujica.MailReader(c, m); r.Link != "" {
			c.OpenURL(r.Link)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMailComposer() {
	draft := mujica.MailDraft{To: []string{"uika@example.com"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.MailComposer(c, &draft, mujica.MailComposerOptions{
			OnSaveDraft: func(d mujica.MailDraft) {}, // store d
		})
		if r.Sent {
			draft = mujica.MailDraft{}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRecipientInput() {
	var to []string
	contacts := []mujica.Contact{{Name: "Uika Misumi", Email: "uika@example.com"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.RecipientInput(c, &to, mujica.RecipientInputOptions{Contacts: contacts, Label: "To"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuotedText() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.QuotedText(c, "> Can you send the score?", &open)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSnoozePicker() {
	var until time.Time
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SnoozePicker(c, &until, mujica.SnoozePickerOptions{}).Changed() {
			_ = until // hide the mail until then
		}
	}
	ui.Render(view, 400, 300, 1)
}
