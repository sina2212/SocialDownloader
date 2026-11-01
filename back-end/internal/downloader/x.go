package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"back-end/internal/models"

	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
)

// isXURL routes x.com and twitter.com URLs through X-specific photo handling.
func isXURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	h := strings.ToLower(u.Hostname())
	return h == "x.com" || h == "www.x.com" || h == "twitter.com" || h == "www.twitter.com" || h == "mobile.twitter.com"
}

func xStatusID(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	for i := 0; i+1 < len(parts); i++ {
		if parts[i] == "status" && parts[i+1] != "" {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("could not find X status ID in URL")
}

// extractXPhotos reads the public syndication representation, which exposes
// photo URLs separately from video formats. It deliberately does not treat a
// tweet thumbnail as the original image.
func (d *Downloader) extractXPhotos(ctx context.Context, tweetURL string) (*models.MediaInfo, error) {
	id, err := xStatusID(tweetURL)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("id", id)
	q.Set("lang", "en")
	q.Set("token", "0")
	endpoint := "https://cdn.syndication.twimg.com/tweet-result?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MediaDownloader/1.0)")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("X syndication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X syndication returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf("read X syndication response: %w", err)
	}

	var generic map[string]interface{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return nil, fmt.Errorf("decode X syndication response: %w", err)
	}

	if details, ok := generic["mediaDetails"].([]interface{}); ok {
		for _, item := range details {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			typ, _ := m["type"].(string)
			if typ == "video" || typ == "animated_gif" {
				return nil, fmt.Errorf("post contains video/gif (type=%s), not pure photo", typ)
			}
		}
	}

	if _, hasVideo := generic["video"]; hasVideo {
		return nil, fmt.Errorf("post contains video object")
	}

	info := &models.MediaInfo{
		ID:         id,
		Title:      getString(generic, "text"),
		Uploader:   "", // بعداً پر می‌کنیم
		WebpageURL: tweetURL,
		MediaType:  "image",
	}

	if user, ok := generic["user"].(map[string]interface{}); ok {
		info.Uploader = getString(user, "name")
	}

	if photos, ok := generic["photos"].([]interface{}); ok {
		for _, p := range photos {
			if pm, ok := p.(map[string]interface{}); ok {
				if u := getString(pm, "url"); validRemoteMediaURL(u) {
					info.MediaURLs = append(info.MediaURLs, u)
				}
			}
		}
	}

	if len(info.MediaURLs) == 0 {
		if details, ok := generic["mediaDetails"].([]interface{}); ok {
			for _, item := range details {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				typ, _ := m["type"].(string)
				if typ != "" && typ != "photo" {
					continue
				}
				for _, key := range []string{"media_url_https", "media_url", "url"} {
					if candidate, ok := m[key].(string); ok && validRemoteMediaURL(candidate) {
						info.MediaURLs = append(info.MediaURLs, candidate)
						break
					}
				}
			}
		}
	}

	if len(info.MediaURLs) == 0 {
		return nil, fmt.Errorf("X syndication returned no photo URLs")
	}

	info.Thumbnail = info.MediaURLs[0]
	return info, nil
}

func validRemoteMediaURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != ""
}

func (d *Downloader) downloadXPhotos(ctx context.Context, urls []string) ([]string, string, error) {
	if len(urls) == 0 {
		return nil, "", fmt.Errorf("X post contains no downloadable photos")
	}
	jobID := uuid.NewString()
	client := &http.Client{Timeout: 45 * time.Second}
	files := make([]string, 0, len(urls))
	for i, raw := range urls {
		if !validRemoteMediaURL(raw) {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
		if err != nil {
			return files, "image", err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36")
		req.Header.Set("Referer", "https://x.com/")
		resp, err := client.Do(req)
		if err != nil {
			return files, "image", fmt.Errorf("download X photo %d: %w", i+1, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return files, "image", fmt.Errorf("download X photo %d: HTTP %d", i+1, resp.StatusCode)
		}
		ext := filepath.Ext(resp.Request.URL.Path)
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			ext = ".jpg"
		}
		name := filepath.Join(d.DownloadDir, fmt.Sprintf("%s_%02d%s", jobID, i+1, ext))
		out, err := os.Create(name)
		if err == nil {
			_, err = io.Copy(out, io.LimitReader(resp.Body, 40<<20))
		}
		closeErr := resp.Body.Close()
		if out != nil {
			if e := out.Close(); err == nil {
				err = e
			}
		}
		if err == nil {
			err = closeErr
		}
		if err != nil {
			os.Remove(name)
			return files, "image", fmt.Errorf("save X photo %d: %w", i+1, err)
		}
		files = append(files, name)
	}
	if len(files) == 0 {
		return nil, "", fmt.Errorf("X returned no valid photo URLs")
	}
	sort.Strings(files)
	return files, "image", nil
}

// extractXVideo
func (d *Downloader) extractXVideo(ctx context.Context, tweetURL string) (*models.MediaInfo, error) {
	dl := ytdlp.New().
		DumpSingleJSON().
		NoPlaylist().
		NoWarnings().
		SkipDownload().
		IgnoreNoFormatsError()

	result, err := dl.Run(ctx, tweetURL)
	if err != nil {
		return nil, fmt.Errorf("X video extract failed: %w", err)
	}

	stdout := ""
	if result != nil {
		stdout = result.Stdout
	}
	if stdout == "" {
		return nil, fmt.Errorf("empty response from yt-dlp for X")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp json: %w", err)
	}

	info := &models.MediaInfo{
		ID:          getString(raw, "id"),
		Title:       getString(raw, "title"),
		Uploader:    getString(raw, "uploader"),
		Thumbnail:   getString(raw, "thumbnail"),
		Description: getString(raw, "description"),
		WebpageURL:  getString(raw, "webpage_url"),
		MediaType:   "video",
	}

	if duration, ok := raw["duration"].(float64); ok {
		info.Duration = duration
	}

	if info.Duration == 0 {
		info.MediaType = detectMediaType(raw)
	}

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
					URL:        getString(fm, "url"),
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

// downloadXVideo
func (d *Downloader) downloadXVideo(ctx context.Context, tweetURL, quality, formatID string) ([]string, string, error) {
	jobID := uuid.NewString()
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
		RecodeVideo("mp4")

	_, err := dl.Run(ctx, tweetURL)
	if err != nil {
		return nil, "", fmt.Errorf("X video download failed: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if err != nil || len(files) == 0 {
		return nil, "", fmt.Errorf("no files downloaded for X video")
	}

	sort.Strings(files)
	return files, "video", nil
}

func (d *Downloader) downloadXVideoWithProgress(ctx context.Context, tweetURL, quality, formatID string, progressCb func(percent float64, status string),) (string, string, error) {

	jobID := uuid.NewString()
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

	_, err := dl.Run(ctx, tweetURL)
	if err != nil {
		return "", "", fmt.Errorf("X video download failed: %w", err)
	}

	files, _ := filepath.Glob(filepath.Join(d.DownloadDir, jobID+"*"))
	if len(files) == 0 {
		return "", "", fmt.Errorf("downloaded file not found")
	}

	return files[0], "video", nil
}