package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"back-end/internal/models"

	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
)

type Downloader struct {
	DownloadDir string
}

func New(downloadDir string) *Downloader {
	_ = os.MkdirAll(downloadDir, 0755)
	return &Downloader{DownloadDir: downloadDir}
}

// ExtractInfo
func (d *Downloader) ExtractInfo(ctx context.Context, url string) (*models.MediaInfo, error) {
	dl := ytdlp.New().
		DumpSingleJSON().
		NoPlaylist().
		NoWarnings().
		SkipDownload().
		IgnoreNoFormatsError()

	result, err := dl.Run(ctx, url)
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "There is no video") &&
			!strings.Contains(errStr, "No video formats found") &&
			!strings.Contains(errStr, "no video") {
			return nil, fmt.Errorf("extract failed: %w", err)
		}
	}

	stdout := ""
	if result != nil {
		stdout = result.Stdout
	}
	if stdout == "" {
		return nil, fmt.Errorf("empty response from yt-dlp")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	info := &models.MediaInfo{
		ID:          getString(raw, "id"),
		Title:       getString(raw, "title"),
		Uploader:    getString(raw, "uploader"),
		Thumbnail:   getString(raw, "thumbnail"),
		Description: getString(raw, "description"),
		WebpageURL:  getString(raw, "webpage_url"),
	}

	if duration, ok := raw["duration"].(float64); ok {
		info.Duration = duration
	}

	info.MediaType = detectMediaType(raw)

	if formatsRaw, ok := raw["formats"].([]interface{}); ok {
		for _, f := range formatsRaw {
			if fm, ok := f.(map[string]interface{}); ok {
				format := models.Format{
					FormatID:   getString(fm, "format_id"),
					Ext:        getString(fm, "ext"),
					Resolution: getString(fm, "resolution"),
					Note:       getString(fm, "format_note"),
					VCodec:     getString(fm, "vcodec"),
					ACodec:     getString(fm, "acodec"),
				}
				if size, ok := fm["filesize"].(float64); ok {
					format.Filesize = int64(size)
				}
				if fps, ok := fm["fps"].(float64); ok {
					format.FPS = fps
				}
				if format.FormatID != "" {
					info.Formats = append(info.Formats, format)
				}
			}
		}
	}

	return info, nil
}

func detectMediaType(raw map[string]interface{}) string {
	if duration, ok := raw["duration"].(float64); ok && duration > 0 {
		return "video"
	}

	if formats, ok := raw["formats"].([]interface{}); ok {
		for _, f := range formats {
			if fm, ok := f.(map[string]interface{}); ok {
				vcodec := strings.ToLower(getString(fm, "vcodec"))
				if vcodec != "" && vcodec != "none" {
					return "video"
				}
			}
		}
	}

	return "image"
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func buildFormatSelector(quality, formatID string) string {
	if formatID != "" {
		return formatID
	}

	switch strings.ToLower(quality) {
	case "1080p", "1080":
		return "bestvideo[height<=1080]+bestaudio/best[height<=1080]/best"
	case "720p", "720":
		return "bestvideo[height<=720]+bestaudio/best[height<=720]/best"
	case "480p", "480":
		return "bestvideo[height<=480]+bestaudio/best[height<=480]/best"
	case "360p", "360":
		return "bestvideo[height<=360]+bestaudio/best[height<=360]/best"
	case "worst":
		return "worst"
	default:
		return "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio/best[ext=mp4]/best"
	}
}

// Download
func (d *Downloader) Download(ctx context.Context, url, quality, formatID string) ([]string, string, error) {
	jobID := uuid.New().String()
	outputTemplate := filepath.Join(d.DownloadDir, jobID+"_%(playlist_index)02d_%(title).60B.%(ext)s")

	dl := ytdlp.New().
		NoOverwrites().
		Continue().
		Output(outputTemplate).
		IgnoreNoFormatsError().
		WriteThumbnail().
		ConvertThumbnails("jpg").
		FormatSort("res,ext:mp4:m4a").
		MergeOutputFormat("mp4").
		RecodeVideo("mp4")

	if formatID != "" {
		dl = dl.Format(formatID)
	} else {
		dl = dl.Format(buildFormatSelector(quality, ""))
	}

	_, err := dl.Run(ctx, url)
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "There is no video") &&
			!strings.Contains(errStr, "No video formats found") {
			return nil, "", fmt.Errorf("download failed: %w", err)
		}
	}

	files, err := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if err != nil || len(files) == 0 {
		return nil, "", fmt.Errorf("no files downloaded")
	}

	sort.Strings(files)

	mediaType := "video"
	imageCount := 0
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			imageCount++
		}
	}
	if imageCount == len(files) {
		mediaType = "image"
	} else if imageCount > 0 {
		mediaType = "mixed"
	}

	return files, mediaType, nil
}

