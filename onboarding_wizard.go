package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// WizardStep is one page of an OnboardingWizard.
type WizardStep struct {
	Title, Description string
	Content            func()
	// Ready, when set, must return true before Next is allowed.
	Ready func() bool
}

// OnboardingWizardOptions configure an OnboardingWizard.
type OnboardingWizardOptions struct {
	// NoSkip hides the Skip action.
	NoSkip bool
}

// OnboardingWizardResult reports completion or skipping.
type OnboardingWizardResult struct {
	Element           *ui.Element
	finished, skipped bool
}

// Finished reports Finish on the last step.
func (r OnboardingWizardResult) Finished() bool { return r.finished }

// Skipped reports that the user skipped the rest.
func (r OnboardingWizardResult) Skipped() bool { return r.skipped }

// OnboardingWizard provides progress and step navigation.
func OnboardingWizard(c *ui.Context, step *int, steps []WizardStep, opts OnboardingWizardOptions) OnboardingWizardResult {
	if len(steps) == 0 || *step < 0 || *step >= len(steps) {
		panic(fmt.Sprintf("mujica: OnboardingWizard step %d outside %d steps", *step, len(steps)))
	}
	k := Tokens(c)
	d := Density(c)
	var r OnboardingWizardResult
	from := *step
	r.Element = ui.Column(c).Gap(d.Unit()*5).Width(520).Padding(d.Unit()*6).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border)
	r.Element.Children(func() {
		head := make([]Step, len(steps))
		for i, s := range steps {
			if s.Title == "" {
				panic("mujica: WizardStep needs a Title")
			}
			head[i] = Step{Title: s.Title}
		}
		at := *step
		Steps(c, &at, head, StepsOptions{Label: msg(c, "onboarding_wizard.progress")})
		cur := steps[*step]
		ui.Column(c).Key(*step).Gap(d.Unit() * 3).Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, "onboarding_wizard.count"), *step+1, len(steps))).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize))
			ui.Text(c, cur.Title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H2Size)).FontWeight(600)
			if cur.Description != "" {
				ui.Text(c, cur.Description).TextColor(k.TextMuted)
			}
			if cur.Content != nil {
				cur.Content()
			}
		})
		last := *step == len(steps)-1
		ready := cur.Ready == nil || cur.Ready()
		ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
			if !opts.NoSkip && !last && Button(c, msg(c, "onboarding_wizard.skip"), ButtonOptions{Variant: Ghost}).Clicked() {
				r.skipped = true
			}
			ui.Box(c).Grow(1)
			if *step > 0 && Button(c, msg(c, "onboarding_wizard.back"), ButtonOptions{Variant: Secondary}).Clicked() {
				*step--
			}
			next := msg(c, "onboarding_wizard.next")
			if last {
				next = msg(c, "onboarding_wizard.finish")
			}
			if Button(c, next, ButtonOptions{Disabled: !ready}).AutoFocus().Clicked() && ready {
				if last {
					r.finished = true
				} else {
					*step++
				}
			}
		})
	})
	if *step != from || r.finished || r.skipped {
		slog.Debug("mujica: onboarding_wizard step", "step", *step, "finished", r.finished, "skipped", r.skipped)
	}
	return r
}
