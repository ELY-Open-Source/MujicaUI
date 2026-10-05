package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ChatContainerOptions configure a ChatContainer.
type ChatContainerOptions struct {
	// List configures the message area; its ID is required.
	List MessageListOptions
	// Input builds the input area under the messages.
	Input func()
	// Sidebar builds a side column, as a session list; nil shows none.
	Sidebar func()
	// SidebarWidth is the side column's width in DIPs; zero means 260.
	SidebarWidth float32
}

// ChatContainer lays out a chat: sidebar, messages that follow the newest, and the input.
func ChatContainer(c *ui.Context, s *MessageListState, n int, opts ChatContainerOptions, row func(i int)) MessageListResult {
	if opts.SidebarWidth < 0 {
		panic("mujica: ChatContainer SidebarWidth must not be negative")
	}
	w := opts.SidebarWidth
	if w == 0 {
		w = 260
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r MessageListResult
	shell := ui.Row(c).Fill().AlignItems(ui.Stretch).Background(k.Background)
	shell.Children(func() {
		if opts.Sidebar != nil {
			ui.Column(c).Width(w).Shrink(0).MinHeight(0).Background(k.Surface).
				BorderWidth(0, theme.BorderWidth, 0, 0).BorderColor(k.Border).Children(opts.Sidebar)
		}
		ui.Column(c).Grow(1).MinWidth(0).MinHeight(0).Children(func() {
			r = MessageList(c, s, n, opts.List, row)
			if opts.Input != nil {
				ui.Column(c).Shrink(0).Padding(u*3, u*4).BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border).Children(opts.Input)
			}
		})
	})
	r.Element = shell
	return r
}
