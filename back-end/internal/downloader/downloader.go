package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"back-end/internal/models"

	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
)

type Downloader struct {
	DownloadDir string
}

func New(downloadDir string) *Downloader {
	//Ckeck Existing Downlder Directory
	_ = os.MkdirAll(downloadDir, 0755)
	return &Downloader{DownloadDir: downloadDir}
}

// ExtractInfo اطلاعات پست رو بدون دانلود می‌گیره
func (d *Downloader) ExtractInfo(ctx context.Context, url string) (*models.VideoInfo, error) {
	dl := ytdlp.New().
		DumpSingleJSON().
		NoPlaylist().
		NoWarnings()

	result, err := dl.Run(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}

	// پارس کردن خروجی JSON
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(result.Stdout), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	info := &models.VideoInfo{
		ID:         getString(raw, "id"),
		Title:      getString(raw, "title"),
		Uploader:   getString(raw, "uploader"),
		Thumbnail:  getString(raw, "thumbnail"),
		Description: getString(raw, "description"),
		WebpageURL: getString(raw, "webpage_url"),
	}

	if duration, ok := raw["duration"].(float64); ok {
		info.Duration = duration
	}

	return info, nil
}

// کمک‌کننده برای گرفتن مقدار string از map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// Download File
func (d *Downloader) Download(ctx context.Context, url string, quality string) (string, error) {
	jobID := uuid.New().String()
	outputTemplate := filepath.Join(d.DownloadDir, jobID+"_%(title).50B.%(ext)s")

	dl := ytdlp.New().
		NoPlaylist().
		NoOverwrites().
		Continue().
		Output(outputTemplate).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4")

	if quality != "" && quality != "best" {
		dl = dl.Format(quality)
	} else {
		dl = dl.Format("bestvideo+bestaudio/best")
	}

	_, err := dl.Run(ctx, url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	// Find the downloaded file
	files, _ := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if len(files) == 0 {
		return "", fmt.Errorf("downloaded file not found")
	}

	return files[0], nil
}

// DownloadWithProgress
func (d *Downloader) DownloadWithProgress(ctx context.Context, url string, quality string, progressCb func(percent float64, status string)) (string, error) {
	jobID := uuid.New().String()
	outputTemplate := filepath.Join(d.DownloadDir, jobID+"_%(title).50B.%(ext)s")

	dl := ytdlp.New().
		NoPlaylist().
		NoOverwrites().
		Continue().
		Output(outputTemplate).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		ProgressFunc(500*time.Millisecond, func(prog ytdlp.ProgressUpdate) {
			if progressCb != nil {
				percent := 0.0
				if prog.PercentString() != "" {
					percent = 50
				}
				progressCb(percent, string(prog.Status))
			}
		})

	if quality != "" && quality != "best" {
		dl = dl.Format(quality)
	} else {
		dl = dl.Format("bestvideo+bestaudio/best")
	}

	_, err := dl.Run(ctx, url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	files, _ := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if len(files) == 0 {
		return "", fmt.Errorf("downloaded file not found")
	}

	return files[0], nil
}