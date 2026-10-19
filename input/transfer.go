package input

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// TransferItem is one entry of a Transfer; ID must be unique and stable.
type TransferItem struct {
	ID       string
	Label    string
	Disabled bool
}

// TransferOptions configure a Transfer.
type TransferOptions struct {
	// SourceTitle and TargetTitle replace the built-in "Available" and "Chosen".
	SourceTitle, TargetTitle string
	// Height of each list in DIPs; 240 when zero.
	Height   float32
	Disabled bool
}

type transferState struct {
	checked map[string]bool
	query   [2]string
}

// Transfer moves items between a source list and *target, the chosen IDs in order.
func Transfer(c *ui.Context, target *[]string, items []TransferItem, opts TransferOptions) core.ChoiceResult {
	byID := make(map[string]TransferItem, len(items))
	for _, it := range items {
		if _, dup := byID[it.ID]; dup || it.ID == "" {
			panic(fmt.Sprintf("mujica: Transfer item ID %q is empty or repeated", it.ID))
		}
		byID[it.ID] = it
	}
	for _, id := range *target {
		if _, ok := byID[id]; !ok {
			panic(fmt.Sprintf("mujica: Transfer target ID %q is not an item", id))
		}
	}
	d := core.Density(c)
	parts := core.ChoiceResult{}
	row := ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Disabled(opts.Disabled)
	parts.Element = row
	st := ui.Local(row, "transfer", func() transferState { return transferState{checked: map[string]bool{}} })
	for id := range st.checked {
		if it, ok := byID[id]; !ok || it.Disabled {
			delete(st.checked, id)
		}
	}
	var source, chosen []TransferItem
	for _, it := range items {
		if !slices.Contains(*target, it.ID) {
			source = append(source, it)
		}
	}
	for _, id := range *target {
		chosen = append(chosen, byID[id])
	}
	move := func(from []TransferItem, right bool) {
		var ids []string
		for _, it := range from {
			if st.checked[it.ID] && !it.Disabled {
				ids = append(ids, it.ID)
				delete(st.checked, it.ID)
			}
		}
		if len(ids) == 0 {
			return
		}
		if right {
			*target = append(slices.Clone(*target), ids...)
		} else {
			*target = slices.DeleteFunc(slices.Clone(*target), func(id string) bool { return slices.Contains(ids, id) })
		}
		parts = core.NewChoiceResult(parts.Element, true)
		slog.Debug("mujica: transfer moved", "count", len(ids), "right", right)
	}
	height := opts.Height
	if height == 0 {
		height = 240
	}
	srcTitle, dstTitle := opts.SourceTitle, opts.TargetTitle
	if srcTitle == "" {
		srcTitle = core.Msg(c, "transfer.source")
	}
	if dstTitle == "" {
		dstTitle = core.Msg(c, "transfer.target")
	}
	row.Children(func() {
		transferPanel(c, st, 0, srcTitle, source, height)
		ui.Column(c).Gap(d.Unit() * 2).Children(func() {
			right := transferCount(st, source) > 0
			left := transferCount(st, chosen) > 0
			if Button(c, "", ButtonOptions{Variant: Secondary, Icon: icons.Must("chevron-right"), Label: fmt.Sprintf(core.Msg(c, "transfer.to"), dstTitle), Disabled: !right}).Clicked() {
				move(source, true)
			}
			if Button(c, "", ButtonOptions{Variant: Secondary, Icon: icons.Must("chevron-left"), Label: fmt.Sprintf(core.Msg(c, "transfer.to"), srcTitle), Disabled: !left}).Clicked() {
				move(chosen, false)
			}
		})
		transferPanel(c, st, 1, dstTitle, chosen, height)
	})
	return parts
}

// transferCount counts the checked items that can move.
func transferCount(st *transferState, items []TransferItem) int {
	n := 0
	for _, it := range items {
		if st.checked[it.ID] && !it.Disabled {
			n++
		}
	}
	return n
}

