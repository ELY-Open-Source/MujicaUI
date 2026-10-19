package messaging

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestEmojiPickerPickSearchToneRecent(t *testing.T) {
	var recent []string
	tone := 0
	var picked []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := EmojiPicker(c, EmojiPickerOptions{Recent: &recent, Tone: &tone}); r.Picked != "" {
			picked = append(picked, r.Picked)
		}
	}, 400, 500)
	if err := tt.Click("People"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Skin tone 3"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("thumbs up"); err != nil {
		t.Fatal(err)
	}
	if len(picked) != 1 || picked[0] != "👍\U0001F3FD" || tone != 3 {
		t.Fatalf("picked %q tone %d", picked, tone)
	}
	if err := tt.Click("Search emoji"); err != nil {
		t.Fatal(err)
	}
	tt.Type("wine")
	if err := tt.Click("wine glass"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(recent, []string{"🍷", "👍\U0001F3FD"}) {
		t.Fatalf("recent %q", recent)
	}
	if err := tt.Click("Search emoji"); err != nil {
		t.Fatal(err)
	}
	tt.Type("zzzz")
	if !tt.HasText("No emoji found") {
		t.Fatal("no empty search state")
	}
}

func TestEmojiPickerKeyboardAndDisabled(t *testing.T) {
	var got string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := EmojiPicker(c, EmojiPickerOptions{}); r.Picked != "" {
			got = r.Picked
		}
	}, 400, 500)
	if err := tt.Click("grinning face"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if got != "😄" {
		t.Fatalf("Right, Enter picked %q", got)
	}
	got = ""
	off := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := EmojiPicker(c, EmojiPickerOptions{Disabled: true}); r.Picked != "" {
			got = r.Picked
		}
	}, 400, 500)
	_ = off.Click("grinning face")
	if got != "" {
		t.Fatal("disabled picker picked")
	}
	bad := 9
	mediaPanics(t, "tone", func(c *ui.Context) { EmojiPicker(c, EmojiPickerOptions{Tone: &bad}) })
}

func TestEmojiPickerSnapshot(t *testing.T) {
	recent := []string{"🌹", "🍷"}
	tone := 0
	snapshot(t, "emoji_picker", 360, 400, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { EmojiPicker(c, EmojiPickerOptions{Recent: &recent, Tone: &tone}) })
	})
}

func TestEmojiPickerKeepsCategoryAfterFirstPick(t *testing.T) {
	var recent []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		EmojiPicker(c, EmojiPickerOptions{Recent: &recent})
	}, 400, 400)
	if err := tt.Click("People"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("thumbs up"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("waving hand"); !ok {
		t.Fatal("first pick changed the category")
	}
	if err := tt.Click("Recent"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("👍"); !ok {
		t.Fatal("recent pick missing")
	}
}

func TestEmojiPickerChineseSearch(t *testing.T) {
	picked := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		if r := EmojiPicker(c, EmojiPickerOptions{}); r.Picked != "" {
			picked = r.Picked
		}
	}, 400, 400)
	if err := tt.Click("搜索表情"); err != nil {
		t.Fatal(err)
	}
	tt.Type("小提琴")
	if err := tt.Click("小提琴"); err != nil {
		t.Fatal(err)
	}
	if picked != "🎻" {
		t.Fatal("localized search did not choose violin")
	}
}

func TestEmojiPickerInitialHistory(t *testing.T) {
	recent := []string{"🌹", "🍷"}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		EmojiPicker(c, EmojiPickerOptions{Recent: &recent})
	}, 400, 400)
	if got := mediaTextCount(tt, "🌹"); got != 3 {
		t.Fatalf("initial history rose occurrences=%d, want 3", got)
	}
}
