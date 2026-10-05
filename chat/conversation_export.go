package chat

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ExportFormat is an export's text form; zero is ExportMarkdown.
type ExportFormat int

const (
	ExportMarkdown ExportFormat = iota
	ExportJSON
	ExportText
)

// ExportRange is which messages an export holds; zero is ExportAll.
type ExportRange int

const (
	ExportAll ExportRange = iota
	// ExportLast holds the last State.Last messages.
	ExportLast
	// ExportSelected holds the messages at Options.Selected.
	ExportSelected
)

var exportRoles = map[SessionRole]string{SessionUser: "user", SessionAssistant: "assistant", SessionSystem: "system"}

// ExportConversation writes messages as Markdown, JSON or plain text.
func ExportConversation(title string, messages []SessionMessage, format ExportFormat) string {
	var b strings.Builder
	switch format {
	case ExportMarkdown:
		fmt.Fprintf(&b, "# %s\n", title)
		for _, m := range messages {
			fmt.Fprintf(&b, "\n**%s**%s\n\n%s\n", exportAuthor(m), exportTime(m, " · "), m.Text)
		}
	case ExportJSON:
		type msgJSON struct {
			Role   string `json:"role"`
			Author string `json:"author,omitempty"`
			Text   string `json:"text"`
			Time   string `json:"time,omitempty"`
		}
		doc := struct {
			Title    string    `json:"title"`
			Messages []msgJSON `json:"messages"`
		}{Title: title, Messages: []msgJSON{}}
		for _, m := range messages {
			doc.Messages = append(doc.Messages, msgJSON{exportRoles[m.Role], m.Author, m.Text, exportTime(m, "")})
		}
		out, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			panic(err)
		}
		b.Write(out)
		b.WriteByte('\n')
	case ExportText:
		fmt.Fprintf(&b, "%s\n", title)
		for _, m := range messages {
			fmt.Fprintf(&b, "\n%s%s:\n%s\n", exportAuthor(m), exportTime(m, " "), m.Text)
		}
	default:
		panic("mujica: unknown ExportFormat")
	}
	return b.String()
}

func exportAuthor(m SessionMessage) string {
	if m.Author != "" {
		return m.Author
	}
	r, ok := exportRoles[m.Role]
	if !ok {
		panic("mujica: SessionMessage Role out of range")
	}
	return strings.ToUpper(r[:1]) + r[1:]
}

func exportTime(m SessionMessage, sep string) string {
	if m.Time.IsZero() {
		return ""
	}
	if sep == "" {
		return m.Time.Format(time.RFC3339)
	}
	return sep + "(" + m.Time.Format("2006-01-02 15:04") + ")"
}

// ConversationExportState is the chosen format and range.
type ConversationExportState struct {
	Format ExportFormat
	Range  ExportRange
	// Last is the count for ExportLast; zero is 10.
	Last  float64
	cache conversationExportCache
}

// conversationExportCache holds the last export text until its inputs change.
type conversationExportCache struct {
	key     conversationExportInputs
	content string
}

type conversationExportInputs struct {
	title  string
	format ExportFormat
	rng    ExportRange
	last   float64
	first  *SessionMessage
	n      int
	sel    string
	built  bool
}

// conversationExportContent rebuilds the export only when its inputs change.
func conversationExportContent(s *ConversationExportState, messages []SessionMessage, opts ConversationExportOptions) string {
	in := conversationExportInputs{title: opts.Title, format: s.Format, rng: s.Range, last: s.Last, n: len(messages), sel: fmt.Sprint(opts.Selected), built: true}
	if len(messages) > 0 {
		in.first = &messages[0]
	}
	if s.cache.key != in {
		s.cache = conversationExportCache{in, ExportConversation(opts.Title, conversationExportPick(s, messages, opts.Selected), s.Format)}
	}
	return s.cache.content
}

// ConversationExportOptions configure a ConversationExport. Title is required.
type ConversationExportOptions struct {
	Title string
	// Selected are message indices for ExportSelected.
	Selected []int
	Disabled bool
}

