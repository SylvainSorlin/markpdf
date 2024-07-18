package main

import (
	"fmt"
	"image"
	imgColor "image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/unidoc/unidoc/pdf/core"
)

func debugInfo(message string) {
	if verbose {
		log.Println(message)
	}
}

func fatalIfError(err error, message string) {
	if err != nil {
		fmt.Printf("ERROR: %s \n", message)
		os.Exit(1)
	}
}

func parseAccessPermissions(permStr string) (core.AccessPermissions, error) {
	var permissions core.AccessPermissions

	// Delete the part before '{' and after '}'.
	start := strings.Index(permStr, "{")
	end := strings.LastIndex(permStr, "}")
	if start == -1 || end == -1 {
		return permissions, fmt.Errorf("incorrect string format")
	}

	permStr = permStr[start+1 : end]

	// Divide string into key/value pairs
	pairs := strings.Split(permStr, ", ")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return permissions, fmt.Errorf("incorrect key/value pair format: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1]) == "true"

		switch key {
		case "Printing":
			permissions.Printing = value
		case "Modify":
			permissions.Modify = value
		case "ExtractGraphics":
			permissions.ExtractGraphics = value
		case "Annotate":
			permissions.Annotate = value
		case "FillForms":
			permissions.FillForms = value
		case "DisabilityExtract":
			permissions.DisabilityExtract = value
		case "RotateInsert":
			permissions.RotateInsert = value
		case "FullPrintQuality":
			permissions.FullPrintQuality = value
		default:
			return permissions, fmt.Errorf("unknown key: %s", key)
		}
	}

	return permissions, nil
}

func textToPNG(watermark string) string {
	imgHeight := int(fontSize * 1.2)
	imgWidth := len(watermark) * int(fontSize)

	imgPath := filepath.Join(os.TempDir(), "Watermark.png")

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill image with transparent color
	draw.Draw(img, img.Bounds(), &image.Uniform{imgColor.Transparent}, image.Point{}, draw.Src)

	// Read the policy
	fontBytes, err := os.ReadFile("C:\\Windows\\Fonts\\arial.ttf")
	fatalIfError(err, fmt.Sprintf("Failed to read font. [%s]", err))

	f, err := freetype.ParseFont(fontBytes)
	fatalIfError(err, fmt.Sprintf("Failed to parse font. [%s]", err))

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	// Set initial text position
	pt := freetype.Pt(10, 10+int(c.PointToFixed(48)>>6))

	// Draw text on the image
	_, err = c.DrawString(watermark, pt)
	fatalIfError(err, fmt.Sprintf("Failed to draw string. [%s]", err))

	// Save image as PNG
	outFile, err := os.Create(imgPath)
	fatalIfError(err, fmt.Sprintf("Failed to create file. [%s]", err))
	defer outFile.Close()

	err = png.Encode(outFile, img)
	fatalIfError(err, fmt.Sprintf("Failed to encode image. [%s]", err))

	return imgPath
}
