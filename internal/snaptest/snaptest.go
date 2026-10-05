// Package snaptest compares rendered views with the PNG goldens in testdata/snapshots.
package snaptest

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

// Modes names the two themes a snapshot covers, in capture order.
var Modes = []string{"light", "dark"}

// child marks a process that renders snapshots itself.
const child = "MUJICA_SNAPSHOT_CHILD"

var ran sync.Map

// dir is the module's testdata/snapshots/<os>, whatever package runs the test.
var dir = func() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("snaptest: no caller")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "snapshots", runtime.GOOS)
}()

// Run compares view, themed by theme for each of Modes, with testdata/snapshots/<os>/<name>-<mode>.png.
// drive, when non-nil, sends input before the capture.
func Run(t *testing.T, name string, w, h int, theme func(c *ui.Context, mode string), view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	if os.Getenv(child) == "" {
		inChild(t)
		return
	}
	for _, m := range Modes {
		ready := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme(c, m)
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
		got := Encode(t, tt.Image())
		path := filepath.Join(dir, name+"-"+m+".png")
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

// inChild reruns the test alone in a fresh process; MyGo's path cache is process-wide.
func inChild(t *testing.T) {
	t.Helper()
	if _, done := ran.LoadOrStore(t.Name(), true); done {
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
	cmd.Env = append(os.Environ(), child+"=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("snapshot child failed: %v\n%s", err, out)
	}
}

// Encode returns img as PNG bytes.
func Encode(t *testing.T, img image.Image) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
