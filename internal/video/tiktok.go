package video

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TikTokDownloader struct{}

func NewTikTokDownloader() *TikTokDownloader {
	return &TikTokDownloader{}
}

type tikWmResponse struct {
	Code int `json:"code"`
	Data struct {
		Play string `json:"play"`
	} `json:"data"`
}

func (d *TikTokDownloader) Download(ctx context.Context, videoUrl string) (*DownloadVideo, error) {
	apiURL := "https://www.tikwm.com/api/?url=" + url.QueryEscape(videoUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp tikWmResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != 0 || apiResp.Data.Play == "" {
		return nil, errors.New("failed to get video link from tikwm API")
	}

	vidReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiResp.Data.Play, nil)
	if err != nil {
		return nil, err
	}
	vidResp, err := http.DefaultClient.Do(vidReq)
	if err != nil {
		return nil, err
	}
	defer vidResp.Body.Close()

	tmpFile, err := os.CreateTemp("", "tiktok-*.mp4")
	if err != nil {
		return nil, err
	}
	filePath := tmpFile.Name()
	
	size, err := io.Copy(tmpFile, vidResp.Body)
	tmpFile.Close()

	if err != nil || size == 0 {
		os.Remove(filePath)
		return nil, errors.New("failed to download video file")
	}

	return &DownloadVideo{
		FilePath:    filePath,
		OriginalURL: videoUrl,
		Size:        size,
	}, nil
}
