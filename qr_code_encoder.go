package mujica

import (
	"fmt"
)

// QRLevel is a QR code's error correction level. The zero value is QRMedium.
type QRLevel int

const (
	QRMedium   QRLevel = iota // M: about 15% recoverable
	QRLow                     // L: about 7%
	QRQuartile                // Q: about 25%
	QRHigh                    // H: about 30%
)

func (l QRLevel) String() string {
	switch l {
	case QRLow:
		return "L"
	case QRMedium:
		return "M"
	case QRQuartile:
		return "Q"
	case QRHigh:
		return "H"
	}
	panic(fmt.Sprintf("mujica: unknown QRLevel %d", int(l)))
}

// index is the level's column in qrCodeTable.
func (l QRLevel) index() int { return [4]int{1, 0, 2, 3}[l] }

// formatBits are the level's two format information bits.
func (l QRLevel) formatBits() int { return [4]int{0, 1, 3, 2}[l] }

// qrCodeBlocks is one version and level: EC codewords per block, then two block groups.
type qrCodeBlocks struct{ ec, n1, d1, n2, d2 int }

func (b qrCodeBlocks) data() int { return b.n1*b.d1 + b.n2*b.d2 }

// qrCodeTable is ISO/IEC 18004 Table 9, indexed [version-1][L, M, Q, H].
var qrCodeTable = [40][4]qrCodeBlocks{
	{{7, 1, 19, 0, 0}, {10, 1, 16, 0, 0}, {13, 1, 13, 0, 0}, {17, 1, 9, 0, 0}},
	{{10, 1, 34, 0, 0}, {16, 1, 28, 0, 0}, {22, 1, 22, 0, 0}, {28, 1, 16, 0, 0}},
	{{15, 1, 55, 0, 0}, {26, 1, 44, 0, 0}, {18, 2, 17, 0, 0}, {22, 2, 13, 0, 0}},
	{{20, 1, 80, 0, 0}, {18, 2, 32, 0, 0}, {26, 2, 24, 0, 0}, {16, 4, 9, 0, 0}},
	{{26, 1, 108, 0, 0}, {24, 2, 43, 0, 0}, {18, 2, 15, 2, 16}, {22, 2, 11, 2, 12}},
	{{18, 2, 68, 0, 0}, {16, 4, 27, 0, 0}, {24, 4, 19, 0, 0}, {28, 4, 15, 0, 0}},
	{{20, 2, 78, 0, 0}, {18, 4, 31, 0, 0}, {18, 2, 14, 4, 15}, {26, 4, 13, 1, 14}},
	{{24, 2, 97, 0, 0}, {22, 2, 38, 2, 39}, {22, 4, 18, 2, 19}, {26, 4, 14, 2, 15}},
	{{30, 2, 116, 0, 0}, {22, 3, 36, 2, 37}, {20, 4, 16, 4, 17}, {24, 4, 12, 4, 13}},
	{{18, 2, 68, 2, 69}, {26, 4, 43, 1, 44}, {24, 6, 19, 2, 20}, {28, 6, 15, 2, 16}},
	{{20, 4, 81, 0, 0}, {30, 1, 50, 4, 51}, {28, 4, 22, 4, 23}, {24, 3, 12, 8, 13}},
	{{24, 2, 92, 2, 93}, {22, 6, 36, 2, 37}, {26, 4, 20, 6, 21}, {28, 7, 14, 4, 15}},
	{{26, 4, 107, 0, 0}, {22, 8, 37, 1, 38}, {24, 8, 20, 4, 21}, {22, 12, 11, 4, 12}},
	{{30, 3, 115, 1, 116}, {24, 4, 40, 5, 41}, {20, 11, 16, 5, 17}, {24, 11, 12, 5, 13}},
	{{22, 5, 87, 1, 88}, {24, 5, 41, 5, 42}, {30, 5, 24, 7, 25}, {24, 11, 12, 7, 13}},
	{{24, 5, 98, 1, 99}, {28, 7, 45, 3, 46}, {24, 15, 19, 2, 20}, {30, 3, 15, 13, 16}},
	{{28, 1, 107, 5, 108}, {28, 10, 46, 1, 47}, {28, 1, 22, 15, 23}, {28, 2, 14, 17, 15}},
	{{30, 5, 120, 1, 121}, {26, 9, 43, 4, 44}, {28, 17, 22, 1, 23}, {28, 2, 14, 19, 15}},
	{{28, 3, 113, 4, 114}, {26, 3, 44, 11, 45}, {26, 17, 21, 4, 22}, {26, 9, 13, 16, 14}},
	{{28, 3, 107, 5, 108}, {26, 3, 41, 13, 42}, {30, 15, 24, 5, 25}, {28, 15, 15, 10, 16}},
	{{28, 4, 116, 4, 117}, {26, 17, 42, 0, 0}, {28, 17, 22, 6, 23}, {30, 19, 16, 6, 17}},
	{{28, 2, 111, 7, 112}, {28, 17, 46, 0, 0}, {30, 7, 24, 16, 25}, {24, 34, 13, 0, 0}},
	{{30, 4, 121, 5, 122}, {28, 4, 47, 14, 48}, {30, 11, 24, 14, 25}, {30, 16, 15, 14, 16}},
	{{30, 6, 117, 4, 118}, {28, 6, 45, 14, 46}, {30, 11, 24, 16, 25}, {30, 30, 16, 2, 17}},
	{{26, 8, 106, 4, 107}, {28, 8, 47, 13, 48}, {30, 7, 24, 22, 25}, {30, 22, 15, 13, 16}},
	{{28, 10, 114, 2, 115}, {28, 19, 46, 4, 47}, {28, 28, 22, 6, 23}, {30, 33, 16, 4, 17}},
	{{30, 8, 122, 4, 123}, {28, 22, 45, 3, 46}, {30, 8, 23, 26, 24}, {30, 12, 15, 28, 16}},
	{{30, 3, 117, 10, 118}, {28, 3, 45, 23, 46}, {30, 4, 24, 31, 25}, {30, 11, 15, 31, 16}},
	{{30, 7, 116, 7, 117}, {28, 21, 45, 7, 46}, {30, 1, 23, 37, 24}, {30, 19, 15, 26, 16}},
	{{30, 5, 115, 10, 116}, {28, 19, 47, 10, 48}, {30, 15, 24, 25, 25}, {30, 23, 15, 25, 16}},
	{{30, 13, 115, 3, 116}, {28, 2, 46, 29, 47}, {30, 42, 24, 1, 25}, {30, 23, 15, 28, 16}},
	{{30, 17, 115, 0, 0}, {28, 10, 46, 23, 47}, {30, 10, 24, 35, 25}, {30, 19, 15, 35, 16}},
	{{30, 17, 115, 1, 116}, {28, 14, 46, 21, 47}, {30, 29, 24, 19, 25}, {30, 11, 15, 46, 16}},
	{{30, 13, 115, 6, 116}, {28, 14, 46, 23, 47}, {30, 44, 24, 7, 25}, {30, 59, 16, 1, 17}},
	{{30, 12, 121, 7, 122}, {28, 12, 47, 26, 48}, {30, 39, 24, 14, 25}, {30, 22, 15, 41, 16}},
	{{30, 6, 121, 14, 122}, {28, 6, 47, 34, 48}, {30, 46, 24, 10, 25}, {30, 2, 15, 64, 16}},
	{{30, 17, 122, 4, 123}, {28, 29, 46, 14, 47}, {30, 49, 24, 10, 25}, {30, 24, 15, 46, 16}},
	{{30, 4, 122, 18, 123}, {28, 13, 46, 32, 47}, {30, 48, 24, 14, 25}, {30, 42, 15, 32, 16}},
	{{30, 20, 117, 4, 118}, {28, 40, 47, 7, 48}, {30, 43, 24, 22, 25}, {30, 10, 15, 67, 16}},
	{{30, 19, 118, 6, 119}, {28, 18, 47, 31, 48}, {30, 34, 24, 34, 25}, {30, 20, 15, 61, 16}},
}

