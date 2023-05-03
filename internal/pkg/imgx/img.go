package imgx

import (
	"image"
	"image/draw"
	"os"
)

func Jpg2RGBA(img image.Image) *image.RGBA {
	tmp := image.NewRGBA(img.Bounds())

	draw.Draw(tmp, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return tmp
}

// ReadPic 把样本图片变成image.Image
func ReadPic(path string) image.Image {
	fr, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	img, _, err := image.Decode(fr)
	if err != nil {
		panic(err)
	}

	return img
}
