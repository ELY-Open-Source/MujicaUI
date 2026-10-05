package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDateSeparatorLabels(t *testing.T) {
	today := DateOf(chatNow)
	for _, tc := range []struct {
		locale Locale
		day    Date
		want   string
	}{
		{EnUS, today, "Today"},
		{EnUS, today.AddDays(-1), "Yesterday"},
		{EnUS, Date{2026, 9, 1}, "Tuesday, September 1, 2026"},
		{ZhCN, today, "今天"},
		{ZhCN, today.AddDays(-1), "昨天"},
		{ZhCN, Date{2026, 9, 1}, "2026年9月1日 星期二"},
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Locale: tc.locale})
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
		Use(c, Settings{})
		today := DateOf(chatNow)
		got = []string{dateSeparatorWhen(c, chatNow, today), dateSeparatorWhen(c, chatNow.AddDate(0, 0, -1), today),
			dateSeparatorWhen(c, chatNow.AddDate(0, -2, 0), today)}
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
		Use(c, Settings{})
		DateSeparator(c, Date{2026, 2, 30}, DateSeparatorOptions{})
	}, 100, 100, 1)
}

func TestDateSeparatorSnapshot(t *testing.T) {
	snapshot(t, "date_separator", 520, 110, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			DateSeparator(c, DateOf(chatNow), DateSeparatorOptions{Today: DateOf(chatNow)})
			DateSeparator(c, Date{2026, 9, 1}, DateSeparatorOptions{Today: DateOf(chatNow)})
		})
	})
}
