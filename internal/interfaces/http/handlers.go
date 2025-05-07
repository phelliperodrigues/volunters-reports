package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phelliperodrigues/ccb-report/internal/domain"
)

type ReportHandler struct {
	reportService domain.ReportService
}

func NewReportHandler(reportService domain.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) GenerateReport(c *gin.Context) {
	localidade := c.Param("localidade")
	if localidade == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "localidade is required"})
		return
	}

	report, err := h.reportService.GenerateReport(c.Request.Context(), localidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	localidade := c.Param("localidade")
	if localidade == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "localidade is required"})
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), localidade)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetAllReports(c *gin.Context) {
	reports, err := h.reportService.GetAllReports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}
