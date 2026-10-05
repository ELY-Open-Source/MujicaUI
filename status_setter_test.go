package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestStatusUntil(t *testing.T) {
	wed := time.Date(2026, 10, 7, 14, 30, 0, 0, time.UTC)
	for d, want := range map[time.Duration]time.Time{
		0:                {},
		time.Hour:        wed.Add(time.Hour),
		-1:               time.Date(2026, 10, 8, 0, 0, 0, 0, time.UTC),
		-2:               time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC),
		30 * time.Minute: wed.Add(30 * time.Minute),
	} {
		if got := statusUntil(wed, d); !got.Equal(want) {
			t.Errorf("%v: %v, want %v", d, got, want)
		}
	}
}

func TestStatusSetterEditAndClear(t *testing.T) {
	st := CustomStatus{}
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(20).Width(420).Children(func() {
			if StatusSetter(c, &st, StatusSetterOptions{Location: time.UTC}).Changed() {
				n++
			}
		})
	}, 520, 600)
	if err := tt.Click("Choose an emoji"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Activities"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("violin"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Status text"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Rehearsing")
	if err := tt.Click("Don't clear"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("1 hour"); err != nil {
		t.Fatal(err)
	}
	if st.Emoji != "🎻" || st.Text != "Rehearsing" || st.Until.IsZero() || n < 3 {
		t.Fatalf("status %+v changes %d", st, n)
	}
	if err := tt.Click("Clear status"); err != nil {
		t.Fatal(err)
	}
	if st != (CustomStatus{}) {
		t.Fatalf("not cleared: %+v", st)
	}
}

func TestStatusSetterSnapshot(t *testing.T) {
	st := CustomStatus{Emoji: "🎻", Text: "Rehearsing the second act"}
	snapshot(t, "status_setter", 460, 140, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(420).Children(func() { StatusSetter(c, &st, StatusSetterOptions{Location: time.UTC}) })
	})
}

func TestStatusSetterLoadedExpiryAndDisabledPopup(t *testing.T) {
	st := CustomStatus{Until: time.Date(2026, 10, 6, 15, 0, 0, 0, time.UTC)}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { StatusSetter(c, &st, StatusSetterOptions{Location: time.UTC}) })
	}, 600, 400)
	if tt.HasText("Don't clear") {
		t.Fatal("loaded expiry appears as never")
	}
	if err := tt.Click("Tue 15:00"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Don't clear"); err != nil {
		t.Fatal(err)
	}
	if !st.Until.IsZero() {
		t.Fatal("loaded expiry was not cleared")
	}
	if err := tt.Click("Don't clear"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	if tt.HasText("30 minutes") {
		t.Fatal("disabled expiry popup remains visible")
	}
	disabled = false
	tt.Frame()
	if tt.HasText("30 minutes") {
		t.Fatal("expiry popup reopened after enabling")
	}
}
