package video 

import(
	"strings"
	"regexp"
)
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

func ExtractURL(text string) string {
	return urlRegex.FindString(text)
}

func DetectPlatform(text string) (url string, platform Platform, supported bool){
	link := ExtractURL(text)

	if link == "" {
		return "", PlatformUnknown, false
	}

	if strings.Contains(link, "tiktok.com"){
		return link, PlatformTikTok, true
	}else if strings.Contains(link, "instagram.com"){
		return link, PlatformInstagram, true
	}else if strings.Contains(link, "youtube.com") || strings.Contains(link, "youtu.be"){
		return link, PlatformYouTube, true
	}else{
		return "", PlatformUnknown, false
	}
}