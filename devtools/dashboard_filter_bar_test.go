package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func filterTester(f *DashboardFilters, opts DashboardFilterBarOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(8).Children(func() {
			if DashboardFilterBar(c, f, opts).Changed() {
				*changes++
			}
		})
	}, 900, 420)
}

var filterOpts = DashboardFilterBarOptions{Envs: []string{"production", "staging"}, Tags: []string{"api", "web", "db"}}

func TestDashboardFilterBarChooseAndClear(t *testing.T) {
	var f DashboardFilters
	n := 0
	tt := filterTester(&f, filterOpts, &n)
	_ = tt.Click("Time range")
	_ = tt.Click("Last 1 d")
	if f.Range != 24*time.Hour || n != 1 {
		t.Fatalf("range %v changes %d", f.Range, n)
	}
	_ = tt.Click("Environment")
	_ = tt.Click("staging")
	if f.Env != "staging" {
		t.Fatalf("env %q", f.Env)
	}
	_ = tt.Click("Tags")
	_ = tt.Click("db")
	if len(f.Tags) != 1 || f.Tags[0] != "db" {
		t.Fatalf("tags %v", f.Tags)
	}
	tt.Key(0, ui.KeyEscape)
	_ = tt.Click("Clear filters")
	if f.Range != 0 || f.Env != "" || len(f.Tags) != 0 || !tt.HasText("All environments") {
		t.Fatalf("clear: %+v", f)
	}
}

func TestDashboardFilterBarDisabledAndRules(t *testing.T) {
	f := DashboardFilters{Env: "production"}
	n := 0
	tt := filterTester(&f, DashboardFilterBarOptions{Envs: filterOpts.Envs, Disabled: true}, &n)
	_ = tt.Click("Clear filters")
	_ = tt.Click("Time range")
	_ = tt.Click("Last 1 h")
	if n != 0 || f.Env != "production" || f.Range != 0 {
		t.Fatalf("disabled changed: %+v", f)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unoffered env did not panic")
		}
	}()
	bad := DashboardFilters{Env: "moon"}
	filterTester(&bad, filterOpts, &n)
}

func TestDashboardFilterBarSnapshot(t *testing.T) {
	snapshot(t, "dashboard_filter_bar", 860, 90, func(c *ui.Context) {
		f := DashboardFilters{Range: time.Hour, Env: "production", Tags: []string{"api"}}
		ui.Column(c).Padding(16).Children(func() { DashboardFilterBar(c, &f, filterOpts) })
	})
}

func TestDashboardFilterBarKeyboardAndAncestorDisabled(t *testing.T) {
	var f DashboardFilters
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { DashboardFilterBar(c, &f, filterOpts) })
	}, 900, 420)
	_ = tt.Click("Environment")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if f.Env != "production" {
		t.Fatal("keyboard environment")
	}
	_ = tt.Click("Environment")
	disabled = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Clear filters")
	if f.Env != "production" || tt.HasText("staging") {
		t.Fatal("disabled filter or popup changed")
	}
}

func TestDashboardFilterBarRejectsUnknownTags(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unoffered tag accepted")
		}
	}()
	f := DashboardFilters{Tags: []string{"unknown"}}
	n := 0
	filterTester(&f, filterOpts, &n)
}
