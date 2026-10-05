package agent

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

type approvalApp struct {
	open            bool
	allowed, denied int
}

func (a *approvalApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Run", input.ButtonOptions{}).Clicked() {
		a.open = true
	}
	r := ToolApprovalDialog(c, &a.open, ToolApprovalOptions{Tool: "delete_files", Args: `{"glob":"*.tmp"}`, Reason: "Clean the build folder."})
	if r.Allowed {
		a.allowed++
	}
	if r.Denied {
		a.denied++
	}
}

func newApprovalApp() (*approvalApp, *ui.Tester) {
	a := &approvalApp{}
	tt := ui.NewTester(a.view, 600, 500)
	_ = tt.Click("Run")
	return a, tt
}

func TestToolApprovalDialogDenyFirst(t *testing.T) {
	a, tt := newApprovalApp()
	for _, s := range []string{"delete_files", "Clean the build folder.", `  "glob": "*.tmp"`} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	if !tt.Focused("Deny") {
		t.Fatal("Deny does not have the focus")
	}
	tt.Key(0, ui.KeyEnter)
	if a.denied != 1 || a.allowed != 0 || a.open {
		t.Fatalf("Enter: denied %d allowed %d open %v", a.denied, a.allowed, a.open)
	}
	if !tt.Focused("Run") {
		t.Fatal("focus not restored")
	}
}

func TestToolApprovalDialogAllowByTabAndBackdrop(t *testing.T) {
	a, tt := newApprovalApp()
	tt.ClickAt(5, 495)
	if !a.open || a.denied != 0 {
		t.Fatal("backdrop click acted")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Allow") {
		t.Fatal("Tab did not reach Allow")
	}
	tt.Key(0, ui.KeyEnter)
	if a.allowed != 1 || a.denied != 0 || a.open {
		t.Fatalf("allowed %d denied %d", a.allowed, a.denied)
	}
}

func TestToolApprovalDialogEscapeDenies(t *testing.T) {
	a, tt := newApprovalApp()
	tt.Key(0, ui.KeyEscape)
	if a.denied != 1 || a.allowed != 0 || a.open {
		t.Fatalf("Escape: denied %d allowed %d open %v", a.denied, a.allowed, a.open)
	}
}

func TestToolApprovalDialogNeedsTool(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := true
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolApprovalDialog(c, &open, ToolApprovalOptions{})
	}, 200, 200)
}

func TestToolApprovalDialogSnapshot(t *testing.T) {
	snapshot(t, "tool_approval_dialog", 560, 380, func(c *ui.Context) {
		open := true
		ToolApprovalDialog(c, &open, ToolApprovalOptions{Tool: "delete_files", Args: `{"glob":"build/*.tmp","dry_run":false}`,
			Reason: "The agent wants to clean temporary files from the build folder."})
	})
}
