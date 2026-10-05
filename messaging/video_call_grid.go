package messaging

import (
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/media"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CallParticipant supplies identity, frames and call state.
type CallParticipant struct {
	ID, Name string
	// Frame is their video; nil shows their avatar.
	Frame     ui.ImageSource
	Image     *ui.Bitmap
	Muted     bool
	CameraOff bool
	Speaking  bool
}

// VideoCallGridOptions configure a VideoCallGrid.
type VideoCallGridOptions struct {
	// Columns fixes the column count; zero fits a near-square grid.
	Columns int
}

// VideoCallGrid displays participant frames and speaker marks.
func VideoCallGrid(c *ui.Context, people []CallParticipant, opts VideoCallGridOptions) *ui.Element {
	if opts.Columns < 0 {
		panic("mujica: VideoCallGrid Columns must not be negative")
	}
	seen := map[string]bool{}
	for _, p := range people {
		if p.ID == "" || p.Name == "" || seen[p.ID] {
			panic("mujica: VideoCallGrid participants need unique IDs and names: " + p.ID)
		}
		seen[p.ID] = true
	}
	cols := opts.Columns
	if cols == 0 {
		cols = max(1, int(math.Ceil(math.Sqrt(float64(len(people))))))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	tracks := make([]ui.Track, cols)
	for i := range tracks {
		tracks[i] = ui.Fr(1)
	}
	grid := ui.Box(c).Grid().ColumnTracks(tracks...).Gap(u * 2).FocusGroup(ui.Horizontal | ui.Vertical).
		Role(ui.RoleGroup).Label(core.Msg(c, "video_call_grid.label"))
	grid.Children(func() {
		for _, p := range people {
			callTile(c, k, p)
		}
	})
	return grid
}

func callTile(c *ui.Context, k theme.Tokens, p CallParticipant) {
	u := core.Density(c).Unit()
	var notes []string
	if p.Speaking {
		notes = append(notes, core.Msg(c, "video_call_grid.speaking"))
	}
	if p.Muted {
		notes = append(notes, core.Msg(c, "video_call_grid.muted"))
	}
	if p.CameraOff {
		notes = append(notes, core.Msg(c, "video_call_grid.cameraOff"))
	}
	label := p.Name
	if len(notes) > 0 {
		label += ", " + strings.Join(notes, ", ")
	}
	tile := ui.Box(c).Key(p.ID).AspectRatio(16.0 / 9).Radius(theme.CardRadius).Clip().Background(media.VideoBlack).Center().
		Focusable().Role(ui.RoleImage).Label(label)
	internal.Ring(tile, k.Focus, theme.CardRadius)
	light := ui.Hex("#F2EBDD")
	tile.Children(func() {
		if p.Frame != nil && !p.CameraOff {
			ui.Image(c, p.Frame).Fit(ui.Cover).Absolute().Left(0).Top(0).Right(0).Bottom(0)
		} else {
			display.Avatar(c, p.Name, display.AvatarOptions{Image: p.Image, Size: 56})
		}
		ui.Row(c).Absolute().Left(u*2).Bottom(u*2).Gap(u).AlignItems(ui.Center).Padding(2, u*1.5).
			Radius(theme.ControlRadius).Background(ui.RGBA(0, 0, 0, 0.6)).TextColor(light).Children(func() {
			if p.Muted {
				ui.Icon(c, msgMicOff).Size(13, 13).TextColor(ui.Hex("#F0A1AA"))
			}
			if p.Speaking {
				ui.Icon(c, msgSpeaking).Size(13, 13)
			}
			ui.Text(c, p.Name).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
		})
	})
	if p.Speaking {
		tile.DrawOver(func(pt *ui.Painter, r ui.Rect) {
			pt.Stroke(ui.Rect{X: r.X + 1.5, Y: r.Y + 1.5, W: r.W - 3, H: r.H - 3}, k.AccentText, theme.CardRadius, 3)
			if tile.FocusVisible() {
				internal.FocusRing(pt, r, theme.CardRadius, k.Focus)
			}
		})
	}
}