// QRMatrix is an encoded QR code: Modules[y][x] is true for dark.
type QRMatrix struct {
	Version int
	Level   QRLevel
	Mask    int
	Modules [][]bool
}

// Size is the side of the symbol in modules, without the quiet zone.
func (m QRMatrix) Size() int { return len(m.Modules) }

// EncodeQR encodes data in byte mode at the smallest version that fits, choosing the mask by penalty.
func EncodeQR(data []byte, level QRLevel) (QRMatrix, error) {
	_ = level.String()
	for v := 1; v <= 40; v++ {
		cc := 8
		if v >= 10 {
			cc = 16
		}
		if len(data) >= 1<<cc {
			continue
		}
		if 4+cc+8*len(data) <= 8*qrCodeTable[v-1][level.index()].data() {
			return qrCodeEncodeAt(data, level, v), nil
		}
	}
	return QRMatrix{}, fmt.Errorf("mujica: %d bytes do not fit a QR code at level %v", len(data), level)
}

// qrCodeEncodeAt builds version v with the lowest-penalty mask.
func qrCodeEncodeAt(data []byte, level QRLevel, v int) QRMatrix {
	b := qrCodeTable[v-1][level.index()]
	var w qrCodeBits
	w.put(4, 4)
	if v >= 10 {
		w.put(len(data), 16)
	} else {
		w.put(len(data), 8)
	}
	for _, d := range data {
		w.put(int(d), 8)
	}
	capBits := 8 * b.data()
	w.put(0, min(4, capBits-w.n))
	w.put(0, (8-w.n%8)%8)
	for pad := 0xEC; w.n < capBits; pad ^= 0xEC ^ 0x11 {
		w.put(pad, 8)
	}
	words := qrCodeInterleave(w.bytes, b)
	g := newQRCodeGrid(v)
	g.place(words)
	best, bestPenalty := 0, -1
	for m := range 8 {
		t := g.clone()
		t.apply(m)
		if p := t.penalty(); bestPenalty < 0 || p < bestPenalty {
			best, bestPenalty = m, p
		}
	}
	g.apply(best)
	g.format(level.formatBits(), best)
	g.version(true)
	return QRMatrix{Version: v, Level: level, Mask: best, Modules: g.dark}
}

