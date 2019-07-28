package main

import (
	"io"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type controller struct {
	LEDCount uint8
	w        io.ReadWriteCloser
	buf      []byte
}

func newController(portName string, totalLED uint8) (*controller, error) {
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

	return &controller{
		LEDCount: totalLED,
		w:        port,
		buf:      make([]byte, 3),
	}, nil
}

func (c *controller) WriteColor(r, g, b uint8) {
	c.buf[0] = r
	c.buf[1] = g
	c.buf[2] = b
	c.w.Write(c.buf)
}

func (c *controller) Close() {
	defer func() {
		if err := c.w.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for i := 0; i < totalLED; i++ {
		c.WriteColor(0, 0, 0)
		time.Sleep(time.Millisecond)
	}
}
