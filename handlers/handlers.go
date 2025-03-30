package handlers

import (
	"net/http"
	"time"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"yt-downloader/models"
	"yt-downloader/services"
)

type Handler struct {
	store     *models.DownloadStore
	downloader *services.Downloader
}

func NewHandler() *Handler {
	store := models.NewDownloadStore()
	downloader := services.NewDownloader(store)
	return &Handler{
		store:     store,
		downloader: downloader,
	}
}

func (h *Handler) IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (h *Handler) SubmitHandler(c echo.Context) error {
	urls := c.FormValue("urls")
	useProxy := c.FormValue("use_proxy") == "on"
	proxyURL := c.FormValue("proxy_url")

	// Генерация уникального ID для запроса
	id := generateID()

	// Разбиваем текстовое поле на отдельные URLs
	urlList := strings.Split(urls, "\n")
	cleanURLs := []string{}
	videos := []models.VideoInfo{}

	for _, url := range urlList {
		url = strings.TrimSpace(url)
		if url != "" {
			cleanURLs = append(cleanURLs, url)
			videos = append(videos, models.VideoInfo{
				URL:       url,
				Status:    models.StatusPending,
				CreatedAt: time.Now(),
			})
		}
	}

	// Создаем запрос на скачивание
	request := &models.DownloadRequest{
		ID:        id,
		URLs:      cleanURLs,
		Videos:    videos,
		UseProxy:  useProxy,
		ProxyURL:  proxyURL,
		CreatedAt: time.Now(),
	}

	// Сохраняем запрос
	h.store.Add(request)

	// Запускаем процесс скачивания в фоне
	go h.downloader.ProcessRequest(request)

	// Перенаправляем на страницу статуса
	return c.Redirect(http.StatusSeeOther, "/status/"+id)
}

func (h *Handler) StatusHandler(c echo.Context) error {
	id := c.Param("id")
	request, ok := h.store.Get(id)
	if !ok {
		return c.String(http.StatusNotFound, "Download request not found")
	}

	// Проверка, завершены ли все скачивания
	allCompleted := true
	for _, video := range request.Videos {
		if video.Status != models.StatusCompleted && video.Status != models.StatusFailed {
			allCompleted = false
			break
		}
	}

	return c.Render(http.StatusOK, "status.html", map[string]interface{}{
		"ID":           id,
		"Videos":       request.Videos,
		"AllCompleted": allCompleted,
	})
}

func (h *Handler) DownloadHandler(c echo.Context) error {
	id := c.Param("id")
	request, ok := h.store.Get(id)
	if !ok {
		return c.String(http.StatusNotFound, "Download request not found")
	}

	// Фильтруем только успешно скачанные видео
	completedVideos := []models.VideoInfo{}
	for _, video := range request.Videos {
		if video.Status == models.StatusCompleted {
			completedVideos = append(completedVideos, video)
		}
	}

	return c.Render(http.StatusOK, "download.html", map[string]interface{}{
		"ID":     id,
		"Videos": completedVideos,
	})
}

// Генерация уникального ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}