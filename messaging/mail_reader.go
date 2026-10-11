package messaging

import (
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MailFormat selects plain text or safe Markdown.
type MailFormat int

const (
	MailPlain MailFormat = iota
	// MailMarkdown uses MarkdownView’s safe subset.
	MailMarkdown
)

// MailAttachment is a file sent with a mail.
type MailAttachment struct {
	Name string
	Size int64
	Type string
}

// Mail supplies headers, body and attachments.
type Mail struct {
	From        string
	To, Cc      []string
	Subject     string
	Date        time.Time
	Body        string
	Format      MailFormat
	Attachments []MailAttachment
}

// MailReaderResult is a MailReader's element and this frame's events.
type MailReaderResult struct {
	Element *ui.Element
	// Link contains the activated URL, or an empty string.
	Link string
	// Opened and Downloaded are attachment indexes, -1 for none.
	Opened, Downloaded int
}

type mailQuotesKey struct{}

// MailReader displays headers, body, folded quotes and attachments.
func MailReader(c *ui.Context, m Mail) MailReaderResult {
	if m.From == "" || m.Subject == "" {
		panic("mujica: MailReader needs From and Subject")
	}
	if m.Format != MailPlain && m.Format != MailMarkdown {
		panic("mujica: unknown MailFormat")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	r := MailReaderResult{Opened: -1, Downloaded: -1}
	col := ui.Column(c).Gap(u * 4).MinWidth(0).Role(ui.RoleGroup).Label(m.Subject)
	r.Element = col
	small := core.FontSize(c, theme.CaptionSize)
	col.Children(func() {
		ui.Text(c, m.Subject).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).FontWeight(600)
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Children(func() {
			display.Avatar(c, m.From, display.AvatarOptions{Size: 36})
			ui.Column(c).Gap(2).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, m.From).FontWeight(600).Selectable().FocusRing(false)
				for _, line := range [][2]string{{"mail_reader.to", strings.Join(m.To, ", ")}, {"mail_reader.cc", strings.Join(m.Cc, ", ")}} {
					if line[1] != "" {
						ui.Text(c, core.Msg(c, line[0])+" "+line[1]).FontSize(small).TextColor(k.TextMuted).Selectable().FocusRing(false)
					}
				}
			})
			if !m.Date.IsZero() {
				ui.Text(c, core.MessagingDate(c, m.Date, "mail_reader.date")).FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
			}
		})
		ui.Box(c).Height(1).Background(k.Border)
		if m.Format == MailMarkdown {
			if url, ok := chat.MarkdownView(c, m.Body).Link(); ok {
				r.Link = url
			}
		} else {
			mailPlainBody(c, ui.Local(col, mailQuotesKey{}, func() map[string]bool { return map[string]bool{} }), m.Body)
		}
		if len(m.Attachments) > 0 {
			ui.Text(c, core.Msg(c, "mail_reader.attachments")).FontSize(small).FontWeight(600).TextColor(k.TextMuted)
			ui.Row(c).Gap(u * 2).Wrap().Children(func() {
				for i, a := range m.Attachments {
					f := chat.FileMessage(c, chat.FileMessageOptions(a))
					f.Element.Key(i)
					if f.Opened() {
						r.Opened = i
					}
					if f.Downloaded() {
						r.Downloaded = i
					}
				}
			})
		}
	})
	if r.Link != "" || r.Opened >= 0 || r.Downloaded >= 0 {
		slog.Debug("mujica: mail_reader action", "opened", r.Opened, "downloaded", r.Downloaded)
	}
	return r
}

// mailPlainBody shows plain text, folding each run of quoted lines.
func mailPlainBody(c *ui.Context, open *map[string]bool, body string) {
	lines := strings.Split(body, "\n")
	ui.Column(c).Gap(core.Density(c).Unit() * 2).MinWidth(0).Children(func() {
		for i := 0; i < len(lines); {
			j := i
			quoted := strings.HasPrefix(lines[i], ">")
			for j < len(lines) && strings.HasPrefix(lines[j], ">") == quoted {
				j++
			}
			chunk := strings.Join(lines[i:j], "\n")
			ui.Box(c).Key(i).Children(func() {
				if quoted {
					// Fold state is kept per quote text.
					o := (*open)[chunk]
					QuotedText(c, chunk, &o)
					(*open)[chunk] = o
				} else if strings.TrimSpace(chunk) != "" {
					ui.Text(c, strings.Trim(chunk, "\n")).Selectable().FocusRing(false).LineHeight(theme.LineHeight)
				}
			})
			i = j
		}
	})
}
