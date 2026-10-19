package chat_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chat"
	"github.com/ELY-Open-Source/MujicaUI/core"
)

func ExampleChatContainer() {
	var thread chat.MessageListState
	msgs := []string{"Hello", "Bonjour"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.ChatContainer(c, &thread, len(msgs), chat.ChatContainerOptions{
			List:  chat.MessageListOptions{ID: func(i int) string { return msgs[i] }},
			Input: func() { ui.Text(c, "Composer goes here") },
		}, func(i int) { ui.Text(c, msgs[i]) })
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageList() {
	var list chat.MessageListState
	ids := []string{"m1", "m2", "m3"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.MessageList(c, &list, len(ids), chat.MessageListOptions{
			ID: func(i int) string { return ids[i] },
		}, func(i int) { ui.Text(c, "Message "+ids[i]) })
		fmt.Println("unread", r.Unread())
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageBubble() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.MessageBubble(c, chat.MessageAssistant, chat.MessageBubbleOptions{Name: "Herald"}, func() {
			ui.Text(c, "The ball begins at nine.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageHeader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.MessageHeader(c, chat.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: chat.MessageSent})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStreamingText() {
	reply := "The decree reads"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.StreamingText(c, reply, chat.StreamingTextOptions{Streaming: true, Smooth: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMarkdownView() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if url, ok := chat.MarkdownView(c, "# Decree\n\nSee **the guide** at [the court](https://example.com).").Link(); ok {
			fmt.Println("open", url)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCodeBlock() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.CodeBlock(c, "func bow() {}", chat.CodeBlockOptions{Language: "go", LineNumbers: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTableBlock() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.TableBlock(c, []string{"Guest", "Seat"}, [][]string{{"Marie", "1"}}, chat.TableBlockOptions{}).Copied() {
			fmt.Println("copied as TSV")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageGrid() {
	pic := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="3"/>`))
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := chat.ImageGrid(c, []chat.GridImage{{Source: pic, Alt: "Portrait"}}).Retried(); ok {
			fmt.Println("reload", i)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileMessage() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.FileMessage(c, chat.FileMessageOptions{Name: "treaty.pdf", Size: 2621440})
		if r.Downloaded() {
			fmt.Println("save treaty.pdf")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAudioMessage() {
	pos := 5 * time.Second
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chat.AudioMessage(c, chat.AudioMessageOptions{Duration: time.Minute, Position: pos})
		if d, ok := r.Seek(); ok {
			pos = d
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLinkPreviewCard() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		p := chat.LinkPreview{URL: "https://example.com", Title: "The Palace", Site: "example.com"}
		if chat.LinkPreviewCard(c, p).Clicked() {
			fmt.Println("open", p.URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuoteReply() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.QuoteReply(c, "Marie", "Let them eat cake.").Clicked() {
			fmt.Println("jump to the quoted message")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageActions() {
	rating := chat.FeedbackNone
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if a, ok := chat.MessageActions(c, chat.MessageAssistant, chat.MessageActionsOptions{Feedback: &rating}).Action(); ok {
			fmt.Println("action", a)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageEditor() {
	draft := "Hello"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.MessageEditor(c, &draft, chat.MessageEditorOptions{Original: "Hello"}).Saved() {
			fmt.Println("new branch from", draft)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchNavigator() {
	branch := 1
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if chat.BranchNavigator(c, &branch, 3).Changed() {
			fmt.Println("show branch", branch)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleErrorMessage() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := chat.ErrorMessageOptions{RetryAt: time.Now().Add(30 * time.Second)}
		if chat.ErrorMessage(c, "429 Too Many Requests", opts).Retried() {
			fmt.Println("send again")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThinkingBlock() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.ThinkingBlock(c, &open, chat.ThinkingBlockOptions{Elapsed: 12 * time.Second}, func() {
			ui.Text(c, "Weigh the guest list.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScrollToBottomButton() {
	var list chat.MessageListState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Fill().Children(func() {
			if chat.ScrollToBottomButton(c, !list.AtEnd(), 3).Clicked() {
				list.ScrollToEnd()
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTypingIndicator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chat.TypingIndicator(c, "Herald")
	}
	ui.Render(view, 400, 300, 1)
}
