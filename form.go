package mujica

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	"github.com/egoist/mygo/ui"
)

// FormState is a form's validation and submit state; keep it beside the form's values.
// Its methods run on the UI thread: deliver async results through Window.Update.
type FormState struct {
	fields     map[string]*formFieldState
	order      []string
	resets     []func()
	seq        uint64
	tried      bool
	submitting bool
	awaiting   bool
	requested  bool
	focus      string
	formError  string
}

type formFieldState struct {
	snapshot   string
	version    uint64
	touched    bool
	focused    bool
	err        string
	asyncFor   uint64
	asyncErr   string
	pending    bool
	check      func() string
	startAsync func()
}

// Version returns the version of name's value, which changes whenever the value does.
func (s *FormState) Version(name string) uint64 {
	if f := s.fields[name]; f != nil {
		return f.version
	}
	return 0
}

// Resolve delivers the async check of name's value at version; "" means valid.
// It reports false and drops the result when the value has changed since.
func (s *FormState) Resolve(name string, version uint64, err string) bool {
	f := s.fields[name]
	if f == nil || f.version != version {
		slog.Debug("mujica: form async result dropped", "field", name, "version", version)
		return false
	}
	f.pending, f.asyncFor, f.asyncErr = false, version, err
	slog.Debug("mujica: form async result", "field", name, "valid", err == "")
	return true
}

// Submit asks the form to submit in its next frame, as the submit button does.
func (s *FormState) Submit() { s.requested = true }

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
	f := s.fields[name]
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
	c  *ui.Context
	st *FormState
}

// FormParts are the parts of a Form.
type FormParts struct {
	Element          *ui.Element
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

// Form lays out the fields that fields registers, validates them on blur and submit,
// and shows submit and reset buttons. It never stores or sends the values itself.
func Form(c *ui.Context, st *FormState, opts FormOptions, fields func(f *FormScope)) FormParts {
	if st.fields == nil {
		st.fields = map[string]*formFieldState{}
	}
	st.order, st.resets = st.order[:0], nil
	d := Density(c)
	k := Tokens(c)
	parts := FormParts{}
	col := ui.Column(c).Gap(d.Unit() * 5)
	parts.Element = col
	scope := &FormScope{c: c, st: st}
	col.Children(func() {
		if opts.Horizontal {
			ui.Form(c, func() { fields(scope) })
		} else {
			ui.Column(c).Gap(d.Unit() * 4).Children(func() { fields(scope) })
		}
		if st.formError != "" {
			ui.Text(c, st.formError).TextColor(k.Danger).Role(ui.RoleStatus)
		}
		ui.Row(c).Gap(d.Unit() * 3).Children(func() {
			submit, reset := opts.SubmitLabel, opts.ResetLabel
			if submit == "" {
				submit = msg(c, "form.submit")
			}
			if reset == "" {
				reset = msg(c, "form.reset")
			}
			if Button(c, submit, ButtonOptions{Loading: st.submitting || st.awaiting}).Clicked() {
				st.requested = true
			}
			if Button(c, reset, ButtonOptions{Variant: Secondary, Disabled: st.submitting}).Clicked() {
				formReset(st)
				parts.reset = true
			}
		})
	})
	if st.requested {
		st.requested = false
		if !st.submitting && !st.awaiting {
			formStartSubmit(st)
		}
	}
	if st.awaiting {
		parts.submitted = formSettle(c, st)
	}
	return parts
}

// formStartSubmit validates every field and starts the async checks still owed.
func formStartSubmit(st *FormState) {
	st.tried, st.awaiting, st.formError = true, true, ""
	for _, name := range st.order {
		f := st.fields[name]
		f.err = f.check()
		if f.err == "" {
			f.startAsync()
		}
	}
	slog.Debug("mujica: form submit requested", "fields", len(st.order))
}

// formSettle finishes a waiting submit once no check is pending; it reports success.
func formSettle(c *ui.Context, st *FormState) bool {
	first, pending := "", false
	for _, name := range st.order {
		if st.Error(name) != "" && first == "" {
			first = name
		}
		pending = pending || st.fields[name].pending
	}
	switch {
	case first != "":
		st.awaiting, st.focus = false, first
		c.Announce(msg(c, "form.failed"))
		c.Invalidate()
		slog.Debug("mujica: form submit failed", "first", first)
	case pending:
	default:
		st.awaiting, st.submitting = false, true
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
	st.fields = map[string]*formFieldState{}
	st.tried, st.awaiting, st.requested, st.focus, st.formError = false, false, false, "", ""
	slog.Debug("mujica: form reset")
}

// FormItem registers field name, labels the control that control returns, and shows its error.
// control must return the element to focus when the field has the first error.
func FormItem[T any](f *FormScope, name string, value *T, opts FormItemOptions[T], control func() *ui.Element) *ui.Element {
	c, st := f.c, f.st
	if name == "" || opts.Label == "" {
		panic("mujica: FormItem needs a name and a Label")
	}
	if slices.Contains(st.order, name) {
		panic("mujica: FormItem " + name + " registered twice")
	}
	st.order = append(st.order, name)
	snap := fmt.Sprintf("%#v", *value)
	fs := st.fields[name]
	if fs == nil {
		st.seq++
		fs = &formFieldState{snapshot: snap, version: st.seq}
		st.fields[name] = fs
	}
	fs.check = func() string {
		if opts.Required && formEmpty(*value) {
			return msg(c, "form.required")
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
	if snap != fs.snapshot {
		st.seq++
		fs.snapshot, fs.version, fs.pending, fs.asyncErr = snap, st.seq, false, ""
		if fs.touched || st.tried {
			fs.err = fs.check()
		}
	}
	initial := opts.Initial
	st.resets = append(st.resets, func() { *value = initial })
	desc := opts.Description
	if fs.pending {
		desc = msg(c, "form.checking")
	}
	var ctrl *ui.Element
	field := formFieldBuild(c, opts.Label, FormFieldOptions{Required: opts.Required, Description: desc, Error: st.Error(name)}, func() {
		ctrl = control()
	})
	if ctrl == nil {
		panic("mujica: FormItem " + name + ": control returned nil")
	}
	focused := ctrl.FocusWithin()
	if fs.focused && !focused {
		fs.touched = true
		fs.err = fs.check()
		if fs.err == "" {
			fs.startAsync()
		}
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
