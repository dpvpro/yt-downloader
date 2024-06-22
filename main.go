package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	yt_path string = "/tmp/yt_downloader/"
	files   string = "/files/"
)

func yt(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("method:", r.Method) //get request method
		fmt.Println("url:", r.URL)       //get request method

		html, err := template.ParseFiles("yt.html")
		check(err)
		err = html.Execute(w, nil)
		check(err)
	}

	if r.Method == "POST" {
		r.ParseForm()
		// logic part of log in
		var values string = r.FormValue("message")
		fmt.Println("page message:")
		fmt.Println(values)
		fmt.Println(":end page message")
		splitValues := strings.Split(values, "\r\n")
		filterValues := filterUrlStrings(splitValues)
		fmt.Println(splitValues)
		fmt.Println("-------")
		fmt.Println(filterValues)
		fmt.Println("-------")

		item, err := process(filterValues)
		if err != nil {
			fmt.Fprintf(w, "Ошибка: '%s' со ссылкой: '%s'", err, item) // write data to response
			return
		}

		serve(w, r)

	}
}

func main() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(yt_path))
	mux.Handle(files, http.StripPrefix(files, fileServer))

	mux.HandleFunc("/hello/", sayHelloName) // setting router rule
	mux.HandleFunc("/yt/", yt)
	//mux.HandleFunc("/process/", process)
	mux.HandleFunc("/serve/", serve)

	err := http.ListenAndServe(":10542", mux) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
