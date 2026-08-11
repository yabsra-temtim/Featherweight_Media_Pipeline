package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
