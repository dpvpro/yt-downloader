package main

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"yt-downloader/handlers"
)

// Шаблонизатор
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Создаем директорию для загрузок, если ее нет
	os.MkdirAll("downloads", os.ModePerm)

	e := echo.New()

	// Настройка шаблонизатора
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Статические файлы
	e.Static("/public", "public")
	e.Static("/downloads", "downloads")

	// Инициализация обработчиков
	h := handlers.NewHandler()

	// Маршруты
	e.GET("/", h.IndexHandler)
	e.POST("/submit", h.SubmitHandler)
	e.GET("/status/:id", h.StatusHandler)
	e.GET("/download/:id", h.DownloadHandler)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}

