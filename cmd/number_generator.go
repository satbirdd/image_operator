package main

import (
	"fmt"

	"github.com/fogleman/gg"
)

func main() {
	var i int

	for i = 1; i < 10; i++ {
		drawNumber(i)
	}
}

func drawNumber(i int) {
	const S = 200
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0.8, 0, 0)
	if err := dc.LoadFontFace("src/fonts/simsun.ttf", 96); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprint(i), S/2, S/2, 0.5, 0.5)
	dc.SetRGB(0.8, 0, 0)
	dc.DrawLine(0, S-1, S-1, S-1)
	dc.DrawLine(S-1, S-1, S-1, 0)
	dc.Stroke()
	err := dc.SaveJPG(fmt.Sprintf("src/images/%v.jpg", i), 70)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
