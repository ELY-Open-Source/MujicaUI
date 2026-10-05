package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func processListSample() []Process {
	return []Process{
		{PID: 101, Name: "chancellor", CPU: 3.5, Memory: 200 << 20},
		{PID: 7, Name: "archivist", CPU: 42.1, Memory: 50 << 20},
		{PID: 330, Name: "Herald", CPU: 0.2, Memory: 2 << 30},
	}
}

func processListTester(procs []Process, s *ProcessListState, disabled bool, ended *[]int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ProcessList(c, procs, s, ProcessListOptions{Disabled: disabled})
		r.Element.Height(300)
		if r.End != nil {
			*ended = append(*ended, r.End.PID)
		}
	}, 560, 340)
}

func processListOrder(procs []Process, s *ProcessListState) []int {
	var out []int
	for _, i := range s.order {
		out = append(out, procs[i].PID)
	}
	return out
}

func TestProcessListSortsByCPUThenMemory(t *testing.T) {
	procs := processListSample()
	s := &ProcessListState{}
	var ended []int
	tt := processListTester(procs, s, false, &ended)
	if got := processListOrder(procs, s); got[0] != 7 || got[2] != 330 {
		t.Fatalf("default CPU order %v", got)
	}
	if err := tt.Click("Memory"); err != nil {
		t.Fatal(err)
	}
	if got := processListOrder(procs, s); got[0] != 7 || got[2] != 330 || s.Sort.Descending {
		t.Fatalf("memory ascending %v", got)
	}
	_ = tt.Click("Memory")
	if got := processListOrder(procs, s); got[0] != 330 {
		t.Fatalf("memory descending %v", got)
	}
	if !tt.HasText("2.0 GB") || !tt.HasText("42.1%") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestProcessListEndAsksFirst(t *testing.T) {
	procs := processListSample()
	s := &ProcessListState{}
	var ended []int
	tt := processListTester(procs, s, false, &ended)
	if err := tt.Click("End Herald"); err != nil {
		t.Fatal(err)
	}
	if len(ended) != 0 || !tt.HasText("End “Herald” (PID 330)?") {
		t.Fatalf("ended %v before confirmation", ended)
	}
	if err := tt.Click("End process"); err != nil {
		t.Fatal(err)
	}
	if len(ended) != 1 || ended[0] != 330 {
		t.Fatalf("ended %v", ended)
	}
}

func TestProcessListKeyboardEnd(t *testing.T) {
	procs := processListSample()
	s := &ProcessListState{}
	var ended []int
	tt := processListTester(procs, s, false, &ended)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDelete)
	if !tt.HasText("End “archivist” (PID 7)?") {
		t.Fatalf("Delete did not ask for the chosen row: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if len(ended) != 0 {
		t.Fatal("Escape ended the process")
	}
}

func TestProcessListDisabledAndEmpty(t *testing.T) {
	var ended []int
	tt := processListTester(processListSample(), &ProcessListState{}, true, &ended)
	_ = tt.Click("End Herald")
	if tt.HasText("End process") {
		t.Fatal("disabled list asked to end a process")
	}
	tt = processListTester(nil, &ProcessListState{}, false, &ended)
	if !tt.HasText("No processes") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestProcessListBadSortPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unknown sort column did not panic")
		}
	}()
	processListSort(processListSample(), &ProcessListState{Sort: ui.SortOrder{Column: "owner"}})
}

func TestProcessListSnapshot(t *testing.T) {
	snapshot(t, "process_list", 560, 220, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *ProcessListState { return &ProcessListState{} })
		ui.Column(c).Padding(12).Children(func() {
			ProcessList(c, processListSample(), *s, ProcessListOptions{}).Element.Height(196)
		})
	})
}
