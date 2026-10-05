package feedback

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestResultActions(t *testing.T) {
	got := []int{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if a := Result(c, ResultOptions{Kind: ResultFailure, Description: "The payment was declined.", Actions: []string{"Try again", "Back to cart"}}).Action; a >= 0 {
			got = append(got, a)
		}
	}, 600, 500)
	if !tt.HasText("Something went wrong") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Back to cart"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 2 || got[0] != 1 || got[1] != 0 {
		t.Fatalf("actions %v", got)
	}
}

func TestResultTitles(t *testing.T) {
	for kind, want := range map[ResultKind]string{ResultSuccess: "操作成功", ResultFailure: "操作失败", ResultForbidden: "无访问权限"} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Locale: core.ZhCN})
			Result(c, ResultOptions{Kind: kind})
		}, 500, 400)
		if !tt.HasText(want) {
			t.Errorf("%d: texts %q", kind, tt.Texts())
		}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Result(c, ResultOptions{Title: "Order placed"})
	}, 500, 400)
	if !tt.HasText("Order placed") || tt.HasText("Done") {
		t.Fatalf("custom title: %q", tt.Texts())
	}
}

func TestResultPanicsOnUnknownKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); Result(c, ResultOptions{Kind: 7}) }, 300, 200)
}

func TestResultSnapshot(t *testing.T) {
	snapshot(t, "result", 520, 340, func(c *ui.Context) {
		Result(c, ResultOptions{
			Kind:        ResultForbidden,
			Description: "You need the Editor role to change billing settings.",
			Actions:     []string{"Request access", "Go back"},
		})
	})
}
