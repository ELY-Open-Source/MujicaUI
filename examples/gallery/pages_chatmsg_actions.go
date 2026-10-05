package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chat"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

var chatmsgActionPages = []page{
	{ID: "211", Name: "AudioMessage", ZhName: "语音消息", Category: chatmsgCategory, Keys: "Tab 到播放按钮与波形；波形上 Left/Right 跳转，Home、End 到首尾。",
		Code: `r := chat.AudioMessage(c, chat.AudioMessageOptions{Duration: time.Minute, Position: pos})
if d, ok := r.Seek(); ok {
	pos = d
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				pos     time.Duration
				playing bool
				at      time.Time
			}
			s := ui.Local(c.Root(), "211", func() state { return state{pos: 12 * time.Second} })
			const length = 47 * time.Second
			// The demo plays: it moves the position while playing.
			if s.playing {
				s.pos += c.Now().Sub(s.at)
				if s.pos >= length {
					s.pos, s.playing = length, false
					log("playback ended")
				}
				c.After(100 * time.Millisecond)
			}
			s.at = c.Now()
			inputSection(c, "Live", func() {
				r := chat.AudioMessage(c, chat.AudioMessageOptions{Duration: length, Position: s.pos, Playing: s.playing, Waveform: chatmsgWave})
				if r.Toggled() {
					s.playing = !s.playing
					if s.pos >= length {
						s.pos = 0
					}
					log(fmt.Sprintf("playing %v", s.playing))
				}
				if d, ok := r.Seek(); ok {
					s.pos = d
					log(fmt.Sprintf("seek to %v", d.Round(time.Second)))
				}
			})
			inputSection(c, "States", func() {
				ui.Column(c).Disabled(true).Children(func() {
					chat.AudioMessage(c, chat.AudioMessageOptions{Duration: length, Waveform: chatmsgWave})
				})
			})
		}},
	{ID: "212", Name: "LinkPreviewCard", ZhName: "链接预览卡片", Category: chatmsgCategory, Keys: "Tab 聚焦卡片，Enter 或 Space 打开。",
		Code: `p := chat.LinkPreview{URL: "https://example.com", Title: "The Palace", Site: "example.com"}
if chat.LinkPreviewCard(c, p).Clicked() {
	fmt.Println("open", p.URL)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			thumb := ui.Local(c.Root(), "212", func() bool { return true })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Thumbnail": thumb}, "Thumbnail")
			})
			p := chat.LinkPreview{URL: "https://example.com/versailles", Title: "The Palace of Versailles", Site: "example.com",
				Description: "A royal château in the Île-de-France region, seat of the French court from 1682 until the Revolution."}
			if *thumb {
				p.Thumbnail = chatmsgPicture("#AD9164")
			}
			inputSection(c, "Live", func() {
				if chat.LinkPreviewCard(c, p).Clicked() {
					log("open " + p.URL)
				}
				if chat.LinkPreviewCard(c, chat.LinkPreview{URL: "https://example.com/hall", Title: "Hall of Mirrors"}).Clicked() {
					log("open https://example.com/hall")
				}
			})
		}},
	{ID: "213", Name: "QuoteReply", ZhName: "引用回复", Category: chatmsgCategory, Keys: "Tab 聚焦，Enter 或 Space 跳转。",
		Code: `if chat.QuoteReply(c, "Marie", "Let them eat cake.").Clicked() {
	fmt.Println("jump to the quoted message")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			inputSection(c, "Live", func() {
				for _, q := range [][2]string{{"Marie", "Let them eat cake."}, {"Herald", chatmsgMarkdown}} {
					if chat.QuoteReply(c, q[0], q[1]).Clicked() {
						log("jump to " + q[0] + ": " + chatmsgSnippet(q[1], 24))
					}
				}
			})
		}},
	{ID: "214", Name: "MessageActions", ZhName: "消息操作", Category: chatmsgCategory, Keys: "Tab 进入工具栏，Left/Right 移动，Enter 或 Space 触发。",
		Code: `if a, ok := chat.MessageActions(c, chat.MessageAssistant, chat.MessageActionsOptions{Feedback: &rating}).Action(); ok {
	fmt.Println("action", a)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				role   chat.MessageRole
				rating chat.MessageFeedback
				rated  bool
			}
			s := ui.Local(c.Root(), "214", func() state { return state{role: chat.MessageAssistant, rated: true} })
			names := map[chat.MessageAction]string{chat.ActionCopy: "copy", chat.ActionRegenerate: "regenerate", chat.ActionEdit: "edit",
				chat.ActionLike: "like", chat.ActionDislike: "dislike", chat.ActionSpeak: "speak", chat.ActionShare: "share"}
			inputSection(c, "Parameters", func() {
				roles := []input.RadioItem[chat.MessageRole]{{Value: chat.MessageUser, Label: "User"}, {Value: chat.MessageAssistant, Label: "Assistant"},
					{Value: chat.MessageSystem, Label: "System"}, {Value: chat.MessageError, Label: "Error"}}
				if input.RadioGroup(c, &s.role, roles, input.RadioGroupOptions{Horizontal: true, Label: "Role"}).Changed() {
					log(fmt.Sprintf("role %d", s.role))
				}
				inputFlags(c, log, map[string]*bool{"Rating": &s.rated}, "Rating")
			})
			inputSection(c, "Live", func() {
				opts := chat.MessageActionsOptions{}
				if s.rated {
					opts.Feedback = &s.rating
				}
				if a, ok := chat.MessageActions(c, s.role, opts).Action(); ok {
					log(fmt.Sprintf("action %s, rating %d", names[a], s.rating))
				}
			})
		}},
	{ID: "215", Name: "MessageEditor", ZhName: "消息编辑器", Category: chatmsgCategory, Keys: "Cmd/Ctrl+Enter 保存，Escape 取消。",
		Code: `if chat.MessageEditor(c, &draft, chat.MessageEditorOptions{Original: "Hello"}).Saved() {
	fmt.Println("new branch from", draft)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				original, draft string
				editing         bool
				branches        []string
				index           int
			}
			s := ui.Local(c.Root(), "215", func() state {
				return state{original: "Seat the ambassador.", branches: []string{"Seat the ambassador."}}
			})
			inputSection(c, "Live", func() {
				if !s.editing {
					chat.MessageBubble(c, chat.MessageUser, chat.MessageBubbleOptions{}, func() { ui.Text(c, s.branches[s.index]) })
					ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
						if chat.BranchNavigator(c, &s.index, len(s.branches)).Changed() {
							log(fmt.Sprintf("branch %d", s.index+1))
						}
						if input.Button(c, "Edit", input.ButtonOptions{Variant: input.Ghost}).Clicked() {
							s.editing, s.draft = true, s.branches[s.index]
							log("editing")
						}
					})
					return
				}
				r := chat.MessageEditor(c, &s.draft, chat.MessageEditorOptions{Original: s.branches[s.index]})
				if r.Saved() {
					s.branches = append(s.branches, s.draft)
					s.index, s.editing = len(s.branches)-1, false
					log(fmt.Sprintf("saved as branch %d", len(s.branches)))
				}
				if r.Canceled() {
					s.editing = false
					log("edit canceled")
				}
			})
		}},
	{ID: "216", Name: "BranchNavigator", ZhName: "分支切换", Category: chatmsgCategory, Keys: "Tab 聚焦，Left/Right 切换，Home、End 到首末。",
		Code: `if chat.BranchNavigator(c, &branch, 3).Changed() {
	fmt.Println("show branch", branch)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				index    int
				count    float64
				disabled bool
			}
			s := ui.Local(c.Root(), "216", func() state { return state{index: 1, count: 3} })
			inputSection(c, "Parameters", func() {
				if input.NumberInput(c, &s.count, input.NumberInputOptions{Min: 1, Max: 9, StepperOnly: true, Label: "Branches"}).Changed() {
					s.index = min(s.index, int(s.count)-1)
					log(fmt.Sprintf("%v branches", s.count))
				}
				inputFlags(c, log, map[string]*bool{"Disabled": &s.disabled}, "Disabled")
			})
			inputSection(c, "Live", func() {
				ui.Row(c).Disabled(s.disabled).Children(func() {
					if chat.BranchNavigator(c, &s.index, int(s.count)).Changed() {
						log(fmt.Sprintf("branch %d of %v", s.index+1, s.count))
					}
				})
			})
		}},
	{ID: "217", Name: "ErrorMessage", ZhName: "错误消息", Category: chatmsgCategory, Keys: "Tab 到错误文本（可选择复制）与重试按钮。",
		Code: `opts := chat.ErrorMessageOptions{RetryAt: time.Now().Add(30 * time.Second)}
if chat.ErrorMessage(c, "429 Too Many Requests", opts).Retried() {
	fmt.Println("send again")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			until := ui.Local(c.Root(), "217", func() time.Time { return time.Time{} })
			inputSection(c, "Parameters", func() {
				ui.Row(c).Children(func() {
					if input.Button(c, "Rate limit for 10 s", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
						*until = c.Now().Add(10 * time.Second)
						log("rate limited for 10 s")
					}
				})
			})
			inputSection(c, "Live", func() {
				opts := chat.ErrorMessageOptions{RetryAt: *until}
				if chat.ErrorMessage(c, `429 Too Many Requests: {"type":"rate_limit_error","message":"Number of requests exceeded"}`, opts).Retried() {
					log("retry")
				}
				if chat.ErrorMessage(c, "dial tcp 10.0.0.7:443: connection refused", chat.ErrorMessageOptions{Title: "Could not reach the server"}).Retried() {
					log("retry connection")
				}
			})
		}},
	{ID: "218", Name: "ThinkingBlock", ZhName: "思考过程", Category: chatmsgCategory, Keys: "Tab 聚焦标题，Enter 或 Space 展开与折叠。",
		Code: `chat.ThinkingBlock(c, &open, chat.ThinkingBlockOptions{Elapsed: 12 * time.Second}, func() {
	ui.Text(c, "Weigh the guest list.")
})`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				open, thinking bool
				started        time.Time
				elapsed        time.Duration
			}
			s := ui.Local(c.Root(), "218", func() state { return state{elapsed: 14 * time.Second} })
			inputSection(c, "Parameters", func() {
				label := "Start thinking"
				if s.thinking {
					label = "Finish thinking"
				}
				ui.Row(c).Children(func() {
					if input.Button(c, label, input.ButtonOptions{Variant: input.Secondary}).Clicked() {
						if s.thinking {
							s.elapsed = c.Now().Sub(s.started)
						} else {
							s.started = c.Now()
						}
						s.thinking = !s.thinking
						log(fmt.Sprintf("thinking %v", s.thinking))
					}
				})
			})
			inputSection(c, "Live", func() {
				opts := chat.ThinkingBlockOptions{Thinking: s.thinking, Started: s.started, Elapsed: s.elapsed}
				if chat.ThinkingBlock(c, &s.open, opts, func() {
					ui.Text(c, "Forty guests, thirty-six seats. The envoys outrank the counts; four counts dine in the antechamber.")
				}).Changed() {
					log(fmt.Sprintf("open %v", s.open))
				}
			})
		}},
	{ID: "219", Name: "ScrollToBottomButton", ZhName: "回到底部按钮", Category: chatmsgCategory, Keys: "Tab 聚焦，Enter 或 Space 回到底部。",
		Code: `ui.Box(c).Fill().Children(func() {
	if chat.ScrollToBottomButton(c, !list.AtEnd(), 3).Clicked() {
		list.ScrollToEnd()
	}
})`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				away   bool
				unread float64
			}
			s := ui.Local(c.Root(), "219", func() state { return state{away: true, unread: 3} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Away from bottom": &s.away}, "Away from bottom")
				if input.NumberInput(c, &s.unread, input.NumberInputOptions{Min: 0, Max: 99, Label: "Unread"}).Changed() {
					log(fmt.Sprintf("unread %v", s.unread))
				}
			})
			inputSection(c, "Live", func() {
				ui.Box(c).Height(120).Background(core.Tokens(c).Surface).Border(1, core.Tokens(c).Border).Children(func() {
					if chat.ScrollToBottomButton(c, s.away, int(s.unread)).Clicked() {
						s.away, s.unread = false, 0
						log("back to the bottom")
					}
				})
			})
		}},
	{ID: "220", Name: "TypingIndicator", ZhName: "正在输入", Category: chatmsgCategory, Keys: "不获取焦点；减少动画时显示静态文字。",
		Code: `chat.TypingIndicator(c, "Herald")`,
		Demo: func(c *ui.Context, log func(string)) {
			typing := ui.Local(c.Root(), "220", func() bool { return true })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Typing": typing}, "Typing")
			})
			inputSection(c, "Live", func() {
				if *typing {
					chat.TypingIndicator(c, "Herald")
				}
			})
		}},
}
