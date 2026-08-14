package processor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

	outWidth, outHeight, err := getMediaDimensions(
		ctx,
		cfg,
		outputPath,
	)
	if err != nil {
		return domain.Output{}, fmt.Errorf(
			"read AVIF dimensions: %w",
			err,
		)
	}
	return domain.Output{
		Format: "avif",
		Path:   outputPath,
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

func getMediaDimensions(
	ctx context.Context,
	cfg config.Config,
	path string,
) (int, int, error) {
	output, err := runCommand(
		ctx,
		cfg.CommandTimeout,
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		path,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("probe media dimensions: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(output)), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dimensions returned by ffprobe: %q", output)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parse width: %w", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse height: %w", err)
	}

	return width, height, nil
}
