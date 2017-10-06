// gapbuffer project gapbuffer.go
package gapbuffer

import (
	"io"
)

type GapBuffer struct {
	data   []rune
	gapPos int
	gapLen int
}

func New(data []rune) *GapBuffer {
	return &GapBuffer{
		data: data,
	}
}

func (gp *GapBuffer) compress() {
	copy(gp.data[gp.gapPos:], gp.data[gp.gapPos+gp.gapLen:])
	gp.data = gp.data[:len(gp.data)-gp.gapLen]
	gp.gapPos = 0
	gp.gapLen = 0
}

func (gp *GapBuffer) gap(p, l int) {
	var dst []rune
	if cap(gp.data) < len(gp.data)+l {
		dst = make([]rune, len(gp.data)+l)
		copy(dst, gp.data[:p])
	} else {
		dst = gp.data[:len(gp.data)+l]
	}
	copy(dst[p+l:], gp.data[p:])
	gp.data = dst
	gp.gapPos = p
	gp.gapLen = l
}

func (gp *GapBuffer) Insert(pos int, char rune) {
	if gp.gapPos != pos {
		gp.compress()
	}
	if gp.gapLen <= 0 {
		gp.gap(pos, 100)
	}
	gp.data[pos] = char
	gp.gapPos++
	gp.gapLen--
}

func (gp *GapBuffer) Delete(pos int) {
	if gp.gapPos != pos {
		if pos == gp.gapPos+1 {
			gp.gapLen++
			return
		}
		if gp.gapLen > 0 {
			gp.compress()
		}
		gp.gapPos = pos
		gp.gapLen = 0
	}
	gp.gapPos--
	gp.gapLen++
}

func (gp *GapBuffer) Replace(pos int, char rune) {
	if pos >= gp.gapPos {
		pos += gp.gapLen
	}
	gp.data[pos] = char
}

func (gp *GapBuffer) String() string {
	if gp.gapLen > 0 {
		gp.compress()
	}
	return string(gp.data)
}

func (gp *GapBuffer) Get(pos int) rune {
	if pos >= gp.gapPos {
		pos += gp.gapLen
	}
	return gp.data[pos]
}

func (gp *GapBuffer) Len() int {
	return len(gp.data) - gp.gapLen
}

func (gp *GapBuffer) WriteTo(w io.Writer) {
	w.Write([]byte(gp.String()))
}