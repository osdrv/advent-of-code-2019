package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Image struct {
	layers []*Layer
}

func NewImage(layers []*Layer) *Image {
	return &Image{
		layers: layers,
	}
}

func (img *Image) ColorAt(x, y int) byte {
	for _, l := range img.layers {
		if c := l.ColorAt(x, y); c < 2 {
			return c
		}
	}
	return 2
}

func (img *Image) ToPNG(dst io.Writer) error {
	w, h := img.layers[0].w, img.layers[0].h
	pane := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			var c color.Color
			switch img.ColorAt(i, j) {
			case 0:
				c = color.Black
			case 1:
				c = color.White
			case 2:
				c = color.Transparent
			}
			pane.Set(i, j, c)
		}
	}
	if err := png.Encode(dst, pane); err != nil {
		return err
	}
	return nil
}

type Layer struct {
	data []byte
	w, h int
}

func NewLayer(data []byte, w, h int) *Layer {
	return &Layer{
		data: data,
		w:    w,
		h:    h,
	}
}

func (l *Layer) ColorAt(x, y int) byte {
	if x >= l.w || y >= l.h {
		log.Fatalf("position out of range: x: %d, y: %d", x, y)
	}
	return l.data[l.w*y+x]
}

func ReadLayers(reader *bytes.Reader, w, h int) []*Layer {
	res := make([]*Layer, 0, 1)
Layers:
	for {
		data := make([]byte, 0, w*h)
		for i := 0; i < w*h; i++ {
			b, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					break Layers
				}
				log.Fatalf("Failed to read images: %s", err)
			}
			data = append(data, b-'0')
		}
		res = append(res, NewLayer(data, w, h))
	}
	return res
}

func CountDigits(l *Layer, dig byte) int {
	res := 0
	for _, d := range l.data {
		if d == dig {
			res++
		}
	}
	return res
}

func main() {
	file, err := os.Open("INPUT")
	if err != nil {
		log.Fatalf("failed to open input file: %s", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read input data: %s", err)
	}
	layers := ReadLayers(bytes.NewReader(data), 25, 6)
	var minLayer *Layer
	min0 := 999999
	for _, l := range layers {
		cnt0 := CountDigits(l, 0)
		if cnt0 < min0 {
			min0 = cnt0
			minLayer = l
		}
	}
	log.Printf("Min layer: %+v", minLayer)
	log.Printf("Result: %d", CountDigits(minLayer, 1)*CountDigits(minLayer, 2))

	img := NewImage(layers)
	dst, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Failed to create result file: %s", err)
	}
	if err := img.ToPNG(dst); err != nil {
		log.Fatalf("Failed to write image: %s", err)
	}
}