// transferPanel is one side: header with select-all, search, and the list.
func transferPanel(c *ui.Context, st *transferState, side int, title string, items []TransferItem, height float32) {
	k := core.Tokens(c)
	d := core.Density(c)
	q := strings.ToLower(strings.TrimSpace(st.query[side]))
	var shown []TransferItem
	enabled := 0
	for _, it := range items {
		if q == "" || strings.Contains(strings.ToLower(it.Label), q) {
			shown = append(shown, it)
			if !it.Disabled {
				enabled++
			}
		}
	}
	checked := transferCount(st, shown)
	ui.Column(c).Width(240).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).Children(func() {
		head := ui.Row(c).Gap(d.Unit()*2).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3, 0, d.Unit()*3-(theme.MinHitSize-16)/2).
			BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border)
		head.Children(func() {
			all := checked > 0 && checked == enabled
			box := ui.CheckboxBase(c, &all).Size(theme.MinHitSize, theme.MinHitSize).Justify(ui.Center).Label(fmt.Sprintf(core.Msg(c, "transfer.all"), title)).Disabled(enabled == 0)
			transferMark(c, box, all, checked > 0 && !all)
			box.Description(transferPartial(c, checked, enabled))
			if box.Changed() {
				for _, it := range shown {
					if !it.Disabled {
						st.checked[it.ID] = all
					}
				}
			}
			ui.Text(c, title).Font(theme.SerifFont).FontWeight(600).TextColor(k.Text).Grow(1)
			ui.Text(c, fmt.Sprintf("%d / %d", checked, len(items))).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
		})
		ui.Box(c).Padding(d.Unit() * 2).Children(func() {
			InputGroup(c, &st.query[side], InputGroupOptions{Placeholder: core.Msg(c, "transfer.search"), Label: fmt.Sprintf(core.Msg(c, "transfer.search_in"), title)})
		})
		list := ui.Scroll(c).Height(height).Padding(0, d.Unit(), d.Unit(), d.Unit()).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(title)
		list.Children(func() {
			if len(shown) == 0 {
				empty := core.Msg(c, "transfer.empty")
				if q != "" {
					empty = core.Msg(c, "transfer.no_match")
				}
				ui.Text(c, empty).TextColor(k.TextMuted).Padding(d.Unit() * 3).TextAlign(ui.Center).FillWidth()
				return
			}
			for _, it := range shown {
				on := st.checked[it.ID]
				ui.Box(c.Key(it.ID)).Children(func() {
					cb := ui.CheckboxBase(c, &on).Gap(d.Unit()*2).MinHeight(max(d.ControlHeight()-d.Unit(), theme.MinHitSize)).Padding(0, d.Unit()*2).
						Radius(theme.ControlRadius).FillWidth().Disabled(it.Disabled)
					if cb.Changed() {
						st.checked[it.ID] = on
					}
					if on {
						cb.Background(k.Selection)
					} else if cb.Hovered() {
						cb.Background(k.SurfaceHover)
					}
					transferMark(c, cb, on, false)
					cb.Children(func() { ui.Text(c, it.Label).SingleLine().Ellipsis("…").MinWidth(0) })
				})
			}
		})
	})
}

// transferPartial describes a select-all box that is neither clear nor full.
func transferPartial(c *ui.Context, checked, enabled int) string {
	if checked == 0 || checked == enabled {
		return ""
	}
	return fmt.Sprintf(core.Msg(c, "transfer.partial"), checked, enabled)
}

// transferMark draws a check box's square inside box, which takes the focus ring.
func transferMark(c *ui.Context, box ui.Element, on, mixed bool) {
	k := core.Tokens(c)
	internal.Ring(box, k.Focus, theme.ControlRadius)
	box.Children(func() {
		sq := ui.Box(c).Size(16, 16).Shrink(0).Center().Radius(theme.ControlRadius).Role(ui.RoleNone)
		if on || mixed {
			sq.Background(k.Accent).TextColor(k.OnAccent)
		} else {
			sq.Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder)
		}
		sq.Children(func() {
			switch {
			case on:
				ui.Icon(c, icons.Must("check")).Size(12, 12)
			case mixed:
				ui.Icon(c, icons.Must("minus")).Size(12, 12)
			}
		})
	})
}
