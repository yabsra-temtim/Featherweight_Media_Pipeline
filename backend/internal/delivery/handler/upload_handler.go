package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"featherweight/internal/config"
	"featherweight/internal/delivery/dto"
	"featherweight/internal/usecase"
	"featherweight/internal/worker"
)

type UploadHandler struct {
	config       config.Config
	mediaUseCase *usecase.MediaUseCase
	workerPool   *worker.Pool
}

func NewUploadHandler(cfg config.Config, mediaUseCase *usecase.MediaUseCase, workerPool *worker.Pool) *UploadHandler {
	return &UploadHandler{
		config:       cfg,
		mediaUseCase: mediaUseCase,
		workerPool:   workerPool,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	// 1. Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("A file is required"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to read file"))
		return
	}
	defer file.Close()

	// 2. Parse width (default 0 means keep original size)
	width, _ := strconv.Atoi(c.DefaultPostForm("width", "0"))

	// 3. Parse quality (default 80)
	quality, _ := strconv.Atoi(c.DefaultPostForm("quality", "80"))

	// 4. Parse formats (frontend sends multiple separate fields, e.g. formats=jpeg&formats=webp)
	formats := c.PostFormArray("formats")
	if len(formats) == 0 {
		formats = []string{
			"jpeg",
			"webp",
			"avif",
			"png",
		}
	}

	// 5. Create the Job
	input := usecase.UploadInput{
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		File:     file,
		Width:    width,
		Quality:  quality,
		Formats:  formats,
	}

	job, err := h.mediaUseCase.CreateJob(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	// 6. Submit the Job to the background worker pool
	if !h.workerPool.Submit(job.ID) {
		c.JSON(http.StatusServiceUnavailable, dto.NewErrorResponse("Server is busy, please try again shortly"))
		return
	}

	// 7. Return the pending job details to the user
	c.JSON(http.StatusAccepted, dto.FromJob(job))
}
