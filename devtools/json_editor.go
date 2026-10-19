package devtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// JsonEditorOptions configure a JsonEditor.
type JsonEditorOptions struct {
	// Label names the editor; empty uses "JSON".
	Label string
	// MinHeight and MaxHeight default to 160 and 400 DIPs.
	MinHeight, MaxHeight float32
	Disabled             bool
}

// JsonEditorResult reports this frame's events and the parse state.
type JsonEditorResult struct {
	Element ui.Element
	Changed bool
	// Valid reports *text parses as one JSON value.
	Valid bool
	// Line and Column locate errors; valid documents use zero.
	Line, Column int
}

type jsonEditorCheck struct {
	src          string
	err          string
	line, column int
}

// JsonEditor edits and formats JSON with error locations.
func JsonEditor(c *ui.Context, text *string, opts JsonEditorOptions) JsonEditorResult {
	if text == nil || opts.MinHeight < 0 || opts.MaxHeight < 0 {
		panic("mujica: JsonEditor needs text and non-negative heights")
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "json_editor.label")
	}
	minH, maxH := opts.MinHeight, opts.MaxHeight
	if minH == 0 {
		minH = 160
	}
	if maxH == 0 {
		maxH = 400
	}
	if minH > maxH {
		panic("mujica: JsonEditor MinHeight exceeds MaxHeight")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r JsonEditorResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 2).Disabled(opts.Disabled)
	chk := ui.Local(r.Element, "mujica.jsoneditor", func() jsonEditorCheck { return jsonEditorCheck{src: "\x00"} })
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		jsonEditorValidate(chk, *text)
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).FontWeight(600).Grow(1)
			if input.Button(c, core.Msg(c, "json_editor.format"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked || chk.err != ""}).Clicked() && !locked {
				var b bytes.Buffer
				if json.Indent(&b, []byte(*text), "", "  ") == nil && b.String() != *text {
					*text, r.Changed = b.String(), true
					slog.Debug("mujica: json_editor format")
				}
			}
		})
		in := input.TextArea(c, text, input.TextAreaOptions{Label: label, MinHeight: minH, MaxHeight: maxH, Disabled: locked, Error: chk.err})
		in.Input.Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize))
		r.Changed = r.Changed || in.Changed()
		jsonEditorValidate(chk, *text)
		if chk.err != "" {
			lines := strings.Split(*text, "\n")
			if chk.line >= 1 && chk.line <= len(lines) {
				jsonEditorExcerpt(c, lines[chk.line-1], chk.line, chk.column)
			}
			return
		}
		ui.Row(c).Gap(u).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
			ui.Icon(c, icons.Must("circle-check")).Size(14, 14).TextColor(k.Success)
			ui.Text(c, core.Msg(c, "json_editor.valid")).TextColor(k.Success).FontSize(core.FontSize(c, theme.CaptionSize+1))
		})
	})
	r.Valid, r.Line, r.Column = chk.err == "", chk.line, chk.column
	return r
}

// jsonEditorValidate parses text once per change and locates any error.
func jsonEditorValidate(chk *jsonEditorCheck, text string) {
	if chk.src == text {
		return
	}
	chk.src, chk.err, chk.line, chk.column = text, "", 0, 0
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	var v any
	err := dec.Decode(&v)
	var offset int64
	if err == nil {
		rest := text[dec.InputOffset():]
		offset = dec.InputOffset() + int64(len(rest)-len(strings.TrimLeft(rest, " \t\r\n"))) + 1
		var trailing any
		if next := dec.Decode(&trailing); next != io.EOF {
			err = errors.New("unexpected data after the value")
			if next != nil {
				err = next
			}
		}
	}
	if err == nil {
		return
	}
	var syn *json.SyntaxError
	var typ *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syn):
		offset = syn.Offset
	case errors.As(err, &typ):
		offset = typ.Offset
	case strings.TrimSpace(text) == "":
		err, offset = errors.New("empty document"), 0
	default:
		if offset == 0 {
			offset = int64(len(text))
		}
	}
	before := text[:min(int(offset), len(text))]
	chk.line = strings.Count(before, "\n") + 1
	chk.column = len([]rune(before[strings.LastIndex(before, "\n")+1:]))
	chk.column = max(chk.column, 1)
	chk.err = fmt.Sprintf("%d:%d %v", chk.line, chk.column, err)
}

// jsonEditorExcerpt shows an error line and column caret.
func jsonEditorExcerpt(c *ui.Context, line string, n, col int) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ui.Column(c).Padding(u*2, u*3).Radius(theme.ControlRadius).Background(k.Background).Border(theme.BorderWidth, k.Border).
		Label(fmt.Sprintf(core.Msg(c, "json_editor.at"), n, col)).Children(func() {
		spans := []ui.Span{{Text: fmt.Sprintf("%4d | ", n), Color: k.TextMuted}}
		toks, _ := internal.Highlight("json", line, false)
		for _, t := range toks {
			spans = append(spans, ui.Span{Text: t.Text, Color: code.CodeTokenColor(k, t.Kind)})
		}
		ui.RichText(c, spans...).Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize)).SingleLine()
		code.CodeMono(c, strings.Repeat(" ", 7+col-1)+"^").TextColor(k.Danger)
	})
}
