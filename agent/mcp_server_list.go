package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// MCPState is a server connection; zero is stopped.
type MCPState int

const (
	MCPStopped MCPState = iota
	MCPConnecting
	MCPConnected
	MCPError
)

// MCPServer is one configured MCP server; the caller owns it.
type MCPServer struct {
	ID        string
	Name      string
	Transport string
	State     MCPState
	Tools     int
	// Error explains MCPError.
	Error string
}

// MCPAction is a list request; zero is none.
type MCPAction int

const (
	MCPNone MCPAction = iota
	MCPStart
	MCPStop
	MCPReconnect
)

// MCPServerListOptions configure an MCPServerList.
type MCPServerListOptions struct {
	// MaxHeight caps the scrolling list; zero means 360.
	MaxHeight float32
	Disabled  bool
	Still     bool
}

// MCPServerListResult reports this frame's request.
type MCPServerListResult struct {
	Element ui.Element
	Action  MCPAction
	// Server is the ID the action is for.
	Server string
}

// MCPServerList lists MCP servers with start, stop and reconnect.
func MCPServerList(c *ui.Context, servers []MCPServer, opts MCPServerListOptions) MCPServerListResult {
	seen := map[string]bool{}
	for i, s := range servers {
		if s.ID == "" || s.Name == "" || seen[s.ID] || s.Tools < 0 || s.State < MCPStopped || s.State > MCPError {
			panic(fmt.Sprintf("mujica: MCPServerList server %d needs a unique ID, a Name, a known State and Tools >= 0", i))
		}
		seen[s.ID] = true
	}
	if opts.MaxHeight < 0 {
		panic("mujica: MCPServerList MaxHeight is negative")
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	var r MCPServerListResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled).Label(core.Msg(c, "mcp_server_list.label"))
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		if len(servers) == 0 {
			AgentCaption(c, core.Msg(c, "mcp_server_list.empty"))
			return
		}
		k := core.Tokens(c)
		AgentRows(c, len(servers), func(i int) {
			if a := mcpServerRow(c, servers[i], locked, opts.Still); a != MCPNone && !locked {
				r.Action, r.Server = a, servers[i].ID
			}
		}).Shrink(1).MaxHeight(maxH).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
			Dividers(theme.BorderWidth, k.Border)
	})
	if r.Action != MCPNone {
		slog.Debug("mujica: mcp_server_list action", "server", r.Server, "action", int(r.Action))
	}
	return r
}

// mcpServerRow builds one server, returning its pressed action.
func mcpServerRow(c *ui.Context, s MCPServer, locked, still bool) MCPAction {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	words := [...]string{"stopped", "connecting", "connected", "error"}
	word := core.Msg(c, "mcp_server_list."+words[s.State])
	act := MCPNone
	ui.Row(c.Key(s.ID)).Gap(u*3).AlignItems(ui.Center).MinHeight(d.ControlHeight()+u*2).Padding(u*2, u*3).
		Label(s.Name + ", " + word).Children(func() {
		switch s.State {
		case MCPConnecting:
			agentMark(c, AgentRunning, 16, still)
		case MCPConnected:
			agentMark(c, AgentDone, 16, still)
		case MCPError:
			agentMark(c, AgentFailed, 16, still)
		default:
			ui.Icon(c, icons.Must("wifi-off")).Size(16, 16).TextColor(k.TextMuted)
		}
		ui.Column(c).Grow(1).Shrink(1).Gap(u / 2).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, s.Name).FontWeight(600).SingleLine()
				AgentCaption(c, s.Transport).Font(theme.MonoFont)
			})
			line := word
			if s.State == MCPConnected {
				line += " · " + fmt.Sprintf(core.Msg(c, "mcp_server_list.tools"), s.Tools)
			}
			AgentCaption(c, line)
			if s.State == MCPError && s.Error != "" {
				ui.Text(c, s.Error).TextColor(k.Danger).FontSize(core.FontSize(c, theme.CaptionSize)).Selectable()
			}
		})
		if s.State == MCPError && input.Button(c, core.Msg(c, "mcp_server_list.reconnect"), input.ButtonOptions{Variant: input.Outline, Disabled: locked}).Clicked() {
			act = MCPReconnect
		}
		if s.State == MCPStopped {
			if input.Button(c, core.Msg(c, "mcp_server_list.start"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() {
				act = MCPStart
			}
		} else if input.Button(c, core.Msg(c, "mcp_server_list.stop"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() {
			act = MCPStop
		}
	})
	return act
}
