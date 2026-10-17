package input

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// FormState is a form's check and submit state; use it on the UI thread only.
type FormState struct {
	Fields     map[string]*FormFieldState
	Order      []string
	resets     []func()
	seq        uint64
	tried      bool
	submitting bool
	Awaiting   bool
	Requested  bool
	focus      string
	formError  string
	generation int
}

type FormFieldState struct {
	snapshot   string
	version    uint64
	touched    bool
	focused    bool
	err        string
	asyncFor   uint64
	asyncErr   string
	pending    bool
	owed       bool // blurred valid; the check starts once no form button is held
	hasAsync   bool
	check      func() string
	refresh    func()
	startAsync func()
}

// Version returns the version of name's value, which changes whenever the value does.
func (s *FormState) Version(name string) uint64 {
	if f := s.Fields[name]; f != nil {
		return f.version
	}
	return 0
}

// Resolve delivers name's async result for version, "" if valid; stale versions return false.
func (s *FormState) Resolve(name string, version uint64, err string) bool {
	f := s.Fields[name]
	if f == nil || f.version != version {
		slog.Debug("mujica: form async result dropped", "field", name, "version", version)
		return false
	}
	f.pending, f.asyncFor, f.asyncErr = false, version, err
	slog.Debug("mujica: form async result", "field", name, "valid", err == "")
	return true
}

// Submit asks the form to submit in its next frame, as the submit button does.
func (s *FormState) Submit() { s.Requested = true }

// Submitting reports whether a submit is in progress.
func (s *FormState) Submitting() bool { return s.submitting }

// EndSubmit finishes a submit; formError, if any, shows above the buttons and keeps all values.
func (s *FormState) EndSubmit(formError string) {
	if !s.submitting {
		panic("mujica: FormState.EndSubmit without a submit in progress")
	}
	s.submitting, s.formError = false, formError
	slog.Debug("mujica: form submit ended", "failed", formError != "")
}

// Error returns the error name shows now, "" when valid, untouched or pending.
func (s *FormState) Error(name string) string {
	f := s.Fields[name]
	if f == nil || !(f.touched || s.tried) {
		return ""
	}
	if f.err != "" {
		return f.err
	}
	if f.asyncFor == f.version {
		return f.asyncErr
	}
	return ""
}

// FormOptions configure a Form.
type FormOptions struct {
	// Horizontal puts labels in a column beside the controls.
	Horizontal bool
	// SubmitLabel and ResetLabel replace the built-in "Submit" and "Reset".
	SubmitLabel, ResetLabel string
}

// FormScope registers a Form's fields with FormItem.
type FormScope struct {
	st         *FormState
	horizontal bool
}

// FormParts are the parts of a Form.
type FormParts struct {
	Element          ui.Element
	submitted, reset bool
}

// Submitted reports a submit that passed every check; call EndSubmit when it finishes.
func (p FormParts) Submitted() bool { return p.submitted }

// Reset reports that the form restored its initial values this frame.
func (p FormParts) Reset() bool { return p.reset }

// FormItemOptions configure one field of a Form.
type FormItemOptions[T any] struct {
	Label       string
	Required    bool
	Description string
	// Initial is the value Reset restores.
	Initial T
	// Validate returns an error message for v, or ""; cross-field rules read other values here.
	Validate func(v T) string
	// Async starts a check of v off the UI thread, finished with FormState.Resolve(name, version, msg).
	Async func(v T, version uint64)
}

// Form lays out and checks its fields and shows submit and reset buttons.
func Form(c *ui.Context, st *FormState, opts FormOptions, fields func(f *FormScope)) FormParts {
	if st.Fields == nil {
		st.Fields = map[string]*FormFieldState{}
	}
	st.Order, st.resets = st.Order[:0], nil
	d := core.Density(c)
	k := core.Tokens(c)
	parts := FormParts{}
	col := ui.Column(c).Gap(d.Unit() * 5)
	parts.Element = col
	scope := &FormScope{st: st, horizontal: opts.Horizontal}
	held := false
	col.Children(func() {
		// A new key after Reset drops the controls' drafts and errors.
		ui.Box(c.Key(st.generation)).Children(func() {
			if opts.Horizontal {
				formFieldGrid(c, func() { fields(scope) })
			} else {
				ui.Column(c).Gap(d.Unit() * 4).Children(func() { fields(scope) })
			}
		})
		if st.formError != "" {
			ui.Text(c, st.formError).TextColor(k.Danger).Role(ui.RoleStatus)
		}
		ui.Row(c).Gap(d.Unit() * 3).Children(func() {
			submit, reset := opts.SubmitLabel, opts.ResetLabel
			if submit == "" {
				submit = core.Msg(c, "form.submit")
			}
			if reset == "" {
				reset = core.Msg(c, "form.reset")
			}
			sb := Button(c, submit, ButtonOptions{Loading: st.submitting || st.Awaiting})
			if sb.Clicked() {
				st.Requested = true
			}
			rb := Button(c, reset, ButtonOptions{Variant: Secondary, Disabled: st.submitting})
			if rb.Clicked() {
				formReset(st)
				parts.reset = true
			}
			// A "Checking…" line would move the buttons under the pointer before the release.
			held = sb.Pressed() || rb.Pressed()
		})
	})
	if st.Requested {
		st.Requested = false
		if !st.submitting && !st.Awaiting {
			formStartSubmit(st)
		}
	}
	if st.Awaiting {
		parts.submitted = formSettle(c, st)
	}
	if !held {
		formStartOwed(c, st)
	}
	return parts
}

