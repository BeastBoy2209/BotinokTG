package video
import(
	"context"
	"os"
	"errors"
)

type VideoService struct{
	ytDownloader Downloader
	tkDownloader Downloader
	maxSize int64
}

func NewService(ytDownloader Downloader, tkDownloader Downloader, maxSize int64) *VideoService{
	return &VideoService{
		ytDownloader: ytDownloader,
		tkDownloader: tkDownloader,
		maxSize:      maxSize,
	}
}

func (s *VideoService) ProcessMessage(ctx context.Context, text string) (*DownloadVideo, error){
	url, platform, supported := DetectPlatform(text)
	if !supported{
		return nil, nil
	}

	var downloader Downloader
	if platform == PlatformTikTok {
		downloader = s.tkDownloader
	} else {
		downloader = s.ytDownloader
	}

	video, err := downloader.Download(ctx, url)
	if err != nil{
		return nil, err
	}
	
	video.Platform = platform
	
	if video.Size > s.maxSize {
		os.Remove(video.FilePath)
		return nil, errors.New("video is too large for Telegram")
	}
	
	return video, nil
}

