package data

import (
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DescriptionItem is a name and its value.
type DescriptionItem struct {
	Label string
	Value string
	// Span is how many columns the item takes; negative spans the whole row.
	Span int
	// Copyable adds a button copying Value to the clipboard.
	Copyable bool
}

// DescriptionListOptions configure a DescriptionList.
type DescriptionListOptions struct {
	// Columns is the most columns shown; 0 means 2.
	Columns int
	// MinColumnWidth drops columns when narrower; 0 means 220 DIPs.
	MinColumnWidth float32
}

// DescriptionListView is a built DescriptionList and what happened to it in this frame.
type DescriptionListView struct {
	Element ui.Element
	copied  int
}

// Copied returns the index of the item whose value was copied in this frame.
func (v DescriptionListView) Copied() (int, bool) { return v.copied, v.copied >= 0 }

// descriptionCopied shows the copied mark for this long.
const descriptionCopied = 1500 * time.Millisecond

type descriptionState struct {
	item  int
	until time.Time
}

// DescriptionList shows names and values in columns, one column when narrow.
func DescriptionList(c *ui.Context, items []DescriptionItem, opts DescriptionListOptions) DescriptionListView {
	if opts.Columns < 0 || opts.MinColumnWidth < 0 {
		panic("mujica: DescriptionList with negative Columns or MinColumnWidth")
	}
	most, minW := opts.Columns, opts.MinColumnWidth
	if most == 0 {
		most = 2
	}
	if minW == 0 {
		minW = 220
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	gap := u * 6
	v := DescriptionListView{Element: ui.Grid(c).GapX(gap).GapY(u * 4).FillWidth(), copied: -1}
	st := ui.Local(v.Element, "state", func() descriptionState { return descriptionState{item: -1} })
	fit := func(w float32) int {
		if w <= 0 {
			return most
		}
		return max(1, min(most, int(math.Floor(float64((w+gap)/(minW+gap))))))
	}
	cols := fit(v.Element.Bounds().W)
	v.Element.Draw(func(p *ui.Painter, r ui.Rect) {
		if fit(r.W) != cols {
			// The new width needs other columns: build again.
			c.Invalidate()
		}
	})
	now := c.Now()
	v.Element.Columns(cols).Children(func() {
		for i, it := range items {
			cell := ui.Column(c).Gap(u).MinWidth(0)
			if it.Span < 0 || it.Span > cols {
				cell.ColumnSpan(-1)
			} else if it.Span > 1 {
				cell.ColumnSpan(it.Span)
			}
			cell.Children(func() {
				ui.Text(c, it.Label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				ui.Row(c).Gap(u).MinHeight(theme.MinHitSize).AlignItems(ui.Center).Children(func() {
					ui.Text(c, it.Value).Shrink(1).MinWidth(0)
					if !it.Copyable {
						return
					}
					done := st.item == i && now.Before(st.until)
					icon, label := icons.Must("copy"), core.Msg(c, "description_list.copy")+" "+it.Label
					if done {
						icon, label = icons.Must("check"), core.Msg(c, "description_list.copied")
						c.After(st.until.Sub(now))
					}
					if input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icon, Label: label}).MinHeight(theme.MinHitSize).Padding(0, u).Clicked() {
						c.WriteClipboard(it.Value)
						c.Announce(core.Msg(c, "description_list.copied"))
						*st = descriptionState{item: i, until: now.Add(descriptionCopied)}
						v.copied = i
						slog.Debug("mujica: description_list copy", "item", i)
					}
				})
			})
		}
	})
	return v
}
