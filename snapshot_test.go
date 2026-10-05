package mujica

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/egoist/mygo/ui"
)

var update = flag.Bool("update", false, "rewrite snapshots")

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snapshotAfter(t, name, w, h, view, nil)
}

// snapshotChild marks a process that renders snapshots itself.
const snapshotChild = "MUJICA_SNAPSHOT_CHILD"

var snapshotRan sync.Map

// snapshotAfter is snapshot with drive sending input before the capture.
func snapshotAfter(t *testing.T, name string, w, h int, view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	if os.Getenv(snapshotChild) == "" {
		snapshotInChild(t)
		return
	}
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
		if drive != nil {
			drive(tt)
		}
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

// snapshotInChild reruns the test alone in a fresh process; MyGo's path cache is process-wide.
func snapshotInChild(t *testing.T) {
	t.Helper()
	if _, done := snapshotRan.LoadOrStore(t.Name(), true); done {
		return
	}
	parts := strings.Split(t.Name(), "/")
	for i, p := range parts {
		parts[i] = "^" + regexp.QuoteMeta(p) + "$"
	}
	args := []string{"-test.run=" + strings.Join(parts, "/"), "-test.count=1"}
	if *update {
		args = append(args, "-update")
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), snapshotChild+"=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("snapshot child failed: %v\n%s", err, out)
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
