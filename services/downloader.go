package services

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"yt-downloader/models"
)

type Downloader struct {
	store *models.DownloadStore
}

func NewDownloader(store *models.DownloadStore) *Downloader {
	return &Downloader{
		store: store,
	}
}

func (d *Downloader) ProcessRequest(request *models.DownloadRequest) {
	for _, video := range request.Videos {
		// Обновляем статус на "скачивание"
		d.store.UpdateVideoStatus(request.ID, video.URL, models.StatusDownloading, "", "")

		// Получаем информацию о видео
		title, err := d.getVideoTitle(video.URL, request.UseProxy, request.ProxyURL)
		if err != nil {
			d.store.UpdateVideoStatus(request.ID, video.URL, models.StatusFailed, 
				fmt.Sprintf("Error getting video title: %v", err), "")
			continue
		}

		// Формируем безопасное имя файла
		safeTitle := sanitizeFilename(title)
		fileName := safeTitle + ".mp3"
		filePath := filepath.Join("downloads", fileName)

		// Скачиваем аудио
		err = d.downloadAudio(video.URL, filePath, request.UseProxy, request.ProxyURL)
		if err != nil {
			d.store.UpdateVideoStatus(request.ID, video.URL, models.StatusFailed, 
				fmt.Sprintf("Error downloading audio: %v", err), "")
			continue
		}

		// Обновляем информацию о видео
		d.store.UpdateVideoStatus(request.ID, video.URL, models.StatusCompleted, "", fileName)
	}
}

func (d *Downloader) getVideoTitle(url string, useProxy bool, proxyURL string) (string, error) {
	args := []string{
		"--get-title",
		url,
	}

	if useProxy && proxyURL != "" {
		args = append(args, "--proxy", proxyURL)
		args = append(args, "--cookies", "env/cookies.txt")
	}

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (d *Downloader) downloadAudio(url string, outputPath string, useProxy bool, proxyURL string) error {
	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"-o", outputPath,
		url,
	}

	if useProxy && proxyURL != "" {
		args = append(args, "--proxy", proxyURL)
		args = append(args, "--cookies", "env/cookies.txt")
		
	}

	cmd := exec.Command("yt-dlp", args...)
	// fmt.Println(cmd.String())
	return cmd.Run()
}

func sanitizeFilename(name string) string {
	// Заменяем недопустимые символы в имени файла
	forbidden := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name

	for _, char := range forbidden {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Ограничиваем длину имени файла
	if len(result) > 200 {
		result = result[:200]
	}

	return result
}