// formStartOwed starts the async checks of fields blurred while a form button was held.
func formStartOwed(c *ui.Context, st *FormState) {
	for _, name := range st.Order {
		if f := st.Fields[name]; f != nil && f.owed {
			f.owed = false
			f.startAsync()
			c.Invalidate()
		}
	}
}

// formStartSubmit validates every field and starts the async checks still owed.
func formStartSubmit(st *FormState) {
	st.tried, st.Awaiting, st.formError = true, true, ""
	for _, name := range st.Order {
		f := st.Fields[name]
		f.err, f.owed = f.check(), false
		if f.err == "" {
			f.startAsync()
		}
	}
	slog.Debug("mujica: form submit requested", "fields", len(st.Order))
}

// formSettle finishes a waiting submit once no check is pending; it reports success.
func formSettle(c *ui.Context, st *FormState) bool {
	// Current values only: later edits and cross-field rules count.
	for _, name := range st.Order {
		f := st.Fields[name]
		f.refresh()
		f.err = f.check()
	}
	first, pending := "", false
	for _, name := range st.Order {
		if st.Error(name) != "" && first == "" {
			first = name
		}
		f := st.Fields[name]
		// A value edited since its check owes a new one.
		if f.hasAsync && f.err == "" && f.asyncFor != f.version {
			f.startAsync()
			pending = true
		}
	}
	switch {
	case first != "":
		st.Awaiting, st.focus = false, first
		c.Announce(core.Msg(c, "form.failed"))
		c.Invalidate()
		slog.Debug("mujica: form submit failed", "first", first)
	case pending:
	default:
		st.Awaiting, st.submitting = false, true
		slog.Debug("mujica: form submitted")
		return true
	}
	return false
}

// formReset restores the initial values and forgets every check.
func formReset(st *FormState) {
	for _, r := range st.resets {
		r()
	}
	st.Fields = map[string]*FormFieldState{}
	st.generation++
	st.tried, st.Awaiting, st.Requested, st.focus, st.formError = false, false, false, "", ""
	slog.Debug("mujica: form reset")
}

// FormItem registers field name around the focusable element control returns.
func FormItem[T any](c *ui.Context, f *FormScope, name string, value *T, opts FormItemOptions[T], control func() ui.Element) ui.Element {
	st := f.st
	if name == "" || opts.Label == "" {
		panic("mujica: FormItem needs a name and a Label")
	}
	if slices.Contains(st.Order, name) {
		panic("mujica: FormItem " + name + " registered twice")
	}
	st.Order = append(st.Order, name)
	fs := st.Fields[name]
	if fs == nil {
		st.seq++
		fs = &FormFieldState{snapshot: fmt.Sprintf("%#v", *value), version: st.seq}
		st.Fields[name] = fs
	}
	fs.hasAsync = opts.Async != nil
	fs.check = func() string {
		if opts.Required && formEmpty(*value) {
			return core.Msg(c, "form.required")
		}
		if opts.Validate != nil {
			return opts.Validate(*value)
		}
		return ""
	}
	fs.startAsync = func() {
		if opts.Async != nil && fs.asyncFor != fs.version && !fs.pending {
			fs.pending = true
			slog.Debug("mujica: form async check started", "field", name, "version", fs.version)
			opts.Async(*value, fs.version)
		}
	}
	// refresh gives an edited value a new version; it runs after the control took its input.
	fs.refresh = func() {
		snap := fmt.Sprintf("%#v", *value)
		if snap == fs.snapshot {
			return
		}
		st.seq++
		fs.snapshot, fs.version, fs.pending, fs.asyncErr = snap, st.seq, false, ""
		if fs.touched || st.tried {
			fs.err = fs.check()
		}
		c.Invalidate()
	}
	initial := opts.Initial
	st.resets = append(st.resets, func() { *value = initial })
	desc := opts.Description
	if fs.pending {
		desc = core.Msg(c, "form.checking")
	}
	var ctrl ui.Element
	field := formFieldBuild(c, name, opts.Label, FormFieldOptions{Required: opts.Required, Description: desc, Error: st.Error(name)}, f.horizontal, func() ui.Element {
		ctrl = control()
		return ctrl
	})
	fs.refresh()
	focused := ctrl.FocusWithin()
	if fs.focused && !focused {
		fs.touched = true
		fs.err = fs.check()
		fs.owed = fs.err == ""
		slog.Debug("mujica: form field blurred", "field", name, "valid", fs.err == "")
	}
	fs.focused = focused
	if st.focus == name {
		st.focus = ""
		ctrl.Focus().ScrollIntoView()
	}
	return field
}

// formEmpty reports whether a required value is missing: zero, or an empty string, slice or map.
func formEmpty(v any) bool {
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return r.Len() == 0
	}
	return r.IsZero()
}
