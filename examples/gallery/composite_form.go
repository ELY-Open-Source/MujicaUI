package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

const formID = "c-form"

// Simulated latencies of the name check and the server.
const (
	formCheckDelay  = 400 * time.Millisecond
	formSubmitDelay = 900 * time.Millisecond
)

// formPhase is what the form composite shows.
type formPhase int

const (
	formEditing formPhase = iota
	formSucceeded
	formFailed
)

type formSubmitDemo struct {
	st       mujica.FormState
	name     string
	password string
	debut    mujica.Date
	tags     []string
	terms    bool
	failNext bool
	phase    formPhase
	failure  string
	checks   chan formAsync
	done     chan string
}

func newFormSubmitDemo() formSubmitDemo {
	return formSubmitDemo{checks: make(chan formAsync, 8), done: make(chan string, 1)}
}

// formTaken are the names the simulated archive already holds.
var formTaken = []string{"admin", "oblivionis"}

func (g *gallery) formView(c *ui.Context, log func(string)) {
	s := &g.cForm
	for drained := false; !drained; {
		select {
		case r := <-s.checks:
			log(fmt.Sprintf("name check v%d: %q, kept %v", r.version, r.err, s.st.Resolve("name", r.version, r.err)))
		case e := <-s.done:
			s.st.EndSubmit("")
			s.phase, s.failure = formSucceeded, e
			if e != "" {
				s.phase = formFailed
			}
			log(fmt.Sprintf("server answered, error %q", e))
		default:
			drained = true
		}
	}
	mujica.Card(c, mujica.CardOptions{Title: g.t("Audition entry", "试演报名")}, func() {
		switch s.phase {
		case formSucceeded:
			g.formSuccess(c, log)
		case formFailed:
			g.formFailure(c, log)
		default:
			ui.Column(c).FillWidth().MaxWidth(640).Gap(16).Children(func() { g.formEdit(c, log) })
		}
	})
}

func (g *gallery) formEdit(c *ui.Context, log func(string)) {
	s := &g.cForm
	if mujica.Switch(c, &s.failNext, g.t("The server rejects the next submit", "服务器拒绝下一次提交"), mujica.SwitchOptions{Disabled: s.st.Submitting()}).Changed() {
		log(fmt.Sprintf("fail next submit %v", s.failNext))
	}
	r := mujica.LoadingOverlay(c, s.st.Submitting(), mujica.LoadingOverlayOptions{Label: g.t("Submitting…", "正在提交…")}, func() {
		res := mujica.Form(c, &s.st, mujica.FormOptions{Horizontal: true, SubmitLabel: g.t("Submit entry", "提交报名")}, func(f *mujica.FormScope) {
			g.formFields(c, f, log)
		})
		res.Element.FillWidth()
		if res.Reset() {
			log("reset to initial values")
		}
		if res.Submitted() {
			fail := ""
			if s.failNext {
				s.failNext = false
				fail = g.t("The archive is closed for the night. Your entry was not sent.", "档案馆今夜已闭馆，报名未送达。")
			}
			// The password is never logged, not even its length.
			log(fmt.Sprintf("submitted %q, debut %s, roles %q", s.name, s.debut.Format("2006-01-02"), s.tags))
			go func() {
				time.Sleep(formSubmitDelay)
				s.done <- fail
				c.Invalidate()
			}()
		}
	})
	r.Element.FillWidth()
}

