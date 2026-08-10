package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"featherweight/internal/config"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "A file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	// 2. Parse width (default 0 means keep original size)
	width, _ := strconv.Atoi(c.DefaultPostForm("width", "0"))

	// 3. Parse quality (default 80)
	quality, _ := strconv.Atoi(c.DefaultPostForm("quality", "80"))

	// 4. Parse formats (e.g. "jpeg,webp")
	formatsStr := c.DefaultPostForm("formats", "jpeg")
	formats := strings.Split(formatsStr, ",")

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 6. Submit the Job to the background worker pool!
	h.workerPool.Submit(job.ID)

	// 7. Return the pending job details to the user
	c.JSON(http.StatusAccepted, job)
}
