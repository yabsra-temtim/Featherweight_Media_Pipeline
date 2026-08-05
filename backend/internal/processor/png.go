package processor

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"featherweight/internal/domain"
)

func writePNG(source image.Image, outputDirectory string) (domain.Output, error) {
	outputPath := filepath.Join(outputDirectory, "optimized.png")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("create PNG output: %w", err)
	}
	defer outputFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	err = encoder.Encode(outputFile, source)
	if err != nil {
		return domain.Output{}, fmt.Errorf("encode PNG: %w", err)
	}

	size, err := fileSize(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("read PNG size: %w", err)
	}

	bounds := source.Bounds()
	return domain.Output{
		Format: "png",
		Path:   "/downloads/" + filepath.Base(outputDirectory) + "/optimized.png",
		Size:   size,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}, nil
}
