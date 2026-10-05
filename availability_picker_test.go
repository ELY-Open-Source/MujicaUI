package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func availabilityGrid(days, slots int) [][]bool {
	v := make([][]bool, days)
	for i := range v {
		v[i] = make([]bool, slots)
	}
	return v
}

var availabilityDays = []Date{oct(12), oct(13), oct(14)}

func availabilityTester(v *[][]bool, opts *AvailabilityPickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		o := *opts
		o.Days, o.From, o.To = availabilityDays, 9, 12
		ui.Column(c).Padding(16).Children(func() {
			if AvailabilityPicker(c, v, o).Changed() {
				*n++
			}
		})
	}, 560, 360)
}

func TestAvailabilityPickerPaintDrag(t *testing.T) {
	v, n := availabilityGrid(3, 6), 0
	tt := availabilityTester(&v, &AvailabilityPickerOptions{}, &n)
	a, _ := tt.Find("Monday, October 12, 2026 09:00")
	b, _ := tt.Find("Monday, October 12, 2026 10:30")
	calDrag(tt, a.X+a.W/2, a.Y+a.H/2, b.X+b.W/2, b.Y+b.H/2)
	for s := range 6 {
		if v[0][s] != (s <= 3) {
			t.Fatalf("painted %v", v[0])
		}
	}
	if n == 0 || !tt.HasText("Monday, October 12, 2026 10:00, you are free") {
		t.Fatal("no change or label")
	}
	// Starting on a free slot erases.
	calDrag(tt, a.X+a.W/2, a.Y+a.H/2, a.X+a.W/2, a.Y+a.H*1.5)
	if v[0][0] || v[0][1] || !v[0][2] {
		t.Fatalf("erase stroke %v", v[0])
	}
}

func TestAvailabilityPickerOthersAndKeyboard(t *testing.T) {
	v, n := availabilityGrid(3, 6), 0
	free := availabilityGrid(3, 6)
	free[1][0] = true
	other := availabilityGrid(3, 6)
	other[1][0], other[1][1] = true, true
	opts := &AvailabilityPickerOptions{Others: []AvailabilityPerson{{Name: "Mutsumi", Free: free}, {Name: "Uika", Free: other}}}
	tt := availabilityTester(&v, opts, &n)
	if !tt.HasText("Tuesday, October 13, 2026 09:00, 2 of 2 others free") || !tt.HasText("Tuesday, October 13, 2026 09:30, 1 of 2 others free") {
		t.Fatalf("overlap counts: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !v[1][0] || !v[1][1] || n != 2 {
		t.Fatalf("keyboard toggles %v", v[1])
	}
}

func TestAvailabilityPickerDisabledAndInvalid(t *testing.T) {
	v, n := availabilityGrid(3, 6), 0
	tt := availabilityTester(&v, &AvailabilityPickerOptions{Disabled: true}, &n)
	a, _ := tt.Find("Monday, October 12, 2026 09:00")
	calDrag(tt, a.X+a.W/2, a.Y+a.H/2, a.X+a.W/2, a.Y+a.H*2.5)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if v[0][0] || n != 0 {
		t.Fatal("disabled picker painted")
	}
	for _, bad := range []*AvailabilityPickerOptions{{SlotMinutes: 7}, {Others: []AvailabilityPerson{{Name: "X", Free: availabilityGrid(2, 6)}}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: expected panic", bad)
				}
			}()
			w := availabilityGrid(3, 6)
			availabilityTester(&w, bad, &n)
		}()
	}
}

func TestAvailabilityPickerSnapshot(t *testing.T) {
	v := availabilityGrid(3, 6)
	v[0][1], v[0][2], v[1][2] = true, true, true
	free := availabilityGrid(3, 6)
	free[0][2], free[1][2], free[2][0] = true, true, true
	snapshot(t, "availability_picker", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			AvailabilityPicker(c, &v, AvailabilityPickerOptions{Days: availabilityDays, From: 9, To: 12,
				Others: []AvailabilityPerson{{Name: "Mutsumi", Free: free}}})
		})
	})
}
