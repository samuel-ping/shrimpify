package main

import (
	"fmt"
	"time"

	"github.com/h2non/bimg"
)

func processImage(imageBuffer []byte, filename string, quality int, dirname string) (string, error) {
	newFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)

	convertedImageBuffer, err := bimg.NewImage(imageBuffer).Convert(bimg.WEBP)
	if err != nil {
		return newFilename, err
	}

	compressedImageBuffer, err := bimg.NewImage(convertedImageBuffer).Process(bimg.Options{Quality: quality})
	if err != nil {
		return newFilename, err
	}

	fileRelativePath := fmt.Sprintf("%s/%s", dirname, newFilename)
	writeErr := bimg.Write(fileRelativePath, compressedImageBuffer)
	if err != nil {
		return newFilename, writeErr
	}

	return fileRelativePath, nil
}
