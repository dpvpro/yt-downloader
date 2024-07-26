package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	yt_path 	string = "/tmp/yt_downloader/"
	fileurl   	string = "/files/"
)


func yt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	fmt.Println("url:", r.URL)       //get request method

	html, err := template.ParseFiles("yt.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)

}

var dowloadedItems *[]string

func prepare(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// logic part of log in
	var values string = r.FormValue("message")
	fmt.Println("page message:")
	fmt.Println(values)
	fmt.Println(":end page message")
	splitValues := strings.Split(values, "\r\n")
	filterValues := filterUrlStrings(splitValues)
	fmt.Println("-------")
	fmt.Println(splitValues)
	fmt.Println("-------")
	fmt.Println(filterValues)
	fmt.Println("-------")

	dowloadedItems = &filterValues

	html, err := template.ParseFiles("banner.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)
	
	// http.Redirect(w, r, "/download/", http.StatusFound)
}


func download(w http.ResponseWriter, r *http.Request) {


	fmt.Println(*dowloadedItems)
	item, err := process(*dowloadedItems)
	if err != nil {
		fmt.Fprintf(w, "Ошибка: '%s' со ссылкой: '%s'", err, item) // write data to response
		return
	}

	// http.Redirect(w, r, "/files/", http.StatusFound)
	serve(w, r)

}

func main() {

	http.HandleFunc("/hello/", sayHelloName)
	http.HandleFunc("/yt/", yt)
	http.HandleFunc("/prepare/", prepare)
	// http.HandleFunc("/banner/", banner)
	http.HandleFunc("/download/", download)
	http.HandleFunc("/serve/", serve)
	
	http.Handle(fileurl, 
		http.StripPrefix(fileurl,
			http.FileServer(
				http.Dir(yt_path))))


	err := http.ListenAndServe(":10542", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
