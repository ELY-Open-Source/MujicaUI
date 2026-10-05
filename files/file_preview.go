package files

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/code"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/layout"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
	Status data.DataStatus
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
	viewer code.CodeViewerState
	text   string
	lines  []string
}

// FilePreview presents images, text, code or unsupported-type copy.
func FilePreview(c *ui.Context, f PreviewFile, opts FilePreviewOptions) FilePreviewResult {
	if f.Name == "" {
		panic("mujica: FilePreview needs a Name")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res FilePreviewResult
	res.Element = code.CodePanel(c).Role(ui.RoleGroup).Label(fmt.Sprintf(core.Msg(c, "file_preview.label"), f.Name))
	codeState := ui.Local(res.Element, "file_preview.code", func() previewCodeState { return previewCodeState{} })
	if res.Element.Shortcut(0, ui.KeyEscape) && !res.Element.IsDisabled() {
		res.Closed = true
	}
	res.Element.Children(func() {
		code.CodeHeader(c, func() {
			FileIcon(c, f.Name, FileIconOptions{})
			ui.Text(c, f.Name).FontWeight(600).SingleLine().MinWidth(0).Shrink(1)
			if f.Size > 0 {
				ui.Text(c, code.ProcessListBytes(uint64(f.Size))).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Shrink(0)
			}
			ui.Spacer(c)
			if input.Button(c, core.Msg(c, "file_preview.open"), input.ButtonOptions{Variant: input.Ghost}).MinHeight(theme.MinHitSize).Clicked() {
				res.Opened = true
			}
			if core.IconAction(c, icons.Must("x"), core.Msg(c, "file_preview.close")).Clicked() {
				res.Closed = true
			}
		})
		body := ui.Column(c).Grow(1).MinHeight(u * 40).Padding(u * 3)
		body.Children(func() {
			if opts.Status != data.DataReady {
				res.Retried = data.DataPlaceholder(c, opts.Status, opts.Error, "", "")
				return
			}
			lang := f.Language
			if lang == "" {
				lang = filePreviewLanguages[fileExt(f.Name)]
			}
			switch {
			case f.Image != nil:
				ui.Box(c).Grow(1).Center().Clip().Children(func() {
					display.Image(c, f.Image, display.ImageOptions{Alt: f.Name}).Element.MaxWidthPercent(100).MaxHeightPercent(100)
				})
			case lang != "" || FileKindOf(f.Name) == FileCode:
				if codeState.lines == nil || codeState.text != f.Text {
					codeState.text, codeState.lines = f.Text, strings.Split(f.Text, "\n")
				}
				if lang == "" {
					lang = fileExt(f.Name)
				}
				code.CodeViewer(c, codeState.lines, &codeState.viewer, code.CodeViewerOptions{Language: lang, Label: f.Name}).Element.Grow(1)
			case f.Text != "" || FileKindOf(f.Name) == FileText:
				layout.ScrollArea(c, nil, layout.ScrollAreaOptions{Label: f.Name}, func() {
					ui.Text(c, f.Text).Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize)).Selectable()
				}).Grow(1)
			default:
				ui.Column(c).Grow(1).Center().Gap(u * 2).Children(func() {
					FileIcon(c, f.Name, FileIconOptions{Size: 48})
					ui.Text(c, fmt.Sprintf(core.Msg(c, "file_preview.unsupported"), strings.ToUpper(fileExt(f.Name)))).TextColor(k.TextMuted).Role(ui.RoleStatus)
				})
			}
		})
	})
	if res.Closed || res.Opened {
		slog.Debug("mujica: file_preview action", "closed", res.Closed, "opened", res.Opened)
	}
	return res
}
