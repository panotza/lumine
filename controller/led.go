package main

import (
	"io"
)

type ledController struct {
	maxLED uint8
	n      uint8
	w      io.Writer
	buf    []byte
}

func (c *ledController) SetColor(r, g, b uint8) {
	c.buf[c.n] = r
	c.buf[c.n+1] = g
	c.buf[c.n+2] = b
	c.n += 3

	if c.n >= c.maxLED*3-1 {
		c.w.Write(c.buf)
		c.n = 0
	}
}
