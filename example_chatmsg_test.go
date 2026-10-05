package mujica_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleChatContainer() {
	var chat mujica.MessageListState
	msgs := []string{"Hello", "Bonjour"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ChatContainer(c, &chat, len(msgs), mujica.ChatContainerOptions{
			List:  mujica.MessageListOptions{ID: func(i int) string { return msgs[i] }},
			Input: func() { ui.Text(c, "Composer goes here") },
		}, func(i int) { ui.Text(c, msgs[i]) })
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageList() {
	var list mujica.MessageListState
	ids := []string{"m1", "m2", "m3"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.MessageList(c, &list, len(ids), mujica.MessageListOptions{
			ID: func(i int) string { return ids[i] },
		}, func(i int) { ui.Text(c, "Message "+ids[i]) })
		fmt.Println("unread", r.Unread())
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageBubble() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MessageBubble(c, mujica.MessageAssistant, mujica.MessageBubbleOptions{Name: "Herald"}, func() {
			ui.Text(c, "The ball begins at nine.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageHeader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MessageHeader(c, mujica.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: mujica.MessageSent})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStreamingText() {
	reply := "The decree reads"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.StreamingText(c, reply, mujica.StreamingTextOptions{Streaming: true, Smooth: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMarkdownView() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if url, ok := mujica.MarkdownView(c, "# Decree\n\nSee **the guide** at [the court](https://example.com).").Link(); ok {
			fmt.Println("open", url)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCodeBlock() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CodeBlock(c, "func bow() {}", mujica.CodeBlockOptions{Language: "go", LineNumbers: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTableBlock() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.TableBlock(c, []string{"Guest", "Seat"}, [][]string{{"Marie", "1"}}, mujica.TableBlockOptions{}).Copied() {
			fmt.Println("copied as TSV")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageGrid() {
	pic := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="3"/>`))
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if i, ok := mujica.ImageGrid(c, []mujica.GridImage{{Source: pic, Alt: "Portrait"}}).Retried(); ok {
			fmt.Println("reload", i)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileMessage() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.FileMessage(c, mujica.FileMessageOptions{Name: "treaty.pdf", Size: 2621440})
		if r.Downloaded() {
			fmt.Println("save treaty.pdf")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAudioMessage() {
	pos := 5 * time.Second
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.AudioMessage(c, mujica.AudioMessageOptions{Duration: time.Minute, Position: pos})
		if d, ok := r.Seek(); ok {
			pos = d
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLinkPreviewCard() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		p := mujica.LinkPreview{URL: "https://example.com", Title: "The Palace", Site: "example.com"}
		if mujica.LinkPreviewCard(c, p).Clicked() {
			fmt.Println("open", p.URL)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuoteReply() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.QuoteReply(c, "Marie", "Let them eat cake.").Clicked() {
			fmt.Println("jump to the quoted message")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageActions() {
	rating := mujica.FeedbackNone
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if a, ok := mujica.MessageActions(c, mujica.MessageAssistant, mujica.MessageActionsOptions{Feedback: &rating}).Action(); ok {
			fmt.Println("action", a)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMessageEditor() {
	draft := "Hello"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.MessageEditor(c, &draft, mujica.MessageEditorOptions{Original: "Hello"}).Saved() {
			fmt.Println("new branch from", draft)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchNavigator() {
	branch := 1
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.BranchNavigator(c, &branch, 3).Changed() {
			fmt.Println("show branch", branch)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleErrorMessage() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.ErrorMessageOptions{RetryAt: time.Now().Add(30 * time.Second)}
		if mujica.ErrorMessage(c, "429 Too Many Requests", opts).Retried() {
			fmt.Println("send again")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThinkingBlock() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ThinkingBlock(c, &open, mujica.ThinkingBlockOptions{Elapsed: 12 * time.Second}, func() {
			ui.Text(c, "Weigh the guest list.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScrollToBottomButton() {
	var list mujica.MessageListState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Box(c).Fill().Children(func() {
			if mujica.ScrollToBottomButton(c, !list.AtEnd(), 3).Clicked() {
				list.ScrollToEnd()
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTypingIndicator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TypingIndicator(c, "Herald")
	}
	ui.Render(view, 400, 300, 1)
}
