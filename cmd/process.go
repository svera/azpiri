package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const backgroundsFolder = "/media/screenshots/"
const foregroundsFolder = "/media/marquees/"
const textContent = "Loading..."
const textSize = 24
const textFont = "Arcade"
const backgroundBrightness = -20
const backgroundKeepAspectRatio = true
const foregroundScale = 2
const targetWidth = 1024
const targetHeight = 768

type config struct {
	BackgroundsFolder         string
	ForegroundsFolder         string
	TextContent               string
	TextSize                  int
	TextFont                  string
	BackgroundBrightness      int
	BackgroundKeepAspectRatio bool
	ForegroundScale           int
	TargetWidth               int
	TargetHeight              int
}

func process() {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	imagick.Initialize()
	defer imagick.Terminate()
	background := imagick.NewMagickWand()
	foreground := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	bar := pb.StartNew(len(files))

	for _, fileInfo := range files {
		if fileInfo.IsDir() || fileInfo.Name()[0:1] == "." {
			bar.Increment()
			continue
		}
		imageName := folder + backgroundsFolder + fileNameWithoutExtension(fileInfo.Name()) + ".png"
		file, err := os.Open(imageName)
		if err != nil {
			log.Printf("Error loading image %s, skipping generation of launching image for %s\n", imageName, fileInfo.Name())
			bar.Increment()
			continue
		}
		background.ReadImageFile(file)

		imageName = folder + foregroundsFolder + fileNameWithoutExtension(fileInfo.Name()) + ".png"
		file, err = os.Open(imageName)
		if err != nil {
			log.Printf("Error loading marquee %s, skipping generation of launching image for %s\n", imageName, fileInfo.Name())
			bar.Increment()
			continue
		}
		foreground.ReadImageFile(file)
		foreground.ResizeImage(foreground.GetImageWidth()*foregroundScale, foreground.GetImageHeight()*foregroundScale, imagick.FILTER_CUBIC, 1)

		if backgroundKeepAspectRatio {
			w := int(background.GetImageWidth())
			h := int(background.GetImageHeight())
			pw.SetColor("black")
			background.SetImageBackgroundColor(pw)
			w = w * targetHeight / h
			h = targetHeight
			background.ScaleImage(uint(w), uint(h))
			background.ExtentImage(targetWidth, targetHeight, -(targetWidth-w)/2, -(targetHeight-h)/2)
		} else {
			background.ScaleImage(targetWidth, targetHeight)
		}

		background.BlurImage(50, 5)
		background.BrightnessContrastImage(backgroundBrightness, 0)

		writeText(dw, pw, foreground)
		// Draw the image on to the background
		background.DrawImage(dw)

		x, y := center(background, foreground)
		background.CompositeImage(foreground, imagick.COMPOSITE_OP_OVER, x, y)
		background.WriteImage(fileNameWithoutExtension(fileInfo.Name()) + "-launching.png")
		dw.Clear()
		background.Clear()
		foreground.Clear()

		bar.Increment()
	}
	bar.Finish()
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func center(background *imagick.MagickWand, foreground *imagick.MagickWand) (int, int) {
	x := int(background.GetImageWidth()/2) - int(foreground.GetImageWidth()/2)
	y := int(background.GetImageHeight()/2) - int(foreground.GetImageHeight()/2)
	return x, y
}

func writeText(dw *imagick.DrawingWand, pw *imagick.PixelWand, foreground *imagick.MagickWand) {
	pw.SetColor("white")
	dw.SetFillColor(pw)
	if err := dw.SetFont(textFont); err != nil {
		log.Printf("Error using font '%s'\n", textFont)
	}
	dw.SetFontSize(textSize)
	// Now draw the text
	dw.SetGravity(imagick.GRAVITY_CENTER)
	dw.Annotation(0, float64((foreground.GetImageHeight()/2)+textSize), textContent)
}
