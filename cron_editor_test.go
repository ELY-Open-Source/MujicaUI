package mujica

import (
	"errors"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func cronRuns(t *testing.T, expr string, from time.Time, loc *time.Location, n int) []CronRun {
	t.Helper()
	s, err := ParseCron(expr)
	if err != nil {
		t.Fatalf("%s: %v", expr, err)
	}
	return s.Next(from, loc, n)
}

func TestParseCronErrors(t *testing.T) {
	for _, c := range []struct {
		expr  string
		field int
		p     CronProblem
	}{{"* * * *", -1, CronFieldCount}, {"60 * * * *", 0, CronOutOfRange}, {"* 9-5 * * *", 1, CronBackwards},
		{"*/0 * * * *", 0, CronBadStep}, {"* * x * *", 2, CronNotNumber}, {"1,,2 * * * *", 0, CronEmptyItem}, {"* * * 13 *", 3, CronOutOfRange}} {
		_, err := ParseCron(c.expr)
		var ce *CronError
		if !errors.As(err, &ce) || ce.Field != c.field || ce.Problem != c.p {
			t.Errorf("%q: %v", c.expr, err)
		}
	}
}

func TestCronNextRuns(t *testing.T) {
	from := time.Date(2026, time.October, 5, 8, 50, 0, 0, calTokyo) // Monday
	got := cronRuns(t, "*/15 9-10 * * MON-FRI", from, calTokyo, 3)
	if len(got) != 3 || !got[0].At.Equal(calAt(t, 9, 0, 5)) || !got[2].At.Equal(calAt(t, 9, 30, 5)) {
		t.Fatalf("weekday runs %v", got)
	}
	got = cronRuns(t, "0 12 1 * SUN", from, calTokyo, 2)
	if !got[0].At.Equal(calAt(t, 12, 0, 11)) || !got[1].At.Equal(calAt(t, 12, 0, 18)) {
		t.Fatalf("day-of-month or day-of-week: %v", got)
	}
	got = cronRuns(t, "0 0 29 2 *", from, calTokyo, 1)
	if len(got) != 1 || got[0].At.Year() != 2028 {
		t.Fatalf("leap day %v", got)
	}
	if got = cronRuns(t, "0 0 31 2 *", from, calTokyo, 1); len(got) != 0 {
		t.Fatalf("Feb 31 never runs: %v", got)
	}
}

func TestCronNextRunsAcrossDST(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	got := cronRuns(t, "30 2 * * *", time.Date(2026, time.March, 7, 3, 0, 0, 0, ny), ny, 2)
	if got[0].At.Day() != 9 {
		t.Fatalf("02:30 on the gap day must not run: %v", got)
	}
	got = cronRuns(t, "30 1 * * *", time.Date(2026, time.October, 31, 3, 0, 0, 0, ny), ny, 1)
	if !got[0].Repeated || got[0].At.Day() != 1 {
		t.Fatalf("01:30 on fall-back is repeated: %v", got)
	}
}

func cronTester(expr *string, opts CronEditorOptions, n *int) *ui.Tester {
	opts.Location, opts.Now = calTokyo, time.Date(2026, time.October, 5, 8, 50, 0, 0, calTokyo)
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if CronEditor(c, expr, opts).Changed() {
				*n++
			}
		})
	}, 760, 520)
}

func TestCronEditorEditAndPreview(t *testing.T) {
	expr, n := "0 9 * * MON-FRI", 0
	tt := cronTester(&expr, CronEditorOptions{Runs: 3}, &n)
	for _, s := range []string{"Mo 2026-10-05 09:00 JST", "Tu 2026-10-06 09:00 JST", "MON–FRI", "Every value"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	_ = tt.Click("Expression")
	tt.Command("selectAll")
	tt.Type("0 25 * * *")
	if n == 0 || !tt.HasText("Hour: 25 is outside 0–23.") || tt.HasText("Next runs") {
		t.Fatalf("error not shown: %v", tt.Texts())
	}
	tt.Command("selectAll")
	tt.Type("*/30 * * * *")
	if !tt.HasText("Every 30") || !tt.HasText("Mo 2026-10-05 09:00 JST") {
		t.Fatalf("fixed expression: %v", tt.Texts())
	}
}

func TestCronEditorDisabledAndLocale(t *testing.T) {
	expr, n := "0 9 * * *", 0
	tt := cronTester(&expr, CronEditorOptions{Disabled: true}, &n)
	_ = tt.Click("Expression")
	tt.Type("x")
	if expr != "0 9 * * *" || n != 0 {
		t.Fatal("disabled editor changed")
	}
	bad := "* *"
	tt = ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		CronEditor(c, &bad, CronEditorOptions{Location: calTokyo})
	}, 760, 400)
	if !tt.HasText("需要 5 个字段，而不是 2 个。") {
		t.Fatalf("zh error: %v", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		CronEditor(c, &expr, CronEditorOptions{})
	}, 10, 10)
}

func TestCronEditorSnapshot(t *testing.T) {
	expr := "*/15 9-17 * * MON-FRI"
	snapshot(t, "cron_editor", 720, 400, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CronEditor(c, &expr, CronEditorOptions{Location: calTokyo, Now: time.Date(2026, time.October, 5, 8, 50, 0, 0, calTokyo)})
		})
	})
}
