package main

import (
	"image"
	"log"
	"time"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/jacobsa/go-serial/serial"
	"github.com/kbinani/screenshot"
)

const (
	leftLEDCount   = 5
	topLEDCount    = 10
	rightLEDCount  = 5
	bottomLEDCount = 10
	maxLED         = leftLEDCount + topLEDCount + rightLEDCount + bottomLEDCount
)

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/cu.usbserial-141210",
		BaudRate:        250000,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	led := ledController{
		maxLED: maxLED,
		w:      port,
		buf:    make([]byte, maxLED*3),
	}
	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	for {
		imgL := c.Left(80, 80)
		processVertical(&led, imgL, leftLEDCount)
		imgT := c.Top(80, 80)
		processHorizontal(&led, imgT, topLEDCount)
		imgR := c.Right(80, 80)
		processVertical(&led, imgR, rightLEDCount)
		imgB := c.Bottom(80, 80)
		processHorizontal(&led, imgB, bottomLEDCount)
		time.Sleep(time.Second)
	}
}

func processVertical(c *ledController, img *image.RGBA, count int) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := height / count

	for i := 0; i < count; i++ {
		img := img.SubImage(image.Rect(0, div*i, width, div*(i+1)))
		c.SetColor(getProminentColor(img))
	}
}

func processHorizontal(c *ledController, img *image.RGBA, count int) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := width / count

	for i := 0; i < count; i++ {
		img := img.SubImage(image.Rect(div*i, 0, div*(i+1), height))
		c.SetColor(getProminentColor(img))
	}
}

var noBGMask []prominentcolor.ColorBackgroundMask

func getProminentColor(img image.Image) (r, g, b uint8) {
	xs, err := prominentcolor.KmeansWithAll(
		prominentcolor.DefaultK,
		img,
		prominentcolor.ArgumentNoCropping,
		prominentcolor.DefaultSize,
		noBGMask,
	)
	if err != nil {
		log.Fatalln(err)
		return 0, 0, 0
	}
	return uint8(xs[0].Color.R), uint8(xs[0].Color.G), uint8(xs[0].Color.B)
}
