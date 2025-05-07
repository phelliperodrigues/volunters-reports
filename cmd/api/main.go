package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phelliperodrigues/ccb-report/internal/application"
	"github.com/phelliperodrigues/ccb-report/internal/infrastructure"
	"github.com/phelliperodrigues/ccb-report/internal/interfaces/http"
	"github.com/phelliperodrigues/ccb-report/pkg/config"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize repository
	repo := infrastructure.NewFileRepository(cfg.ReportsDirPath)

	// Initialize PDF generator
	pdfGenerator := infrastructure.NewPDFGenerator(cfg.PDFOutputDirPath)

	// Initialize service
	reportService := application.NewReportService(repo, pdfGenerator, cfg)

	// Initialize handler
	reportHandler := http.NewReportHandler(reportService)

	// Setup router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		api.POST("/reports/:localidade", reportHandler.GenerateReport)
		api.GET("/reports/:localidade", reportHandler.GetReport)
		api.GET("/reports", reportHandler.GetAllReports)
	}

	// Start server
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
