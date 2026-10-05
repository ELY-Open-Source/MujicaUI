package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// FileOperationKind is what a file operation does.
type FileOperationKind int

const (
	FileCopy FileOperationKind = iota
	FileMove
	FileDelete
)

// FileOperation is the caller's running operation; counts are of items.
type FileOperation struct {
	Kind        FileOperationKind
	Done, Total int
	// BytesDone and BytesTotal show the amount when BytesTotal is set.
	BytesDone, BytesTotal int64
	// Current is the item in progress.
	Current string
	// Conflict names an existing target awaiting a decision.
	Conflict string
}

// ConflictChoice answers a name conflict.
type ConflictChoice int

const (
	ConflictNone ConflictChoice = iota
	ConflictReplace
	ConflictKeepBoth
	ConflictSkip
)

// FileOperationResult is the panel and the user's answer this frame.
type FileOperationResult struct {
	Element  *ui.Element
	Canceled bool
	// Resolved answers a conflict; ApplyAll extends that decision.
	Resolved ConflictChoice
	ApplyAll bool
}

// FileOperationProgress displays progress, cancellation and conflict choices.
func FileOperationProgress(c *ui.Context, op FileOperation, opts FileOperationOptions) FileOperationResult {
	if op.Total < 0 || op.Done < 0 || op.Done > op.Total || op.BytesTotal < 0 || op.BytesDone < 0 || op.BytesTotal > 0 && op.BytesDone > op.BytesTotal {
		panic("mujica: FileOperation counts out of range")
	}
	if op.Kind < FileCopy || op.Kind > FileDelete {
		panic("mujica: FileOperation kind out of range")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var res FileOperationResult
	res.Element = codePanel(c).Padding(u * 4).Gap(u * 3).Disabled(opts.Disabled).Role(ui.RoleGroup)
	all := ui.Local(res.Element, "file_operation_progress.all", func() bool { return false })
	disabled := res.Element.IsDisabled()
	verb := msg(c, [...]string{"file_operation_progress.copying", "file_operation_progress.moving", "file_operation_progress.deleting"}[op.Kind])
	title := fmt.Sprintf(verb, op.Done, op.Total)
	res.Element.Label(title)
	frac := 0.0
	switch {
	case op.BytesTotal > 0:
		frac = float64(op.BytesDone) / float64(op.BytesTotal)
	case op.Total > 0:
		frac = float64(op.Done) / float64(op.Total)
	}
	res.Element.Children(func() {
		Progress(c, frac, ProgressOptions{Label: title, ShowValue: true})
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				if op.Current != "" {
					ui.Text(c, op.Current).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				}
				if op.BytesTotal > 0 {
					ui.Text(c, fmt.Sprintf(msg(c, "file_operation_progress.bytes"), processListBytes(uint64(op.BytesDone)), processListBytes(uint64(op.BytesTotal)))).
						FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				}
			})
			if Button(c, msg(c, "file_operation_progress.cancel"), ButtonOptions{Variant: Secondary}).Clicked() && !disabled {
				res.Canceled = true
			}
		})
		if op.Conflict == "" {
			return
		}
		ui.Column(c).Gap(u*2).Padding(u*3).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.Warning).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.Warning).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(16, 16)
				ui.Text(c, fmt.Sprintf(msg(c, "file_operation_progress.conflict"), op.Conflict)).TextColor(k.Text).Role(ui.RoleStatus)
			})
			Checkbox(c, all, msg(c, "file_operation_progress.all"), CheckboxOptions{Disabled: opts.Disabled})
			ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
				for _, b := range []struct {
					key string
					v   ButtonVariant
					cho ConflictChoice
				}{{"file_operation_progress.skip", Ghost, ConflictSkip}, {"file_operation_progress.keepBoth", Secondary, ConflictKeepBoth}, {"file_operation_progress.replace", Danger, ConflictReplace}} {
					if Button(c, msg(c, b.key), ButtonOptions{Variant: b.v}).Clicked() && !disabled {
						res.Resolved, res.ApplyAll = b.cho, *all
					}
				}
			})
		})
	})
	if res.Canceled || res.Resolved != ConflictNone {
		slog.Debug("mujica: file_operation_progress answer", "canceled", res.Canceled, "choice", res.Resolved, "all", res.ApplyAll)
		*all = false
	}
	return res
}

// FileOperationOptions configure a FileOperationProgress.
type FileOperationOptions struct {
	Disabled bool
}
