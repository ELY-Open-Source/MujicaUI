package input

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type filePickerProbe struct {
	paths            []string
	changes, cancels int
	got              mygo.OpenDialogOptions
	reply            []string
	err              error
	release          chan struct{}
}

func (p *filePickerProbe) stub(t *testing.T) {
	saved := filePickerOpen
	t.Cleanup(func() { filePickerOpen = saved })
	filePickerOpen = func(o mygo.OpenDialogOptions) ([]string, error) {
		p.got = o
		if p.release != nil {
			<-p.release
		}
		return p.reply, p.err
	}
}

func (p *filePickerProbe) tester(opts FilePickerOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := FilePicker(c, &p.paths, opts)
		if r.Changed() {
			p.changes++
		}
		if r.Canceled() {
			p.cancels++
		}
	}, 500, 200)
}

// filePickerWait draws frames until done holds.
func filePickerWait(t *testing.T, tt *ui.Tester, done func() bool) {
	t.Helper()
	for range 200 {
		tt.Frame()
		if done() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("dialog result never arrived")
}

func TestFilePickerPicksFiles(t *testing.T) {
	p := &filePickerProbe{reply: []string{"/court/score.pdf", "/court/cast.png"}, release: make(chan struct{})}
	p.stub(t)
	filters := []mygo.FileFilter{{Name: "Docs", Extensions: []string{"pdf", "png"}}}
	tt := p.tester(FilePickerOptions{Multiple: true, Filters: filters, Title: "Score"})
	if err := tt.Click("Choose file…"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Choose file…")
	close(p.release)
	filePickerWait(t, tt, func() bool { return p.changes > 0 })
	if !slices.Equal(p.paths, p.reply) || p.changes != 1 || !tt.HasText("score.pdf") {
		t.Fatalf("paths %q changes %d", p.paths, p.changes)
	}
	if !p.got.Multiple || p.got.Directory || p.got.Title != "Score" || p.got.Filters[0].Extensions[1] != "png" {
		t.Fatalf("dialog options %+v", p.got)
	}
}

func TestFilePickerCancelIsNotAnError(t *testing.T) {
	p := &filePickerProbe{paths: []string{"/kept"}}
	p.stub(t)
	tt := p.tester(FilePickerOptions{Directory: true})
	if err := tt.Click("Choose folder…"); err != nil {
		t.Fatal(err)
	}
	filePickerWait(t, tt, func() bool { return p.cancels > 0 })
	if p.changes != 0 || p.paths[0] != "/kept" || !p.got.Directory {
		t.Fatalf("cancel changed paths: %q", p.paths)
	}
}

func TestFilePickerReportsDialogError(t *testing.T) {
	p := &filePickerProbe{err: errors.New("no display")}
	p.stub(t)
	tt := p.tester(FilePickerOptions{})
	_ = tt.Click("Choose file…")
	want := "Could not open the file dialog: no display"
	filePickerWait(t, tt, func() bool { return tt.HasText(want) })
	if p.changes != 0 || p.cancels != 0 {
		t.Fatalf("error counted as result")
	}
}

func TestFilePickerDisabledAndBadFilter(t *testing.T) {
	p := &filePickerProbe{}
	p.stub(t)
	called := false
	filePickerOpen = func(mygo.OpenDialogOptions) ([]string, error) { called = true; return nil, nil }
	tt := p.tester(FilePickerOptions{Disabled: true})
	_ = tt.Click("Choose file…")
	tt.Frame()
	if called {
		t.Fatal("disabled picker opened the dialog")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("dotted extension accepted")
		}
	}()
	p.tester(FilePickerOptions{Filters: []mygo.FileFilter{{Name: "PDF", Extensions: []string{".pdf"}}}})
}

func TestFilePickerSnapshot(t *testing.T) {
	a, b := []string{"/Users/sakiko/scores/Imprisoned XII.pdf", "/Users/sakiko/scores/Ave Mujica.mid"}, []string{}
	snapshot(t, "file_picker", 460, 190, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(20).Children(func() {
			FilePicker(c, &a, FilePickerOptions{Multiple: true})
			FilePicker(c, &b, FilePickerOptions{Directory: true, Disabled: true})
		})
	})
}
