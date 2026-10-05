package input

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// testdata/qr: segno 1.6.6 byte-mode symbols, padding fixed to ISO 7.4.10.
func TestEncodeQRMatchesReference(t *testing.T) {
	files, err := filepath.Glob("testdata/qr/*.txt")
	if err != nil || len(files) < 6 {
		t.Fatalf("reference files: %v %v", files, err)
	}
	levels := map[string]QRLevel{"L": QRLow, "M": QRMedium, "Q": QRQuartile, "H": QRHigh}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(string(b)), "\n")
		head := strings.Fields(lines[0])
		version, _ := strconv.Atoi(head[0])
		mask, _ := strconv.Atoi(head[2])
		data, err := hex.DecodeString(head[3])
		if err != nil {
			t.Fatal(err)
		}
		m, err := EncodeQR(data, levels[head[1]])
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if m.Version != version || m.Mask != mask {
			t.Errorf("%s: version %d mask %d, want %d %d", f, m.Version, m.Mask, version, mask)
			continue
		}
		for y, row := range lines[1:] {
			for x, ch := range row {
				if m.Modules[y][x] != (ch == '1') {
					t.Fatalf("%s: module (%d,%d) differs", f, x, y)
				}
			}
		}
	}
}

func TestQRCodeAlignmentPositions(t *testing.T) {
	for v, want := range map[int]string{1: "[]", 2: "[6 18]", 7: "[6 22 38]", 32: "[6 34 60 86 112 138]", 40: "[6 30 58 86 114 142 170]"} {
		if got := fmtInts(qrCodeAlignment(v)); got != want {
			t.Errorf("version %d: %s, want %s", v, got, want)
		}
	}
}

func fmtInts(v []int) string {
	s := make([]string, len(v))
	for i, n := range v {
		s[i] = strconv.Itoa(n)
	}
	return "[" + strings.Join(s, " ") + "]"
}

func TestEncodeQRLimits(t *testing.T) {
	m, err := EncodeQR(make([]byte, 2953), QRLow)
	if err != nil || m.Version != 40 || m.Size() != 177 {
		t.Fatalf("largest L symbol: %v %v", m.Version, err)
	}
	if _, err := EncodeQR(make([]byte, 2954), QRLow); err == nil {
		t.Fatal("2954 bytes fit at L")
	}
	if _, err := EncodeQR(make([]byte, 1274), QRHigh); err == nil {
		t.Fatal("1274 bytes fit at H")
	}
	m, err = EncodeQR(nil, QRHigh)
	if err != nil || m.Version != 1 {
		t.Fatalf("empty data: %v %v", m.Version, err)
	}
}
