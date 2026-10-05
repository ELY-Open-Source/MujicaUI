package overlay

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// overlayUntil draws frames until cond holds, for at most 5 s of wall time.
func overlayUntil(tt *ui.Tester, cond func() bool) bool {
	end := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(end) {
			return false
		}
		time.Sleep(2 * time.Millisecond)
		tt.Frame()
	}
	return true
}
