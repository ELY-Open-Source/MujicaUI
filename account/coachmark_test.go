package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/input"
)

func coachmarkView(c *ui.Context, open *bool, step *int, r *CoachmarkResult) {
	var a, b ui.Element
	ui.Row(c).Gap(40).Padding(20).Children(func() {
		a = input.Button(c, "Compose", input.ButtonOptions{})
		b = input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	})
	*r = Coachmark(c, open, step, []CoachmarkStep{
		{Anchor: a, Title: "Write here", Body: "Start a new letter."},
		{Anchor: b, Title: "Keep it tidy", Body: "Archive finished threads."},
	})
}

func TestCoachmarkStepsAndEscape(t *testing.T) {
	open, step := true, 0
	var r CoachmarkResult
	done, dismissed := 0, 0
	tt := accountTester(600, 400, func(c *ui.Context) {
		coachmarkView(c, &open, &step, &r)
		if r.Finished() {
			done++
		}
		if r.Dismissed() {
			dismissed++
		}
	})
	tt.Frame()
	if !tt.HasText("Write here") || !tt.HasText("1 of 2") {
		t.Fatalf("first: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if step != 1 || !tt.HasText("Keep it tidy") {
		t.Fatalf("next: %d", step)
	}
	_ = tt.Click("Done")
	if done != 1 || open {
		t.Fatal("done not reported")
	}
	open, step = true, 0
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if open || dismissed != 1 || done != 1 {
		t.Fatalf("escape: open %v dismissed %d", open, dismissed)
	}
}

func TestCoachmarkSnapshot(t *testing.T) {
	snapshot(t, "coachmark", 420, 230, func(c *ui.Context) {
		open := ui.Local(c.Root(), "open", func() bool { return true })
		step := ui.Local(c.Root(), "step", func() int { return 0 })
		var r CoachmarkResult
		coachmarkView(c, open, step, &r)
	})
}

func TestCoachmarkDisabledAnchorCloses(t *testing.T) {
	open, step, disabled := true, 0, false
	tt := accountTester(600, 400, func(c *ui.Context) {
		a := input.Button(c, "Compose", input.ButtonOptions{Disabled: disabled})
		Coachmark(c, &open, &step, []CoachmarkStep{{Anchor: a, Title: "Write here"}})
	})
	disabled = true
	tt.Frame()
	if open || tt.HasText("Write here") {
		t.Fatal("disabled anchor retains coachmark")
	}
}

func TestCoachmarkInvalidStep(t *testing.T) {
	open, step := true, 1
	accountPanics(t, func(c *ui.Context) { Coachmark(c, &open, &step, []CoachmarkStep{{Title: "Only step"}}) })
}
