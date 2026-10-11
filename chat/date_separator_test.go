package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestDateSeparatorLabels(t *testing.T) {
	today := core.DateOf(chatNow)
	for _, tc := range []struct {
		locale core.Locale
		day    core.Date
		want   string
	}{
		{core.EnUS, today, "Today"},
		{core.EnUS, today.AddDays(-1), "Yesterday"},
		{core.EnUS, core.Date{Year: 2026, Month: 9, Day: 1}, "Tuesday, September 1, 2026"},
		{core.ZhCN, today, "今天"},
		{core.ZhCN, today.AddDays(-1), "昨天"},
		{core.ZhCN, core.Date{Year: 2026, Month: 9, Day: 1}, "2026年9月1日 星期二"},
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Locale: tc.locale})
			DateSeparator(c, tc.day, DateSeparatorOptions{Today: today})
		}, 500, 80)
		if !tt.HasText(tc.want) {
			t.Errorf("%s %v: texts %q", tc.locale, tc.day, tt.Texts())
		}
	}
}

func TestDateSeparatorShortStamp(t *testing.T) {
	var got []string
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		today := core.DateOf(chatNow)
		got = []string{DateSeparatorWhen(c, chatNow, today), DateSeparatorWhen(c, chatNow.AddDate(0, 0, -1), today),
			DateSeparatorWhen(c, chatNow.AddDate(0, -2, 0), today)}
	}, 100, 100, 1)
	if got[0] != "15:30" || got[1] != "Yesterday" || got[2] != "August 5" {
		t.Fatalf("stamps %q", got)
	}
}

func TestDateSeparatorRejectsInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DateSeparator(c, core.Date{Year: 2026, Month: 2, Day: 30}, DateSeparatorOptions{})
	}, 100, 100, 1)
}

func TestDateSeparatorSnapshot(t *testing.T) {
	snapshot(t, "date_separator", 520, 110, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			DateSeparator(c, core.DateOf(chatNow), DateSeparatorOptions{Today: core.DateOf(chatNow)})
			DateSeparator(c, core.Date{Year: 2026, Month: 9, Day: 1}, DateSeparatorOptions{Today: core.DateOf(chatNow)})
		})
	})
}
