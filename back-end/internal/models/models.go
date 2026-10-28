package models

type MediaInfo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Uploader    string   `json:"uploader"`
	Duration    float64  `json:"duration"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	WebpageURL  string   `json:"webpage_url"`
	MediaType   string   `json:"media_type"`
	Formats     []Format `json:"formats"`
}

type Format struct {
	FormatID   string  `json:"format_id"`
	Ext        string  `json:"ext"`
	Resolution string  `json:"resolution"`
	Filesize   int64   `json:"filesize"`
	Note       string  `json:"note"`
	VCodec     string  `json:"vcodec,omitempty"`
	ACodec     string  `json:"acodec,omitempty"`
	FPS        float64 `json:"fps,omitempty"`
	URL        string  `json:"url,omitempty"`
}

type DownloadRequest struct {
	URL      string `json:"url" binding:"required"`
	FormatID string `json:"format_id"`
	Quality  string `json:"quality"` // "best", "1080p", "720p", "480p", "360p", "worst"
}

type DownloadResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	FilePath  string `json:"file_path,omitempty"`
	MediaType string `json:"media_type,omitempty"`
}