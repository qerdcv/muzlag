package output

import (
	"io"
)

type chOut struct {
	out chan []byte
}

func NewChOut(ch chan []byte) io.Writer {
	return &chOut{ch}
}

func (o *chOut) Write(b []byte) (int, error) {
	o.out <- b
	return len(b), nil
}
