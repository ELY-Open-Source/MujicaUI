package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// PreviewFile holds caller-provided preview content.
type PreviewFile struct {
	Name string
	Size int64
	// Image previews a picture; Text previews text or code.
	Image ui.ImageSource
	Text  string
	// Language highlights code; empty guesses it from the extension.
	Language string
}

// FilePreviewOptions configure a FilePreview.
type FilePreviewOptions struct {
	// Status shows the caller's read in progress or failed.
	Status DataStatus
	Error  string
}

// FilePreviewResult exposes the panel and requested actions.
type FilePreviewResult struct {
	Element *ui.Element
	// Closed, Opened and Retried report requested actions.
	Closed, Opened, Retried bool
}

var filePreviewLanguages = map[string]string{"go": "go", "js": "javascript", "ts": "typescript", "py": "python", "json": "json", "sh": "shell", "sql": "sql"}

type previewCodeState struct {
	viewer CodeViewerState
	text   string
	lines  []string
}

// FilePreview presents images, text, code or unsupported-type copy.
func FilePreview(c *ui.Context, f PreviewFile, opts FilePreviewOptions) FilePreviewResult {
	if f.Name == "" {
		panic("mujica: FilePreview needs a Name")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var res FilePreviewResult
	res.Element = codePanel(c).Role(ui.RoleGroup).Label(fmt.Sprintf(msg(c, "file_preview.label"), f.Name))
	code := ui.Local(res.Element, "file_preview.code", func() previewCodeState { return previewCodeState{} })
	if res.Element.Shortcut(0, ui.KeyEscape) && !res.Element.IsDisabled() {
		res.Closed = true
	}
	res.Element.Children(func() {
		codeHeader(c, func() {
			FileIcon(c, f.Name, FileIconOptions{})
			ui.Text(c, f.Name).FontWeight(600).SingleLine().MinWidth(0).Shrink(1)
			if f.Size > 0 {
				ui.Text(c, processListBytes(uint64(f.Size))).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Shrink(0)
			}
			ui.Spacer(c)
			if Button(c, msg(c, "file_preview.open"), ButtonOptions{Variant: Ghost}).MinHeight(theme.MinHitSize).Clicked() {
				res.Opened = true
			}
			if iconAction(c, icons.Must("x"), msg(c, "file_preview.close")).Clicked() {
				res.Closed = true
			}
		})
		body := ui.Column(c).Grow(1).MinHeight(u * 40).Padding(u * 3)
		body.Children(func() {
			if opts.Status != DataReady {
				res.Retried = dataPlaceholder(c, opts.Status, opts.Error, "", "")
				return
			}
			lang := f.Language
			if lang == "" {
				lang = filePreviewLanguages[fileExt(f.Name)]
			}
			switch {
			case f.Image != nil:
				ui.Box(c).Grow(1).Center().Clip().Children(func() {
					Image(c, f.Image, ImageOptions{Alt: f.Name}).Element.MaxWidthPercent(100).MaxHeightPercent(100)
				})
			case lang != "" || FileKindOf(f.Name) == FileCode:
				if code.lines == nil || code.text != f.Text {
					code.text, code.lines = f.Text, strings.Split(f.Text, "\n")
				}
				if lang == "" {
					lang = fileExt(f.Name)
				}
				CodeViewer(c, code.lines, &code.viewer, CodeViewerOptions{Language: lang, Label: f.Name}).Element.Grow(1)
			case f.Text != "" || FileKindOf(f.Name) == FileText:
				ScrollArea(c, nil, ScrollAreaOptions{Label: f.Name}, func() {
					ui.Text(c, f.Text).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).Selectable()
				}).Grow(1)
			default:
				ui.Column(c).Grow(1).Center().Gap(u * 2).Children(func() {
					FileIcon(c, f.Name, FileIconOptions{Size: 48})
					ui.Text(c, fmt.Sprintf(msg(c, "file_preview.unsupported"), strings.ToUpper(fileExt(f.Name)))).TextColor(k.TextMuted).Role(ui.RoleStatus)
				})
			}
		})
	})
	if res.Closed || res.Opened {
		slog.Debug("mujica: file_preview action", "closed", res.Closed, "opened", res.Opened)
	}
	return res
}
