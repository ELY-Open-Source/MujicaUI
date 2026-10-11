package account

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TwoFactorInputOptions configure a TwoFactorInput.
type TwoFactorInputOptions struct {
	// ResendAt sets when resending becomes available.
	ResendAt time.Time
	// Error displays the caller-provided verification failure.
	Error    string
	Disabled bool
}

// TwoFactorResult reports completion and resend events.
type TwoFactorResult struct {
	Element           *ui.Element
	completed, resend bool
}

// Completed reports a transition to six filled digits.
func (r TwoFactorResult) Completed() bool { return r.completed }

// Resend reports a request for a new code.
func (r TwoFactorResult) Resend() bool { return r.resend }

// TwoFactorInput collects six digits and displays a resend countdown.
func TwoFactorInput(c *ui.Context, code *string, opts TwoFactorInputOptions) TwoFactorResult {
	if code == nil || len(*code) > 6 || strings.Trim(*code, "0123456789") != "" {
		panic("mujica: TwoFactorInput needs up to six digits")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	var r TwoFactorResult
	r.Element = ui.Column(c).Gap(d.Unit() * 3).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		pin := input.PinInput(c, code, input.PinInputOptions{Length: 6, Label: core.Msg(c, "two_factor_input.label"), Disabled: disabled, Invalid: opts.Error != ""})
		r.completed = pin.Completed() && !disabled
		if opts.Error != "" {
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(k.Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(14, 14)
				ui.Text(c, opts.Error).FontSize(core.FontSize(c, theme.CaptionSize))
			})
		}
		wait := opts.ResendAt.Sub(c.Now())
		label := core.Msg(c, "two_factor_input.resend")
		if wait > 0 {
			secs := int((wait + time.Second - 1) / time.Second)
			label = fmt.Sprintf(core.Msg(c, "two_factor_input.resend_in"), secs/60, secs%60)
			// Wake when the shown second changes.
			c.After(wait - time.Duration(secs-1)*time.Second)
		}
		ui.Row(c).Children(func() {
			b := input.Button(c, label, input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("refresh-cw"), Disabled: disabled || wait > 0})
			if b.Clicked() && !disabled && wait <= 0 {
				r.resend = true
				slog.Debug("mujica: two_factor_input resend")
			}
		})
	})
	if r.completed {
		slog.Debug("mujica: two_factor_input completed")
	}
	return r
}
