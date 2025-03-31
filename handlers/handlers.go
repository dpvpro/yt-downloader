package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"yt-downloader/models"
	"yt-downloader/services"

	"github.com/labstack/echo/v4"
)

var once sync.Once

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
	if urls == "" {
		return c.String(http.StatusBadRequest, "No urls provided")
	}
	useProxy := c.FormValue("use_proxy") == "on"
	proxyURL := ProxySettings()

	// генерация уникального ID для запроса
	id := generateID()


	// Разбиваем текстовое поле на отдельные URLs
	cleanUrls := flterUrlStrings(urls)
	videos := []models.VideoInfo{}

	for _, url := range cleanUrls {
		url = strings.TrimSpace(url)
		videos = append(videos, models.VideoInfo{
			URL:       url,
			Status:    models.StatusPending,
			CreatedAt: time.Now(),
		})
	}

	// Создаем запрос на скачивание
	request := &models.DownloadRequest{
		ID:        id,
		URLs:      cleanUrls,
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

	return c.Render(http.StatusOK, "status.html", map[string]any{
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

	return c.Render(http.StatusOK, "download.html", map[string]any{
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

// filter empty strings and strings that begins with http or https prefix
func flterUrlStrings(s string) []string {
	var r []string
	
	urlList := strings.Split(s, "\n")
	
	regExFilter, _ := regexp.Compile("^https?")
	for _, str := range urlList {
		if str != "" && regExFilter.MatchString(str) {
			r = append(r, str)
		}
	}
	return r
}


func ProxySettings() string {
	var ProxySettings []byte
	var err error
	
	once.Do(func() {
		ProxySettings, err = os.ReadFile("env/proxy.txt")
	})
	if err != nil {
		return ""
	}
	return string(ProxySettings)
	
}
