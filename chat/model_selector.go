package chat

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ModelOption is one model a ModelSelector offers.
type ModelOption struct {
	// ID is the value chosen; required.
	ID   string
	Name string
	// Capabilities are short tags, as "Vision" or "Tools".
	Capabilities []string
	// ContextTokens is the context length; zero hides it.
	ContextTokens int
	// Unavailable explains disabled models; empty permits selection.
	Unavailable string
}

// ModelSelectorOptions configure a ModelSelector.
type ModelSelectorOptions struct {
	Disabled bool
	// Label names the selector; empty uses the built-in "Model".
	Label string
}

// ModelSelector chooses models with capabilities and availability details.
func ModelSelector(c *ui.Context, value *string, models []ModelOption, opts ModelSelectorOptions) core.ChoiceResult {
	if value == nil {
		panic("mujica: ModelSelector needs value")
	}
	for _, m := range models {
		if m.ID == "" || m.Name == "" || m.ContextTokens < 0 {
			panic("mujica: ModelOption needs ID, Name and non-negative ContextTokens")
		}
	}
	k, u := core.Tokens(c), core.Density(c).Unit()
	var tr *ui.Element
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() { tr = ui.ButtonBase(c).Role(ui.RolePopUpButton) })
	st := ui.Local(tr, "model-selector", func() modelSelectorState { return modelSelectorState{hi: -1, lastHi: -1} })
	disabled := tr.IsDisabled()
	if disabled {
		st.open = false
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "model_selector.label")
	}
	tr.Label(label)
	old := *value
	choose := func(i int) {
		if !disabled && i >= 0 && i < len(models) && models[i].Unavailable == "" {
			*value = models[i].ID
			st.open = false
			tr.Focus()
		}
	}
	if tr.Clicked() && !disabled {
		if st.open {
			choose(st.hi)
		} else {
			st.open = true
			st.hi = modelSelectorIndex(models, *value)
		}
	}
	if !disabled {
		move := 0
		if tr.Shortcut(0, ui.KeyDown) {
			move = 1
		}
		if tr.Shortcut(0, ui.KeyUp) {
			move = -1
		}
		if move != 0 {
			if !st.open {
				st.open = true
				st.hi = modelSelectorIndex(models, *value)
			} else {
				for i := st.hi + move; i >= 0 && i < len(models); i += move {
					if models[i].Unavailable == "" {
						st.hi = i
						break
					}
				}
			}
		}
		if tr.Shortcut(0, ui.KeyHome) {
			st.open = true
			st.hi = modelSelectorIndex(models, "")
		}
		if tr.Shortcut(0, ui.KeyEnd) {
			st.open = true
			for i := len(models) - 1; i >= 0; i-- {
				if models[i].Unavailable == "" {
					st.hi = i
					break
				}
			}
		}
	}
	if st.hi >= len(models) || st.hi >= 0 && models[st.hi].Unavailable != "" {
		st.hi = modelSelectorIndex(models, *value)
	}
	cur, has := modelSelectorFind(models, *value)
	input.SelectField(c, tr, st.open, func() {
		if has {
			ui.Text(c, cur.Name).SingleLine().FontWeight(600).Grow(1).Shrink(1)
		} else {
			ui.Text(c, core.Msg(c, "model_selector.placeholder")).SingleLine().TextColor(k.TextMuted).Grow(1).Shrink(1)
		}
	})
	if st.open {
		ww, wh := c.Size()
		ui.PopoverBase(c, tr, &st.open, func(panel *ui.Element) {
			core.OverlayStyle(c, panel, overlay.OverlayAnchored(c, panel, tr.Bounds()))
			a := tr.Bounds()
			room := max(a.Y, wh-a.Y-a.H) - 8
			panel.Margin(4, 0, 0, 0).Padding(4).Width(min(float32(360), max(0, ww-8))).Height(min(float32(360), max(0, room))).Label(label)
			if st.hi != st.lastHi {
				st.list.ScrollIntoView(max(st.hi, 0))
				st.lastHi = st.hi
			}
			ui.List(c, &st.list, len(models), func(i int) {
				m := models[i]
				hi, chosen := i == st.hi, m.ID == *value
				row := ui.Row(c).MinHeight(core.Density(c).ControlHeight()).Disabled(m.Unavailable != "").Description(m.Unavailable).
					Padding(u*2, u*3, u*2, u*2).Gap(u * 2).AlignItems(ui.Start).Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(m.Name).Cursor(ui.CursorPointer).Selected(chosen)
				if chosen {
					row.Background(k.Selection).TextColor(k.Text)
				} else if hi || row.Hovered() {
					row.Background(k.SurfaceHover)
				}
				if hi {
					row.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 2, H: r.H - 8}, k.AccentText, 1) })
				}
				row.Children(func() { input.SelectCheck(c, chosen); modelSelectorRow(c, m) })
				if row.Clicked() && !row.IsDisabled() {
					choose(i)
				}
			}).Fill()
		})
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: model_selector changed", "model", *value)
	}
	return core.NewChoiceResult(tr, changed)
}

type modelSelectorState struct {
	open       bool
	hi, lastHi int
	list       ui.ListState
}

func modelSelectorIndex(models []ModelOption, id string) int {
	first := -1
	for i, m := range models {
		if m.Unavailable != "" {
			continue
		}
		if first < 0 {
			first = i
		}
		if m.ID == id {
			return i
		}
	}
	return first
}

// modelSelectorRow shows a model's name, tags, context and unavailability.
func modelSelectorRow(c *ui.Context, m ModelOption) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	ui.Column(c).Gap(u).Grow(1).MinWidth(0).Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, m.Name).FontWeight(600).SingleLine()
			if m.ContextTokens > 0 {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "model_selector.context"), modelSelectorTokens(m.ContextTokens))).
					FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
			}
		})
		if len(m.Capabilities) > 0 {
			ui.Row(c).Gap(u).Wrap().Children(func() {
				for _, cap := range m.Capabilities {
					ui.Text(c, cap).FontSize(small).Padding(0, u*1.5).Radius(theme.ControlRadius).
						Border(theme.BorderWidth, k.Border).TextColor(k.AccentText).SingleLine()
				}
			})
		}
		if m.Unavailable != "" {
			ui.Text(c, m.Unavailable).FontSize(small).TextColor(k.TextMuted).Italic()
		}
	})
}

func modelSelectorFind(models []ModelOption, id string) (ModelOption, bool) {
	for _, m := range models {
		if m.ID == id {
			return m, true
		}
	}
	return ModelOption{}, false
}

// modelSelectorTokens writes a token count as 8K, 200K or 1M.
func modelSelectorTokens(n int) string {
	switch {
	case n >= 1_000_000 && n%100_000 == 0:
		return strconv.FormatFloat(float64(n)/1e6, 'f', -1, 64) + "M"
	case n >= 1000:
		return strconv.Itoa(n/1000) + "K"
	}
	return strconv.Itoa(n)
}
