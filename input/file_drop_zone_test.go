package input

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// fileDropZoneFiles writes files of the given sizes into a temp dir.
func fileDropZoneFiles(t *testing.T, sizes map[string]int) string {
	t.Helper()
	dir := t.TempDir()
	for name, n := range sizes {
		if err := os.WriteFile(filepath.Join(dir, name), make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func fileDropZoneTester(files *[]DroppedFile, opts FileDropZoneOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if FileDropZone(c, files, opts).Changed() {
			*changes++
		}
	}, 560, 600)
}

// fileDropZoneBrowse answers the Browse dialog with paths.
func fileDropZoneBrowse(t *testing.T, paths ...string) *mygo.OpenDialogOptions {
	saved := filePickerOpen
	t.Cleanup(func() { filePickerOpen = saved })
	got := &mygo.OpenDialogOptions{}
	filePickerOpen = func(o mygo.OpenDialogOptions) ([]string, error) {
		*got = o
		return paths, nil
	}
	return got
}

func TestFileDropZoneChecksTypeSizeAndCount(t *testing.T) {
	dir := fileDropZoneFiles(t, map[string]int{"a.pdf": 10, "b.PDF": 3000, "c.pdf": 20, "d.txt": 5, "e.pdf": 30})
	p := func(n string) string { return filepath.Join(dir, n) }
	got := fileDropZoneBrowse(t, p("a.pdf"), p("b.PDF"), p("c.pdf"), p("d.txt"), p("e.pdf"), p("missing.pdf"), dir+"/sub.pdf")
	if err := os.Mkdir(p("sub.pdf"), 0o755); err != nil {
		t.Fatal(err)
	}
	var files []DroppedFile
	n := 0
	tt := fileDropZoneTester(&files, FileDropZoneOptions{Extensions: []string{"pdf"}, MaxFiles: 2, MaxSize: 1024}, &n)
	if err := tt.Click("Browse…"); err != nil {
		t.Fatal(err)
	}
	filePickerWait(t, tt, func() bool { return len(files) == 2 && tt.HasText("More than 2 files") })
	if files[0] != (DroppedFile{p("a.pdf"), 10}) || files[1] != (DroppedFile{p("c.pdf"), 20}) || n != 1 {
		t.Fatalf("files %+v changes %d", files, n)
	}
	for _, want := range []string{"Only pdf files", "Larger than 1.0 KB", "Folders are not accepted", "Cannot read: no such file or directory"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if !got.Multiple || got.Filters[0].Extensions[0] != "pdf" {
		t.Fatalf("dialog options %+v", got)
	}
}

func TestFileDropZoneRemoveAndDedupe(t *testing.T) {
	dir := fileDropZoneFiles(t, map[string]int{"score.mid": 2048})
	path := filepath.Join(dir, "score.mid")
	fileDropZoneBrowse(t, path)
	files := []DroppedFile{{path, 2048}}
	n := 0
	tt := fileDropZoneTester(&files, FileDropZoneOptions{}, &n)
	_ = tt.Click("Browse…")
	for range 20 {
		tt.Frame()
	}
	if len(files) != 1 || n != 0 || !tt.HasText("2.0 KB") {
		t.Fatalf("duplicate added: %+v", files)
	}
	if err := tt.Click("Remove score.mid"); err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 || n != 1 {
		t.Fatalf("remove: %+v changes %d", files, n)
	}
}

func TestFileDropZoneDisabled(t *testing.T) {
	called := false
	saved := filePickerOpen
	t.Cleanup(func() { filePickerOpen = saved })
	filePickerOpen = func(mygo.OpenDialogOptions) ([]string, error) { called = true; return nil, nil }
	files := []DroppedFile{{"/x/kept.pdf", 1}}
	n := 0
	tt := fileDropZoneTester(&files, FileDropZoneOptions{Disabled: true}, &n)
	_ = tt.Click("Browse…")
	_ = tt.Click("Remove kept.pdf")
	tt.Frame()
	if called || len(files) != 1 {
		t.Fatal("disabled zone acted")
	}
	if fileDropZoneSize(1536) != "1.5 KB" || fileDropZoneSize(3<<20) != "3.0 MB" || fileDropZoneSize(12) != "12 B" {
		t.Fatal("size format")
	}
}

func TestFileDropZoneSnapshot(t *testing.T) {
	files := []DroppedFile{{"/scores/Imprisoned XII.pdf", 482133}, {"/scores/cover.png", 1 << 20}}
	snapshot(t, "file_drop_zone", 520, 330, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			FileDropZone(c, &files, FileDropZoneOptions{Extensions: []string{"pdf", "png"}, MaxFiles: 3, MaxSize: 2 << 20})
		})
	})
}
