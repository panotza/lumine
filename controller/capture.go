package main

import (
	"image"
	"log"

	"github.com/kbinani/screenshot"
)

type capture struct {
	width  int
	height int
}

func (d *capture) Left(percent int, size int) *image.RGBA {
	height := d.height * percent / 100
	offset := (d.height - height) / 2

	return captureRect(image.Rect(0, offset, size, d.height-offset))
}

func (d *capture) Right(percent int, size int) *image.RGBA {
	height := d.height * percent / 100
	offset := (d.height - height) / 2

	return captureRect(image.Rect(d.width-size, offset, d.width, d.height-offset))
}

func (d *capture) Top(percent int, size int) *image.RGBA {
	width := d.width * percent / 100
	offset := (d.width - width) / 2

	return captureRect(image.Rect(offset, 0, d.width-offset, size))
}

func (d *capture) Bottom(percent int, size int) *image.RGBA {
	width := d.width * percent / 100
	offset := (d.width - width) / 2

	return captureRect(image.Rect(offset, d.height-size, d.width-offset, d.height))
}

func captureRect(rect image.Rectangle) *image.RGBA {
	img, err := screenshot.CaptureRect(rect)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
