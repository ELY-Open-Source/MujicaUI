package mujica

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
)

// Priority describes urgency; zero means no priority.
type Priority int

const (
	PriorityNone Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

// Member is a person tasks can be assigned to.
type Member struct {
	ID   string
	Name string
}

// TaskLabel is a colored tag tasks carry.
type TaskLabel struct {
	ID    string
	Name  string
	Color EventColor
}

// Subtask is one checklist step of a task.
type Subtask struct {
	ID    string
	Title string
	Done  bool
}

// Task is one unit of project work.
type Task struct {
	// ID is the issue key, as ENG-123.
	ID        string
	Title     string
	Done      bool
	Priority  Priority
	Due       Date
	Status    string
	Assignees []string
	Labels    []string
	Subtasks  []Subtask
}

// StatusGroup is the workflow stage a status belongs to.
type StatusGroup int

const (
	StatusBacklog StatusGroup = iota
	StatusUnstarted
	StatusStarted
	StatusCompleted
	StatusCanceled
)

// WorkflowStatus is one workflow state.
type WorkflowStatus struct {
	ID    string
	Name  string
	Group StatusGroup
}

// ProjectResult exposes the element and task events.
type ProjectResult struct {
	Element *ui.Element
	changed string
	opened  string
}

// Changed returns the edited task ID.
func (r ProjectResult) Changed() (string, bool) { return r.changed, r.changed != "" }

// Opened returns the activated task ID.
func (r ProjectResult) Opened() (string, bool) { return r.opened, r.opened != "" }

// projectCheckPriority panics on an unknown priority.
func projectCheckPriority(p Priority) {
	if p < PriorityNone || p > PriorityUrgent {
		panic(fmt.Sprintf("mujica: Priority %d out of range", p))
	}
}

// projectMember resolves a member ID or panics.
func projectMember(members []Member, id string) Member {
	i := slices.IndexFunc(members, func(m Member) bool { return m.ID == id })
	if i < 0 {
		panic(fmt.Sprintf("mujica: unknown member %q", id))
	}
	return members[i]
}

// projectLabel resolves a label ID or panics.
func projectLabel(labels []TaskLabel, id string) TaskLabel {
	i := slices.IndexFunc(labels, func(l TaskLabel) bool { return l.ID == id })
	if i < 0 {
		panic(fmt.Sprintf("mujica: unknown label %q", id))
	}
	return labels[i]
}

// projectStatus resolves a status ID or panics.
func projectStatus(statuses []WorkflowStatus, id string) WorkflowStatus {
	i := slices.IndexFunc(statuses, func(s WorkflowStatus) bool { return s.ID == id })
	if i < 0 {
		panic(fmt.Sprintf("mujica: unknown status %q", id))
	}
	return statuses[i]
}

// projectDueText formats a date relative to today.
func projectDueText(c *ui.Context, due, today Date, done bool) (string, bool) {
	datetimeCheckDate("project due", due)
	if due.IsZero() {
		return "", false
	}
	days := int(due.civil().Sub(today.civil()).Hours() / 24)
	switch {
	case days < 0 && !done:
		return calendarCount(c, "projects.overdue", -days), true
	case days == 0:
		return msg(c, "projects.due_today"), false
	case days == 1:
		return msg(c, "projects.due_tomorrow"), false
	}
	return fmt.Sprintf(msg(c, "projects.due"), datetimeMonthName(c, due.Month), due.Day), false
}

// projectBar draws a named completion bar.
func projectBar(c *ui.Context, frac float64, label string) *ui.Element {
	k := Tokens(c)
	return ui.Box(c).Height(6).FillWidth().Role(ui.RoleProgress).Label(label).Description(fmt.Sprintf("%d%%", int(frac*100))).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.Border, 3)
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(frac), H: r.H}, k.AccentText, 3)
		})
}