// qrCodeBits appends bits most significant first.
type qrCodeBits struct {
	bytes []byte
	n     int
}

func (w *qrCodeBits) put(v, count int) {
	for i := count - 1; i >= 0; i-- {
		if w.n%8 == 0 {
			w.bytes = append(w.bytes, 0)
		}
		if v>>i&1 == 1 {
			w.bytes[w.n/8] |= 0x80 >> (w.n % 8)
		}
		w.n++
	}
}

// qrCodeInterleave splits data into blocks, adds Reed-Solomon codewords and interleaves both.
func qrCodeInterleave(data []byte, b qrCodeBlocks) []byte {
	var blocks, ecs [][]byte
	gen := qrCodeGenerator(b.ec)
	at := 0
	for i := range b.n1 + b.n2 {
		n := b.d1
		if i >= b.n1 {
			n = b.d2
		}
		blk := data[at : at+n]
		at += n
		blocks = append(blocks, blk)
		ecs = append(ecs, qrCodeRemainder(blk, gen))
	}
	var out []byte
	for i := range max(b.d1, b.d2) {
		for _, blk := range blocks {
			if i < len(blk) {
				out = append(out, blk[i])
			}
		}
	}
	for i := range b.ec {
		for _, e := range ecs {
			out = append(out, e[i])
		}
	}
	return out
}

// qrCodeMul multiplies in GF(256) modulo x^8+x^4+x^3+x^2+1.
func qrCodeMul(x, y byte) byte {
	var z byte
	for i := 7; i >= 0; i-- {
		hi := z & 0x80
		z <<= 1
		if hi != 0 {
			z ^= 0x1D
		}
		if y>>i&1 == 1 {
			z ^= x
		}
	}
	return z
}

// qrCodeGenerator returns the degree-n generator's coefficients, highest first, leading 1 dropped.
func qrCodeGenerator(n int) []byte {
	g := make([]byte, n)
	g[n-1] = 1
	root := byte(1)
	for range n {
		for j := range n {
			g[j] = qrCodeMul(g[j], root)
			if j+1 < n {
				g[j] ^= g[j+1]
			}
		}
		root = qrCodeMul(root, 2)
	}
	return g
}

func qrCodeRemainder(data, gen []byte) []byte {
	r := make([]byte, len(gen))
	for _, d := range data {
		f := d ^ r[0]
		copy(r, r[1:])
		r[len(r)-1] = 0
		for i := range r {
			r[i] ^= qrCodeMul(gen[i], f)
		}
	}
	return r
}

// qrCodeAlignment returns the alignment pattern centers of version v.
func qrCodeAlignment(v int) []int {
	if v == 1 {
		return nil
	}
	n := v/7 + 2
	step := 26
	if v != 32 {
		step = (v*4 + n*2 + 1) / (n*2 - 2) * 2
	}
	pos := make([]int, n)
	pos[0] = 6
	for i, p := n-1, v*4+10; i >= 1; i, p = i-1, p-step {
		pos[i] = p
	}
	return pos
}