// ConversationExportView is a built ConversationExport and the export asked for.
type ConversationExportView struct {
	Element *ui.Element
	// Exported reports Content is ready for the caller to write.
	Exported bool
	Content  string
	Format   ExportFormat
}

// ConversationExport picks format and range, previews, and hands over text.
func ConversationExport(c *ui.Context, s *ConversationExportState, messages []SessionMessage, opts ConversationExportOptions) ConversationExportView {
	if s == nil || opts.Title == "" {
		panic("mujica: ConversationExport needs a state and a Title")
	}
	for _, i := range opts.Selected {
		if i < 0 || i >= len(messages) {
			panic("mujica: ConversationExport Selected index out of range")
		}
	}
	if s.Last == 0 {
		s.Last = 10
	}
	k := core.Tokens(c)
	d := core.Density(c)
	col := ui.Column(c).Gap(d.Unit() * 3).Disabled(opts.Disabled)
	v := ConversationExportView{Element: col, Format: s.Format}
	content := conversationExportContent(s, messages, opts)
	col.Children(func() {
		input.FormField(c, core.Msg(c, "conversation_export.format"), input.FormFieldOptions{}, func() *ui.Element {
			return input.RadioGroup(c, &s.Format, []input.RadioItem[ExportFormat]{
				{Value: ExportMarkdown, Label: "Markdown"}, {Value: ExportJSON, Label: "JSON"},
				{Value: ExportText, Label: core.Msg(c, "conversation_export.text")},
			}, input.RadioGroupOptions{Horizontal: true, Label: core.Msg(c, "conversation_export.format"), Disabled: opts.Disabled}).Element
		})
		input.FormField(c, core.Msg(c, "conversation_export.range"), input.FormFieldOptions{}, func() *ui.Element {
			return ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
				input.RadioGroup(c, &s.Range, []input.RadioItem[ExportRange]{
					{Value: ExportAll, Label: fmt.Sprintf(core.Msg(c, "conversation_export.all"), len(messages))},
					{Value: ExportLast, Label: core.Msg(c, "conversation_export.last")},
					{Value: ExportSelected, Label: fmt.Sprintf(core.Msg(c, "conversation_export.selected"), len(opts.Selected)), Disabled: len(opts.Selected) == 0},
				}, input.RadioGroupOptions{Horizontal: true, Label: core.Msg(c, "conversation_export.range"), Disabled: opts.Disabled})
				if s.Range == ExportLast {
					input.NumberInput(c, &s.Last, input.NumberInputOptions{Min: 1, Max: float64(max(1, len(messages))), Step: 1,
						Label: core.Msg(c, "conversation_export.last"), Disabled: opts.Disabled}).Element.Width(120)
				}
			})
		})
		ui.Text(c, core.Msg(c, "conversation_export.preview")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		ui.Scroll(c).MaxHeight(220).Padding(d.Unit()*3).Background(k.Background).Border(theme.BorderWidth, k.Border).
			Radius(theme.ControlRadius).Children(func() {
			ui.Text(c, content).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text).Selectable()
		})
		ui.Row(c).Justify(ui.End).Children(func() {
			v.Exported = input.Button(c, core.Msg(c, "conversation_export.export"), input.ButtonOptions{Icon: icons.Must("upload"), Disabled: opts.Disabled}).Clicked()
		})
	})
	if v.Exported && !col.IsDisabled() {
		v.Content = content
		slog.Debug("mujica: conversation_export export", "format", s.Format, "bytes", len(content))
	} else {
		v.Exported = false
	}
	return v
}

// conversationExportPick returns the messages the range holds, in order.
func conversationExportPick(s *ConversationExportState, messages []SessionMessage, selected []int) []SessionMessage {
	switch s.Range {
	case ExportLast:
		return messages[max(0, len(messages)-int(s.Last)):]
	case ExportSelected:
		var out []SessionMessage
		for i := range messages {
			for _, j := range selected {
				if i == j {
					out = append(out, messages[i])
					break
				}
			}
		}
		return out
	}
	return messages
}