// downloadImageWithYtDlp
func (d *Downloader) downloadImageWithYtDlp(ctx context.Context, url, jobID string) (string, string, error) {
	outputTemplate := filepath.Join(d.DownloadDir, jobID+"_%(title).80B")
	dl := ytdlp.New().
		NoPlaylist().
		SkipDownload().
		WriteThumbnail().
		IgnoreNoFormatsError().
		ConvertThumbnails("jpg").
		Output(outputTemplate)

	_, err := dl.Run(ctx, url)
	if err != nil {
		return "", "", fmt.Errorf("image download failed: %w", err)
	}

	patterns := []string{
		filepath.Join(d.DownloadDir, jobID+"*.jpg"),
		filepath.Join(d.DownloadDir, jobID+"*.jpeg"),
		filepath.Join(d.DownloadDir, jobID+"*.webp"),
		filepath.Join(d.DownloadDir, jobID+"*.png"),
		filepath.Join(d.DownloadDir, jobID+"*"),
	}

	var found string
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			found = matches[0]
			break
		}
	}

	if found == "" {
		return "", "", fmt.Errorf("thumbnail file not found after download")
	}

	return found, "image", nil
}

// DownloadWithProgress
func (d *Downloader) DownloadWithProgress(ctx context.Context, url, quality, formatID string, progressCb func(percent float64, status string)) (string, string, error) {
	info, err := d.ExtractInfo(ctx, url)
	if err != nil {
		return "", "", err
	}

	jobID := uuid.New().String()

	if info.MediaType == "image" {
		if progressCb != nil {
			progressCb(30, "downloading_image")
		}
		path, mediaType, err := d.downloadImageWithYtDlp(ctx, url, jobID)
		if progressCb != nil && err == nil {
			progressCb(100, "finished")
		}
		return path, mediaType, err
	}

	outputTemplate := filepath.Join(d.DownloadDir, jobID+"_%(title).80B.%(ext)s")
	formatSelector := buildFormatSelector(quality, formatID)

	dl := ytdlp.New().
		NoPlaylist().
		NoOverwrites().
		Continue().
		Output(outputTemplate).
		Format(formatSelector).
		FormatSort("res,ext:mp4:m4a").
		MergeOutputFormat("mp4").
		RecodeVideo("mp4").
		ProgressFunc(300*time.Millisecond, func(prog ytdlp.ProgressUpdate) {
			if progressCb == nil {
				return
			}
			percent := 0.0
			if p := prog.Percent(); p > 0 {
				percent = p
			} else {
				switch prog.Status {
				case ytdlp.ProgressStatusDownloading:
					percent = 45
				case ytdlp.ProgressStatusPostProcessing:
					percent = 85
				case ytdlp.ProgressStatusFinished:
					percent = 100
				}
			}
			progressCb(percent, string(prog.Status))
		})

	_, err = dl.Run(ctx, url)
	if err != nil {
		return "", "", fmt.Errorf("download failed: %w", err)
	}

	files, _ := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if len(files) == 0 {
		return "", "", fmt.Errorf("downloaded file not found")
	}

	return files[0], "video", nil
}