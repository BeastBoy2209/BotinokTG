package video

import "context"

type Downloader interface {
	Download(ctx context.Context, url string) (*DownloadVideo, error)
}