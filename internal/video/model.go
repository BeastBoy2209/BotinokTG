package video

type Platform string

const (
	PlatformTikTok    Platform = "TikTok"
	PlatformInstagram Platform = "Instagram"
	PlatformYouTube   Platform = "YouTube"
	PlatformUnknown   Platform = "Unknown"
)

type DownloadVideo struct {
	FilePath    string
	Platform    Platform
	OriginalURL string
	Size        int64
}
