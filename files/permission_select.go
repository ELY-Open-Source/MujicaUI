package files

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// Permission defines access to a shared item.
type Permission int

const (
	PermissionView Permission = iota
	PermissionComment
	PermissionEdit
)

// PermissionSelectOptions configure a PermissionSelect.
type PermissionSelectOptions struct {
	// Label names the select, as the person it is for.
	Label    string
	Disabled bool
}

// permissionName is a permission's built-in name.
func permissionName(c *ui.Context, p Permission) string {
	if p < PermissionView || p > PermissionEdit {
		panic("mujica: Permission out of range")
	}
	return core.Msg(c, [...]string{"permission_select.view", "permission_select.comment", "permission_select.edit"}[p])
}

// PermissionSelect chooses can view, can comment or can edit.
func PermissionSelect(c *ui.Context, value *Permission, opts PermissionSelectOptions) core.ChoiceResult {
	if value == nil {
		panic("mujica: PermissionSelect needs a value")
	}
	permissionName(c, *value)
	options := make([]input.SelectOption[Permission], 3)
	for p := PermissionView; p <= PermissionEdit; p++ {
		options[p] = input.SelectOption[Permission]{Value: p, Label: permissionName(c, p)}
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "permission_select.label")
	}
	r := input.Select(c, value, options, input.SelectOptions{Label: label, Disabled: opts.Disabled})
	if r.Changed() {
		slog.Debug("mujica: permission_select change", "permission", *value)
	}
	return r
}
