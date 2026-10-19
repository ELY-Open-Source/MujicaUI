package agent

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// MemoryEntry is one memory; the caller persists it.
type MemoryEntry struct {
	ID   string
	Text string
	// Source says where it came from, as "Chat, 3 Oct".
	Source string
}

// MemoryPanelOptions configure a MemoryPanel.
type MemoryPanelOptions struct {
	// MaxHeight caps the scrolling list; zero means 360.
	MaxHeight float32
	Disabled  bool
}

// MemoryPanelResult reports this frame's in-place edits.
type MemoryPanelResult struct {
	Element ui.Element
	// Edited and Deleted are changed IDs, else "".
	Edited  string
	Deleted string
}

type memoryPanelState struct {
	query, editing, draft string
	confirm               string
	// The filter reruns only when query or entries change.
	key   string
	shown []int
}

// MemoryPanel lists memories with search, inline edit and confirmed delete.
func MemoryPanel(c *ui.Context, entries *[]MemoryEntry, opts MemoryPanelOptions) MemoryPanelResult {
	if entries == nil {
		panic("mujica: MemoryPanel needs entries")
	}
	if opts.MaxHeight < 0 {
		panic("mujica: MemoryPanel MaxHeight is negative")
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	u := core.Density(c).Unit()
	var r MemoryPanelResult
	r.Element = AgentCard(c).Disabled(opts.Disabled).Label(core.Msg(c, "memory_panel.title"))
	st := ui.Local(r.Element, "mujica.memory", func() memoryPanelState { return memoryPanelState{} })
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "memory_panel.title")).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			AgentCaption(c, fmt.Sprintf(core.Msg(c, "memory_panel.count"), len(*entries)))
		})
		input.SearchInput(c, &st.query, input.SearchInputOptions{Placeholder: core.Msg(c, "memory_panel.search"), Disabled: locked})
		memoryFilter(st, *entries)
		if len(st.shown) == 0 {
			AgentCaption(c, core.Msg(c, "memory_panel.empty"))
			return
		}
		AgentRows(c, len(st.shown), func(i int) {
			idx := st.shown[i]
			ui.Column(c.Key((*entries)[idx].ID)).Children(func() { memoryRow(c, entries, idx, st, locked, &r) })
		}).Shrink(1).MaxHeight(maxH).Dividers(theme.BorderWidth, core.Tokens(c).Border)
	})
	if r.Deleted != "" {
		for i, e := range *entries {
			if e.ID == r.Deleted {
				*entries = append((*entries)[:i:i], (*entries)[i+1:]...)
				break
			}
		}
		st.key = "" // refilter
	}
	if r.Edited != "" || r.Deleted != "" {
		slog.Debug("mujica: memory_panel change", "edited", r.Edited, "deleted", r.Deleted)
	}
	return r
}

// memoryFilter recomputes matches when inputs change.
func memoryFilter(st *memoryPanelState, entries []MemoryEntry) {
	key := fmt.Sprintf("%s\x00%d", st.query, len(entries))
	if len(entries) > 0 {
		key += fmt.Sprintf("\x00%p", &entries[0])
	}
	if key == st.key && st.shown != nil {
		return
	}
	st.key = key
	q := strings.ToLower(strings.TrimSpace(st.query))
	st.shown = st.shown[:0]
	for i, e := range entries {
		if q == "" || strings.Contains(strings.ToLower(e.Text), q) || strings.Contains(strings.ToLower(e.Source), q) {
			st.shown = append(st.shown, i)
		}
	}
	if st.shown == nil {
		st.shown = []int{}
	}
}

// memoryRow shows one entry or its editor.
func memoryRow(c *ui.Context, entries *[]MemoryEntry, idx int, st *memoryPanelState, locked bool, r *MemoryPanelResult) {
	e := &(*entries)[idx]
	if e.ID == "" {
		panic(fmt.Sprintf("mujica: MemoryPanel entry %d needs an ID", idx))
	}
	u := core.Density(c).Unit()
	ui.Column(c).Gap(u).Padding(u*2, 0).Children(func() {
		if st.editing == e.ID && !locked {
			core.AutoFocus(c, input.TextArea(c, &st.draft, input.TextAreaOptions{Label: core.Msg(c, "memory_panel.edit"), MinHeight: 60}).Input)
			ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
				if input.Button(c, core.Msg(c, "memory_panel.cancel"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					st.editing = ""
				}
				blank := strings.TrimSpace(st.draft) == ""
				if input.Button(c, core.Msg(c, "memory_panel.save"), input.ButtonOptions{Disabled: blank}).Clicked() && !blank {
					e.Text, st.editing, st.key = strings.TrimSpace(st.draft), "", ""
					r.Edited = e.ID
				}
			})
			return
		}
		ui.Text(c, e.Text).Selectable()
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			AgentCaption(c, e.Source).Grow(1)
			if AgentToggle(c, core.Msg(c, "memory_panel.edit"), "", locked) {
				st.editing, st.draft = e.ID, e.Text
			}
			del := ui.Row(c)
			del.Children(func() {
				if AgentToggle(c, core.Msg(c, "memory_panel.delete"), "trash-2", locked) {
					st.confirm = e.ID
				}
			})
			open := st.confirm == e.ID && !locked
			res := overlay.Popconfirm(c, del, &open, overlay.PopconfirmOptions{Title: core.Msg(c, "memory_panel.confirm"), Description: e.Text,
				ConfirmLabel: core.Msg(c, "memory_panel.delete"), Danger: true})
			if st.confirm == e.ID && !open {
				st.confirm = ""
			}
			if res.Confirmed {
				r.Deleted = e.ID
			}
		})
	})
}
