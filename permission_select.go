package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
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
	return msg(c, [...]string{"permission_select.view", "permission_select.comment", "permission_select.edit"}[p])
}

// PermissionSelect chooses can view, can comment or can edit.
func PermissionSelect(c *ui.Context, value *Permission, opts PermissionSelectOptions) ChoiceResult {
	if value == nil {
		panic("mujica: PermissionSelect needs a value")
	}
	permissionName(c, *value)
	options := make([]SelectOption[Permission], 3)
	for p := PermissionView; p <= PermissionEdit; p++ {
		options[p] = SelectOption[Permission]{Value: p, Label: permissionName(c, p)}
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "permission_select.label")
	}
	r := Select(c, value, options, SelectOptions{Label: label, Disabled: opts.Disabled})
	if r.Changed() {
		slog.Debug("mujica: permission_select change", "permission", *value)
	}
	return r
}
