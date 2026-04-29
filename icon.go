package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

func buildAppIcon() ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	bgColor := color.RGBA{R: 99, G: 102, B: 241, A: 255}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, bgColor)
		}
	}

	cursorColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for i := 0; i < 20; i++ {
		img.Set(20+i, 22+i, cursorColor)
	}
	for i := 0; i < 12; i++ {
		img.Set(20+i, 22+i+8, cursorColor)
	}
	for i := 0; i < 8; i++ {
		img.Set(20, 22+i, cursorColor)
		img.Set(20+i, 22, cursorColor)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
