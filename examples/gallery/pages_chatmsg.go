package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

const chatmsgCategory = "Chat Messages"

var chatmsgPages = append(append(chatmsgListPages, chatmsgBlockPages...), chatmsgActionPages...)

// chatmsgReplies are the lines the demo assistant sends.
var chatmsgReplies = []string{
	"The ball begins at nine, in the Hall of Mirrors.",
	"Guests arrive by the **east door**; carriages wait in the *lower court*.",
	"Seating follows rank: dukes nearest the throne, then the envoys.",
	"Shall I send the invitations tonight?",
}

// chatmsgDemo is a chat whose messages the demo adds.
type chatmsgDemo struct {
	ids     []string
	texts   []string
	list    mujica.MessageListState
	next    int
	sidebar bool
	dated   bool
}

func newChatmsgDemo(n int) *chatmsgDemo {
	d := &chatmsgDemo{sidebar: true, dated: true}
	for range n {
		d.add()
	}
	return d
}

func (d *chatmsgDemo) add() {
	d.ids = append(d.ids, fmt.Sprintf("msg-%d", d.next))
	d.texts = append(d.texts, chatmsgReplies[d.next%len(chatmsgReplies)])
	d.next++
}

// day puts three messages on each day, ending today.
func (d *chatmsgDemo) day(i int) time.Time {
	last := (len(d.ids) - 1) / 3
	return time.Now().AddDate(0, 0, i/3-last)
}

// row shows message i as a bubble, users and the assistant taking turns.
func (d *chatmsgDemo) row(c *ui.Context, i int) {
	role := mujica.MessageAssistant
	if i%2 == 0 {
		role = mujica.MessageUser
	}
	mujica.MessageBubble(c, role, mujica.MessageBubbleOptions{}, func() { mujica.MarkdownView(c, d.texts[i]) })
}

// chatmsgControls adds messages to a demo chat and logs it.
func chatmsgControls(c *ui.Context, log func(string), d *chatmsgDemo) {
	ui.Row(c).Gap(8).Children(func() {
		if mujica.Button(c, "Add message", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			d.add()
			log(fmt.Sprintf("message %s added", d.ids[len(d.ids)-1]))
		}
		if mujica.Button(c, "Add 10", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			for range 10 {
				d.add()
			}
			log(fmt.Sprintf("10 messages added, %d in all", len(d.ids)))
		}
		if mujica.Button(c, "Jump to latest", mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
			d.list.ScrollToEnd()
			log("scrolled to end")
		}
	})
}

