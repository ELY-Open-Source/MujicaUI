package files

import (
	"log/slog"
	"path"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RenameInlineOptions configure a RenameInline.
type RenameInlineOptions struct {
	// Validate returns an error message, or empty for acceptance.
	Validate func(name string) string
	// Folder selects the whole name: folders have no extension.
	Folder   bool
	Disabled bool
}

// RenameInlineResult is a RenameInline and its outcome in this frame.
type RenameInlineResult struct {
	Element ui.Element
	// Renamed reports a saved name; Canceled reports an abandoned edit.
	Renamed, Canceled bool
}

// renameState is one editing session.
type renameState struct {
	active bool
	draft  string
	// gen rebuilds the editor with its caret at the end.
	gen int
	// tail appends text next frame behind the caret.
	tail    string
	shown   bool
	stemSel bool
	focused bool
	err     string
	folder  bool
}

// renameStem is the part of name before its extension.
func renameStem(name string, folder bool) string {
	ext := path.Ext(name)
	if folder || ext == name || ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}

// RenameInline edits names with validation and initial stem replacement.
// Enter validates; Escape cancels; blur confirms valid names.
func RenameInline(c *ui.Context, name *string, editing *bool, opts RenameInlineOptions) RenameInlineResult {
	if name == nil || editing == nil {
		panic("mujica: RenameInline needs a name and editing")
	}
	var res RenameInlineResult
	res.Element = ui.Row(c).MinWidth(0).MinHeight(theme.MinHitSize).AlignItems(ui.Center).Disabled(opts.Disabled)
	st := ui.Local(res.Element, "rename_inline", func() renameState { return renameState{} })
	if res.Element.IsDisabled() && *editing {
		*editing = false
		res.Canceled = st.active
		st.active = false
	}
	res.Element.Children(func() {
		if !*editing {
			st.active = false
			t := ui.Text(c, *name).SingleLine().MinWidth(0)
			if t.DoubleClicked() && !res.Element.IsDisabled() {
				*editing = true
				c.Invalidate()
			}
			return
		}
		stem := renameStem(*name, opts.Folder)
		if !st.active {
			// Start with the stem and its caret at the end.
			*st = renameState{active: true, draft: stem, gen: st.gen + 1, tail: (*name)[len(stem):], stemSel: true, folder: opts.Folder}
			slog.Debug("mujica: rename_inline start")
		}
		switch {
		case st.tail != "" && st.shown:
			st.draft, st.tail, st.shown = st.draft+st.tail, "", false
		case st.tail != "":
			// Show the head for one frame to place its caret.
			st.shown = true
			c.Invalidate()
		}
		finish := func(ok bool) {
			if ok {
				if why := renameCheck(c, st.draft, opts.Validate); why != "" {
					st.err = why
					return
				}
				res.Renamed = st.draft != *name
				*name = st.draft
			} else {
				res.Canceled = true
			}
			*editing, st.active = false, false
			slog.Debug("mujica: rename_inline end", "renamed", res.Renamed)
		}
		ui.Box(c.Key(st.gen)).Grow(1).MinWidth(0).Children(func() {
			in := input.TextInput(c, &st.draft, input.TextInputOptions{Label: core.Msg(c, "rename_inline.label"), Error: st.err})
			in.Element.FillWidth()
			core.AutoFocus(c, in.Input)
			if in.Changed() {
				st.err = ""
				renameTyped(st, *name, opts.Folder)
			}
			renameSelection(c, in.Input, st)
			switch {
			case in.Input.Shortcut(0, ui.KeyEscape):
				finish(false)
			case in.Submitted():
				finish(true)
			case in.Input.Focused():
				st.focused = true
			case st.focused && st.tail == "":
				// Confirm valid names when focus leaves.
				finish(renameCheck(c, st.draft, opts.Validate) == "")
			}
		})
	})
	return res
}

func renameCheck(c *ui.Context, name string, validate func(string) string) string {
	if strings.TrimSpace(name) == "" {
		return core.Msg(c, "rename_inline.empty")
	}
	if validate != nil {
		return validate(name)
	}
	return ""
}

// renameSelection marks the stem for replacement on initial input.
func renameSelection(c *ui.Context, in ui.Element, st *renameState) {
	k := core.Tokens(c)
	stem := renameStem(st.draft, st.folder)
	if st.stemSel && st.tail == "" {
		w, _ := c.MeasureText(0, ui.Span{Text: stem, Size: core.FontSize(c, theme.BodySize)})
		in.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 2, W: w, H: r.H - 4}, k.Selection, 2)
		})
	}
	disabled := in.IsDisabled()
	in.HandleInput(func(ev ui.InputEvent) bool {
		if !st.stemSel || disabled {
			return false
		}
		switch {
		case ev.Kind == ui.InputKeyDown && ev.Mods == 0:
			switch ev.Key {
			case ui.KeyBackspace, ui.KeyDelete:
				renameReplace(st, "")
				return true
			case ui.KeyLeft, ui.KeyHome, ui.KeyUp:
				// An empty head clamps the caret to the start.
				st.draft, st.tail, st.stemSel, st.shown = "", st.draft, false, false
				return true
			case ui.KeyRight:
				st.stemSel = false
				return true
			}
		}
		if ev.Kind == ui.InputPointerDown || ev.Kind == ui.InputKeyDown &&
			(ev.Key == ui.KeyLeft || ev.Key == ui.KeyRight || ev.Key == ui.KeyHome || ev.Key == ui.KeyEnd || ev.Key == ui.KeyUp || ev.Key == ui.KeyDown || ev.Key == ui.KeyA && ev.Mods&ui.Cmd != 0) {
			st.stemSel = false
		}
		return false
	})
}

// renameReplace substitutes the stem and rebuilds the caret.
func renameReplace(st *renameState, text string) {
	ext := st.draft[len(renameStem(st.draft, st.folder)):]
	st.draft, st.tail, st.gen, st.stemSel, st.shown = text, ext, st.gen+1, false, false
}

// renameTyped replaces the stem with newly inserted text.
func renameTyped(st *renameState, name string, folder bool) {
	if !st.stemSel {
		return
	}
	stem := renameStem(name, folder)
	ext := name[len(stem):]
	st.stemSel = false
	if strings.HasPrefix(st.draft, stem) && strings.HasSuffix(st.draft, ext) && len(st.draft) >= len(name) {
		typed := st.draft[len(stem) : len(st.draft)-len(ext)]
		st.draft = stem + ext
		renameReplace(st, typed)
		slog.Debug("mujica: rename_inline replace stem")
	}
}
