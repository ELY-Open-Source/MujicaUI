package agent

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func sandboxSample(expires time.Time) SandboxInfo {
	return SandboxInfo{Name: "sandbox-7", State: SandboxRunning, Expires: expires, Resources: []SandboxResource{
		{Name: "CPU", Used: 1.5, Limit: 2, Unit: "cores"},
		{Name: "Memory", Used: 3.8, Limit: 4, Unit: "GB"},
	}}
}

func TestSandboxStatusResourcesAndCountdown(t *testing.T) {
	info := sandboxSample(time.Now().Add(90 * time.Second))
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SandboxStatus(c, info, SandboxStatusOptions{Still: true})
	}, 480, 300)
	for _, s := range []string{"sandbox-7", "Running", "1.5 / 2 cores", "3.8 / 4 GB", "Near the limit"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	left := ""
	for _, s := range tt.Texts() {
		if strings.HasSuffix(s, " left") {
			left = s
		}
	}
	if !strings.HasPrefix(left, "1 m ") {
		t.Fatalf("countdown %q", left)
	}
	info.Expires = time.Now().Add(-time.Second)
	tt.Frame()
	if !tt.HasText("Time is up") {
		t.Fatalf("expired: %q", tt.Texts())
	}
	info.State = SandboxStopped
	tt.Frame()
	if !tt.HasText("Stopped") || tt.HasText("Time is up") {
		t.Fatalf("stopped: %q", tt.Texts())
	}
}

func TestSandboxStatusRejectsBadResource(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SandboxStatus(c, SandboxInfo{Name: "x", Resources: []SandboxResource{{Name: "CPU", Used: 3, Limit: 2}}}, SandboxStatusOptions{})
	}, 200, 200)
}

func TestSandboxStatusSnapshot(t *testing.T) {
	snapshot(t, "sandbox_status", 440, 230, func(c *ui.Context) {
		info := sandboxSample(time.Time{})
		info.Resources = append(info.Resources, SandboxResource{Name: "Disk", Used: 12, Limit: 40, Unit: "GB"})
		ui.Column(c).Padding(16).Children(func() { SandboxStatus(c, info, SandboxStatusOptions{Still: true}) })
	})
}
