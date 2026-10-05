package mujica

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func storageSample() []StorageSegment {
	return []StorageSegment{{"Documents", 24 << 30}, {"Photos", 40 << 30}, {"Apps", 12 << 30}}
}

func storageTester(segs []StorageSegment, capacity int64) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		StorageUsage(c, segs, StorageUsageOptions{Capacity: capacity, Label: "Court disk"}).Width(420)
	}, 460, 300)
}

func TestStorageUsageLegend(t *testing.T) {
	tt := storageTester(storageSample(), 128<<30)
	for _, s := range []string{"76.0 GB of 128.0 GB used", "Photos", "40.0 GB", "31.2%", "Free", "52.0 GB", "40.6%"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
}

func TestStorageUsageKeyboardReadout(t *testing.T) {
	tt := storageTester(storageSample(), 128<<30)
	tt.Key(0, ui.KeyTab)
	if !tt.HasText("Documents: 24.0 GB (18.8%)") {
		t.Fatalf("focused readout missing: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Free: 52.0 GB (40.6%)") {
		t.Fatalf("readout after three Rights: %q", tt.Texts())
	}
	ann := tt.Announcements()
	if len(ann) == 0 || ann[len(ann)-1] != "Free: 52.0 GB (40.6%)" {
		t.Fatalf("announcements %q", ann)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Documents: 24.0 GB (18.8%)") {
		t.Fatal("Right did not wrap to the first segment")
	}
}

func TestStorageUsageInvalid(t *testing.T) {
	for _, f := range []func(){
		func() { storageTester(storageSample(), 0) },
		func() { storageTester(storageSample(), 10) },
		func() { storageTester([]StorageSegment{{"x", -1}}, 10) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("invalid storage did not panic")
				}
			}()
			f()
		}()
	}
}

func TestStorageUsageSnapshot(t *testing.T) {
	snapshot(t, "storage_usage", 460, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			StorageUsage(c, storageSample(), StorageUsageOptions{Capacity: 128 << 30, Label: "Court disk"})
		})
	})
}

func TestStorageUsageOverflow(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("overflowing storage was accepted")
		}
	}()
	storageTester([]StorageSegment{{"a", 1 << 62}, {"b", 1 << 62}}, 1<<62)
}

func TestStorageUsageLargeLegend(t *testing.T) {
	segments := make([]StorageSegment, 100000)
	for i := range segments {
		segments[i] = StorageSegment{Label: fmt.Sprintf("Category %d", i), Bytes: 1}
	}
	tt := storageTester(segments, int64(len(segments)))
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "Category ") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d legend labels", n)
	}
}
