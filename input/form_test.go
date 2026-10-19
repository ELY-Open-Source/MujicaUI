package input

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

type formProbe struct {
	st              FormState
	name, email     string
	tags            []string
	asyncCalls      []uint64
	asyncValues     []string
	submits, resets int
	horizontal      bool
	useAsync        bool
}

func (p *formProbe) run(c *ui.Context) {
	core.Use(c, core.Settings{})
	p.view(c)
}

func (p *formProbe) view(c *ui.Context) {
	res := Form(c, &p.st, FormOptions{Horizontal: p.horizontal}, func(f *FormScope) {
		FormItem(c, f, "name", &p.name, FormItemOptions[string]{Label: "Name", Required: true, Initial: "Sakiko"}, func() ui.Element {
			return InputGroup(c, &p.name, InputGroupOptions{Invalid: p.st.Error("name") != ""}).Input
		})
		opts := FormItemOptions[string]{Label: "Email", Validate: func(v string) string {
			if v != "" && !strings.Contains(v, "@") {
				return "Enter an email address."
			}
			return ""
		}}
		if p.useAsync {
			opts.Async = func(v string, version uint64) {
				p.asyncCalls = append(p.asyncCalls, version)
				p.asyncValues = append(p.asyncValues, v)
			}
		}
		FormItem(c, f, "email", &p.email, opts, func() ui.Element {
			return InputGroup(c, &p.email, InputGroupOptions{Invalid: p.st.Error("email") != ""}).Input
		})
		FormItem(c, f, "tags", &p.tags, FormItemOptions[[]string]{Label: "Tags", Required: true, Initial: []string{"rose"}}, func() ui.Element {
			return TagsInput(c, &p.tags, TagsInputOptions{}).Input
		})
	})
	if res.Submitted() {
		p.submits++
	}
	if res.Reset() {
		p.resets++
	}
}

func TestFormBlurValidationAndUntouchedFields(t *testing.T) {
	p := &formProbe{name: "Sakiko", tags: []string{"rose"}}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Email")
	tt.Type("mutsumi")
	if tt.HasText("Enter an email address.") {
		t.Fatal("error before blur")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.HasText("Enter an email address.") {
		t.Fatal("no error after blur")
	}
	if tt.HasText("Required") {
		t.Fatal("untouched field shows an error")
	}
	_ = tt.Click("Email")
	tt.Type("@court")
	if tt.HasText("Enter an email address.") || p.email != "mutsumi@court" {
		t.Fatalf("touched field did not revalidate: %q", p.email)
	}
}

func TestFormSubmitFailureFocusesFirstErrorAndKeepsValues(t *testing.T) {
	p := &formProbe{email: "bad"}
	tt := ui.NewTester(p.run, 600, 500)
	if err := tt.Click("Submit"); err != nil {
		t.Fatal(err)
	}
	if p.submits != 0 || !tt.HasText("Required") || !tt.HasText("Enter an email address.") {
		t.Fatalf("submits %d texts %q", p.submits, tt.Texts())
	}
	if !tt.Focused("Name") || p.email != "bad" {
		t.Fatal("first error not focused or value lost")
	}
	if got := tt.Announcements(); len(got) == 0 || got[len(got)-1] != "Some fields need attention" {
		t.Fatalf("announcements %q", got)
	}
}

func TestFormSubmitOnceAndEndSubmit(t *testing.T) {
	p := &formProbe{name: "Uika", tags: []string{"x"}}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Submit")
	if p.submits != 1 || !p.st.Submitting() {
		t.Fatalf("submits %d", p.submits)
	}
	_ = tt.Click("Submit")
	p.st.Submit()
	tt.Frame()
	if p.submits != 1 {
		t.Fatalf("re-entered: %d", p.submits)
	}
	p.st.EndSubmit("The server is away.")
	tt.Frame()
	if !tt.HasText("The server is away.") || p.name != "Uika" {
		t.Fatal("form error not shown")
	}
	_ = tt.Click("Submit")
	if p.submits != 2 || tt.HasText("The server is away.") {
		t.Fatalf("second submit: %d", p.submits)
	}
}

func TestFormAsyncStaleResultsAndPendingSubmit(t *testing.T) {
	p := &formProbe{name: "Uika", tags: []string{"x"}, useAsync: true}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Email")
	tt.Type("a@b")
	tt.Key(0, ui.KeyTab)
	if len(p.asyncCalls) != 1 || p.asyncValues[0] != "a@b" || !tt.HasText("Checking…") {
		t.Fatalf("async calls %v", p.asyncCalls)
	}
	old := p.asyncCalls[0]
	_ = tt.Click("Email")
	tt.Type("c")
	if p.st.Resolve("email", old, "Taken") {
		t.Fatal("stale result accepted")
	}
	tt.Frame()
	if tt.HasText("Taken") {
		t.Fatal("stale result shown")
	}
	_ = tt.Click("Submit")
	if len(p.asyncCalls) != 2 || p.asyncValues[1] != "a@bc" || p.submits != 0 {
		t.Fatalf("submit did not wait for the new check: %v %d", p.asyncValues, p.submits)
	}
	if !p.st.Resolve("email", p.st.Version("email"), "") {
		t.Fatal("current result refused")
	}
	tt.Frame()
	if p.submits != 1 {
		t.Fatalf("submit did not finish after the check: %d", p.submits)
	}
}

func TestFormAsyncErrorBlocksSubmit(t *testing.T) {
	p := &formProbe{name: "Uika", tags: []string{"x"}, email: "a@b", useAsync: true}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Submit")
	p.st.Resolve("email", p.st.Version("email"), "Already registered")
	tt.Frame()
	tt.Frame()
	if p.submits != 0 || !tt.HasText("Already registered") || !tt.Focused("Email") {
		t.Fatalf("submits %d texts %q", p.submits, tt.Texts())
	}
}

func TestFormResetRestoresInitialValues(t *testing.T) {
	p := &formProbe{name: "", email: "bad", tags: nil}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Submit")
	if err := tt.Click("Reset"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if p.name != "Sakiko" || p.email != "" || len(p.tags) != 1 || p.tags[0] != "rose" || p.resets != 1 {
		t.Fatalf("reset values %q %q %q", p.name, p.email, p.tags)
	}
	if tt.HasText("Required") || p.st.Error("email") != "" {
		t.Fatal("errors survive reset")
	}
}

func TestFormHorizontalLayout(t *testing.T) {
	p := &formProbe{horizontal: true}
	tt := ui.NewTester(p.run, 600, 500)
	star, _ := tt.Find("*")
	email, _ := tt.Find("Email")
	input, _ := tt.Find("Name")
	if star.X+star.W != email.X+email.W || input.X >= email.X {
		t.Fatalf("labels not right-aligned: %v %v", star, email)
	}
}

// An edit during a pending check must be checked again before the submit goes out.
func TestFormEditWhilePendingRechecks(t *testing.T) {
	p := &formProbe{name: "Uika", tags: []string{"x"}, email: "a@b", useAsync: true}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Submit")
	if len(p.asyncCalls) != 1 {
		t.Fatalf("calls %v", p.asyncCalls)
	}
	_ = tt.Click("Email")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("taken@b")
	tt.Frame()
	if p.submits != 0 || len(p.asyncCalls) != 2 || p.asyncValues[1] != "taken@b" {
		t.Fatalf("submits %d values %q", p.submits, p.asyncValues)
	}
	if p.st.Resolve("email", p.asyncCalls[0], "") {
		t.Fatal("stale result accepted")
	}
	p.st.Resolve("email", p.asyncCalls[1], "Taken")
	tt.Frame()
	tt.Frame()
	if p.submits != 0 || !tt.HasText("Taken") {
		t.Fatalf("submitted unchecked value: %d", p.submits)
	}
}

// Reset also drops drafts and errors kept inside the controls.
func TestFormResetClearsControlState(t *testing.T) {
	var st FormState
	tags, seen := []string{}, ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Form(c, &st, FormOptions{}, func(f *FormScope) {
			FormItem(c, f, "tags", &tags, FormItemOptions[[]string]{Label: "Tags"}, func() ui.Element {
				return TagsInput(c, &tags, TagsInputOptions{Validate: func(s string) string {
					seen = s
					return "Invalid tag"
				}}).Input
			})
		})
	}, 600, 400)
	_ = tt.Click("Tags")
	tt.Type("bad tag")
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Invalid tag") {
		t.Fatal("no tag error")
	}
	_ = tt.Click("Reset")
	tt.Frame()
	if tt.HasText("Invalid tag") {
		t.Fatal("tag error survived reset")
	}
	seen = ""
	_ = tt.Click("Tags")
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if seen != "x" {
		t.Fatalf("draft survived reset: %q", seen)
	}
}

