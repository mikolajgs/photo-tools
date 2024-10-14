package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikolajgs/broccli"
	"github.com/nfnt/resize"
	"github.com/strukturag/libheif/go/heif"
)

func isSupportedImage(fileName string) bool {
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return true
		}
	}
	return false
}

func processImage(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening image file %s: %v", inputPath, err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(inputPath))

	var img image.Image

	if ext == ".heif" || ext == ".heic" {
		// Decode HEIF images
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("error reading HEIF file %s: %v", inputPath, err)
		}

		ctx, err := heif.NewContext()
		if err != nil {
			return fmt.Errorf("error creating HEIF context: %v", err)
		}
		if err := ctx.ReadFromMemory(data); err != nil {
			return fmt.Errorf("error reading HEIF data: %v", err)
		}
		imgHandle, err := ctx.GetPrimaryImageHandle()
		if err != nil {
			return fmt.Errorf("error getting primary HEIF image: %v", err)
		}
		heifImg, err := imgHandle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
		if err != nil {
			return fmt.Errorf("error decoding HEIF image: %v", err)
		}
		img, err = heifImg.GetImage()
		if err != nil {
			return fmt.Errorf("error getting image from HEIF image: %v", err)
		}
	} else {
		// Decode JPG/PNG images
		switch ext {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(file)
			if err != nil {
				return fmt.Errorf("error decoding JPEG file %s: %v", inputPath, err)
			}
		case ".png":
			img, err = png.Decode(file)
			if err != nil {
				return fmt.Errorf("error decoding PNG file %s: %v", inputPath, err)
			}
		}
	}

	// Resize the image while maintaining aspect ratio
	newImage := resize.Resize(200, 0, img, resize.Lanczos3)

	// Create output directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Save the resized image
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %v", outputPath, err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, newImage, nil)
	if err != nil {
		return fmt.Errorf("error saving resized image to %s: %v", outputPath, err)
	}

	fmt.Printf("Resized image saved to %s\n", outputPath)
	return nil
}

func createThumbsHandler(c *broccli.CLI) int {
	photos := c.Flag("photos")
	thumbs := c.Flag("thumbs")

	if err := findAndResizeImages(photos, thumbs); err != nil {
		fmt.Printf("Error creating photo thumbnails: %v\n", err)
	}

	return 0
}

func findAndResizeImages(dir, outputDir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a supported image file
		if !info.IsDir() && isSupportedImage(info.Name()) {
			// Construct output path
			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("error constructing relative path: %v", err)
			}
			outputPath := filepath.Join(outputDir, relativePath+"_th.jpg")

			if fileExists(outputPath) {
				fmt.Printf("Skipping %s, resized image already exists.\n", outputPath)
				return nil
			}

			// Process and resize the image
			if err := processImage(path, outputPath); err != nil {
				return fmt.Errorf("error processing image %s: %v", path, err)
			}
			log.Printf("%v -> %v", path, outputPath)
		}

		return nil
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
