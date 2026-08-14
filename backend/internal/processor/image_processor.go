package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"featherweight/internal/config"
	"featherweight/internal/domain"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

var (
	ErrUnsupportedFormat = errors.New(
		"unsupported output format",
	)

	ErrInvalidImage = errors.New(
		"the uploaded file is not a valid supported image",
	)
)

type ProcessInput struct {
	JobID string

	OriginalPath string

	OutputDirectory string

	Width int

	Quality int

	Formats []string
}

type ImageProcessor struct {
	config config.Config

	webpEnabled bool

	avifEnabled bool
}

func NewImageProcessor(
	cfg config.Config,
) *ImageProcessor {
	_, webpError := exec.LookPath(
		"cwebp",
	)

	_, avifError := exec.LookPath(
		"ffmpeg",
	)

	return &ImageProcessor{
		config: cfg,

		webpEnabled: webpError == nil,

		avifEnabled: avifError == nil,
	}

}

func (p *ImageProcessor) WebPEnabled() bool {
	return p.webpEnabled
}

func (p *ImageProcessor) AVIFEnabled() bool {
	return p.avifEnabled
}

func (p *ImageProcessor) Process(
	ctx context.Context,
	input ProcessInput,
) (
	[]domain.Output,
	[]string,
	error,
) {
	imageData, err := os.ReadFile(
		input.OriginalPath,
	)

	if err != nil {
		return nil,
			nil,
			fmt.Errorf(
				"read original image: %w",
				err,
			)
	}

	sourceImage, err := decodeImage(
		imageData,
	)

	if err != nil {
		return nil,
			nil,
			err
	}

	resizedImage := resizeImage(
		sourceImage,
		input.Width,
	)

	jobDirectory := filepath.Join(
		input.OutputDirectory,
		input.JobID,
	)

	if err := os.MkdirAll(
		jobDirectory,
		0755,
	); err != nil {
		return nil,
			nil,
			fmt.Errorf(
				"create output directory: %w",
				err,
			)
	}

	outputs := make(
		[]domain.Output,
		0,
	)

	notes := make(
		[]string,
		0,
	)

	for _, requestedFormat := range input.Formats {
		format := normalizeFormat(
			requestedFormat,
		)

		switch format {
		case "jpeg":
			output, err :=
				writeJPEG(
					resizedImage,
					jobDirectory,
					input.Quality,
				)

			if err != nil {
				return nil,
					nil,
					err
			}

			outputs = append(
				outputs,
				output,
			)

		case "png":
			output, err :=
				writePNG(
					resizedImage,
					jobDirectory,
				)

			if err != nil {
				return nil,
					nil,
					err
			}

			outputs = append(
				outputs,
				output,
			)

		case "webp":
			if !p.webpEnabled {
				notes = append(
					notes,
					"WebP was skipped because cwebp is not installed",
				)

				continue
			}

			output, err :=
				writeWebP(
					ctx,
					p.config,
					input.OriginalPath,
					jobDirectory,
					input.Width,
					input.Quality,
				)

			if err != nil {
				notes = append(
					notes,
					fmt.Sprintf(
						"WebP conversion failed: %v",
						err,
					),
				)

				continue
			}

			outputs = append(
				outputs,
				output,
			)

		case "avif":
			if !p.avifEnabled {
				notes = append(
					notes,
					"AVIF was skipped because ffmpeg is not installed",
				)

				continue
			}

			output, err :=
				writeAVIF(
					ctx,
					p.config,
					input.OriginalPath,
					jobDirectory,
					input.Width,
					input.Quality,
				)

			if err != nil {
				notes = append(
					notes,
					fmt.Sprintf(
						"AVIF conversion failed: %v",
						err,
					),
				)

				continue
			}

			outputs = append(
				outputs,
				output,
			)

		default:
			notes = append(
				notes,
				fmt.Sprintf(
					"%s is not a supported output format",
					requestedFormat,
				),
			)
		}
	}

	if len(outputs) == 0 {
		return nil,
			notes,
			errors.New(
				"no output files were created",
			)
	}

	return outputs,
		notes,
		nil

}

func decodeImage(
	data []byte,
) (image.Image, error) {
	if decodedImage, _, err :=
		image.Decode(
			bytes.NewReader(data),
		); err == nil {
		return decodedImage,
			nil
	}

	if decodedImage, err :=
		bmp.Decode(
			bytes.NewReader(data),
		); err == nil {
		return decodedImage,
			nil
	}

	if decodedImage, err :=
		tiff.Decode(
			bytes.NewReader(data),
		); err == nil {
		return decodedImage,
			nil
	}

	if decodedImage, err :=
		webp.Decode(
			bytes.NewReader(data),
		); err == nil {
		return decodedImage,
			nil
	}

	return nil,
		ErrInvalidImage

}

func resizeImage(
	source image.Image,
	targetWidth int,
) image.Image {
	if targetWidth <= 0 {
		return source
	}

	currentWidth :=
		source.Bounds().Dx()

	if currentWidth <= targetWidth {
		return source
	}

	return imaging.Resize(
		source,
		targetWidth,
		0,
		imaging.Lanczos,
	)

}

func normalizeFormat(
	format string,
) string {
	format = strings.ToLower(
		strings.TrimSpace(
			format,
		),
	)

	switch format {
	case "jpg":
		return "jpeg"

	case "tif":
		return "tiff"

	default:
		return format
	}

}

func copyToWriter(
	destination io.Writer,
	source io.Reader,
) error {
	_, err := io.Copy(
		destination,
		source,
	)

	return err

}
