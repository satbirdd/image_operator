package main

import (
	"fmt"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(1024, 860)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	err := dc.SaveJPG("src/images/blank.jpg", 70)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
