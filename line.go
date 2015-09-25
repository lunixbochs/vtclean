package vtclean

type lineEdit struct {
	buf       []byte
	pos, size int
}

func (l *lineEdit) Move(x int) {
	if x < 0 && l.pos <= -x {
		l.pos = 0
	} else if x > 0 && l.pos+x > l.size {
		l.pos = l.size
	} else {
		l.pos += x
	}
}

func (l *lineEdit) Write(p []byte) {
	if len(l.buf)-l.pos < len(p) {
		l.buf = append(l.buf[:l.pos], p...)
	} else {
		copy(l.buf[l.pos:], p)
	}
	l.pos += len(p)
	l.size += len(p)
}

func (l *lineEdit) Insert(p []byte) {
	left := append(l.buf[:l.pos], p...)
	l.buf = append(left, l.buf[l.pos:]...)
	l.size += len(p)
}

func (l *lineEdit) Delete(n int) {
	most := l.size - l.pos
	if n > most {
		n = most
	}
	copy(l.buf[l.pos:], l.buf[l.pos+n:])
	l.size -= n
}

func (l *lineEdit) Clear() {
	for i := 0; i < len(l.buf); i++ {
		l.buf[i] = ' '
	}
}
func (l *lineEdit) ClearLeft() {
	for i := 0; i < l.pos; i++ {
		l.buf[i] = ' '
	}
}
func (l *lineEdit) ClearRight() {
	l.size = l.pos
}

func (l *lineEdit) Bytes() []byte {
	return l.buf[:l.size]
}
