package video

import(
	"context"
	"errors"
	"os"
	"os/exec"
) 

type YTDLPDownloader struct{
	executablePath string
}

func NewYTDLPDownloader(path string) *YTDLPDownloader{
	return &YTDLPDownloader{executablePath: path}
}

func (d *YTDLPDownloader)Download(ctx context.Context, url string) (*DownloadVideo, error){
	tmpFile, err := os.CreateTemp("", "video-*.mp4")
	if err != nil{
		return nil, err
	}
	filePath := tmpFile.Name()
	tmpFile.Close() 
	cmd := exec.CommandContext(ctx, d.executablePath, 
		"-f", "bestvideo[ext=mp4][vcodec^=avc1]+bestaudio[ext=m4a]/best[ext=mp4]/best", 
		"--merge-output-format", "mp4",
		"--force-overwrites",
		"-o", filePath, 
		url,
	)
	err = cmd.Run()

	if err != nil {
		os.Remove(filePath)
		return nil, err
	}
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		os.Remove(filePath)
		return nil, err
	}
	size := fileInfo.Size()
	if size == 0 {
		os.Remove(filePath)
		return nil, errors.New("yt-dlp returned an empty file (possibly missing ffmpeg or format issue)")
	}

	return &DownloadVideo{
    FilePath:    filePath,
    OriginalURL: url,
    Size:        size,
	}, nil
}