package input

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func tagsInputTester(tags *[]string, opts TagsInputOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if TagsInput(c, tags, opts).Changed() {
				*changes++
			}
		})
	}, 520, 200)
}

func TestTagsInputEnterSeparatorsAndBackspace(t *testing.T) {
	tags, n := []string{"rose"}, 0
	tt := tagsInputTester(&tags, TagsInputOptions{Label: "Tags"}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Tags") {
		t.Fatal("input not after the chip's remove button")
	}
	tt.Type("thorn")
	tt.Key(0, ui.KeyEnter)
	tt.Type("lily, ivy；moss")
	if !slices.Equal(tags, []string{"rose", "thorn", "lily", "ivy"}) {
		t.Fatalf("tags %q", tags)
	}
	tt.Key(0, ui.KeyEnter)
	if tags[4] != "moss" || n != 3 {
		t.Fatalf("tags %q changes %d", tags, n)
	}
	tt.Key(0, ui.KeyBackspace)
	if !slices.Equal(tags, []string{"rose", "thorn", "lily", "ivy"}) {
		t.Fatalf("backspace: %q", tags)
	}
}

func TestTagsInputPasteDedupeAndMax(t *testing.T) {
	tags, n := []string{"a"}, 0
	tt := tagsInputTester(&tags, TagsInputOptions{Label: "Tags", Max: 3}, &n)
	_ = tt.Click("Tags")
	tt.SetClipboard("b,a,c,d,")
	tt.Command("paste")
	if !slices.Equal(tags, []string{"a", "b", "c"}) {
		t.Fatalf("tags %q", tags)
	}
	if !tt.HasText("“a” is already added") && !tt.HasText("At most 3 tags") {
		t.Fatalf("no error shown: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if len(tags) != 3 {
		t.Fatalf("max exceeded: %q", tags)
	}
}

func TestTagsInputValidateKeepsText(t *testing.T) {
	tags, n, seen := []string{}, 0, ""
	opts := TagsInputOptions{Label: "Tags", Validate: func(s string) string {
		seen = s
		if strings.ContainsAny(s, " ") {
			return "No spaces"
		}
		return ""
	}}
	tt := tagsInputTester(&tags, opts, &n)
	_ = tt.Click("Tags")
	tt.Type("two words")
	tt.Key(0, ui.KeyEnter)
	if len(tags) != 0 || !tt.HasText("No spaces") {
		t.Fatalf("tags %q texts %q", tags, tt.Texts())
	}
	tt.Type("x")
	if tt.HasText("No spaces") {
		t.Fatal("error stays after editing")
	}
	tt.Key(0, ui.KeyEnter)
	if seen != "two wordsx" {
		t.Fatalf("typed text lost: %q", seen)
	}
}

func TestTagsInputRemoveButtonAndDisabled(t *testing.T) {
	tags, n := []string{"rose", "ivy"}, 0
	tt := tagsInputTester(&tags, TagsInputOptions{}, &n)
	if err := tt.Click("Remove rose"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tags, []string{"ivy"}) || n != 1 {
		t.Fatalf("tags %q", tags)
	}
	tags = []string{"rose"}
	tt = tagsInputTester(&tags, TagsInputOptions{Disabled: true}, &n)
	_ = tt.Click("Remove rose")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyBackspace)
	if len(tags) != 1 {
		t.Fatal("disabled input removed a tag")
	}
}

func TestTagsInputSnapshot(t *testing.T) {
	a, b := []string{"Ave Mujica", "Oblivionis", "Doloris"}, []string{"rose"}
	snapshot(t, "tags_input", 520, 190, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			TagsInput(c, &a, TagsInputOptions{Placeholder: "Add a member"})
			TagsInput(c, &b, TagsInputOptions{Invalid: true, Placeholder: "Add a flower"})
		})
	})
}

func TestTagsInputRemoveHitArea(t *testing.T) {
	tags, n := []string{"rose"}, 0
	tt := tagsInputTester(&tags, TagsInputOptions{}, &n)
	r, ok := tt.Find("Remove rose")
	if !ok || r.W < 28 || r.H < 28 {
		t.Fatalf("remove button %v", r)
	}
}

func TestTagsInputDisabledAfterFocus(t *testing.T) {
	tags, n := []string{"rose"}, 0
	opts := TagsInputOptions{Label: "Tags"}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TagsInput(c, &tags, opts).Changed() {
			n++
		}
	}, 520, 200)
	_ = tt.Click("Tags")
	opts.Disabled = true
	tt.Key(0, ui.KeyBackspace)
	tt.Type("ivy,")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !slices.Equal(tags, []string{"rose"}) || n != 0 {
		t.Fatalf("tags %q changes %d", tags, n)
	}
}

// Cut and paste after disabling with a selected draft change neither draft nor clipboard.
func TestTagsInputDisabledEditCommands(t *testing.T) {
	tags, n, seen := []string{}, 0, ""
	opts := TagsInputOptions{Label: "Tags", Validate: func(s string) string { seen = s; return "" }}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TagsInput(c, &tags, opts).Changed() {
			n++
		}
	}, 520, 200)
	_ = tt.Click("Tags")
	tt.Type("ivy")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.SetClipboard("clip")
	opts.Disabled = true
	tt.Frame()
	tt.Key(ui.Cmd, ui.KeyX)
	tt.Command("cut")
	tt.Command("paste")
	tt.Frame()
	if tt.Clipboard() != "clip" || n != 0 {
		t.Fatalf("clipboard %q changes %d", tt.Clipboard(), n)
	}
	opts.Disabled = false
	tt.Frame()
	_ = tt.Click("Tags")
	tt.Key(0, ui.KeyEnter)
	if seen != "ivy" || tt.Clipboard() != "clip" || len(tags) != 1 {
		t.Fatalf("draft %q clipboard %q tags %q", seen, tt.Clipboard(), tags)
	}
}
