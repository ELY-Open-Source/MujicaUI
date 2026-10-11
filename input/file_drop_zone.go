package input

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DroppedFile is a file a FileDropZone accepted.
type DroppedFile struct {
	Path string
	Size int64
}

// FileDropZoneOptions configure a FileDropZone. Zero limits accept anything.
type FileDropZoneOptions struct {
	// Extensions lists the accepted extensions without the dot, matched case-insensitively.
	Extensions []string
	// MaxFiles caps how many files the list holds.
	MaxFiles int
	// MaxSize caps each file's size in bytes.
	MaxSize int64
	// Hint replaces the built-in "Drop files here".
	Hint     string
	Disabled bool
}

// fileDropZoneReject is a file that failed a check, shown until dismissed.
type fileDropZoneReject struct {
	path, reason string
}

type fileDropZoneStat struct {
	path string
	size int64
	dir  bool
	err  error
}

type fileDropZoneState struct {
	pending  []string
	stats    chan []fileDropZoneStat
	browse   chan filePickerResult
	rejected []fileDropZoneReject
}

// FileDropZone checks dropped or browsed files and lists them; it never reads or uploads them.
func FileDropZone(c *ui.Context, files *[]DroppedFile, opts FileDropZoneOptions) core.ChoiceResult {
	if opts.MaxFiles < 0 || opts.MaxSize < 0 {
		panic("mujica: FileDropZone limits must not be negative")
	}
	for _, ext := range opts.Extensions {
		if ext == "" || strings.HasPrefix(ext, ".") {
			panic(fmt.Sprintf("mujica: FileDropZone extension %q: give it without the dot", ext))
		}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	parts := core.ChoiceResult{}
	col := ui.Column(c).Gap(d.Unit() * 2)
	parts.Element = col
	st := ui.Local(col, "zone", func() fileDropZoneState { return fileDropZoneState{} })
	reject := func(path, key string, args ...any) {
		st.rejected = append(st.rejected, fileDropZoneReject{path, fmt.Sprintf(core.Msg(c, key), args...)})
		slog.Debug("mujica: file drop zone rejected", "reason", key)
	}
	// take checks types now and sizes off the UI thread.
	take := func(paths []string) {
		var ok []string
		for _, p := range paths {
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(p)), ".")
			switch {
			case slices.ContainsFunc(*files, func(f DroppedFile) bool { return f.Path == p }), slices.Contains(st.pending, p), slices.Contains(ok, p):
			case len(opts.Extensions) > 0 && !slices.ContainsFunc(opts.Extensions, func(e string) bool { return strings.EqualFold(e, ext) }):
				reject(p, "file_drop_zone.type", strings.Join(opts.Extensions, ", "))
			default:
				ok = append(ok, p)
			}
		}
		if len(ok) == 0 {
			return
		}
		st.pending = append(st.pending, ok...)
		if st.stats != nil {
			return // the running check picks them up next
		}
		fileDropZoneCheck(c, st)
	}
	if st.stats != nil {
		select {
		case res := <-st.stats:
			st.stats = nil
			st.pending = st.pending[len(res):]
			for _, s := range res {
				switch {
				case s.err != nil:
					err := s.err
					if pe, ok := errors.AsType[*fs.PathError](err); ok {
						err = pe.Err
					}
					reject(s.path, "file_drop_zone.unreadable", err)
				case s.dir:
					reject(s.path, "file_drop_zone.folder")
				case opts.MaxSize > 0 && s.size > opts.MaxSize:
					reject(s.path, "file_drop_zone.size", fileDropZoneSize(opts.MaxSize))
				case opts.MaxFiles > 0 && len(*files) >= opts.MaxFiles:
					reject(s.path, "file_drop_zone.count", opts.MaxFiles)
				default:
					*files = append(slices.Clone(*files), DroppedFile{s.path, s.size})
					parts = core.NewChoiceResult(parts.Element, true)
				}
			}
			if len(st.pending) > 0 {
				fileDropZoneCheck(c, st)
			}
		default:
		}
	}
	if st.browse != nil {
		select {
		case r := <-st.browse:
			st.browse = nil
			if r.err != nil {
				reject("", "file_picker.failed", r.err)
			} else {
				take(r.paths)
			}
		default:
		}
	}
	col.Children(func() {
		zone := ui.Column(c).AlignItems(ui.Center).Justify(ui.Center).Gap(d.Unit() * 2).MinHeight(140).Padding(d.Unit() * 5).
			Radius(theme.CardRadius).Disabled(opts.Disabled)
		disabled := zone.IsDisabled()
		if dropped := zone.DroppedFiles(); dropped != nil && !disabled {
			slog.Debug("mujica: file drop zone dropped", "count", len(dropped))
			take(dropped)
		}
		over := zone.FileDragOver() && !disabled
		bg, border := k.Surface, k.ControlBorder
		if over {
			bg, border = k.Selection, k.AccentText
		}
		zone.Background(bg).Draw(func(p *ui.Painter, r ui.Rect) {
			p.StrokeDashed(r, border, theme.CardRadius, theme.BorderWidth*1.5)
		})
		hint := opts.Hint
		if hint == "" {
			hint = core.Msg(c, "file_drop_zone.hint")
		}
		zone.Children(func() {
			ui.Icon(c, icons.Must("upload")).Size(24, 24).TextColor(k.AccentText)
			ui.Text(c, hint).TextAlign(ui.Center)
			ui.Text(c, fileDropZoneRules(c, opts)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).TextAlign(ui.Center)
			b := Button(c, core.Msg(c, "file_drop_zone.browse"), ButtonOptions{Variant: Outline, Loading: st.browse != nil, Disabled: opts.Disabled})
			if b.Clicked() {
				ch := make(chan filePickerResult, 1)
				st.browse = ch
				dlg := mygo.OpenDialogOptions{Multiple: opts.MaxFiles != 1}
				if len(opts.Extensions) > 0 {
					dlg.Filters = []mygo.FileFilter{{Name: hint, Extensions: opts.Extensions}}
				}
				go func() {
					p, err := filePickerOpen(dlg)
					ch <- filePickerResult{p, err}
					c.Invalidate()
				}()
			}
		})
		fileDropZoneList(c, st, files, &parts, opts.Disabled)
	})
	return parts
}

