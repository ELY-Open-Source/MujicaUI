package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/input"
)

func TestOnboardingWizardStepsSkipFinish(t *testing.T) {
	step, name := 0, ""
	var finished, skipped int
	steps := []WizardStep{
		{Title: "Welcome", Description: "A short tour."},
		{Title: "Your name", Ready: func() bool { return name != "" }, Content: func() {}},
		{Title: "Done"},
	}
	tt := accountTester(600, 500, func(c *ui.Context) {
		steps[1].Content = func() { input.TextInput(c, &name, input.TextInputOptions{Label: "Name"}) }
		r := OnboardingWizard(c, &step, steps, OnboardingWizardOptions{})
		if r.Finished() {
			finished++
		}
		if r.Skipped() {
			skipped++
		}
	})
	if !tt.HasText("Step 1 of 3") || !tt.Focused("Next") {
		t.Fatalf("first step: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Next")
	if step != 1 {
		t.Fatalf("not ready but moved: %d", step)
	}
	_ = tt.Click("Name")
	tt.Type("Sakiko")
	_ = tt.Click("Next")
	_ = tt.Click("Back")
	_ = tt.Click("Next")
	if step != 2 || tt.HasText("Skip") {
		t.Fatalf("step %d", step)
	}
	_ = tt.Click("Finish")
	if finished != 1 || skipped != 0 {
		t.Fatalf("finished %d", finished)
	}
	step = 0
	tt.Frame()
	_ = tt.Click("Skip")
	if skipped != 1 {
		t.Fatal("skip not reported")
	}
}

func TestOnboardingWizardSnapshot(t *testing.T) {
	snapshot(t, "onboarding_wizard", 580, 330, func(c *ui.Context) {
		step := ui.Local(c.Root(), "step", func() int { return 1 })
		ui.Column(c).Padding(16).Children(func() {
			OnboardingWizard(c, step, []WizardStep{{Title: "Welcome"}, {Title: "Pick a theme", Description: "You can change it later in Settings."}, {Title: "Done"}}, OnboardingWizardOptions{})
		})
	})
}

func TestOnboardingWizardDisabledAndSingleStep(t *testing.T) {
	step, done := 0, 0
	disabled := true
	tt := accountTester(600, 500, func(c *ui.Context) {
		ui.Column(c).Disabled(disabled).Children(func() {
			if OnboardingWizard(c, &step, []WizardStep{{Title: "Only step"}}, OnboardingWizardOptions{NoSkip: true}).Finished() {
				done++
			}
		})
	})
	if err := tt.Click("Finish"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if done != 0 || step != 0 || tt.HasText("Skip") || tt.HasText("Back") {
		t.Fatal("disabled single step acts")
	}
	disabled = false
	tt.Frame()
	if err := tt.Click("Finish"); err != nil {
		t.Fatal(err)
	}
	if done != 1 {
		t.Fatal("single step cannot finish")
	}
	step = 1
	accountPanics(t, func(c *ui.Context) {
		OnboardingWizard(c, &step, []WizardStep{{Title: "Only step"}}, OnboardingWizardOptions{})
	})
}
