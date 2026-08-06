package processor

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"featherweight/internal/config"
	"featherweight/internal/domain"
)

func writeAVIF(
	ctx context.Context,
	cfg config.Config,
	inputPath string,
	outputDirectory string,
	width int,
	quality int,
) (domain.Output, error) {
	outputPath := filepath.Join(outputDirectory, "optimized.avif")
	crf := qualityToCRF(quality)

	filter := ""
	if width > 0 {
		filter = "scale=" + strconv.Itoa(width) + ":-2"
	}

	args := []string{
		"-y",
		"-i",
		inputPath,
	}

	if filter != "" {
		args = append(args, "-vf", filter)
	}

	args = append(
		args,
		"-c:v",
		"libaom-av1",
		"-still-picture",
		"1",
		"-crf",
		strconv.Itoa(crf),
		"-b:v",
		"0",
		outputPath,
	)

	_, err := runCommand(ctx, cfg.CommandTimeout, "ffmpeg", args...)
	if err != nil {
		return domain.Output{}, err
	}

	size, err := fileSize(outputPath)
	if err != nil {
		return domain.Output{}, fmt.Errorf("read AVIF size: %w", err)
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
		Format: "avif",
		Path:   "/downloads/" + filepath.Base(outputDirectory) + "/optimized.avif",
		Size:   size,
		Width:  outWidth,
		Height: outHeight,
	}, nil
}

func runCommand(ctx context.Context, timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("run command %s failed: %w\n%s", name, err, string(output))
	}

	return output, nil
}

func fileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func qualityToCRF(quality int) int {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	// Maps quality (1-100) down to CRF scale (63-0)
	return 63 - ((quality - 1) * 63 / 99)
}
