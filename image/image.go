package image

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/graphics-go/graphics"
)

// DoImageHandler 产生缩略图,等比例缩放
func DoImageHandler(url string, newdx int) {
	src, err := LoadImage("." + url)
	//bound := src.Bounds()
	//dx := bound.Dx()
	//dy := bound.Dy()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(dx, dy, newdx)
	// 缩略图的大小
	dst := image.NewRGBA(image.Rect(640, 640, 200, 200))
	//dst := image.NewRGBA(image.Rect(0, 0, newdx, newdx*dy/dx))
	// 产生缩略图,等比例缩放
	//err = graphics.Scale(dst, src)
	err = graphics.Thumbnail(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	filen := strings.Replace(url, ".", "-cropper.", -1)
	file, err := os.Create("." + filen)
	defer file.Close()

	err = jpeg.Encode(file, dst, &jpeg.Options{Quality: 100}) //图像质量值为100，是最好的图像显示

	//header := w.Header()
	//header.Add("Content-Type", "image/jpeg")

	//png.Encode(w, dst)
}

// LoadImage 解码 an image from a file of image.
func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}
