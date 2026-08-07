package processor

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"featherweight/internal/domain"
)

func writeJPEG(source image.Image, outputDirectory string, quality int) (domain.Output, error) {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}

	outputPath := filepath.Join(outputDirectory, "optimized.jpg")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("create JPEG output: %w", err)
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, source, &jpeg.Options{Quality: quality})
	if err != nil {
		return domain.Output{}, fmt.Errorf("encode JPEG: %w", err)
	}

	size, err := fileSize(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("read JPEG size: %w", err)
	}

	bounds := source.Bounds()
	return domain.Output{
		Format: "jpeg",
		Path:   "/downloads/" + filepath.Base(outputDirectory) + "/optimized.jpg",
		Size:   size,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}, nil
}
