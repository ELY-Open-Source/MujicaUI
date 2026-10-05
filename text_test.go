package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestTextSelectAndCopy(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Text(c, "ID-4821", TextOptions{Selectable: true})
	}, 300, 100)
	if err := tt.Click("ID-4821"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Command("copy")
	if got := tt.Clipboard(); got != "ID-4821" {
		t.Fatalf("clipboard %q", got)
	}
}

func TestTextSelectableTakesTab(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Text(c, "ID-4821", TextOptions{Selectable: true})
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("ID-4821") {
		t.Fatal("selectable text not in Tab order")
	}
	tt.Command("selectAll")
	tt.Command("copy")
	if got := tt.Clipboard(); got != "ID-4821" {
		t.Fatalf("clipboard %q", got)
	}
}

func TestEditableTextBlurConfirms(t *testing.T) {
	name := "a"
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if EditableText(c, &name, EditableTextOptions{}).Changed() {
			changed++
		}
		Button(c, "Other", ButtonOptions{})
	}, 300, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Type("b")
	if err := tt.Click("Other"); err != nil {
		t.Fatal(err)
	}
	if name != "ab" || changed != 1 {
		t.Fatalf("name %q, changed %d", name, changed)
	}
}

func TestTextNotSelectableByDefault(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Text(c, "plain", TextOptions{})
	}, 300, 100)
	_ = tt.Click("plain")
	tt.Command("selectAll")
	tt.Command("copy")
	if got := tt.Clipboard(); got != "" {
		t.Fatalf("clipboard %q, want empty", got)
	}
}

func TestTextMaxLinesCuts(t *testing.T) {
	long := "The court gathers in the hall of mirrors while candles burn low and the music fades away."
	var free, cut float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(120).Children(func() {
			free = Text(c, long, TextOptions{}).Bounds().H
			cut = Text(c, long, TextOptions{MaxLines: 1}).Bounds().H
		})
	}, 300, 400)
	tt.Frame()
	if cut <= 0 || free < cut*3 {
		t.Fatalf("free %v, one line %v", free, cut)
	}
}

func TestTextNegativeMaxLinesPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Text(c, "x", TextOptions{MaxLines: -1})
	}, 100, 100)
	tt.Frame()
}

func TestEditableTextConfirm(t *testing.T) {
	name := "report.pdf"
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if EditableText(c, &name, EditableTextOptions{Label: "File name"}).Changed() {
			changed++
		}
	}, 300, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.Focused("File name") {
		t.Fatal("editing input not named File name")
	}
	tt.Command("selectAll")
	tt.Type("summary.pdf")
	tt.Key(0, ui.KeyEnter)
	if name != "summary.pdf" || changed != 1 {
		t.Fatalf("name %q, changed %d", name, changed)
	}
	if !tt.Focused("File name") {
		t.Fatal("focus not back on the text after confirm")
	}
}

func TestEditableTextCancel(t *testing.T) {
	name := "report.pdf"
	changed, canceled := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		v := EditableText(c, &name, EditableTextOptions{})
		if v.Changed() {
			changed++
		}
		if v.Canceled {
			canceled++
		}
	}, 300, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Type("draft")
	tt.Key(0, ui.KeyEscape)
	if name != "report.pdf" || changed != 0 || canceled != 1 {
		t.Fatalf("name %q, changed %d, canceled %d", name, changed, canceled)
	}
	if !tt.HasText("report.pdf") {
		t.Fatal("text not shown after cancel")
	}
}

func TestEditableTextSavingBlocksAndErrorShows(t *testing.T) {
	name := "notes"
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		EditableText(c, &name, EditableTextOptions{Saving: true, Error: "disk full"})
	}, 300, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if name != "notes" {
		t.Fatalf("edited while saving: %q", name)
	}
	if !tt.HasText("disk full") {
		t.Fatal("save error not shown")
	}
	if _, ok := tt.Find("Saving"); !ok {
		t.Fatal("saving spinner not labeled")
	}
}

func TestTextSnapshot(t *testing.T) {
	name := "Coronation.md"
	snapshot(t, "text", 420, 280, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Width(380).Children(func() {
			Text(c, "Velvet curtains fall over the long gallery, and every candle leans toward the open door.", TextOptions{})
			Text(c, "A single line that is far too long to fit and so ends in an ellipsis.", TextOptions{MaxLines: 1})
			Text(c, "Selectable: ID-4821", TextOptions{Selectable: true})
			EditableText(c, &name, EditableTextOptions{Error: "The name is already taken."})
			EditableText(c, &name, EditableTextOptions{Saving: true})
		})
	})
}

func TestEditableTextInputFocusRing(t *testing.T) {
	name := "report.pdf"
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Column(c).Padding(20).Gap(20).AlignItems(ui.Start).Children(func() {
			EditableText(c, &name, EditableTextOptions{Label: "File name"})
			Button(c, "Other", ButtonOptions{})
		})
	}, 300, 160)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	r, ok := tt.Find("File name")
	if !ok || !tt.Focused("File name") {
		t.Fatal("editing input not found or not focused")
	}
	// The ring's middle: focus gap plus half its width outside the box.
	x, y := int(r.X-theme.FocusGap-theme.FocusWidth/2), int(r.Y+r.H/2)
	focused := tt.Image().RGBAAt(x, y)
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyTab)
	if tt.Focused("File name") {
		t.Fatal("focus did not leave")
	}
	plain := tt.Image().RGBAAt(x, y)
	want := theme.Light().Focus
	if d := int(focused.R) - int(want.R); d < -24 || d > 24 || focused == plain {
		t.Fatalf("ring pixel %v focused, %v unfocused, want near %v", focused, plain, want)
	}
}
