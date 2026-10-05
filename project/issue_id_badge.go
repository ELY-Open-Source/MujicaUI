package project

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

var issueIDPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]*-[1-9][0-9]*$`)

// IssueIdBadgeOptions configure an IssueIdBadge.
type IssueIdBadgeOptions struct {
	// Copyable makes a click or Enter copy the key.
	Copyable bool
	Disabled bool
}

// issueCopiedFor is how long "Copied" shows after a copy.
const issueCopiedFor = 1500 * time.Millisecond

// IssueIdBadge displays a monospace issue key with optional copying.
func IssueIdBadge(c *ui.Context, id string, opts IssueIdBadgeOptions) core.ChoiceResult {
	if !issueIDPattern.MatchString(id) {
		panic(fmt.Sprintf("mujica: IssueIdBadge %q is not a key like ENG-123", id))
	}
	k := core.Tokens(c)
	d := core.Density(c)
	small := core.FontSize(c, theme.CaptionSize)
	if !opts.Copyable {
		e := ui.Row(c).Padding(d.Unit()/2, d.Unit()*1.5).Radius(theme.ControlRadius).Background(k.Background).
			Border(theme.BorderWidth, k.Border).Shrink(0).Label(id)
		e.Children(func() { ui.Text(c, id).Font(theme.MonoFont).FontSize(small).TextColor(k.TextMuted).SingleLine() })
		return core.ChoiceResult{Element: e}
	}
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Gap(d.Unit()).PaddingX(d.Unit()*1.5).Radius(theme.ControlRadius).
		Background(k.Background).Border(theme.BorderWidth, k.Border).Disabled(opts.Disabled)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	copied := ui.Local(b, "issue-copied", func() time.Time { return time.Time{} })
	changed := false
	if b.Clicked() && !b.IsDisabled() {
		c.WriteClipboard(id)
		c.Announce(fmt.Sprintf(core.Msg(c, "issue_id_badge.copied_label"), id))
		*copied, changed = c.Now(), true
		slog.Debug("mujica: issue_id_badge copy", "id", id)
	}
	showing := !copied.IsZero() && c.Now().Sub(*copied) < issueCopiedFor
	if showing {
		c.After(issueCopiedFor - c.Now().Sub(*copied))
	}
	b.Label(fmt.Sprintf(core.Msg(c, "issue_id_badge.copy"), id))
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Transition(core.ColorFade).Children(func() {
		ui.Text(c, id).Font(theme.MonoFont).FontSize(small).TextColor(k.Text).SingleLine()
		icon, col := "copy", k.TextMuted
		if showing {
			icon, col = "check", k.Success
		}
		ui.Icon(c, icons.Must(icon)).Size(12, 12).TextColor(col)
		if showing {
			ui.Text(c, core.Msg(c, "issue_id_badge.copied")).FontSize(small).TextColor(k.Success)
		}
	})
	return core.NewChoiceResult(b, changed)
}
