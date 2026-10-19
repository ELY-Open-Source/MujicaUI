package agent

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var mcpServersSample = []MCPServer{
	{ID: "fs", Name: "filesystem", Transport: "stdio", State: MCPConnected, Tools: 12},
	{ID: "gh", Name: "github", Transport: "http", State: MCPError, Error: "401 Unauthorized"},
	{ID: "db", Name: "postgres", Transport: "stdio", State: MCPStopped},
	{ID: "web", Name: "browser", Transport: "sse", State: MCPConnecting},
}

func mcpTester(servers []MCPServer, opts MCPServerListOptions, got *[]MCPServerListResult) *ui.Tester {
	opts.Still = true
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := MCPServerList(c, servers, opts)
		if r.Action != MCPNone {
			*got = append(*got, r)
		}
	}, 600, 460)
}

func TestMCPServerListStatesAndActions(t *testing.T) {
	var got []MCPServerListResult
	tt := mcpTester(mcpServersSample, MCPServerListOptions{}, &got)
	for _, s := range []string{"Connected · 12 tools", "Connection failed", "401 Unauthorized", "Stopped", "Connecting"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Reconnect")
	_ = tt.Click("Start")
	if len(got) != 2 || got[0].Action != MCPReconnect || got[0].Server != "gh" || got[1].Action != MCPStart || got[1].Server != "db" {
		t.Fatalf("actions %+v", got)
	}
	tt = mcpTester(mcpServersSample, MCPServerListOptions{}, &got)
	tt.Key(0, ui.KeyTab) // filesystem's Stop
	tt.Key(0, ui.KeyEnter)
	if len(got) != 3 || got[2].Action != MCPStop || got[2].Server != "fs" {
		t.Fatalf("keyboard stop: %+v", got)
	}
}

func TestMCPServerListDisabledAndEmpty(t *testing.T) {
	var got []MCPServerListResult
	tt := mcpTester(mcpServersSample, MCPServerListOptions{Disabled: true}, &got)
	_ = tt.Click("Reconnect")
	_ = tt.Click("Start")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 0 {
		t.Fatalf("disabled list acted: %+v", got)
	}
	tt = mcpTester(nil, MCPServerListOptions{}, &got)
	if !tt.HasText("No MCP servers") {
		t.Fatalf("empty: %q", tt.Texts())
	}
}

func TestMCPServerListManyServersVirtualized(t *testing.T) {
	servers := make([]MCPServer, 100000)
	for i := range servers {
		servers[i] = MCPServer{ID: fmt.Sprint(i), Name: fmt.Sprintf("server-%d", i), State: MCPConnected, Tools: i % 30}
	}
	var got []MCPServerListResult
	tt := mcpTester(servers, MCPServerListOptions{}, &got)
	if n := len(tt.Texts()); n > 300 {
		t.Fatalf("%d texts built for 100000 servers", n)
	}
}

func TestMCPServerListRejectsDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	var got []MCPServerListResult
	mcpTester([]MCPServer{{ID: "a", Name: "x"}, {ID: "a", Name: "y"}}, MCPServerListOptions{}, &got)
}

func TestMCPServerListSnapshot(t *testing.T) {
	snapshot(t, "mcp_server_list", 560, 330, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { MCPServerList(c, mcpServersSample, MCPServerListOptions{Still: true}) })
	})
}
