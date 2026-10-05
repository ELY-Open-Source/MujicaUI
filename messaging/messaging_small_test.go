package messaging

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// mediaTextCount counts texts equal to s.
func mediaTextCount(tt *ui.Tester, s string) int {
	n := 0
	for _, x := range tt.Texts() {
		if x == s {
			n++
		}
	}
	return n
}

func TestMessagingChineseWeekdays(t *testing.T) {
	date := time.Date(2026, 10, 6, 15, 4, 0, 0, time.UTC)
	st := CustomStatus{Until: date}
	var until time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		StatusSetter(c, &st, StatusSetterOptions{Location: time.UTC})
		SnoozePicker(c, &until, SnoozePickerOptions{Now: time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC), Location: time.UTC})
		MailReader(c, Mail{Subject: "Test", From: "a@example.com", Date: date})
	}, 900, 900)
	if !tt.HasText("2026年10月6日 星期二 15:04") || !tt.HasText("星期二 15:04") || !tt.HasText("星期二 09:00") {
		t.Fatalf("weekday texts=%v", tt.Texts())
	}
}
