package main

import (
	"io"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type ledController struct {
	LEDCount uint8
	n        uint8
	w        io.ReadWriteCloser
	buf      []byte
}

func newController(portName string, totalLED uint8) (*ledController, error) {
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        250000,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	return &ledController{
		LEDCount: totalLED,
		w:        port,
		buf:      make([]byte, totalLED*3),
	}, nil
}

func (c *ledController) WriteColor(r, g, b uint8) {
	c.buf[c.n] = r
	c.buf[c.n+1] = g
	c.buf[c.n+2] = b
	c.n += 3

	if c.n >= c.LEDCount*3-1 {
		_, err := c.w.Write(c.buf)
		if err != nil {
			log.Fatal(err)
		}
		c.n = 0
	}
}

func (c *ledController) Close() {
	for i := 0; i < totalLED; i++ {
		c.WriteColor(0, 0, 0)
	}
	time.Sleep(10 * time.Millisecond)
	err := c.w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
