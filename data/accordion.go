package data

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// AccordionItem is one section of an Accordion.
type AccordionItem struct {
	// ID identifies the section in the open list; required, unique.
	ID       string
	Title    string
	Content  func()
	Disabled bool
}

// AccordionOptions configure an Accordion.
type AccordionOptions struct {
	// Multiple lets several sections stay open; otherwise opening one closes the rest.
	Multiple bool
	// KeepOne stops the last open section from closing.
	KeepOne bool
}

// Accordion shows sections whose open IDs live in *open, which the caller controls.
func Accordion(c *ui.Context, open *[]string, items []AccordionItem, opts AccordionOptions) core.ChoiceResult {
	seen := map[string]bool{}
	for _, it := range items {
		if it.ID == "" || seen[it.ID] || it.Content == nil {
			panic("mujica: Accordion item needs a unique ID and Content: " + it.ID)
		}
		seen[it.ID] = true
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	var v core.ChoiceResult
	v.Element = ui.Column(c).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius)
	focusTo := ui.Local(v.Element, "focus", func() int { return -1 })
	v.Element.Children(func() {
		for i, it := range items {
			ui.Column(c.Key(it.ID)).Children(func() {
				if i > 0 {
					ui.Box(c).Height(1).FillWidth().Background(k.Border)
				}
				isOpen := slices.Contains(*open, it.ID)
				was := isOpen
				p := ui.CollapsibleBase(c, &isOpen)
				tr := p.Trigger.Disabled(it.Disabled).AlignItems(ui.Center).Gap(u*2).MinHeight(d.ControlHeight()+u*2).
					Padding(0, u*4).Label(it.Title)
				internal.Ring(tr, k.Focus, theme.ControlRadius)
				if tr.Hovered() && !it.Disabled {
					tr.Background(k.SurfaceHover)
				}
				if *focusTo == i {
					tr.Focus()
					*focusTo = -1
				}
				if j := accordionKey(tr, items, i); j >= 0 {
					*focusTo = j
				}
				if isOpen != was && !it.Disabled && accordionSet(open, it.ID, isOpen, opts) {
					v = core.NewChoiceResult(v.Element, true)
				}
				isOpen = slices.Contains(*open, it.ID)
				prog := tr.Animate("mujica.open", collapsibleTarget(isOpen), core.Motion(c, theme.StateDuration))
				tr.Children(func() {
					ui.Text(c, it.Title).FontWeight(600).Grow(1).MinWidth(0)
					ui.Icon(c, icons.Must("chevron-down")).Size(16, 16).Rotate(180 * prog).TextColor(k.TextMuted).Shrink(0)
				})
				if panel := collapsiblePanel(c, isOpen, prog, it.Content); panel.Valid() {
					panel.Padding(0, u*4, u*4).Gap(u * 2)
				}
			})
		}
	})
	return v
}

// accordionKey returns the header the arrows, Home or End move to from header i, or -1.
func accordionKey(tr ui.Element, items []AccordionItem, i int) int {
	step, from := 0, i
	switch {
	case tr.Shortcut(0, ui.KeyDown):
		step = 1
	case tr.Shortcut(0, ui.KeyUp):
		step = -1
	case tr.Shortcut(0, ui.KeyHome):
		step, from = 1, -1
	case tr.Shortcut(0, ui.KeyEnd):
		step, from = -1, len(items)
	default:
		return -1
	}
	for j := from + step; j >= 0 && j < len(items); j += step {
		if !items[j].Disabled {
			return j
		}
	}
	return -1
}

// accordionSet applies a header toggle under the options and reports a change.
func accordionSet(open *[]string, id string, on bool, opts AccordionOptions) bool {
	switch {
	case on && opts.Multiple:
		*open = append(slices.Clone(*open), id)
	case on:
		*open = []string{id}
	case opts.KeepOne && len(*open) == 1:
		return false
	default:
		*open = slices.DeleteFunc(slices.Clone(*open), func(s string) bool { return s == id })
	}
	slog.Debug("mujica: accordion toggle", "id", id, "open", on)
	return true
}
