package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// mediaSVG wraps stroked 24×24 path data as an icon.
func mediaSVG(body string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
}

// Media glyphs drawn for this group.
var (
	mediaPlay       = mediaSVG(`<path d="M7 4.5v15l12-7.5z"/>`)
	mediaPause      = mediaSVG(`<path d="M8 5v14M16 5v14"/>`)
	mediaPrev       = mediaSVG(`<path d="M18 5v14L8 12zM6 5v14"/>`)
	mediaNext       = mediaSVG(`<path d="M6 5v14l10-7zM18 5v14"/>`)
	mediaRepeat     = mediaSVG(`<path d="M17 3l3 3-3 3M4 11V9a3 3 0 0 1 3-3h13M7 21l-3-3 3-3M20 13v2a3 3 0 0 1-3 3H4"/>`)
	mediaRepeatOne  = mediaSVG(`<path d="M17 3l3 3-3 3M4 11V9a3 3 0 0 1 3-3h13M7 21l-3-3 3-3M20 13v2a3 3 0 0 1-3 3H4M11 10h1v4"/>`)
	mediaShuffle    = mediaSVG(`<path d="M16 3h5v5M4 20 21 3M21 16v5h-5M15 15l6 6M4 4l5 5"/>`)
	mediaVolume     = mediaSVG(`<path d="M11 5 6 9H3v6h3l5 4zM15.5 8.5a5 5 0 0 1 0 7M18.5 5.5a9 9 0 0 1 0 13"/>`)
	mediaMute       = mediaSVG(`<path d="M11 5 6 9H3v6h3l5 4zM16 9l6 6M22 9l-6 6"/>`)
	mediaRotate     = mediaSVG(`<path d="M21 12a9 9 0 1 1-3-6.7L21 8M21 3v5h-5"/>`)
	mediaRotateLeft = mediaSVG(`<path d="M3 12a9 9 0 1 0 3-6.7L3 8M3 3v5h5"/>`)
	mediaMirror     = mediaSVG(`<path d="M12 3v18M8 7 3 12l5 5zM16 7l5 5-5 5z"/>`)
	mediaPiP        = mediaSVG(`<rect x="2" y="4" width="20" height="16" rx="2"/><rect x="12" y="12" width="7" height="5" rx="1"/>`)
	mediaFullscreen = mediaSVG(`<path d="M3 8V3h5M16 3h5v5M21 16v5h-5M8 21H3v-5"/>`)
	mediaCaptions   = mediaSVG(`<rect x="2" y="5" width="20" height="14" rx="2"/><path d="M10 10.5a2 2 0 1 0 0 3M17 10.5a2 2 0 1 0 0 3"/>`)
	mediaRecord     = mediaSVG(`<circle cx="12" cy="12" r="6" fill="currentColor"/>`)
	mediaStop       = mediaSVG(`<rect x="6" y="6" width="12" height="12" rx="1" fill="currentColor"/>`)
	mediaMic        = mediaSVG(`<rect x="9" y="2" width="6" height="12" rx="3"/><path d="M5 10a7 7 0 0 0 14 0M12 17v5M8 22h8"/>`)
	mediaCamera     = mediaSVG(`<rect x="2" y="6" width="14" height="12" rx="2"/><path d="m16 10 6-4v12l-6-4z"/>`)
	mediaSpeaker    = mediaSVG(`<rect x="5" y="2" width="14" height="20" rx="2"/><circle cx="12" cy="14" r="4"/><path d="M12 6h.01"/>`)
	mediaRect       = mediaSVG(`<rect x="4" y="6" width="16" height="12" rx="1"/>`)
	mediaArrow      = mediaSVG(`<path d="M5 19 19 5M10 5h9v9"/>`)
	mediaText       = mediaSVG(`<path d="M5 6V4h14v2M12 4v16M9 20h6"/>`)
	mediaPointer    = mediaSVG(`<path d="M5 3l14 7-6 2-2 6z"/>`)
)

// mediaClock writes d as m:ss, or h:mm:ss past an hour.
func mediaClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	s := int(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
