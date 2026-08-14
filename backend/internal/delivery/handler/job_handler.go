package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"featherweight/internal/domain"
	"featherweight/internal/usecase"
)

type JobHandler struct {
	jobUseCase *usecase.JobUseCase
}

func NewJobHandler(jobUseCase *usecase.JobUseCase) *JobHandler {
	return &JobHandler{
		jobUseCase: jobUseCase,
	}
}

func (h *JobHandler) GetByID(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.jobUseCase.GetJob(c.Request.Context(), jobID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		}
		return
	}

	c.JSON(http.StatusOK, job)
}
