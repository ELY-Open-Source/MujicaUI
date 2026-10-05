package main

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/account"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func accountClick(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func TestAccountPagesInteractiveEvents(t *testing.T) {
	actions := map[string]func(*testing.T, *ui.Tester){
		"461": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Privacy") },
		"462": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Search settings")
			tt.Type("audio")
			tt.Key(0, ui.KeyEnter)
		},
		"463": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Dark") },
		"464": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "#355E80") },
		"465": func(t *testing.T, tt *ui.Tester) {
			value := account.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}
			accountClick(t, tt, "Command palette: "+value.String())
			tt.Key(ui.Cmd, ui.KeyK)
		},
		"466": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Search shortcuts"); tt.Type("Save") },
		"467": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Email or username")
			tt.Type("test@court")
			accountClick(t, tt, "Password")
			tt.Type("test-secret-467")
			accountClick(t, tt, "Sign in")
		},
		"468": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Continue with Git host") },
		"469": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Verification code"); tt.Type("381942") },
		"470": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Name")
			tt.Type(" edited")
			accountClick(t, tt, "Save")
		},
		"471": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Sakiko Togawa")
			accountClick(t, tt, "Oblivionis · ob@mask")
		},
		"472": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Ave Mujica")
			tt.Type("arch")
			tt.Key(0, ui.KeyEnter)
		},
		"473": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Sakiko Togawa, Online")
			accountClick(t, tt, "Away")
		},
		"474": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Sign out of iPhone · App") },
		"475": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Create key")
			tt.Type("CI")
			tt.Key(0, ui.KeyEnter)
		},
		"476": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "API calls used"); tt.Key(0, ui.KeyRight) },
		"477": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Next") },
		"478": func(t *testing.T, tt *ui.Tester) { accountClick(t, tt, "Start tour"); accountClick(t, tt, "Next") },
		"479": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Show what's new")
			accountClick(t, tt, "Don't show again")
		},
		"480": func(t *testing.T, tt *ui.Tester) {
			accountClick(t, tt, "Message")
			tt.Type("This is a feedback test.")
			accountClick(t, tt, "Send feedback")
		},
	}
	if len(accountPages) != 20 {
		t.Fatalf("%d account pages", len(accountPages))
	}
	for _, p := range accountPages {
		t.Run(p.ID, func(t *testing.T) {
			if p.Category != "Settings & Account" || p.Code == "" || p.Keys == "" {
				t.Fatal("incomplete gallery metadata")
			}
			var events []string
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Column(c).Padding(20).Gap(16).Children(func() { p.Demo(c, func(s string) { events = append(events, s) }) })
			}, 1100, 800)
			actions[p.ID](t, tt)
			if len(events) == 0 {
				t.Fatal("interaction has no gallery event")
			}
			logs := strings.Join(events, "\n")
			if strings.Contains(logs, "test-secret-467") || strings.Contains(logs, "381942") {
				t.Fatal("gallery event exposes a secret")
			}
		})
	}
}
