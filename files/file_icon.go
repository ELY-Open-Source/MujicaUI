package files

import (
	"path"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// FileKind is the family a file belongs to.
type FileKind int

const (
	FileOther FileKind = iota
	FileFolder
	FileImage
	FileText
	FileCode
	FileDocument
	FileArchive
	FileAudio
	FileVideo
)

var fileKinds = map[string]FileKind{
	"png": FileImage, "jpg": FileImage, "jpeg": FileImage, "gif": FileImage, "svg": FileImage, "webp": FileImage, "bmp": FileImage, "heic": FileImage,
	"txt": FileText, "md": FileText, "log": FileText, "csv": FileText, "rtf": FileText,
	"go": FileCode, "js": FileCode, "ts": FileCode, "py": FileCode, "json": FileCode, "sh": FileCode, "sql": FileCode, "html": FileCode, "css": FileCode, "rs": FileCode, "c": FileCode, "h": FileCode, "java": FileCode, "yaml": FileCode, "yml": FileCode, "toml": FileCode,
	"pdf": FileDocument, "doc": FileDocument, "docx": FileDocument, "xls": FileDocument, "xlsx": FileDocument, "ppt": FileDocument, "pptx": FileDocument, "pages": FileDocument, "key": FileDocument,
	"zip": FileArchive, "tar": FileArchive, "gz": FileArchive, "tgz": FileArchive, "7z": FileArchive, "rar": FileArchive, "xz": FileArchive,
	"mp3": FileAudio, "wav": FileAudio, "flac": FileAudio, "aac": FileAudio, "ogg": FileAudio, "m4a": FileAudio,
	"mp4": FileVideo, "mov": FileVideo, "mkv": FileVideo, "avi": FileVideo, "webm": FileVideo,
}

// FileKindOf classifies a filename by its extension.
func FileKindOf(name string) FileKind {
	return fileKinds[fileExt(name)]
}

// fileExt is the lower-case extension without the dot.
func fileExt(name string) string {
	return strings.ToLower(strings.TrimPrefix(path.Ext(name), "."))
}

// FileIconOptions configure a FileIcon.
type FileIconOptions struct {
	// Folder draws a folder; Open draws it opened.
	Folder, Open bool
	// Size is the icon's edge in DIPs; zero is 20.
	Size float32
}

// fileFolderOpen is an opened folder glyph in Lucide's style.
var fileFolderOpen = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 14 1.5-2.9A2 2 0 0 1 9.24 10H20a2 2 0 0 1 1.94 2.5l-1.54 6a2 2 0 0 1-1.95 1.5H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h3.9a2 2 0 0 1 1.69.9l.81 1.2a2 2 0 0 0 1.67.9H18a2 2 0 0 1 2 2v2"/></svg>`))

// fileKindColor is a kind's tint, readable on Surface.
func fileKindColor(k theme.Tokens, kind FileKind) ui.Color {
	switch kind {
	case FileFolder:
		return k.Warning
	case FileImage, FileVideo:
		return k.Info
	case FileCode:
		return k.AccentText
	case FileArchive, FileAudio:
		return k.Warning
	case FileDocument:
		return k.Success
	}
	return k.TextMuted
}

// FileIcon shows file type or folder opening state.
func FileIcon(c *ui.Context, name string, opts FileIconOptions) ui.Element {
	if opts.Size < 0 {
		panic("mujica: FileIcon Size is negative")
	}
	if opts.Open && !opts.Folder {
		panic("mujica: FileIcon Open needs Folder")
	}
	k := core.Tokens(c)
	size := opts.Size
	if size == 0 {
		size = 20
	}
	kind := FileKindOf(name)
	label := core.Msg(c, "file_icon.file")
	if ext := fileExt(name); ext != "" {
		label = strings.ToUpper(ext) + " " + label
	}
	if opts.Folder {
		kind, label = FileFolder, core.Msg(c, "file_icon.folder")
		if opts.Open {
			label = core.Msg(c, "file_icon.open")
		}
	}
	col := fileKindColor(k, kind)
	box := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleImage).Label(label).TextColor(col)
	return box.Children(func() {
		switch {
		case opts.Open:
			ui.Icon(c, fileFolderOpen).Size(size, size)
		case opts.Folder:
			ui.Icon(c, icons.Must("folder")).Size(size, size)
		case kind == FileImage:
			ui.Icon(c, icons.Must("image")).Size(size, size)
		case kind == FileText:
			ui.Icon(c, icons.Must("file-text")).Size(size, size)
		default:
			ui.Icon(c, icons.Must("file")).Size(size, size)
			ext := strings.ToUpper(fileExt(name))
			if ext != "" && size >= 20 {
				// The extension tag sits across the page's lower half.
				if len(ext) > 3 {
					ext = ext[:3]
				}
				ui.Text(c, ext).Font(theme.MonoFont).FontSize(core.FontSize(c, size*0.3)).FontWeight(700).TextColor(k.Surface).
					Background(col).Radius(1).Padding(0, 1).Absolute().Bottom(size * 0.12).Left(0).Role(ui.RoleNone)
			}
		}
	})
}
