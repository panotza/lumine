package main

import (
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/kbinani/screenshot"
)

const (
	leftLEDCount   = 10
	topLEDCount    = 20
	rightLEDCount  = 10
	bottomLEDCount = 20
	totalLED       = leftLEDCount + topLEDCount + rightLEDCount + bottomLEDCount
)

func main() {
	ctrl, err := newController("/dev/cu.usbserial-14210", totalLED)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer ctrl.Close()

	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigterm:
			return
		default:
			imgL := c.Left(90, 50)
			processVertical(ctrl, imgL, leftLEDCount, true)
			imgT := c.Top(100, 50)
			processHorizontal(ctrl, imgT, topLEDCount, false)
			imgR := c.Right(90, 50)
			processVertical(ctrl, imgR, rightLEDCount, false)
			imgB := c.Bottom(100, 50)
			processHorizontal(ctrl, imgB, bottomLEDCount, true)
		}
	}
}

func processVertical(lc *controller, img *image.RGBA, count int, swap bool) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := height / count

	for i := 0; i < count; i++ {
		n := i
		if swap {
			n = count - i - 1
		}
		img := img.SubImage(image.Rect(0, div*n, width, div*(n+1)))
		lc.WriteColor(getProminentColor(img))
	}
}

func processHorizontal(lc *controller, img *image.RGBA, count int, swap bool) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := width / count

	for i := 0; i < count; i++ {
		n := i
		if swap {
			n = count - i - 1
		}
		img := img.SubImage(image.Rect(div*n, 0, div*(n+1), height))
		lc.WriteColor(getProminentColor(img))
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
