package datetime

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CronSchedule is a parsed five-field cron expression.
type CronSchedule struct {
	fields [5]uint64
	tokens [5]string
	// domAny and dowAny note a "*" day field.
	domAny, dowAny bool
}

// CronProblem is what is wrong with a cron field.
type CronProblem int

const (
	CronFieldCount CronProblem = iota
	CronEmptyItem
	CronBadStep
	CronNotNumber
	CronOutOfRange
	CronBackwards
)

// CronError locates a problem in a cron expression.
type CronError struct {
	// Field is 0 to 4, or -1 for the count.
	Field   int
	Problem CronProblem
	// Text is the offending part.
	Text string
}

var cronProblemKeys = [...]string{"cron_editor.err_count", "cron_editor.err_empty", "cron_editor.err_step", "cron_editor.err_number",
	"cron_editor.err_range", "cron_editor.err_backwards"}

func (e *CronError) Error() string {
	return fmt.Sprintf("cron field %d: problem %d at %q", e.Field+1, e.Problem, e.Text)
}

var cronBounds = [5][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 7}}

var cronNames = [5][]string{3: {"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"},
	4: {"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}}

// ParseCron parses standard five-field cron with names.
func ParseCron(expr string) (CronSchedule, error) {
	parts := strings.Fields(expr)
	var s CronSchedule
	if len(parts) != 5 {
		return s, &CronError{-1, CronFieldCount, strconv.Itoa(len(parts))}
	}
	for f, part := range parts {
		s.tokens[f] = part
		for _, item := range strings.Split(part, ",") {
			bits, err := cronItem(f, strings.ToUpper(item))
			if err != nil {
				err.Field = f
				return CronSchedule{}, err
			}
			s.fields[f] |= bits
		}
	}
	if s.fields[4]&(1<<7) != 0 {
		s.fields[4] |= 1
	}
	s.domAny, s.dowAny = parts[2] == "*", parts[4] == "*"
	return s, nil
}

// cronItem parses one list item into a bit set.
func cronItem(f int, item string) (uint64, *CronError) {
	lo, hi := cronBounds[f][0], cronBounds[f][1]
	if item == "" {
		return 0, &CronError{Problem: CronEmptyItem}
	}
	rng, step := item, 1
	if i := strings.IndexByte(item, '/'); i >= 0 {
		n, err := strconv.Atoi(item[i+1:])
		if err != nil || n < 1 {
			return 0, &CronError{Problem: CronBadStep, Text: item[i+1:]}
		}
		rng, step = item[:i], n
	}
	a, b := lo, hi
	switch {
	case rng == "*":
	case strings.Contains(rng, "-"):
		x, y, _ := strings.Cut(rng, "-")
		var e *CronError
		if a, e = cronValue(f, x); e != nil {
			return 0, e
		}
		if b, e = cronValue(f, y); e != nil {
			return 0, e
		}
		if a > b {
			return 0, &CronError{Problem: CronBackwards, Text: rng}
		}
	default:
		var e *CronError
		if a, e = cronValue(f, rng); e != nil {
			return 0, e
		}
		if step == 1 {
			b = a
		}
	}
	var bits uint64
	for v := a; v <= b; v += step {
		bits |= 1 << v
	}
	return bits, nil
}

func cronValue(f int, s string) (int, *CronError) {
	for i, n := range cronNames[f] {
		if s == n {
			return i + cronBounds[f][0], nil
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, &CronError{Problem: CronNotNumber, Text: s}
	}
	if v < cronBounds[f][0] || v > cronBounds[f][1] {
		return 0, &CronError{Problem: CronOutOfRange, Text: s}
	}
	return v, nil
}

// CronRun is one upcoming run; Repeated marks a DST repeat.
type CronRun struct {
	At       time.Time
	Repeated bool
}

// Next returns up to n runs within five years.
func (s CronSchedule) Next(t time.Time, loc *time.Location, n int) []CronRun {
	var out []CronRun
	w := t.In(loc)
	day := core.DateOf(w)
	for i := 0; i < 5*366 && len(out) < n; i, day = i+1, day.AddDays(1) {
		if !s.dayMatches(day) {
			continue
		}
		for h := 0; h < 24 && len(out) < n; h++ {
			if s.fields[1]&(1<<h) == 0 {
				continue
			}
			for m := 0; m < 60 && len(out) < n; m++ {
				if s.fields[0]&(1<<m) == 0 {
					continue
				}
				at, k := dateTimeResolve(day, TimeOfDay{h, m, 0}, loc)
				if k == 0 || !at.After(t) {
					continue
				}
				out = append(out, CronRun{at, k > 1})
			}
		}
	}
	return out
}

// dayMatches applies cron's either-day-field rule.
func (s CronSchedule) dayMatches(d core.Date) bool {
	if s.fields[3]&(1<<int(d.Month)) == 0 {
		return false
	}
	dom := s.fields[2]&(1<<d.Day) != 0
	dow := s.fields[4]&(1<<int(d.Weekday())) != 0
	if !s.domAny && !s.dowAny {
		return dom || dow
	}
	return dom && dow
}

// CronEditorOptions configure a CronEditor.
type CronEditorOptions struct {
	// Location is the zone runs are computed in; required.
	Location *time.Location
	// Now starts the preview; zero uses the system clock.
	Now time.Time
	// Runs is how many next runs show; zero means 5.
	Runs     int
	Hour12   bool
	Disabled bool
	Label    string
}

// CronEditor edits, explains and previews a cron expression.
func CronEditor(c *ui.Context, expr *string, opts CronEditorOptions) core.ChoiceResult {
	loc := eventLocation("CronEditor", opts.Location)
	if opts.Runs < 0 {
		panic("mujica: CronEditor Runs is negative")
	}
	runs := opts.Runs
	if runs == 0 {
		runs = 5
	}
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "cron_editor.label")
	}
	old := *expr
	col := ui.Column(c).Gap(d.Unit() * 3).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	cache := ui.Local(col, "cron-parse", func() *cronParsed { return &cronParsed{} })
	col.Children(func() {
		cp := *cache
		if !cp.ok || cp.expr != *expr || !cp.from.Equal(now.Truncate(time.Minute)) {
			cp.s, cp.err = ParseCron(*expr)
			cp.expr, cp.ok, cp.from = *expr, true, now.Truncate(time.Minute)
			if cp.err == nil {
				cp.runs = cp.s.Next(now, loc, runs)
			}
		}
		problem := ""
		if cp.err != nil {
			problem = cronProblem(c, cp.err.(*CronError))
		}
		in := input.TextInput(c, expr, input.TextInputOptions{Label: core.Msg(c, "cron_editor.expression"), Placeholder: "*/15 9-17 * * MON-FRI", Disabled: opts.Disabled, Error: problem})
		in.Input.Font(theme.MonoFont)
		bad := -2
		if cp.err != nil {
			bad = cp.err.(*CronError).Field
		}
		fields := strings.Fields(*expr)
		ui.Row(c).Gap(d.Unit() * 2).Wrap().Children(func() {
			for f := range 5 {
				tok := ""
				if f < len(fields) {
					tok = fields[f]
				}
				cronFieldCard(c, f, tok, f == bad || bad == -1)
			}
		})
		if cp.err != nil {
			return
		}
		ui.Text(c, core.Msg(c, "cron_editor.next")).Font(theme.SerifFont).TextColor(k.Text)
		ui.Column(c).Gap(d.Unit()).Role(ui.RoleList).Children(func() {
			list := cp.runs
			if len(list) == 0 {
				ui.Text(c, core.Msg(c, "cron_editor.never")).TextColor(k.TextMuted)
			}
			for _, run := range list {
				w := run.At.In(loc)
				text := core.Msg(c, "calendar.weekday_"+strconv.Itoa(int(w.Weekday()))) + " " + core.DateOf(w).String() + " " +
					timeText(c, TimeOfDay{w.Hour(), w.Minute(), 0}, timeConfig{hour12: opts.Hour12}) + " " + w.Format("MST")
				if run.Repeated {
					text += " · " + core.Msg(c, "cron_editor.repeated")
				}
				ui.Text(c, text).FontFeatures("tnum").TextColor(k.Text).Role(ui.RoleListItem)
			}
		})
	})
	changed := *expr != old && !col.IsDisabled()
	if changed {
		slog.Debug("mujica: cron_editor change", "valid", (*cache).err == nil)
	}
	return core.NewChoiceResult(col, changed)
}

// cronParsed caches the parse and the runs from one minute.
type cronParsed struct {
	expr string
	from time.Time
	s    CronSchedule
	err  error
	runs []CronRun
	ok   bool
}

var cronFieldKeys = [5]string{"cron_editor.minute", "cron_editor.hour", "cron_editor.dom", "cron_editor.month", "cron_editor.dow"}

func cronProblem(c *ui.Context, e *CronError) string {
	reason := core.Msg(c, cronProblemKeys[e.Problem])
	switch e.Problem {
	case CronFieldCount, CronBadStep, CronNotNumber, CronBackwards:
		reason = fmt.Sprintf(reason, e.Text)
	case CronOutOfRange:
		reason = fmt.Sprintf(reason, e.Text, cronBounds[e.Field][0], cronBounds[e.Field][1])
	}
	if e.Field < 0 {
		return reason
	}
	return fmt.Sprintf(core.Msg(c, "cron_editor.field_error"), core.Msg(c, cronFieldKeys[e.Field]), reason)
}

// cronFieldCard shows one field's token and meaning.
func cronFieldCard(c *ui.Context, f int, tok string, bad bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	border := k.Border
	if bad {
		border = k.Danger
	}
	ui.Column(c).MinWidth(d.ControlHeight()*3).Padding(d.Unit()*2).Gap(d.Unit()).Radius(theme.CardRadius).
		Border(theme.BorderWidth, border).Background(k.Surface).Children(func() {
		ui.Text(c, core.Msg(c, cronFieldKeys[f])).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		ui.Text(c, tok).Font(theme.MonoFont).TextColor(k.Text)
		desc := cronDescribe(c, f, tok)
		if bad {
			desc = core.Msg(c, "cron_editor.invalid")
		}
		ui.Text(c, desc).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
}

// cronDescribe words a field token, as "every 15" or "9–17".
func cronDescribe(c *ui.Context, f int, tok string) string {
	if tok == "*" {
		return core.Msg(c, "cron_editor.any")
	}
	var parts []string
	for _, item := range strings.Split(tok, ",") {
		rng, step, hasStep := strings.Cut(item, "/")
		rng = strings.ReplaceAll(rng, "-", "–")
		switch {
		case hasStep && rng == "*":
			parts = append(parts, fmt.Sprintf(core.Msg(c, "cron_editor.every"), step))
		case hasStep:
			parts = append(parts, fmt.Sprintf(core.Msg(c, "cron_editor.every_in"), step, rng))
		default:
			parts = append(parts, rng)
		}
	}
	return strings.Join(parts, ", ")
}
