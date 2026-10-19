package input

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// FilePickerOptions configure a FilePicker.
type FilePickerOptions struct {
	// Directory picks folders instead of files.
	Directory bool
	Multiple  bool
	// Filters limit the files offered by extension, without the dot.
	Filters []mygo.FileFilter
	// Title is the system dialog's title.
	Title string
	// Label is the button text; empty uses the built-in "Choose file…" or "Choose folder…".
	Label string
	// Parent attaches the dialog to a window, as a sheet on macOS.
	Parent   *mygo.Window
	Disabled bool
}

// FilePickerParts are the parts of a FilePicker.
type FilePickerParts struct {
	Element  ui.Element
	changed  bool
	canceled bool
}

// Changed reports whether the user picked new paths this frame.
func (p FilePickerParts) Changed() bool { return p.changed }

// Canceled reports whether the user closed the dialog without picking, which is not an error.
func (p FilePickerParts) Canceled() bool { return p.canceled }

// filePickerOpen shows the system dialog; tests replace it.
var filePickerOpen = mygo.Dialog.Open

type filePickerResult struct {
	paths []string
	err   error
}

type filePickerState struct {
	result chan filePickerResult
	err    string
}

// FilePicker opens the system file dialog off the UI thread and stores the chosen paths.
func FilePicker(c *ui.Context, paths *[]string, opts FilePickerOptions) FilePickerParts {
	for _, f := range opts.Filters {
		for _, ext := range f.Extensions {
			if ext == "" || strings.HasPrefix(ext, ".") {
				panic(fmt.Sprintf("mujica: FilePicker extension %q: give it without the dot", ext))
			}
		}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	parts := FilePickerParts{}
	col := ui.Column(c).Gap(d.Unit() * 2).AlignItems(ui.Start)
	parts.Element = col
	st := ui.Local(col, "picker", func() filePickerState { return filePickerState{} })
	if st.result != nil {
		select {
		case r := <-st.result:
			st.result = nil
			switch {
			case r.err != nil:
				st.err = fmt.Sprintf(core.Msg(c, "file_picker.failed"), r.err)
				slog.Debug("mujica: file picker failed", "err", r.err)
			case r.paths == nil:
				parts.canceled = true
				slog.Debug("mujica: file picker canceled")
			default:
				*paths = r.paths
				parts.changed = true
				slog.Debug("mujica: file picker picked", "count", len(r.paths))
			}
		default:
		}
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "file_picker.choose_file")
		if opts.Directory {
			label = core.Msg(c, "file_picker.choose_folder")
		}
	}
	icon := icons.Must("file")
	if opts.Directory {
		icon = icons.Must("folder")
	}
	col.Children(func() {
		b := Button(c, label, ButtonOptions{Variant: Secondary, Icon: icon, Loading: st.result != nil, Disabled: opts.Disabled})
		if b.Clicked() {
			st.err = ""
			ch := make(chan filePickerResult, 1)
			st.result = ch
			dlg := mygo.OpenDialogOptions{Parent: opts.Parent, Title: opts.Title, Filters: opts.Filters,
				Directory: opts.Directory, Multiple: opts.Multiple}
			slog.Debug("mujica: file picker opened", "directory", opts.Directory, "multiple", opts.Multiple)
			go func() {
				p, err := filePickerOpen(dlg)
				ch <- filePickerResult{p, err}
				c.Invalidate()
			}()
		}
		switch {
		case st.err != "":
			ui.Text(c, st.err).TextColor(k.Danger).FontSize(core.FontSize(c, theme.CaptionSize)).Role(ui.RoleStatus)
		case len(*paths) == 0:
			ui.Text(c, core.Msg(c, "file_picker.none")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
		default:
			ui.Column(c).Gap(d.Unit()).Children(func() {
				for _, p := range *paths {
					ui.Row(c.Key(p)).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon).Size(14, 14).TextColor(k.TextMuted)
						ui.Text(c, filepath.Base(p)).SingleLine()
						ui.Text(c, filepath.Dir(p)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine().Ellipsis("…").MinWidth(0).Shrink(1)
					})
				}
			})
		}
	})
	return parts
}
