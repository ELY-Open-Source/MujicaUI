package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
	k := core.Tokens(c)
	d := core.Density(c)
	var r OnboardingWizardResult
	from := *step
	r.Element = ui.Column(c).Gap(d.Unit()*5).Width(520).Padding(d.Unit()*6).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border)
	r.Element.Children(func() {
		head := make([]navigation.Step, len(steps))
		for i, s := range steps {
			if s.Title == "" {
				panic("mujica: WizardStep needs a Title")
			}
			head[i] = navigation.Step{Title: s.Title}
		}
		at := *step
		navigation.Steps(c, &at, head, navigation.StepsOptions{Label: core.Msg(c, "onboarding_wizard.progress")})
		cur := steps[*step]
		ui.Column(c).Key(*step).Gap(d.Unit() * 3).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "onboarding_wizard.count"), *step+1, len(steps))).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
			ui.Text(c, cur.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).FontWeight(600)
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
			if !opts.NoSkip && !last && input.Button(c, core.Msg(c, "onboarding_wizard.skip"), input.ButtonOptions{Variant: input.Ghost}).Clicked() {
				r.skipped = true
			}
			ui.Box(c).Grow(1)
			if *step > 0 && input.Button(c, core.Msg(c, "onboarding_wizard.back"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				*step--
			}
			next := core.Msg(c, "onboarding_wizard.next")
			if last {
				next = core.Msg(c, "onboarding_wizard.finish")
			}
			if input.Button(c, next, input.ButtonOptions{Disabled: !ready}).AutoFocus().Clicked() && ready {
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
