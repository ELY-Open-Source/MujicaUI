package mujica

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func workloadSample() []Workload {
	return []Workload{{Member: projMembers[0], Assigned: 32, Capacity: 40}, {Member: projMembers[1], Assigned: 46.5, Capacity: 40},
		{Member: projMembers[2], Assigned: 12, Capacity: 20}}
}

func workloadTester(loads []Workload, opts WorkloadViewOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Column(c).Height(300).Width(600).Children(func() { WorkloadView(c, loads, opts) })
	}, 600, 300)
}

func TestWorkloadViewOverload(t *testing.T) {
	tt := workloadTester(workloadSample(), WorkloadViewOptions{})
	for _, s := range []string{"1 person is over capacity", "46.5 / 40 h", "Mutsumi Wakaba, 46.5 / 40 h over by 6.5 h", "Sakiko Togawa, 32 / 40 h "} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	r, _ := tt.Find("Mutsumi Wakaba, 46.5 / 40 h over by 6.5 h")
	img := tt.Image()
	found := false
	for x := int(r.X); x < int(r.X+r.W); x++ {
		if datetimeNear(img.RGBAAt(x, int(r.Y+r.H/2)), theme.Light().Danger) {
			found = true
		}
	}
	if !found {
		t.Fatal("overload is not drawn in Danger")
	}
	if tt := workloadTester(workloadSample()[:1], WorkloadViewOptions{Unit: "pts"}); !tt.HasText("32 / 40 pts") || tt.HasText("over capacity") {
		t.Fatal("custom unit or false alarm")
	}
}

func TestWorkloadViewLargeAndInvalid(t *testing.T) {
	var loads []Workload
	for i := range 100000 {
		loads = append(loads, Workload{Member: Member{ID: fmt.Sprint(i), Name: "M"}, Assigned: float64(i % 50), Capacity: 40})
	}
	if n := len(workloadTester(loads, WorkloadViewOptions{}).Texts()); n > 200 {
		t.Fatalf("built %d labels", n)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	workloadTester([]Workload{{Member: projMembers[0], Capacity: 0}}, WorkloadViewOptions{})
}

func TestWorkloadViewSnapshot(t *testing.T) {
	loads := workloadSample()
	snapshot(t, "workload_view", 600, 200, func(c *ui.Context) {
		ui.Column(c).Height(200).Padding(16).Children(func() { WorkloadView(c, loads, WorkloadViewOptions{}) })
	})
}

func TestWorkloadViewRejectsInfinite(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected infinite load panic")
		}
	}()
	workloadTester([]Workload{{Member: projMembers[0], Assigned: math.Inf(1), Capacity: 40}}, WorkloadViewOptions{})
}
