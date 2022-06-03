package main

import (
	"github.com/h2non/bimg"
)

func processImage(imageBuffer []byte, quality int) ([]byte, error) {
	convertedImageBuffer, err := bimg.NewImage(imageBuffer).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	compressedImageBuffer, err := bimg.NewImage(convertedImageBuffer).Process(bimg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return compressedImageBuffer, nil
}
