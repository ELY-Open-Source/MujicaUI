package mujica

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

var update = flag.Bool("update", false, "rewrite snapshots")

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	for _, m := range []Mode{Light, Dark} {
		ready := false
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Mode: m})
			if ready {
				view(c)
			}
		}, w, h)
		// Animations must not paint before ReduceMotion applies.
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		ready = true
		tt.Frame()
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

func encode(t *testing.T, img image.Image) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
