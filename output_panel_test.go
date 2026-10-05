package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

type outputProbe struct {
	channels []OutputChannel
	s        OutputPanelState
	disabled bool
	changed  int
	cleared  []string
}

func (p *outputProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := OutputPanel(c, p.channels, &p.s, OutputPanelOptions{Disabled: p.disabled})
	r.Element.Height(300)
	if r.Changed {
		p.changed++
	}
	if r.Cleared != "" {
		p.cleared = append(p.cleared, r.Cleared)
		for i := range p.channels {
			if p.channels[i].Name == r.Cleared {
				p.channels[i].Lines = nil
			}
		}
	}
}

func newOutputProbe() *outputProbe {
	var build []string
	for i := range 200 {
		build = append(build, fmt.Sprintf("compile unit %d", i))
	}
	return &outputProbe{channels: []OutputChannel{{"Build", build}, {"Tests", []string{"\x1b[32mPASS\x1b[0m court", "\x1b[31mFAIL\x1b[0m ledger"}}}}
}

func TestOutputPanelFollowsTail(t *testing.T) {
	p := newOutputProbe()
	tt := ui.NewTester(p.view, 520, 320)
	if !tt.HasText("compile unit 199") || tt.HasText("compile unit 0") {
		t.Fatal("panel did not start at the tail")
	}
	p.channels[0].Lines = append(p.channels[0].Lines, "linked")
	tt.Frame()
	if !tt.HasText("linked") {
		t.Fatal("new output not followed")
	}
}

func TestOutputPanelSwitchAndClear(t *testing.T) {
	p := newOutputProbe()
	tt := ui.NewTester(p.view, 520, 320)
	if err := tt.Click("Build"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Tests"); err != nil {
		t.Fatal(err)
	}
	if p.s.Channel != 1 || p.changed != 1 || !tt.HasText("PASS court") {
		t.Fatalf("channel %d changed %d texts %q", p.s.Channel, p.changed, tt.Texts())
	}
	if err := tt.Click("Clear output"); err != nil {
		t.Fatal(err)
	}
	if len(p.cleared) != 1 || p.cleared[0] != "Tests" || !tt.HasText("No output") {
		t.Fatalf("cleared %v", p.cleared)
	}
}

func TestOutputPanelAnsiColors(t *testing.T) {
	var spans []ui.Span
	var ok ui.Color
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ok = Tokens(c).Success
		spans = codeAnsiSpans(c, "\x1b[1;32mPASS\x1b[0m court")
	}, 10, 10)
	if len(spans) != 2 || spans[0].Color != ok || spans[0].Weight != 700 || spans[1].Color.A != 0 {
		t.Fatalf("spans %+v", spans)
	}
}

func TestOutputPanelDisabledAndBadChannel(t *testing.T) {
	p := newOutputProbe()
	p.disabled = true
	tt := ui.NewTester(p.view, 520, 320)
	_ = tt.Click("Clear output")
	if len(p.cleared) != 0 {
		t.Fatal("disabled panel cleared")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("channel out of range did not panic")
		}
	}()
	p = newOutputProbe()
	p.s.Channel = 5
	ui.NewTester(p.view, 520, 320)
}

func TestOutputPanelSnapshot(t *testing.T) {
	snapshot(t, "output_panel", 520, 200, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *outputProbe {
			p := newOutputProbe()
			p.s.Channel = 1
			return p
		})
		ui.Column(c).Padding(12).Children(func() {
			OutputPanel(c, (*p).channels, &(*p).s, OutputPanelOptions{}).Element.Height(176)
		})
	})
}
