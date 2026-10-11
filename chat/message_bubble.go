package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MessageRole is who a chat message comes from.
type MessageRole int

const (
	MessageUser MessageRole = iota
	MessageAssistant
	MessageSystem
	MessageError
)

func (r MessageRole) key() string {
	switch r {
	case MessageUser:
		return "message_bubble.user"
	case MessageAssistant:
		return "message_bubble.assistant"
	case MessageSystem:
		return "message_bubble.system"
	case MessageError:
		return "message_bubble.error"
	}
	panic("mujica: unknown MessageRole")
}

// MessageBubbleOptions configure a MessageBubble.
type MessageBubbleOptions struct {
	// Name shows above the bubble; empty shows none.
	Name string
	// Avatar builds the picture beside the bubble.
	Avatar func()
	// MaxWidth caps the bubble in DIPs; zero means 560.
	MaxWidth float32
}

// MessageBubble frames one chat message by role: users right, others left, system centered.
func MessageBubble(c *ui.Context, role MessageRole, opts MessageBubbleOptions, content func()) *ui.Element {
	if opts.MaxWidth < 0 {
		panic("mujica: MessageBubble MaxWidth must not be negative")
	}
	if content == nil {
		panic("mujica: MessageBubble needs content")
	}
	maxW := opts.MaxWidth
	if maxW == 0 {
		maxW = 560
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := core.Msg(c, role.key())
	if opts.Name != "" {
		label += ", " + opts.Name
	}
	row := ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Role(ui.RoleGroup).Label(label)
	switch role {
	case MessageUser:
		row.Justify(ui.End)
	case MessageSystem:
		row.Justify(ui.Center)
	}
	row.Children(func() {
		if role == MessageSystem {
			ui.Row(c).MaxWidth(maxW).Padding(u, u*3).Gap(u * 2).Radius(theme.ControlRadius).
				Background(k.SurfaceHover).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).
				Children(content)
			return
		}
		if opts.Avatar != nil && role != MessageUser {
			opts.Avatar()
		}
		ui.Column(c).Gap(u).MaxWidth(maxW).Shrink(1).MinWidth(0).AlignItems(messageBubbleAlign(role)).Children(func() {
			if opts.Name != "" {
				ui.Text(c, opts.Name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			}
			b := ui.Column(c).Padding(u*2, u*3).Gap(u * 2).Radius(theme.CardRadius).MinWidth(0)
			switch role {
			case MessageUser:
				b.Background(k.Selection)
			case MessageError:
				b.Background(k.Surface).Border(theme.BorderWidth, k.Danger)
			default:
				b.Background(k.Surface).Border(theme.BorderWidth, k.Border)
			}
			b.Children(func() {
				if role != MessageError {
					content()
					return
				}
				ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
					ui.Icon(c, icons.Must("circle-alert")).Size(16, 16).Shrink(0).TextColor(k.Danger)
					ui.Column(c).Gap(u * 2).Grow(1).MinWidth(0).Children(content)
				})
			})
		})
		if opts.Avatar != nil && role == MessageUser {
			opts.Avatar()
		}
	})
	return row
}

func messageBubbleAlign(r MessageRole) ui.Align {
	if r == MessageUser {
		return ui.End
	}
	return ui.Start
}
