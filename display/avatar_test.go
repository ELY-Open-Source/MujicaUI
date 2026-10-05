package display

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestAvatarInitials(t *testing.T) {
	for name, want := range map[string]string{
		"Sakiko Togawa":        "ST",
		"uika":                 "U",
		"Anon Chihaya Mortis":  "AM",
		"  ":                   "",
		"张伟":                   "张",
		"mutsumi (wakaba) 若叶睦": "M若",
	} {
		if got := avatarInitials(name); got != want {
			t.Errorf("%q: %q, want %q", name, got, want)
		}
	}
}

func TestAvatarFallbacks(t *testing.T) {
	var size ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			size = Avatar(c, "Sakiko Togawa", AvatarOptions{Size: 48}).Bounds()
			Avatar(c, "", AvatarOptions{})
			Avatar(c, "Photo Person", AvatarOptions{Image: imageTestBitmap()})
		})
	}, 300, 100)
	tt.Frame()
	if size.W != 48 || size.H != 48 {
		t.Fatalf("size %v", size)
	}
	if !tt.HasText("ST") {
		t.Fatal("initials missing")
	}
	if _, ok := tt.Find("Unnamed"); !ok {
		t.Fatal("unnamed avatar not labeled")
	}
	if _, ok := tt.Find("Photo Person"); !ok {
		t.Fatal("image avatar not labeled")
	}
	if tt.HasText("PP") {
		t.Fatal("initials shown over image")
	}
}

func TestAvatarSnapshot(t *testing.T) {
	snapshot(t, "avatar", 320, 90, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(12).AlignItems(ui.Center).Children(func() {
			Avatar(c, "Sakiko Togawa", AvatarOptions{Size: 48})
			Avatar(c, "Uika", AvatarOptions{})
			Avatar(c, "", AvatarOptions{Size: 40})
			Avatar(c, "Photo", AvatarOptions{Image: imageTestBitmap(), Size: 40})
			Avatar(c, "Mutsumi Wakaba", AvatarOptions{Size: 24})
		})
	})
}
