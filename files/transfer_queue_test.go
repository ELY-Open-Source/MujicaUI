package files

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func transferSample() []QueuedTransfer {
	return []QueuedTransfer{
		{ID: "a", Name: "ledger.zip", Upload: true, Status: TransferActive, Done: 30 << 20, Total: 120 << 20},
		{ID: "b", Name: "crest.png", Status: TransferPaused, Done: 1 << 20, Total: 4 << 20},
		{ID: "c", Name: "deed.pdf", Upload: true, Status: TransferFailed, Done: 0, Total: 2 << 20, Error: "network lost"},
		{ID: "d", Name: "hymn.mp3", Status: TransferDone, Done: 5 << 20, Total: 5 << 20},
	}
}

type transferProbe struct {
	items    []QueuedTransfer
	s        TransferQueueState
	disabled bool
	got      []string
}

func (p *transferProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := TransferQueue(c, p.items, &p.s, TransferQueueOptions{Disabled: p.disabled})
	r.Element.Height(380)
	if r.Action != TransferNone {
		p.got = append(p.got, fmt.Sprintf("%d:%s", r.Action, r.ID))
	}
}

func TestTransferQueueActions(t *testing.T) {
	p := &transferProbe{items: transferSample()}
	tt := ui.NewTester(p.view, 560, 400)
	for _, s := range []string{"1 in progress, 1 failed", "25% · 30.0 MB / 120.0 MB", "Paused at 25%", "Failed: network lost", "Done"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	for _, name := range []string{"Pause ledger.zip", "Resume crest.png", "Retry deed.pdf", "Remove hymn.mp3"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	want := "1:a 2:b 3:c 4:d"
	if strings.Join(p.got, " ") != want {
		t.Fatalf("actions %v, want %s", p.got, want)
	}
	if tt.HasText("Pause hymn.mp3") || tt.HasText("Retry ledger.zip") {
		t.Fatal("actions shown for the wrong status")
	}
}

func TestTransferQueueKeyboardAndDisabled(t *testing.T) {
	p := &transferProbe{items: transferSample()}
	tt := ui.NewTester(p.view, 560, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(p.got) != 1 {
		t.Fatalf("keyboard actions %v", p.got)
	}
	p = &transferProbe{items: transferSample(), disabled: true}
	tt = ui.NewTester(p.view, 560, 400)
	_ = tt.Click("Pause ledger.zip")
	if len(p.got) != 0 {
		t.Fatal("disabled queue acted")
	}
}

func TestTransferQueueEmptyAndInvalid(t *testing.T) {
	tt := ui.NewTester((&transferProbe{}).view, 560, 400)
	if !tt.HasText("No transfers") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Done above Total did not panic")
		}
	}()
	ui.NewTester((&transferProbe{items: []QueuedTransfer{{ID: "x", Done: 9, Total: 1}}}).view, 560, 400)
}

func TestTransferQueueSnapshot(t *testing.T) {
	snapshot(t, "transfer_queue", 520, 330, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *transferProbe { return &transferProbe{items: transferSample()} })
		ui.Column(c).Padding(12).Children(func() {
			TransferQueue(c, (*p).items, &(*p).s, TransferQueueOptions{}).Element.Height(306)
		})
	})
}

func TestTransferQueueLargeAndRevision(t *testing.T) {
	items := make([]QueuedTransfer, 100000)
	for i := range items {
		items[i] = QueuedTransfer{ID: fmt.Sprint(i), Name: fmt.Sprintf("f-%d.txt", i), Status: TransferActive, Total: 100}
	}
	s := TransferQueueState{}
	revision := uint64(0)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TransferQueue(c, items, &s, TransferQueueOptions{Revision: revision}).Element.Height(220)
	}, 500, 250)
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "f-") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d transfers", n)
	}
	items[0].Status = TransferFailed
	revision++
	tt.Frame()
	if !tt.HasText("99999 in progress, 1 failed") {
		t.Fatalf("revised summary %d active %d failed", s.active, s.failed)
	}
}
