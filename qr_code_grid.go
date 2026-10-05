package mujica

// qrCodeGrid is a symbol under construction; fixed marks function modules.
type qrCodeGrid struct {
	v           int
	dark, fixed [][]bool
}

func newQRCodeGrid(v int) *qrCodeGrid {
	n := 4*v + 17
	g := &qrCodeGrid{v: v, dark: make([][]bool, n), fixed: make([][]bool, n)}
	for i := range n {
		g.dark[i], g.fixed[i] = make([]bool, n), make([]bool, n)
	}
	for i := range n {
		g.set(6, i, i%2 == 0)
		g.set(i, 6, i%2 == 0)
	}
	g.finder(3, 3)
	g.finder(n-4, 3)
	g.finder(3, n-4)
	al := qrCodeAlignment(v)
	for i, x := range al {
		for j, y := range al {
			if i == 0 && j == 0 || i == 0 && j == len(al)-1 || i == len(al)-1 && j == 0 {
				continue
			}
			for dy := -2; dy <= 2; dy++ {
				for dx := -2; dx <= 2; dx++ {
					g.set(x+dx, y+dy, max(qrCodeAbs(dx), qrCodeAbs(dy)) != 1)
				}
			}
		}
	}
	// Reserve format and version areas; apply fills them.
	g.format(-1, 0)
	g.version(false)
	return g
}

func qrCodeAbs(v int) int { return max(v, -v) }

func (g *qrCodeGrid) set(x, y int, dark bool) {
	g.dark[y][x], g.fixed[y][x] = dark, true
}

// finder draws a finder pattern and its separator centered at (cx, cy).
func (g *qrCodeGrid) finder(cx, cy int) {
	n := len(g.dark)
	for dy := -4; dy <= 4; dy++ {
		for dx := -4; dx <= 4; dx++ {
			x, y := cx+dx, cy+dy
			if x < 0 || y < 0 || x >= n || y >= n {
				continue
			}
			d := max(qrCodeAbs(dx), qrCodeAbs(dy))
			g.set(x, y, d != 2 && d != 4)
		}
	}
}

// format writes the 15 format bits and the dark module; levelBits -1 leaves them light.
func (g *qrCodeGrid) format(levelBits, mask int) {
	data := levelBits<<3 | mask
	rem := data
	for range 10 {
		rem = rem<<1 ^ (rem>>9)*0x537
	}
	b := (data<<10 | rem) ^ 0x5412
	bit := func(i int) bool { return levelBits >= 0 && b>>i&1 == 1 }
	n := len(g.dark)
	for i := 0; i <= 5; i++ {
		g.set(8, i, bit(i))
	}
	g.set(8, 7, bit(6))
	g.set(8, 8, bit(7))
	g.set(7, 8, bit(8))
	for i := 9; i < 15; i++ {
		g.set(14-i, 8, bit(i))
	}
	for i := range 8 {
		g.set(n-1-i, 8, bit(i))
	}
	for i := 8; i < 15; i++ {
		g.set(8, n-15+i, bit(i))
	}
	g.set(8, n-8, levelBits >= 0)
}

// version writes the 18 version bits from version 7 on, or light ones.
func (g *qrCodeGrid) version(on bool) {
	if g.v < 7 {
		return
	}
	rem := g.v
	for range 12 {
		rem = rem<<1 ^ (rem>>11)*0x1F25
	}
	b := g.v<<12 | rem
	n := len(g.dark)
	for i := range 18 {
		dark := on && b>>i&1 == 1
		a, c := n-11+i%3, i/3
		g.set(a, c, dark)
		g.set(c, a, dark)
	}
}

// place lays the codewords in the zigzag from the bottom right; remainder bits stay light.
func (g *qrCodeGrid) place(words []byte) {
	n := len(g.dark)
	i := 0
	for right := n - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for vert := range n {
			for j := range 2 {
				x := right - j
				up := (right+1)&2 == 0
				y := vert
				if up {
					y = n - 1 - vert
				}
				if g.fixed[y][x] || i >= len(words)*8 {
					continue
				}
				g.dark[y][x] = words[i/8]>>(7-i%8)&1 == 1
				i++
			}
		}
	}
}

func (g *qrCodeGrid) clone() *qrCodeGrid {
	t := &qrCodeGrid{v: g.v, dark: make([][]bool, len(g.dark)), fixed: g.fixed}
	for i, row := range g.dark {
		t.dark[i] = append([]bool(nil), row...)
	}
	return t
}

// apply XORs mask over the data modules; the information areas stay light, as 7.8.3 evaluates them.
func (g *qrCodeGrid) apply(mask int) {
	for y, row := range g.dark {
		for x := range row {
			if g.fixed[y][x] {
				continue
			}
			var flip bool
			switch mask {
			case 0:
				flip = (x+y)%2 == 0
			case 1:
				flip = y%2 == 0
			case 2:
				flip = x%3 == 0
			case 3:
				flip = (x+y)%3 == 0
			case 4:
				flip = (x/3+y/2)%2 == 0
			case 5:
				flip = x*y%2+x*y%3 == 0
			case 6:
				flip = (x*y%2+x*y%3)%2 == 0
			case 7:
				flip = ((x+y)%2+x*y%3)%2 == 0
			}
			row[x] = row[x] != flip
		}
	}
}

// penalty scores the symbol by the four rules of ISO/IEC 18004 7.8.3.
func (g *qrCodeGrid) penalty() int {
	n := len(g.dark)
	p, dark := 0, 0
	row, col := make([]bool, n), make([]bool, n)
	for i := range n {
		for j := range n {
			row[j], col[j] = g.dark[i][j], g.dark[j][i]
			if row[j] {
				dark++
			}
			if i > 0 && j > 0 && row[j] == row[j-1] && row[j] == g.dark[i-1][j] && row[j] == g.dark[i-1][j-1] {
				p += 3
			}
		}
		p += qrCodeRuns(row) + qrCodeFinders(row) + qrCodeRuns(col) + qrCodeFinders(col)
	}
	return p + 10*(qrCodeAbs(dark*100-50*n*n)/(5*n*n))
}

// qrCodeRuns scores runs of five or more equal modules.
func qrCodeRuns(seq []bool) int {
	p, run := 0, 1
	for j := 1; j <= len(seq); j++ {
		if j < len(seq) && seq[j] == seq[j-1] {
			run++
			continue
		}
		if run >= 5 {
			p += run - 2
		}
		run = 1
	}
	return p
}

// qrCodeFinders scores 1:1:3:1:1 patterns with four light modules on a side.
func qrCodeFinders(seq []bool) int {
	core := [7]bool{true, false, true, true, true, false, true}
	n, p := len(seq), 0
	light := func(a, b int) bool {
		for i := max(a, 0); i < min(b, n); i++ {
			if seq[i] {
				return false
			}
		}
		return true
	}
	for idx := 0; idx+7 <= n; {
		match := true
		for i, d := range core {
			if seq[idx+i] != d {
				match = false
				break
			}
		}
		if !match {
			idx++
			continue
		}
		if idx == 0 || idx == n-7 || light(idx-4, idx) || light(idx+7, idx+11) {
			p += 40
			idx += 7
		} else {
			idx += 4
		}
	}
	return p
}
