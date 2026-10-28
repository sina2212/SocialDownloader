package api

import (
	"back-end/internal/downloader"
	"back-end/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Downloader *downloader.Downloader
}

func NewHandler(dl *downloader.Downloader) *Handler {
	return &Handler{Downloader: dl}
}

func (h *Handler) Extract(c *gin.Context) {
	var req models.DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := h.Downloader.ExtractInfo(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *Handler) Download(c *gin.Context) {
	var req models.DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quality := req.Quality
	if quality == "" {
		quality = "best"
	}

	filePath, err := h.Downloader.Download(c.Request.Context(), req.URL, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.DownloadResponse{
		Status:   "completed",
		Message:  "Download finished",
		FilePath: filePath,
	})
}