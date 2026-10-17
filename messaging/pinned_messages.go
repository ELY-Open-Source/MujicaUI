package messaging

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PinnedMessagesResult is the list's element and this frame's events.
type PinnedMessagesResult struct {
	Element ui.Element
	// Jump requests navigation; Unpinned requests removing a pin.
	Jump, Unpinned string
}

// PinnedMessagesOptions configure a PinnedMessages.
type PinnedMessagesOptions struct {
	// Height caps the list; zero means 320.
	Height   float32
	Disabled bool
}

// PinnedMessages virtualizes pins and emits navigation or removal requests.
func PinnedMessages(c *ui.Context, pins []ChatEntry, opts PinnedMessagesOptions) PinnedMessagesResult {
	if opts.Height < 0 {
		panic("mujica: PinnedMessages Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r PinnedMessagesResult
	col := ui.Column(c).Gap(u * 2).Role(ui.RoleGroup).Label(fmt.Sprintf(core.Msg(c, "pinned_messages.title"), len(pins))).
		Disabled(opts.Disabled)
	r.Element = col
	col.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, msgPin).Size(16, 16).TextColor(k.Ornament)
			ui.Text(c, fmt.Sprintf(core.Msg(c, "pinned_messages.title"), len(pins))).Font(theme.SerifFont).FontWeight(600)
		})
		if len(pins) == 0 {
			ui.Text(c, core.Msg(c, "pinned_messages.empty")).TextColor(k.TextMuted)
			return
		}
		state := ui.Local(col, "pinned.rows", func() ui.ListState { return ui.ListState{} })
		state.Key = func(i int) any { return pins[i].ID }
		ui.List(c, state, len(pins), func(i int) {
			m := pins[i]
			if m.ID == "" || m.Author == "" {
				panic("mujica: PinnedMessages entries need ID and Author")
			}
			ui.Row(c).MinHeight(core.Density(c).ControlHeight()).Gap(u).AlignItems(ui.Center).Children(func() {
				if chat.QuoteReply(c, m.Author, m.Text).Grow(1).MaxWidth(10000).Clicked() && !col.IsDisabled() {
					r.Jump = m.ID
				}
				b := core.IconAction(c, msgPin, fmt.Sprintf(core.Msg(c, "pinned_messages.unpin"), m.Author))
				if b.Clicked() && !b.IsDisabled() {
					r.Unpinned = m.ID
				}
			})
		}).Height(min(h, float32(len(pins))*core.Density(c).ControlHeight())).Gap(u * 2).FocusGroup(ui.Vertical)

	})
	if r.Jump != "" || r.Unpinned != "" {
		slog.Debug("mujica: pinned_messages", "jump", r.Jump, "unpin", r.Unpinned)
	}
	return r
}
