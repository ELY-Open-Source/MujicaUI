package main

const messagingCategory = "Messaging & Mail"

var messagingPages = []page{
	{ID: "401", Name: "ChannelList", ZhName: "频道列表", Category: messagingCategory, Demo: channelListDemo,
		Keys: "Tab 进入；上下方向键切换行；Enter / Space 选择或折叠；分组标题上的左、右方向键折叠、展开。", Code: `groups := []messaging.ChannelGroup{{ID: "ch", Title: "Channels", Channels: []messaging.Channel{
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
ui.Render(view, 400, 300, 1)`},
	{ID: "402", Name: "ChatMessage", ZhName: "聊天消息", Category: messagingCategory, Demo: chatMessageDemo,
		Keys: "可选文字支持选择复制；插槽控件按 Tab 访问。", Code: `prev := messaging.ChatEntry{ID: "1", Author: "Uika", Time: time.Now(), Text: "Rehearsal at nine."}
m := messaging.ChatEntry{ID: "2", Author: "Uika", Time: time.Now(), Text: "Bring the masks.", Edited: true}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.ChatMessage(c, prev, messaging.ChatMessageOptions{})
	messaging.ChatMessage(c, m, messaging.ChatMessageOptions{
		Continued: messaging.ChatContinues(prev, m),
		Reactions: func() { display.Tag(c, "👍 2", display.TagOptions{}) },
	})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "403", Name: "UnreadDivider", ZhName: "未读分隔线", Category: messagingCategory, Demo: unreadDividerDemo,
		Keys: "展示组件保持自然焦点序列。", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.UnreadDivider(c, 3)
}
ui.Render(view, 400, 300, 1)`},
	{ID: "404", Name: "ReadReceipt", ZhName: "已读回执", Category: messagingCategory, Demo: readReceiptDemo,
		Keys: "展示组件保持自然焦点序列。", Code: `readers := []display.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.ReadReceipt(c, messaging.ReceiptRead, messaging.ReadReceiptOptions{Readers: readers})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "405", Name: "MemberList", ZhName: "成员列表", Category: messagingCategory, Demo: memberListDemo,
		Keys: "Tab 聚焦列表；上下方向键选择；Home / End 首尾；Enter 或双击打开成员。", Code: `var s messaging.MemberListState
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
ui.Render(view, 400, 300, 1)`},
	{ID: "406", Name: "UserProfileCard", ZhName: "用户资料卡", Category: messagingCategory, Demo: userProfileCardDemo,
		Keys: "Tab 访问操作插槽；Enter / Space 激活按钮。", Code: `profile := messaging.UserProfile{Name: "Sakiko Togawa", Presence: messaging.PresenceBusy, Role: "Keyboard", Status: "In rehearsal"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.UserProfileCard(c, profile, messaging.UserProfileCardOptions{Actions: func() {
		input.Button(c, "Message", input.ButtonOptions{})
	}})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "407", Name: "OnlineStatus", ZhName: "在线状态", Category: messagingCategory, Demo: onlineStatusDemo,
		Keys: "展示组件保持自然焦点序列。", Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.OnlineStatus(c, messaging.PresenceAway, messaging.OnlineStatusOptions{})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "408", Name: "StatusSetter", ZhName: "状态设置", Category: messagingCategory, Demo: statusSetterDemo,
		Keys: "Tab 访问表情、文字、期限和清除按钮；选择器用方向键移动、Enter 确认、Escape 关闭。", Code: `status := messaging.CustomStatus{Emoji: "🎻", Text: "Rehearsing"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if messaging.StatusSetter(c, &status, messaging.StatusSetterOptions{}).Changed() {
		_ = status // publish it; clear it once status.Until passes
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "409", Name: "CallControls", ZhName: "通话控制", Category: messagingCategory, Demo: callControlsDemo,
		Keys: "Tab 进入工具栏；左右方向键移动；Enter / Space 激活。Cmd/Ctrl+D 静音，Cmd/Ctrl+E 摄像头，Cmd/Ctrl+Shift+S 共享，Cmd/Ctrl+Shift+H 挂断。", Code: `st := messaging.CallState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.CallControls(c, st, messaging.CallControlsOptions{
		OnMute:   func(on bool) { st.Muted = on },
		OnCamera: func(off bool) { st.CameraOff = off },
		OnShare:  func(on bool) { st.Sharing = on },
		OnHangUp: func() {}, // leave the call
	})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "410", Name: "VideoCallGrid", ZhName: "视频通话宫格", Category: messagingCategory, Demo: videoCallGridDemo,
		Keys: "Tab 进入宫格；方向键移动参与者焦点并读取状态。", Code: `people := []messaging.CallParticipant{
	{ID: "1", Name: "Sakiko", Speaking: true}, // Frame: the latest decoded frame
	{ID: "2", Name: "Uika", Muted: true},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.VideoCallGrid(c, people, messaging.VideoCallGridOptions{})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "411", Name: "IncomingCallDialog", ZhName: "来电对话框", Category: messagingCategory, Demo: incomingCallDialogDemo,
		Keys: "默认焦点为拒绝；Tab 循环访问操作；Enter / Space 激活；Escape 拒绝；关闭后恢复焦点。", Code: `ringing := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.IncomingCallDialog(c, &ringing, messaging.IncomingCall{Caller: "Sakiko", Video: true}, messaging.IncomingCallOptions{
		OnAccept:  func(video bool) {}, // join, with the camera on when video
		OnDecline: func() {},
	})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "412", Name: "EmojiPicker", ZhName: "表情选择器", Category: messagingCategory, Demo: emojiPickerDemo,
		Keys: "Tab 访问搜索、肤色、分类、表情；方向键在组内移动；Enter / Space 选择。", Code: `var recent []string
tone := 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.EmojiPicker(c, messaging.EmojiPickerOptions{Recent: &recent, Tone: &tone}); r.Picked != "" {
		_ = r.Picked // insert it
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "413", Name: "PinnedMessages", ZhName: "置顶消息", Category: messagingCategory, Demo: pinnedMessagesDemo,
		Keys: "Tab 访问跳转与取消置顶；方向键移动；Enter / Space 激活。", Code: `pins := []messaging.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := messaging.PinnedMessages(c, pins, messaging.PinnedMessagesOptions{})
	if r.Unpinned != "" {
		pins = pins[:0] // unpin it
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "414", Name: "ThreadPanel", ZhName: "话题侧栏", Category: messagingCategory, Demo: threadPanelDemo,
		Keys: "Tab 访问关闭、回复和发送；Enter 发送回复；Escape 关闭面板。", Code: `open := true
root := messaging.ChatEntry{ID: "r", Author: "Sakiko", Text: "Set list for Friday?"}
var replies []messaging.ChatEntry
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.ThreadPanel(c, &open, root, replies, messaging.ThreadPanelOptions{}); r.Sent != "" {
		replies = append(replies, messaging.ChatEntry{ID: r.Sent, Author: "Me", Text: r.Sent, Own: true})
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "415", Name: "MailList", ZhName: "邮件列表", Category: messagingCategory, Demo: mailListDemo,
		Keys: "上下方向键、Home / End 移动选择；Enter 或双击打开；S 请求星标；Cmd/Ctrl+A 全选；Cmd/Ctrl 点击切换单项，Shift 点击选择范围。", Code: `var s data.ListState[string]
mails := []messaging.MailItem{{ID: "1", From: "Sakiko", Subject: "Second act", Unread: true, Attachments: 1}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.MailList(c, &s, mails, messaging.MailListOptions{}); r.Starred != "" {
		mails[0].Starred = !mails[0].Starred
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "416", Name: "MailReader", ZhName: "邮件阅读器", Category: messagingCategory, Demo: mailReaderDemo,
		Keys: "Tab 访问链接、引用开关和附件；Enter / Space 激活；正文支持选择复制。", Code: `m := messaging.Mail{From: "Sakiko", Subject: "Second act", Body: "The score is attached.\n> Can you send it?",
	Attachments: []messaging.MailAttachment{{Name: "score.pdf", Size: 2048}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.MailReader(c, m); r.Link != "" {
		c.OpenURL(r.Link)
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "417", Name: "MailComposer", ZhName: "邮件编辑器", Category: messagingCategory, Demo: mailComposerDemo,
		Keys: "Tab 访问字段与操作；Cmd/Ctrl+S 保存草稿；组件内 Cmd/Ctrl+Enter 发送；正文 Enter 换行。", Code: `draft := messaging.MailDraft{To: []string{"uika@example.com"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := messaging.MailComposer(c, &draft, messaging.MailComposerOptions{
		OnSaveDraft: func(d messaging.MailDraft) {}, // store d
	})
	if r.Sent {
		draft = messaging.MailDraft{}
	}
}
ui.Render(view, 400, 300, 1)`},
	{ID: "418", Name: "RecipientInput", ZhName: "收件人输入", Category: messagingCategory, Demo: recipientInputDemo,
		Keys: "Enter 或分隔符提交标签；空输入 Backspace 删除末项；上下方向键移动联系人建议，Enter 选择，Escape 收起建议。", Code: `var to []string
contacts := []messaging.Contact{{Name: "Uika Misumi", Email: "uika@example.com"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.RecipientInput(c, &to, messaging.RecipientInputOptions{Contacts: contacts, Label: "To"})
}
ui.Render(view, 400, 300, 1)`},
	{ID: "419", Name: "QuotedText", ZhName: "引用文字", Category: messagingCategory, Demo: quotedTextDemo,
		Keys: "Tab 聚焦开关；Enter / Space 展开或折叠；展开内容支持选择复制。", Code: `open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.QuotedText(c, "> Can you send the score?", &open)
}
ui.Render(view, 400, 300, 1)`},
	{ID: "420", Name: "SnoozePicker", ZhName: "稍后提醒", Category: messagingCategory, Demo: snoozePickerDemo,
		Keys: "Tab 进入；上下方向键选择快捷项；Enter / Space 激活；自定义时间沿用 DateTimePicker 的输入、确认和取消操作。", Code: `var until time.Time
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if messaging.SnoozePicker(c, &until, messaging.SnoozePickerOptions{}).Changed() {
		_ = until // hide the mail until then
	}
}
ui.Render(view, 400, 300, 1)`},
}
