package files

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestFileKindOf(t *testing.T) {
	for name, want := range map[string]FileKind{"a.PNG": FileImage, "main.go": FileCode, "notes.md": FileText, "x.tar.gz": FileArchive,
		"song.flac": FileAudio, "film.mov": FileVideo, "deed.pdf": FileDocument, "README": FileOther, "blob.xyz": FileOther} {
		if got := FileKindOf(name); got != want {
			t.Errorf("%s: kind %d, want %d", name, got, want)
		}
	}
}

func TestFileIconLabels(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Children(func() {
			FileIcon(c, "main.go", FileIconOptions{})
			FileIcon(c, "src", FileIconOptions{Folder: true})
			FileIcon(c, "lib", FileIconOptions{Folder: true, Open: true, Size: 32})
		})
	}, 300, 100)
	for _, s := range []string{"GO file", "Folder", "Open folder"} {
		if _, ok := tt.Find(s); !ok {
			t.Fatalf("no icon named %q", s)
		}
	}
	if r, _ := tt.Find("Open folder"); r.W != 32 {
		t.Fatalf("size %v", r.W)
	}
}

func TestFileIconColorsReadable(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for kind := FileOther; kind <= FileVideo; kind++ {
			if r := theme.Contrast(fileKindColor(k, kind), k.Surface); r < 3 {
				t.Errorf("dark=%v kind %d contrast %.2f", k.Dark, kind, r)
			}
		}
	}
}

func TestFileIconBadOptionsPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Open without Folder did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileIcon(c, "a", FileIconOptions{Open: true})
	}, 100, 100)
}

func TestFileIconSnapshot(t *testing.T) {
	snapshot(t, "file_icon", 500, 70, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(14).Children(func() {
			FileIcon(c, "src", FileIconOptions{Folder: true, Size: 32})
			FileIcon(c, "src", FileIconOptions{Folder: true, Open: true, Size: 32})
			for _, n := range []string{"seal.go", "notes.md", "crest.png", "deed.pdf", "ledger.zip", "hymn.mp3", "rite.mov", "blob"} {
				FileIcon(c, n, FileIconOptions{Size: 32})
			}
		})
	})
}
