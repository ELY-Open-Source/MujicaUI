package devtools

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// KeyValue holds an editable pair and enable flag.
type KeyValue struct {
	Key, Value string
	Enabled    bool
}

// KeyValueInputOptions configure a KeyValueInput.
type KeyValueInputOptions struct {
	// KeyPlaceholder and ValuePlaceholder default to "Key" and "Value".
	KeyPlaceholder, ValuePlaceholder string
	// Label names the editor for assistive technology.
	Label    string
	Disabled bool
	// MaxHeight caps the rows; zero means 320.
	MaxHeight float32
}

// KeyValueInput edits virtualized pairs with enable, add and remove.
func KeyValueInput(c *ui.Context, rows *[]KeyValue, opts KeyValueInputOptions) core.ChoiceResult {
	if rows == nil || opts.MaxHeight < 0 {
		panic("mujica: KeyValueInput needs rows")
	}
	kp, vp := opts.KeyPlaceholder, opts.ValuePlaceholder
	if kp == "" {
		kp = core.Msg(c, "key_value_input.key")
	}
	if vp == "" {
		vp = core.Msg(c, "key_value_input.value")
	}
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 2).Disabled(opts.Disabled).Label(opts.Label)
	locked := r.Element.IsDisabled()
	if locked {
		copy := slices.Clone(*rows)
		rows = &copy
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 320
	}
	remove := -1
	r.Element.Children(func() {
		agent.AgentRows(c, len(*rows), func(i int) {
			kv := &(*rows)[i]
			ui.Row(c.Key(i)).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				name := fmt.Sprintf(core.Msg(c, "key_value_input.row"), i+1)
				on := kv.Enabled
				box := input.Checkbox(c, &on, "", input.CheckboxOptions{Disabled: locked}).Label(name)
				if box.Changed() && !locked {
					kv.Enabled, r = on, core.NewChoiceResult(r.Element, true)
				}
				k := input.TextInput(c, &kv.Key, input.TextInputOptions{Placeholder: kp, Label: name + " " + kp, Disabled: locked})
				k.Element.Grow(1).Basis(0)
				v := input.TextInput(c, &kv.Value, input.TextInputOptions{Placeholder: vp, Label: name + " " + vp, Disabled: locked})
				v.Element.Grow(2).Basis(0)
				r = core.NewChoiceResult(r.Element, r.Changed() || k.Changed() || v.Changed())
				if core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "key_value_input.remove"), i+1)).Disabled(locked).Clicked() && !locked {
					remove = i
				}
			})
		}).Shrink(1).MaxHeight(maxH)
		if input.Button(c, core.Msg(c, "key_value_input.add"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("plus"), Disabled: locked}).AlignSelf(ui.Start).Clicked() && !locked {
			*rows = append(*rows, KeyValue{Enabled: true})
			r = core.NewChoiceResult(r.Element, true)
		}
	})
	if remove >= 0 {
		*rows = append((*rows)[:remove:remove], (*rows)[remove+1:]...)
		r = core.NewChoiceResult(r.Element, true)
	}
	if r.Changed() {
		slog.Debug("mujica: key_value_input change", "rows", len(*rows))
	}
	return r
}

// keyValueEnabled returns the enabled pairs with a key.
func keyValueEnabled(rows []KeyValue) []KeyValue {
	var out []KeyValue
	for _, kv := range rows {
		if kv.Enabled && kv.Key != "" {
			out = append(out, kv)
		}
	}
	return out
}
