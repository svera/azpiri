package cmd

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/disintegration/imaging"
)

type config struct {
	BackgroundBrightness      float64
	BackgroundBlur            float64
	BackgroundKeepAspectRatio bool
	ForegroundWidth           int
	TargetWidth               int
	TargetHeight              int
}

func process() {
	files, err := ioutil.ReadDir(romsFolder)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	filtered := filterFiles(files)
	bar := pb.StartNew(len(filtered))

	outputFolder := romsFolder + "/images/"
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, 0755)
	}

	for _, fileInfo := range filtered {
		imageName := backgroundsFolder + "/" + fileNameWithoutExtension(fileInfo.Name()) + ".png"
		background, err := imaging.Open(imageName)
		if err != nil {
			log.Printf("Error loading background image %s, skipping generation of launching image for %s\n", imageName, fileInfo.Name())
			bar.Increment()
			continue
		}

		imageName = foregroundsFolder + "/" + fileNameWithoutExtension(fileInfo.Name()) + ".png"
		foreground, err := imaging.Open(imageName)
		if err != nil {
			log.Printf("Error loading foreground image %s, skipping generation of launching image for %s\n", imageName, fileInfo.Name())
			bar.Increment()
			continue
		}

		dst := composeImage(background, foreground, cfg)
		err = imaging.Save(dst, outputFolder+fileNameWithoutExtension(fileInfo.Name())+"-launching.png")
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}

		bar.Increment()
	}
	bar.Finish()
}

func filterFiles(files []os.FileInfo) []os.FileInfo {
	filtered := []os.FileInfo{}
	for _, fileInfo := range files {
		if filepath.Ext(fileInfo.Name()) != ".cfg" && filepath.Ext(fileInfo.Name()) != ".bsv" && !fileInfo.IsDir() && fileInfo.Name()[0:1] != "." {
			filtered = append(filtered, fileInfo)
		}
	}
	return filtered
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func composeImage(background image.Image, foreground image.Image, cfg config) *image.NRGBA {
	foreground = imaging.Resize(foreground, cfg.ForegroundWidth, 0, imaging.Lanczos)

	canvas := imaging.New(cfg.TargetWidth, cfg.TargetHeight, color.Black)
	background = imaging.Resize(background, cfg.TargetWidth, 0, imaging.Lanczos)
	background = imaging.PasteCenter(canvas, background)

	background = imaging.Blur(background, cfg.BackgroundBlur)
	background = imaging.AdjustBrightness(background, cfg.BackgroundBrightness)

	return imaging.OverlayCenter(background, foreground, 1)
}

func loadConfig() (config, error) {
	cfg := config{}
	f, err := os.Open("azpiri.json")
	defer f.Close()
	if err != nil {
		return cfg, err
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(byteValue, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
