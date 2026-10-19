package navigation

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestPaginationItems(t *testing.T) {
	for _, tc := range []struct {
		cur, pages int
		want       []int
	}{
		{1, 1, []int{1}},
		{3, 7, []int{1, 2, 3, 4, 5, 6, 7}},
		{1, 20, []int{1, 2, 3, 4, 5, 0, 20}},
		{10, 20, []int{1, 0, 9, 10, 11, 0, 20}},
		{20, 20, []int{1, 0, 16, 17, 18, 19, 20}},
	} {
		if got := paginationItems(tc.cur, tc.pages, 1); !slices.Equal(got, tc.want) {
			t.Errorf("items(%d, %d) = %v, want %v", tc.cur, tc.pages, got, tc.want)
		}
	}
}

func TestPaginationNavigate(t *testing.T) {
	page, size, changes := 1, 10, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Pagination(c, &page, &size, PaginationOptions{Total: 95}).Changed() {
			changes++
		}
	}, 700, 120)
	// Previous is disabled on page 1, so Tab lands on "1".
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if page != 2 || changes != 1 {
		t.Fatalf("keyboard: page %d changes %d", page, changes)
	}
	if err := tt.Click("Next page"); err != nil {
		t.Fatal(err)
	}
	if page != 3 || changes != 2 {
		t.Fatalf("next: page %d changes %d", page, changes)
	}
	if err := tt.Click("10"); err != nil {
		t.Fatal(err)
	}
	if page != 10 {
		t.Fatalf("click 10: page %d", page)
	}
	_ = tt.Click("Next page")
	if page != 10 {
		t.Fatalf("next past the end: page %d", page)
	}
}

func TestPaginationCorrectsOnTotalChange(t *testing.T) {
	page, size, total, changes := 9, 10, 95, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Pagination(c, &page, &size, PaginationOptions{Total: total}).Changed() {
			changes++
		}
	}, 700, 120)
	total = 31
	tt.Frame()
	if page != 4 || changes != 1 {
		t.Fatalf("page %d changes %d, want 4 and 1", page, changes)
	}
	total = 0
	tt.Frame()
	if page != 1 {
		t.Fatalf("empty total: page %d", page)
	}
}

func TestPaginationPageSizeKeepsRecord(t *testing.T) {
	page, size := 3, 10
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Pagination(c, &page, &size, PaginationOptions{Total: 200, PageSizes: []int{10, 20, 50}})
	}, 800, 300)
	if err := tt.Click("10 / page"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("50 / page"); err != nil {
		t.Fatal(err)
	}
	// Record 21 was first on page 3; with 50 per page it is on page 1.
	if size != 50 || page != 1 {
		t.Fatalf("size %d page %d", size, page)
	}
}

func TestPaginationJump(t *testing.T) {
	page, size := 1, 10
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		Pagination(c, &page, &size, PaginationOptions{Total: 100, Jump: true})
	}, 900, 120)
	if err := tt.Click("页码"); err != nil {
		t.Fatal(err)
	}
	tt.Type("42")
	tt.Key(0, ui.KeyEnter)
	if page != 1 || !tt.HasText("1–10") {
		t.Fatalf("out of range accepted: page %d texts %q", page, tt.Texts())
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	tt.Type("7")
	tt.Key(0, ui.KeyEnter)
	if page != 7 || tt.HasText("1–10") {
		t.Fatalf("jump: page %d texts %q", page, tt.Texts())
	}
}

func TestPaginationDisabled(t *testing.T) {
	page, size := 2, 10
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Pagination(c, &page, &size, PaginationOptions{Total: 50, Disabled: true})
	}, 700, 120)
	_ = tt.Click("Next page")
	_ = tt.Click("4")
	if page != 2 {
		t.Fatalf("disabled pagination moved to %d", page)
	}
}

func TestPaginationInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		size int
		opts PaginationOptions
	}{
		"total": {10, PaginationOptions{Total: -1}},
		"size":  {0, PaginationOptions{}},
		"sizes": {15, PaginationOptions{PageSizes: []int{10, 20}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			page := 1
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Pagination(c, &page, &tc.size, tc.opts)
			}, 300, 100, 1)
		})
	}
}

func TestPaginationSnapshot(t *testing.T) {
	snapshot(t, "pagination", 760, 140, func(c *ui.Context) {
		a, b, s := 10, 1, 20
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			Pagination(c, &a, &s, PaginationOptions{Total: 400, PageSizes: []int{10, 20, 50}, Jump: true})
			Pagination(c, &b, &s, PaginationOptions{Total: 45, Disabled: true})
		})
	})
}
