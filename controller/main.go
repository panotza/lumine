package main

import (
	"image"
	"log"
	"sync"
	"time"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/jacobsa/go-serial/serial"
	"github.com/kbinani/screenshot"
)

const (
	leftLEDCount   = 10
	topLEDCount    = 20
	rightLEDCount  = 10
	bottomLEDCount = 20
	maxLED         = leftLEDCount + topLEDCount + rightLEDCount + bottomLEDCount
)

type color struct {
	r, g, b uint8
}

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/cu.usbserial-14120",
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

	sig := make(chan struct{}, 4)
	leftRet := make(chan []color)
	topRet := make(chan []color)
	rightRet := make(chan []color)
	bottomRet := make(chan []color)
	var wg sync.WaitGroup

	go leftWorker(&wg, sig, leftRet)
	go topWorker(&wg, sig, topRet)
	go rightWorker(&wg, sig, rightRet)
	go bottomWorker(&wg, sig, bottomRet)

	writeColor := func(cs []color) {
		for _, c := range cs {
			led.WriteColor(c.r, c.b, c.b)
		}
	}

	for {
		sig <- struct{}{}
		sig <- struct{}{}
		sig <- struct{}{}
		sig <- struct{}{}

		wg.Wait()
		writeColor(<-leftRet)
		writeColor(<-topRet)
		writeColor(<-rightRet)
		writeColor(<-bottomRet)
		time.Sleep(100 * time.Millisecond)
	}
}

func leftWorker(wg *sync.WaitGroup, sig chan struct{}, ret chan []color) {
	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	for {
		select {
		case <-sig:
			wg.Add(1)
			img := c.Left(100, 80)
			ret <- processVertical(img, leftLEDCount, true)
			wg.Done()
		}
	}
}

func topWorker(wg *sync.WaitGroup, sig chan struct{}, ret chan []color) {
	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	for {
		select {
		case <-sig:
			wg.Add(1)
			img := c.Top(100, 80)
			ret <- processHorizontal(img, topLEDCount, false)
			wg.Done()
		}
	}
}

func rightWorker(wg *sync.WaitGroup, sig chan struct{}, ret chan []color) {
	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	for {
		select {
		case <-sig:
			wg.Add(1)
			img := c.Right(100, 80)
			ret <- processVertical(img, rightLEDCount, false)
			wg.Done()
		}
	}
}

func bottomWorker(wg *sync.WaitGroup, sig chan struct{}, ret chan []color) {
	bounds := screenshot.GetDisplayBounds(0)
	c := capture{bounds.Dx(), bounds.Dy()}

	for {
		select {
		case <-sig:
			wg.Add(1)
			img := c.Bottom(100, 80)
			ret <- processHorizontal(img, bottomLEDCount, true)
			wg.Done()
		}
	}
}

func processVertical(img *image.RGBA, count int, swap bool) []color {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := height / count

	cs := make([]color, count)
	for i := 0; i < count; i++ {
		n := i
		if swap {
			n = count - i - 1
		}
		img := img.SubImage(image.Rect(0, div*n, width, div*(n+1)))
		r, g, b := getProminentColor(img)
		cs[i] = color{r, g, b}
	}
	return cs
}

func processHorizontal(img *image.RGBA, count int, swap bool) []color {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	div := width / count

	cs := make([]color, count)
	for i := 0; i < count; i++ {
		n := i
		if swap {
			n = count - i - 1
		}
		img := img.SubImage(image.Rect(div*n, 0, div*(n+1), height))
		r, g, b := getProminentColor(img)
		cs[i] = color{r, g, b}
	}
	return cs
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
