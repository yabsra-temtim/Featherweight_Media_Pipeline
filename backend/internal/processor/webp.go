package processor

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"

	"featherweight/internal/config"
	"featherweight/internal/domain"
)

func writeWebP(
	ctx context.Context,
	cfg config.Config,
	inputPath string,
	outputDirectory string,
	width int,
	quality int,
) (domain.Output, error) {
	outputPath := filepath.Join(outputDirectory, "optimized.webp")

	args := []string{
		"-quiet",
		"-q",
		strconv.Itoa(quality),
	}

	if width > 0 {
		args = append(args, "-resize", strconv.Itoa(width), "0")
	}

	args = append(args, inputPath, "-o", outputPath)

	_, err := runCommand(ctx, cfg.CommandTimeout, "cwebp", args...)
	if err != nil {
		return domain.Output{}, err
	}

	size, err := fileSize(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("read WebP size: %w", err)
	}

	var outWidth, outHeight int
	if file, err := os.Open(outputPath); err == nil {
		if img, _, err := image.DecodeConfig(file); err == nil {
			outWidth = img.Width
			outHeight = img.Height
		}
		file.Close()
	}

	return domain.Output{
		Format: "webp",
		Path:   outputPath,
		Size:   size,
		Width:  outWidth,
		Height: outHeight,
	}, nil
}
