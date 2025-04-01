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

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// настройка шаблонизатора
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// создаем директорию для загрузок, если ее нет
	os.MkdirAll("downloads", os.ModePerm)

	// middleware

	// create log file
	logFile, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		e.Logger.Fatal("error opening file: %v", err)
	}
	defer logFile.Close()

	// output in log file
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{Output: logFile},
	))

	// output in console
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Статические файлы
	e.Static("/templates", "templates")
	e.Static("/downloads", "downloads")

	// Инициализация обработчиков
	h := handlers.NewHandler()

	// Маршруты
	e.GET("/", h.IndexHandler)
	e.POST("/submit", h.SubmitHandler)
	e.GET("/status/:id", h.StatusHandler)
	e.GET("/download/:id", h.DownloadHandler)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":10542"))
}
