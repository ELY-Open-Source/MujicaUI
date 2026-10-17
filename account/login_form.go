package account

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

// LoginFormState holds the form; only OnSubmit sees the password.
type LoginFormState struct {
	Account    string
	Remember   bool
	password   string
	tried      bool
	submitting bool
	err        string
}

// EndSubmit finishes a sign-in, shows err and wipes the password.
func (s *LoginFormState) EndSubmit(err string) {
	if !s.submitting {
		panic("mujica: LoginFormState.EndSubmit without a submit in progress")
	}
	s.submitting, s.err, s.password = false, err, ""
	slog.Debug("mujica: login_form submit ended", "failed", err != "")
}

// LoginFormOptions configure a LoginForm. OnSubmit is required.
type LoginFormOptions struct {
	// OnSubmit alone receives the credentials; then call EndSubmit.
	OnSubmit func(account, password string, remember bool)
	// AccountLabel replaces the built-in "Email or username".
	AccountLabel string
	Disabled     bool
}

// LoginForm collects credentials; it never checks or stores them.
func LoginForm(c *ui.Context, s *LoginFormState, opts LoginFormOptions) ui.Element {
	if opts.OnSubmit == nil {
		panic("mujica: LoginForm needs OnSubmit")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	accountLabel := opts.AccountLabel
	if accountLabel == "" {
		accountLabel = core.Msg(c, "login_form.account")
	}
	e := ui.Column(c).Gap(d.Unit() * 4).Width(340).Disabled(opts.Disabled)
	disabled := e.IsDisabled()
	missing := func(v string) string {
		if s.tried && strings.TrimSpace(v) == "" {
			return core.Msg(c, "form.required")
		}
		return ""
	}
	submit := false
	var acc, pw input.InputResult
	e.Children(func() {
		input.FormField(c, accountLabel, input.FormFieldOptions{Required: true}, func() ui.Element {
			acc = input.TextInput(c, &s.Account, input.TextInputOptions{Disabled: disabled || s.submitting, Error: missing(s.Account)})
			return acc.Input
		})
		input.FormField(c, core.Msg(c, "login_form.password"), input.FormFieldOptions{Required: true}, func() ui.Element {
			pw = input.PasswordInput(c, &s.password, input.PasswordInputOptions{Disabled: disabled || s.submitting, Error: missing(s.password)})
			return pw.Input
		})
		if acc.Changed() || pw.Changed() {
			s.err = ""
		}
		input.Checkbox(c, &s.Remember, core.Msg(c, "login_form.remember"), input.CheckboxOptions{Disabled: disabled || s.submitting})
		if s.err != "" {
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).TextColor(k.Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16).Shrink(0)
				ui.Text(c, s.err).Grow(1)
			})
		}
		b := input.Button(c, core.Msg(c, "login_form.submit"), input.ButtonOptions{Loading: s.submitting, Disabled: disabled})
		b.FillWidth()
		submit = b.Clicked() || acc.Submitted() || pw.Submitted()
	})
	if submit && !disabled && !s.submitting {
		s.tried = true
		switch {
		case strings.TrimSpace(s.Account) == "":
			acc.Input.Focus()
			return e
		case s.password == "":
			pw.Input.Focus()
			return e
		}
		s.submitting, s.err = true, ""
		slog.Debug("mujica: login_form submitted")
		opts.OnSubmit(strings.TrimSpace(s.Account), s.password, s.Remember)
	}
	return e
}
