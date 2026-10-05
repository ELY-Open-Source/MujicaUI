package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// KeyValue is one editable pair; disabled pairs are kept but not used.
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
}

// KeyValueInput edits rows of pairs with an enable switch, add and remove.
func KeyValueInput(c *ui.Context, rows *[]KeyValue, opts KeyValueInputOptions) ChoiceResult {
	if rows == nil {
		panic("mujica: KeyValueInput needs rows")
	}
	kp, vp := opts.KeyPlaceholder, opts.ValuePlaceholder
	if kp == "" {
		kp = msg(c, "key_value_input.key")
	}
	if vp == "" {
		vp = msg(c, "key_value_input.value")
	}
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 2).Disabled(opts.Disabled).Label(opts.Label)
	locked := r.Element.IsDisabled()
	remove := -1
	r.Element.Children(func() {
		for i := range *rows {
			kv := &(*rows)[i]
			ui.Row(c).Key(i).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				name := fmt.Sprintf(msg(c, "key_value_input.row"), i+1)
				on := kv.Enabled
				box := Checkbox(c, &on, "", CheckboxOptions{Disabled: locked}).Label(name)
				if box.Changed() && !locked {
					kv.Enabled, r.changed = on, true
				}
				k := TextInput(c, &kv.Key, TextInputOptions{Placeholder: kp, Label: name + " " + kp, Disabled: locked})
				k.Element.Grow(1).Basis(0)
				v := TextInput(c, &kv.Value, TextInputOptions{Placeholder: vp, Label: name + " " + vp, Disabled: locked})
				v.Element.Grow(2).Basis(0)
				r.changed = r.changed || k.Changed() || v.Changed()
				if iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "key_value_input.remove"), i+1)).Disabled(locked).Clicked() && !locked {
					remove = i
				}
			})
		}
		if Button(c, msg(c, "key_value_input.add"), ButtonOptions{Variant: Ghost, Icon: icons.Must("plus"), Disabled: locked}).AlignSelf(ui.Start).Clicked() && !locked {
			*rows = append(*rows, KeyValue{Enabled: true})
			r.changed = true
		}
	})
	if remove >= 0 {
		*rows = append((*rows)[:remove:remove], (*rows)[remove+1:]...)
		r.changed = true
	}
	if r.changed {
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
