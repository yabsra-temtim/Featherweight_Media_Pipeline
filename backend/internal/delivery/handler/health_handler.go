package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"featherweight/internal/processor"
)

type HealthHandler struct {
	imageProcessor *processor.ImageProcessor
}

func NewHealthHandler(imageProcessor *processor.ImageProcessor) *HealthHandler {
	return &HealthHandler{imageProcessor: imageProcessor}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"webp_enabled": h.imageProcessor.WebPEnabled(),
		"avif_enabled": h.imageProcessor.AVIFEnabled(),
	})
}
