package mujica

import (
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// FileMessageOptions describe a file sent in a chat.
type FileMessageOptions struct {
	// Name is the file name; required.
	Name string
	// Size is the length in bytes.
	Size int64
	// Type is the MIME type; empty guesses from the name's extension.
	Type string
}

// FileMessageResult is a FileMessage's element and this frame's actions.
type FileMessageResult struct {
	Element            *ui.Element
	opened, downloaded bool
}

// Opened reports a click or Enter on the file.
func (r FileMessageResult) Opened() bool { return r.opened }

// Downloaded reports a click on the download button.
func (r FileMessageResult) Downloaded() bool { return r.downloaded }

// FileMessage shows a file's type icon, name and size, with open and download actions.
func FileMessage(c *ui.Context, opts FileMessageOptions) FileMessageResult {
	if opts.Name == "" {
		panic("mujica: FileMessage needs a Name")
	}
	if opts.Size < 0 {
		panic("mujica: FileMessage Size must not be negative")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	size := fileMessageSize(c, opts.Size)
	var r FileMessageResult
	r.Element = ui.Row(c).Gap(u).AlignItems(ui.Center).Padding(u, u, u, u*2).MaxWidth(360).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		open := ui.ButtonBase(c).Grow(1).MinWidth(0).Gap(u * 3).Padding(u).Justify(ui.Start).Radius(theme.ControlRadius).
			Label(fmt.Sprintf(msg(c, "file_message.open"), opts.Name, size))
		internal.Ring(open, k.Focus, theme.ControlRadius)
		if open.Hovered() {
			open.Background(k.SurfaceHover)
		}
		open.Children(func() {
			ui.Icon(c, icons.Must(fileMessageIcon(opts.Name, opts.Type))).Size(24, 24).Shrink(0).TextColor(k.AccentText)
			ui.Column(c).MinWidth(0).Grow(1).Children(func() {
				ui.Text(c, opts.Name).SingleLine().FontWeight(600)
				ui.Text(c, size).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			})
		})
		if open.Clicked() && !disabled {
			r.opened = true
			slog.Debug("mujica: file_message open")
		}
		if iconAction(c, icons.Must("download"), fmt.Sprintf(msg(c, "file_message.download"), opts.Name)).Clicked() && !disabled {
			r.downloaded = true
			slog.Debug("mujica: file_message download")
		}
	})
	return r
}

// fileMessageSize formats n bytes in binary units.
func fileMessageSize(c *ui.Context, n int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	v, i := float64(n), 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf(msg(c, "file_message.bytes"), n)
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

// fileMessageKinds maps extensions to type icons.
var fileMessageKinds = func() map[string]string {
	m := map[string]string{}
	for icon, exts := range map[string]string{
		"image":        ".png .jpg .jpeg .gif .webp .svg",
		"file-archive": ".zip .tar .gz .tgz .7z .rar",
		"file-code":    ".go .js .ts .py .json .sh .sql .html .css .rs .java .c .h",
		"file-text":    ".txt .md .pdf .doc .docx .rtf",
	} {
		for _, e := range strings.Fields(exts) {
			m[e] = icon
		}
	}
	return m
}()

// fileMessageIcon picks a type icon from the MIME type, else the extension.
func fileMessageIcon(name, mime string) string {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return "image"
	case strings.Contains(mime, "zip"):
		return "file-archive"
	case strings.HasPrefix(mime, "text/"):
		return "file-text"
	}
	if icon, ok := fileMessageKinds[strings.ToLower(path.Ext(name))]; ok {
		return icon
	}
	return "file"
}