func (g *gallery) formFields(c *ui.Context, f *mujica.FormScope, log func(string)) {
	s := &g.cForm
	mujica.FormItem(c, f, "name", &s.name, mujica.FormItemOptions[string]{Label: g.t("Stage name", "艺名"), Required: true,
		Description: g.t("Checked against the archive; “Oblivionis” is taken.", "将与档案核对；“Oblivionis”已被占用。"),
		Validate: func(v string) string {
			if len([]rune(strings.TrimSpace(v))) < 2 {
				return g.t("Use at least 2 characters.", "至少 2 个字符。")
			}
			return ""
		},
		Async: func(v string, version uint64) {
			log(fmt.Sprintf("name check started v%d", version))
			go func() {
				time.Sleep(formCheckDelay)
				msg := ""
				for _, t := range formTaken {
					if strings.EqualFold(strings.TrimSpace(v), t) {
						msg = g.t("That name is taken.", "该名称已被占用。")
					}
				}
				s.checks <- formAsync{version, msg}
				c.Invalidate()
			}()
		}}, func() *ui.Element {
		return mujica.TextInput(c, &s.name, mujica.TextInputOptions{}).Input
	})
	mujica.FormItem(c, f, "password", &s.password, mujica.FormItemOptions[string]{Label: g.t("Password", "密码"), Required: true,
		Validate: func(v string) string {
			if len([]rune(v)) < 8 {
				return g.t("Use at least 8 characters.", "至少 8 个字符。")
			}
			return ""
		}}, func() *ui.Element {
		r := mujica.PasswordInput(c, &s.password, mujica.PasswordInputOptions{})
		if r.Submitted() {
			s.st.Submit()
		}
		return r.Input
	})
	mujica.FormItem(c, f, "debut", &s.debut, mujica.FormItemOptions[mujica.Date]{Label: g.t("Debut date", "首演日期"), Required: true}, func() *ui.Element {
		r := mujica.DatePicker(c, &s.debut, mujica.DatePickerOptions{Clearable: true})
		if r.Changed() {
			log("debut " + s.debut.Format("2006-01-02"))
		}
		return r.Input
	})
	mujica.FormItem(c, f, "roles", &s.tags, mujica.FormItemOptions[[]string]{Label: g.t("Roles", "担当"), Required: true,
		Description: g.t("Up to 4; press Enter after each.", "最多 4 个；每个后按 Enter。")}, func() *ui.Element {
		r := mujica.TagsInput(c, &s.tags, mujica.TagsInputOptions{Max: 4, Invalid: s.st.Error("roles") != "", Placeholder: g.t("keyboard, bass…", "键盘、贝斯…")})
		if r.Changed() {
			log(fmt.Sprintf("roles %q", s.tags))
		}
		return r.Input
	})
	mujica.FormItem(c, f, "terms", &s.terms, mujica.FormItemOptions[bool]{Label: g.t("Terms", "条款"),
		Validate: func(v bool) string {
			if !v {
				return g.t("Accept the terms to continue.", "请先同意条款。")
			}
			return ""
		}}, func() *ui.Element {
		e := mujica.Checkbox(c, &s.terms, g.t("I accept the covenant of the masquerade", "我接受假面舞会的盟约"), mujica.CheckboxOptions{})
		if e.Changed() {
			log(fmt.Sprintf("terms %v", s.terms))
		}
		return e
	})
}

func (g *gallery) formSuccess(c *ui.Context, log func(string)) {
	s := &g.cForm
	r := mujica.Result(c, mujica.ResultOptions{Kind: mujica.ResultSuccess, Title: g.t("Entry received", "报名已收到"),
		Description: fmt.Sprintf(g.t("%s, your debut on %s is on the ledger.", "%s，你在 %s 的首演已登记。"), s.name, s.debut.Format("2006-01-02")),
		Actions:     []string{g.t("New entry", "新的报名")}})
	if r.Action == 0 {
		*s = newFormSubmitDemo()
		log("new entry")
	}
}

func (g *gallery) formFailure(c *ui.Context, log func(string)) {
	s := &g.cForm
	r := mujica.Result(c, mujica.ResultOptions{Kind: mujica.ResultFailure, Title: g.t("Entry not sent", "报名未送达"), Description: s.failure,
		Actions: []string{g.t("Retry", "重试"), g.t("Back to the form", "返回表单")}})
	switch r.Action {
	case 0:
		s.phase = formEditing
		s.st.Submit()
		log("retry")
	case 1:
		s.phase = formEditing
		log("back to the form")
	}
}