var chatmsgListPages = []page{
	{ID: "201", Name: "ChatContainer", ZhName: "对话容器", Category: chatmsgCategory,
		Keys: "Tab 在侧栏、消息与输入区之间移动；消息区滚轮与 PageUp/PageDown 滚动；上翻后出现回到最新按钮。",
		Code: `mujica.ChatContainer(c, &chat, len(msgs), mujica.ChatContainerOptions{
	List:  mujica.MessageListOptions{ID: func(i int) string { return msgs[i] }},
	Input: func() { ui.Text(c, "Composer goes here") },
}, func(i int) { ui.Text(c, msgs[i]) })`,
		Demo: func(c *ui.Context, log func(string)) {
			d := *ui.Local(c.Root(), "201", func() *chatmsgDemo { return newChatmsgDemo(6) })
			draft := ui.Local(c.Root(), "201.draft", func() string { return "" })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Sidebar": &d.sidebar}, "Sidebar")
				chatmsgControls(c, log, d)
			})
			inputSection(c, "Live", func() {
				ui.Box(c).Height(420).Border(1, mujica.Tokens(c).Border).Clip().Children(func() {
					opts := mujica.ChatContainerOptions{
						List: mujica.MessageListOptions{ID: func(i int) string { return d.ids[i] }},
						Input: func() {
							r := mujica.TextInput(c, draft, mujica.TextInputOptions{Label: "Message", Placeholder: "Write to the herald"})
							if r.Submitted() && strings.TrimSpace(*draft) != "" {
								d.add()
								d.texts[len(d.texts)-1] = *draft
								log("sent a message")
								*draft = ""
							}
						},
					}
					if d.sidebar {
						opts.Sidebar = func() {
							ui.Column(c).Padding(12).Gap(8).Children(func() {
								for _, s := range []string{"The ball", "Seating", "Envoys"} {
									ui.Text(c, s)
								}
							})
						}
					}
					r := mujica.ChatContainer(c, &d.list, len(d.ids), opts, func(i int) { d.row(c, i) })
					if r.Unread() > 0 {
						ui.Text(c, fmt.Sprintf("%d unread", r.Unread())).Absolute().Top(8).Right(12)
					}
				})
			})
		}},
	{ID: "202", Name: "MessageList", ZhName: "消息列表", Category: chatmsgCategory,
		Keys: "滚轮与方向键滚动；上翻后停止跟随；Tab 到回到最新按钮，Enter 回到底部。",
		Code: `r := mujica.MessageList(c, &list, len(ids), mujica.MessageListOptions{
	ID: func(i int) string { return ids[i] },
}, func(i int) { ui.Text(c, "Message "+ids[i]) })
fmt.Println("unread", r.Unread())`,
		Demo: func(c *ui.Context, log func(string)) {
			d := *ui.Local(c.Root(), "202", func() *chatmsgDemo { return newChatmsgDemo(30) })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Date separators": &d.dated}, "Date separators")
				chatmsgControls(c, log, d)
				ui.Row(c).Children(func() {
					if mujica.Button(c, "Load 100,000", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
						for len(d.ids) < 100000 {
							d.add()
						}
						log("100,000 messages loaded")
					}
				})
			})
			inputSection(c, "Live", func() {
				opts := mujica.MessageListOptions{ID: func(i int) string { return d.ids[i] }, Label: func(i int) string { return d.texts[i] }}
				if d.dated {
					opts.Date = d.day
				}
				ui.Column(c).Height(360).Border(1, mujica.Tokens(c).Border).Children(func() {
					r := mujica.MessageList(c, &d.list, len(d.ids), opts, func(i int) { d.row(c, i) })
					ui.Text(c, fmt.Sprintf("%d messages, %d unread, at end %v", len(d.ids), r.Unread(), d.list.AtEnd())).Padding(8)
				})
			})
		}},
	{ID: "203", Name: "MessageBubble", ZhName: "消息气泡", Category: chatmsgCategory,
		Keys: "气泡不获取焦点；内容中的链接与按钮按 Tab 顺序访问。",
		Code: `mujica.MessageBubble(c, mujica.MessageAssistant, mujica.MessageBubbleOptions{Name: "Herald"}, func() {
	ui.Text(c, "The ball begins at nine.")
})`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				role          mujica.MessageRole
				avatar, named bool
				width         float64
			}
			s := ui.Local(c.Root(), "203", func() state { return state{role: mujica.MessageAssistant, avatar: true, named: true, width: 360} })
			roles := []mujica.RadioItem[mujica.MessageRole]{{Value: mujica.MessageUser, Label: "User"}, {Value: mujica.MessageAssistant, Label: "Assistant"},
				{Value: mujica.MessageSystem, Label: "System"}, {Value: mujica.MessageError, Label: "Error"}}
			inputSection(c, "Parameters", func() {
				if mujica.RadioGroup(c, &s.role, roles, mujica.RadioGroupOptions{Horizontal: true, Label: "Role"}).Changed() {
					log(fmt.Sprintf("role %d", s.role))
				}
				inputFlags(c, log, map[string]*bool{"Avatar": &s.avatar, "Name": &s.named}, "Avatar", "Name")
				if mujica.NumberInput(c, &s.width, mujica.NumberInputOptions{Min: 160, Max: 640, Step: 40, Label: "MaxWidth"}).Changed() {
					log(fmt.Sprintf("MaxWidth %v", s.width))
				}
			})
			inputSection(c, "Live", func() {
				opts := mujica.MessageBubbleOptions{MaxWidth: float32(s.width)}
				if s.named {
					opts.Name = "Herald"
				}
				if s.avatar {
					opts.Avatar = func() { mujica.Avatar(c, "Herald", mujica.AvatarOptions{}) }
				}
				mujica.MessageBubble(c, s.role, opts, func() { ui.Text(c, chatmsgReplies[0]+" "+chatmsgReplies[2]) })
			})
			inputSection(c, "States", func() {
				for _, r := range roles {
					mujica.MessageBubble(c, r.Value, mujica.MessageBubbleOptions{Name: r.Label}, func() { ui.Text(c, r.Label+" message") })
				}
			})
		}},
	{ID: "204", Name: "MessageHeader", ZhName: "消息头", Category: chatmsgCategory, Keys: "不获取焦点。",
		Code: `mujica.MessageHeader(c, mujica.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: mujica.MessageSent})`,
		Demo: func(c *ui.Context, log func(string)) {
			status := ui.Local(c.Root(), "204", func() mujica.MessageStatus { return mujica.MessageSent })
			items := []mujica.RadioItem[mujica.MessageStatus]{{Value: mujica.MessageNoStatus, Label: "None"}, {Value: mujica.MessageSending, Label: "Sending"},
				{Value: mujica.MessageStreaming, Label: "Streaming"}, {Value: mujica.MessageSent, Label: "Sent"}, {Value: mujica.MessageFailed, Label: "Failed"}}
			inputSection(c, "Parameters", func() {
				if mujica.RadioGroup(c, status, items, mujica.RadioGroupOptions{Horizontal: true, Label: "Status"}).Changed() {
					log(fmt.Sprintf("status %d", *status))
				}
			})
			inputSection(c, "Live", func() {
				mujica.MessageHeader(c, mujica.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: *status})
			})
			inputSection(c, "States", func() {
				for _, it := range items[1:] {
					mujica.MessageHeader(c, mujica.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: it.Value})
				}
			})
		}},
	{ID: "205", Name: "StreamingText", ZhName: "流式文本", Category: chatmsgCategory, Keys: "不获取焦点；减少动画时直接显示全部文本。",
		Code: `mujica.StreamingText(c, reply, mujica.StreamingTextOptions{Streaming: true, Smooth: true})`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				text                    string
				smooth, streaming, auto bool
				rate                    float64
				at                      time.Time
			}
			s := ui.Local(c.Root(), "205", func() state { return state{smooth: true, streaming: true, rate: 80} })
			const full = "By order of the Queen, the ball begins at nine. Guests arrive by the east door, and carriages wait in the lower court until midnight."
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Smooth": &s.smooth, "Streaming cursor": &s.streaming, "Stream": &s.auto}, "Smooth", "Streaming cursor", "Stream")
				if mujica.NumberInput(c, &s.rate, mujica.NumberInputOptions{Min: 10, Max: 400, Step: 10, Label: "Rate"}).Changed() {
					log(fmt.Sprintf("Rate %v", s.rate))
				}
				ui.Row(c).Children(func() {
					if mujica.Button(c, "Restart", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
						s.text = ""
						log("stream restarted")
					}
				})
			})
			// The demo's own source appends a word every 150 ms while streaming.
			if s.auto && len(s.text) < len(full) && c.Now().Sub(s.at) > 150*time.Millisecond {
				rest := full[len(s.text):]
				n := strings.IndexByte(rest[1:], ' ') + 2
				if n <= 1 {
					n = len(rest)
				}
				s.text += rest[:n]
				s.at = c.Now()
				if len(s.text) == len(full) {
					log("stream finished")
				}
			}
			if s.auto && len(s.text) < len(full) {
				c.After(150 * time.Millisecond)
			}
			inputSection(c, "Live", func() {
				mujica.StreamingText(c, s.text, mujica.StreamingTextOptions{Streaming: s.streaming && len(s.text) < len(full), Smooth: s.smooth, Rate: s.rate})
			})
		}},
}
