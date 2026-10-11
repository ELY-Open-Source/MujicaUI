package devtools

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func kvTester(rows *[]KeyValue, opts KeyValueInputOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if KeyValueInput(c, rows, opts).Changed() {
			*changes++
		}
	}, 600, 300)
}

func TestKeyValueInputEditToggleAddRemove(t *testing.T) {
	rows := []KeyValue{{Key: "page", Value: "1", Enabled: true}, {Key: "debug", Value: "true"}}
	n := 0
	tt := kvTester(&rows, KeyValueInputOptions{}, &n)
	_ = tt.Click("Row 2") // its enable box
	if !rows[1].Enabled || n != 1 {
		t.Fatalf("toggle: %+v", rows)
	}
	_ = tt.Click("Add row")
	if len(rows) != 3 || !rows[2].Enabled {
		t.Fatalf("add: %+v", rows)
	}
	_ = tt.Click("Row 3 Key")
	tt.Type("sort")
	tt.Key(0, ui.KeyTab)
	tt.Type("name")
	if rows[2].Key != "sort" || rows[2].Value != "name" {
		t.Fatalf("typing: %+v", rows[2])
	}
	_ = tt.Click("Remove row 1")
	if len(rows) != 2 || rows[0].Key != "debug" {
		t.Fatalf("remove: %+v", rows)
	}
	if got := keyValueEnabled(rows); len(got) != 2 {
		t.Fatalf("enabled pairs: %+v", got)
	}
}

func TestKeyValueInputDisabled(t *testing.T) {
	rows := []KeyValue{{Key: "page", Value: "1"}}
	n := 0
	tt := kvTester(&rows, KeyValueInputOptions{Disabled: true}, &n)
	_ = tt.Click("Row 1")
	_ = tt.Click("Add row")
	_ = tt.Click("Remove row 1")
	_ = tt.Click("Row 1 Key")
	tt.Type("x")
	if n != 0 || len(rows) != 1 || rows[0].Enabled || rows[0].Key != "page" {
		t.Fatalf("disabled editor changed: %+v", rows)
	}
}

func TestKeyValueInputSnapshot(t *testing.T) {
	snapshot(t, "key_value_input", 560, 200, func(c *ui.Context) {
		rows := []KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}, {Key: "X-Trace", Value: "on"}}
		ui.Column(c).Padding(16).Children(func() { KeyValueInput(c, &rows, KeyValueInputOptions{}) })
	})
}

func TestKeyValueInputManyRows(t *testing.T) {
	rows := make([]KeyValue, 100000)
	for i := range rows {
		rows[i] = KeyValue{Key: "item", Value: "value", Enabled: true}
	}
	n := 0
	tt := kvTester(&rows, KeyValueInputOptions{}, &n)
	if got := len(tt.Texts()); got > 500 {
		t.Fatalf("built %d texts for 100000 rows", got)
	}
}

func TestKeyValueInputAncestorDisabledAndEmpty(t *testing.T) {
	var rows []KeyValue
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() { KeyValueInput(c, &rows, KeyValueInputOptions{}) })
	}, 600, 300)
	_ = tt.Click("Add row")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(rows) != 0 {
		t.Fatal("disabled empty editor changed")
	}
}
