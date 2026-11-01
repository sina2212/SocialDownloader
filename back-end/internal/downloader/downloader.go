package downloader

import (
	"context"
	"fmt"
	"os"
	"strings"

	"back-end/internal/models"
)

// Downloader is the platform-independent entry point. Platform-specific
// extraction and download logic lives in instagram.go and twitter.go.
type Downloader struct{ DownloadDir string }

func New(downloadDir string) *Downloader {
	_ = os.MkdirAll(downloadDir, 0755)
	return &Downloader{DownloadDir: downloadDir}
}

func (d *Downloader) ExtractInfo(ctx context.Context, rawURL string) (*models.MediaInfo, error) {
	if isXURL(rawURL) {
		info, photoErr := d.extractXPhotos(ctx, rawURL)
		if photoErr == nil && len(info.MediaURLs) > 0 {
			return info, nil
		}

		videoInfo, videoErr := d.extractXVideo(ctx, rawURL)
		if videoErr != nil {
			if photoErr != nil {
				return nil, fmt.Errorf("X photo extraction failed: %v; video extraction failed: %w", photoErr, videoErr)
			}
			return nil, videoErr
		}
		return videoInfo, nil
	}
	return d.instaExtractInfo(ctx, rawURL)
}

func (d *Downloader) Download(ctx context.Context, rawURL, quality, formatID string) ([]string, string, error) {
	if isXURL(rawURL) {
		info, err := d.extractXPhotos(ctx, rawURL)
		if err == nil && len(info.MediaURLs) > 0 {
			return d.downloadXPhotos(ctx, info.MediaURLs)
		}

		files, typ, dlErr := d.downloadXVideo(ctx, rawURL, quality, formatID)
		if dlErr != nil && err != nil {
			return nil, "", fmt.Errorf("X photo extraction: %v; X video download: %w", err, dlErr)
		}
		return files, typ, dlErr
	}
	return d.instaDownload(ctx, rawURL, quality, formatID)
}

func (d *Downloader) DownloadWithProgress(ctx context.Context, rawURL, quality, formatID string, cb func(float64, string)) (string, string, error) {
	if isXURL(rawURL) {
		info, err := d.extractXPhotos(ctx, rawURL)
		if err == nil && len(info.MediaURLs) > 0 {
			if cb != nil {
				cb(5, "downloading_images")
			}
			files, typ, err := d.downloadXPhotos(ctx, info.MediaURLs)
			if err != nil {
				return "", "", err
			}
			if cb != nil {
				cb(100, "finished")
			}
			return files[0], typ, nil
		}
		return d.downloadXVideoWithProgress(ctx, rawURL, quality, formatID, cb)
	}
	return d.instaDownloadWithProgress(ctx, rawURL, quality, formatID, cb)
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
