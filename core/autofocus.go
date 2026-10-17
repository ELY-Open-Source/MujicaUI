package core

import "github.com/egoist/mygo/ui"

type autoFocusKey struct{}

// AutoFocus focuses e the frame after it appears, once MyGo has recorded an overlay's opener from the focus it had.
func AutoFocus(c *ui.Context, e ui.Element) ui.Element {
	step := ui.Local(e, autoFocusKey{}, func() int { return 0 })
	switch *step {
	case 0:
		*step = 1
		c.Invalidate()
	case 1:
		*step = 2
		e.Focus()
	}
	return e
}
