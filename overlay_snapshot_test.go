package mujica

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// overlaySnapshot is snapshot, with drive opening the real component through input before the capture.
func overlaySnapshot(t *testing.T, name string, w, h int, view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	for _, m := range []Mode{Light, Dark} {
		ready := false
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Mode: m})
			if ready {
				view(c)
			}
		}, w, h)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		ready = true
		tt.Frame()
		drive(tt)
		got := encode(t, tt.Image())
		path := filepath.Join("testdata", "snapshots", runtime.GOOS, name+"-"+strings.ToLower(m.String())+".png")
		if *update {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, got, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}
		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%v (run go test -update)", err)
		}
		if !bytes.Equal(got, want) {
			if err := os.WriteFile(path+".got.png", got, 0o644); err != nil {
				t.Fatal(err)
			}
			t.Errorf("%s differs; see %s.got.png", path, path)
		}
	}
}

// overlayUntil draws frames until cond holds, for at most 5 s of wall time.
func overlayUntil(tt *ui.Tester, cond func() bool) bool {
	end := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(end) {
			return false
		}
		time.Sleep(2 * time.Millisecond)
		tt.Frame()
	}
	return true
}
