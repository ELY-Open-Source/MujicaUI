package mujica

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
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
	// Unavailable explains why the model cannot be chosen; empty means available.
	Unavailable string
}

// ModelSelectorOptions configure a ModelSelector.
type ModelSelectorOptions struct {
	Disabled bool
	// Label names the selector; empty uses the built-in "Model".
	Label string
}

// ModelSelector chooses a model ID into *value from a list showing capabilities, context length and why some are unavailable.
func ModelSelector(c *ui.Context, value *string, models []ModelOption, opts ModelSelectorOptions) ChoiceResult {
	for _, m := range models {
		if m.ID == "" || m.Name == "" || m.ContextTokens < 0 {
			panic("mujica: ModelOption needs ID, Name and a non-negative ContextTokens")
		}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	old := *value
	pick := old
	var s *ui.SelectParts[string]
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() { s = ui.SelectBase(c, &pick) })
	tr := s.Trigger
	label := opts.Label
	if label == "" {
		label = msg(c, "model_selector.label")
	}
	tr.Label(label)
	cur, has := modelSelectorFind(models, *value)
	selectField(c, tr, s.Open(), func() {
		if has {
			ui.Text(c, cur.Name).SingleLine().FontWeight(600).Grow(1).Shrink(1)
		} else {
			ui.Text(c, msg(c, "model_selector.placeholder")).SingleLine().TextColor(k.TextMuted).Grow(1).Shrink(1)
		}
	})
	if opts.Disabled {
		return ChoiceResult{Element: tr}
	}
	s.Popup(func(panel *ui.Element) {
		selectPanel(c, panel)
		panel.MinWidth(320)
		ui.Scroll(c).MaxHeight(360).Children(func() {
			for _, m := range models {
				var row *ui.Element
				if m.Unavailable != "" {
					row = ui.Row(c).Disabled(true).Description(m.Unavailable)
				} else {
					row = s.Item(m.ID)
				}
				hi := m.Unavailable == "" && row.Highlighted()
				chosen := m.ID == *value
				row.Padding(u*2, u*3, u*2, u*2).Gap(u * 2).AlignItems(ui.Start).Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(m.Name)
				switch {
				case chosen:
					row.Background(k.Selection)
				case hi:
					row.Background(k.SurfaceHover)
				}
				if hi {
					row.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 2, H: r.H - 8}, k.AccentText, 1) })
				}
				row.Children(func() {
					selectCheck(c, chosen)
					modelSelectorRow(c, m)
				})
			}
		})
	})
	if pick != *value && !tr.IsDisabled() {
		if m, ok := modelSelectorFind(models, pick); ok && m.Unavailable == "" {
			*value = pick
		}
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: model_selector changed", "model", *value)
	}
	return ChoiceResult{Element: tr, changed: changed}
}

// modelSelectorRow shows a model's name, tags, context and unavailability.
func modelSelectorRow(c *ui.Context, m ModelOption) {
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	ui.Column(c).Gap(u).Grow(1).MinWidth(0).Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, m.Name).FontWeight(600).SingleLine()
			if m.ContextTokens > 0 {
				ui.Text(c, fmt.Sprintf(msg(c, "model_selector.context"), modelSelectorTokens(m.ContextTokens))).
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
