package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func stepsSample() []Step {
	return []Step{
		{Title: "Account", Clickable: true},
		{Title: "Profile", Clickable: true},
		{Title: "Payment"},
		{Title: "Confirm"},
	}
}

func TestStepsClickable(t *testing.T) {
	cur, changes := 2, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Steps(c, &cur, stepsSample(), StepsOptions{}).Changed() {
			changes++
		}
	}, 700, 120)
	_ = tt.Click("Confirm")
	if cur != 2 || changes != 0 {
		t.Fatalf("non-clickable step chosen: %d", cur)
	}
	if err := tt.Click("Account"); err != nil {
		t.Fatal(err)
	}
	if cur != 0 || changes != 1 {
		t.Fatalf("click: cur %d changes %d", cur, changes)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if cur != 1 || changes != 2 {
		t.Fatalf("keyboard: cur %d changes %d", cur, changes)
	}
}

func TestStepsStatesNamed(t *testing.T) {
	cur := 1
	steps := stepsSample()
	steps[2].Error = true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Steps(c, &cur, steps, StepsOptions{Vertical: true})
	}, 400, 400)
	for _, name := range []string{
		"Step 1 of 4: Account, done",
		"Step 2 of 4: Profile, current",
		"Step 3 of 4: Payment, failed",
		"Step 4 of 4: Confirm, waiting",
	} {
		if _, ok := tt.Find(name); !ok {
			t.Errorf("missing %q in %q", name, tt.Texts())
		}
	}
	// All done: current past the last step.
	cur = 4
	tt.Frame()
	if _, ok := tt.Find("Step 4 of 4: Confirm, done"); !ok {
		t.Errorf("finished steps: %q", tt.Texts())
	}
}

func TestStepsInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		cur   int
		steps []Step
	}{
		"empty": {0, nil},
		"range": {5, stepsSample()},
		"title": {0, []Step{{}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Steps(c, &tc.cur, tc.steps, StepsOptions{})
			}, 300, 100, 1)
		})
	}
}

func TestStepsSnapshot(t *testing.T) {
	snapshot(t, "steps", 640, 380, func(c *ui.Context) {
		a, b := 2, 1
		steps := []Step{
			{Title: "Account", Description: "Name and seal"},
			{Title: "Profile", Description: "Titles and house"},
			{Title: "Payment", Description: "Treasury"},
			{Title: "Confirm"},
		}
		ui.Column(c).Padding(16).Gap(24).Children(func() {
			Steps(c, &a, steps, StepsOptions{})
			v := append([]Step(nil), steps...)
			v[2].Error = true
			v[2].Description = "The card was declined"
			Steps(c, &b, v, StepsOptions{Vertical: true})
		})
	})
}
