package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestQRCodeDrawsModulesAndCaches(t *testing.T) {
	text := "HELLO"
	var box ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Dark})
		ui.Row(c).Children(func() { box = QRCode(c, text, QRCodeOptions{Level: QRLow, Size: 116}).Bounds() })
	}, 200, 200)
	tt.Frame()
	if box.W != 116 {
		t.Fatalf("size %v", box)
	}
	m, _ := EncodeQR([]byte(text), QRLow)
	// 21 modules plus 8 of quiet zone in 116 DIPs: 4 DIPs each.
	px := func(x, y int) ui.Color {
		r, g, b, _ := tt.Image().At(int(box.X)+(x+4)*4+2, int(box.Y)+(y+4)*4+2).RGBA()
		return ui.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	}
	for _, at := range [][2]int{{0, 0}, {3, 3}, {7, 7}, {10, 6}, {12, 12}} {
		want := theme.Light().Surface
		if m.Modules[at[1]][at[0]] {
			want = theme.Light().Text
		}
		if got := px(at[0], at[1]); got != want {
			t.Errorf("module %v is %v, want %v", at, got, want)
		}
	}
	if _, ok := tt.Find("QR code"); !ok {
		t.Fatal("no accessible name")
	}
}

func TestQRCodeTooLong(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		QRCode(c, strings.Repeat("x", 1300), QRCodeOptions{Level: QRHigh})
	}, 300, 300)
	if !tt.HasText("Text too long for a QR code") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestQRCodeRejectsBadOptions(t *testing.T) {
	for name, o := range map[string]QRCodeOptions{"level": {Level: 9}, "size": {Size: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); QRCode(c, "x", o) }, 200, 200)
		}()
	}
}

func TestQRCodeSnapshot(t *testing.T) {
	snapshot(t, "qr_code", 420, 200, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			QRCode(c, "https://example.com/mujica", QRCodeOptions{})
			QRCode(c, "MujicaUI", QRCodeOptions{Level: QRHigh, Size: 100})
			QRCode(c, strings.Repeat("x", 1300), QRCodeOptions{Level: QRHigh, Size: 100})
		})
	})
}