func TestFormSnapshot(t *testing.T) {
	p := &formProbe{name: "Sakiko Togawa", email: "oblivionis", tags: []string{"Ave Mujica"}}
	snapshot(t, "form", 560, 420, func(c *ui.Context) {
		ui.Column(c).Padding(24).Children(func() {
			p.view(c)
			if !p.st.tried {
				p.st.Submit()
				c.Invalidate()
			}
		})
	})
}

// Text typed in the frame the check passes gets its own check before submitting.
func TestFormInputAfterAsyncPassIsRechecked(t *testing.T) {
	p := &formProbe{name: "Uika", tags: []string{"x"}, email: "good@x", useAsync: true}
	tt := ui.NewTester(p.run, 600, 500)
	_ = tt.Click("Email")
	p.st.Submit()
	tt.Frame()
	if len(p.asyncCalls) != 1 {
		t.Fatalf("calls %v", p.asyncCalls)
	}
	p.st.Resolve("email", p.asyncCalls[0], "")
	tt.Type("bad")
	if p.submits != 0 || len(p.asyncValues) != 2 || p.asyncValues[1] != "good@xbad" {
		t.Fatalf("submits %d checked %q", p.submits, p.asyncValues)
	}
}

// A cross-field rule sees the other field's latest value when the submit settles.
func TestFormCrossFieldRuleAtSettle(t *testing.T) {
	var st FormState
	user, pw, confirm := "uika", "a", "a"
	var version uint64
	submits := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res := Form(c, &st, FormOptions{}, func(f *FormScope) {
			FormItem(c, f, "user", &user, FormItemOptions[string]{Label: "User", Async: func(_ string, v uint64) { version = v }}, func() ui.Element {
				return InputGroup(c, &user, InputGroupOptions{}).Input
			})
			FormItem(c, f, "pw", &pw, FormItemOptions[string]{Label: "Password"}, func() ui.Element {
				return InputGroup(c, &pw, InputGroupOptions{}).Input
			})
			FormItem(c, f, "confirm", &confirm, FormItemOptions[string]{Label: "Confirm", Validate: func(v string) string {
				if v != pw {
					return "Passwords differ."
				}
				return ""
			}}, func() ui.Element {
				return InputGroup(c, &confirm, InputGroupOptions{}).Input
			})
		})
		if res.Submitted() {
			submits++
		}
	}, 600, 500)
	_ = tt.Click("Submit")
	pw = "b"
	tt.Frame()
	st.Resolve("user", version, "")
	tt.Frame()
	tt.Frame()
	if submits != 0 || !tt.HasText("Passwords differ.") {
		t.Fatalf("mismatched pair submitted: %d", submits)
	}
}