// fileDropZoneCheck stats the pending files on a goroutine.
func fileDropZoneCheck(c *ui.Context, st *fileDropZoneState) {
	ch := make(chan []fileDropZoneStat, 1)
	st.stats = ch
	paths := slices.Clone(st.pending)
	go func() {
		res := make([]fileDropZoneStat, len(paths))
		for i, p := range paths {
			info, err := os.Stat(p)
			res[i] = fileDropZoneStat{path: p, err: err}
			if err == nil {
				res[i].size, res[i].dir = info.Size(), info.IsDir()
			}
		}
		ch <- res
		c.Invalidate()
	}()
}

// fileDropZoneList shows accepted, checking and rejected files.
func fileDropZoneList(c *ui.Context, st *fileDropZoneState, files *[]DroppedFile, parts *core.ChoiceResult, disabled bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	row := func(key string, icon string, color ui.Color, name, detail string, remove func()) {
		ui.Row(c).Key(key).Gap(d.Unit()*2).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).
			BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
			ui.Icon(c, icons.Must(icon)).Size(16, 16).TextColor(color)
			ui.Text(c, name).SingleLine().Ellipsis("…").MinWidth(0).Shrink(1)
			ui.Text(c, detail).TextColor(color).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine().Grow(1)
			if remove != nil {
				x := ui.ButtonBase(c).Size(28, 28).Radius(theme.ControlRadius).TextColor(k.TextMuted).
					Label(fmt.Sprintf(core.Msg(c, "file_drop_zone.remove"), name)).Disabled(disabled)
				internal.Ring(x, k.Focus, theme.ControlRadius)
				if x.Hovered() {
					x.Background(k.SurfaceHover).TextColor(k.Text)
				}
				x.Children(func() { ui.Icon(c, icons.Must("x")).Size(14, 14) })
				if x.Clicked() {
					remove()
				}
			}
		})
	}
	ui.Column(c).Children(func() {
		for i, f := range *files {
			row("ok"+f.Path, "file", k.TextMuted, filepath.Base(f.Path), fileDropZoneSize(f.Size), func() {
				*files = slices.Delete(slices.Clone(*files), i, i+1)
				*parts = core.NewChoiceResult(parts.Element, true)
				slog.Debug("mujica: file drop zone removed", "count", len(*files))
			})
		}
		for _, p := range st.pending {
			row("pending"+p, "loader-circle", k.TextMuted, filepath.Base(p), core.Msg(c, "file_drop_zone.checking"), nil)
		}
		for i, r := range st.rejected {
			name := filepath.Base(r.path)
			if r.path == "" {
				name = core.Msg(c, "file_drop_zone.dialog")
			}
			row(fmt.Sprint("bad", i, r.path), "circle-alert", k.Danger, name, r.reason, func() {
				st.rejected = slices.Delete(slices.Clone(st.rejected), i, i+1)
			})
		}
	})
}

// fileDropZoneRules describes the limits, such as "PDF, PNG · up to 3 files · 2 MB each".
func fileDropZoneRules(c *ui.Context, opts FileDropZoneOptions) string {
	var rules []string
	if len(opts.Extensions) > 0 {
		rules = append(rules, strings.ToUpper(strings.Join(opts.Extensions, ", ")))
	}
	if opts.MaxFiles > 0 {
		rules = append(rules, fmt.Sprintf(core.Msg(c, "file_drop_zone.max_files"), opts.MaxFiles))
	}
	if opts.MaxSize > 0 {
		rules = append(rules, fmt.Sprintf(core.Msg(c, "file_drop_zone.max_size"), fileDropZoneSize(opts.MaxSize)))
	}
	if len(rules) == 0 {
		return core.Msg(c, "file_drop_zone.any")
	}
	return strings.Join(rules, " · ")
}

// fileDropZoneSize formats n bytes in binary units.
func fileDropZoneSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	v, i := float64(n)/unit, 0
	for v >= unit && i < 3 {
		v /= unit
		i++
	}
	return fmt.Sprintf("%.1f %s", v, []string{"KB", "MB", "GB", "TB"}[i])
}
