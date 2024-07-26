package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	yt_path        string = "/tmp/yt_downloader/"
	fileurl        string = "/mp3s/"
	dowloadedItems *[]string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func filterUrlStrings(s []string) []string {
	// filter empty strings and strings that begins with http or https prefix
	var r []string
	for _, str := range s {
		if str != "" && strings.HasPrefix(str, "http") || strings.HasPrefix(str, "https") {
			r = append(r, str)
		}
	}
	return r
}

func process(arr_clips []string) (item string, error error) {

	pwd, err := os.Getwd()
	check(err)
	err = os.RemoveAll(yt_path)
	check(err)
	err = os.Mkdir(yt_path, 0755)
	check(err)
	err = os.Chdir(yt_path)
	check(err)

	defer os.Chdir(pwd)

	for key, value := range arr_clips { // range over []string

		fmt.Println("Processing ", key, value)

		//yt-dlp -x --audio-format mp3 --audio-quality 0 https://youtu.be/BS5N_lAIohQ
		cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "0", value)
		//if err := cmd.Run(); err != nil {
		//	fmt.Println("Error: ", err)
		//}
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error: ", err)
			//time.Sleep(5 * time.Second)
			//log.Fatal(err)
			return value, err
		}
		fmt.Printf("%s\n", out)

	}

	return "", nil
}

func yt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	fmt.Println("url:", r.URL)       //get request method

	html, err := template.ParseFiles("yt.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)

}

func banner(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
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

}

func download(w http.ResponseWriter, r *http.Request) {

	fmt.Println(*dowloadedItems)
	item, err := process(*dowloadedItems)
	if err != nil {
		fmt.Fprintf(w, "Ошибка: '%s' со ссылкой: '%s'", err, item) // write data to response
		return
	}

	http.Redirect(w, r, "/serve/", http.StatusFound)

}

func serve(w http.ResponseWriter, r *http.Request) {

	// list directory
	fmt.Println("Listing", yt_path, "directory")
	c, err := os.ReadDir(yt_path)
	check(err)
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// redirect to directory listing
	w.Header().Set("Content-Type", "audio/mpeg")
	http.Redirect(w, r, fileurl, http.StatusSeeOther)

}

func main() {

	http.HandleFunc("/yt/", yt)
	http.HandleFunc("/prepare/", banner)
	http.HandleFunc("/download/", download)
	http.HandleFunc("/serve/", serve)

	http.Handle(fileurl,
		http.StripPrefix(fileurl,
			http.FileServer(
				http.Dir(yt_path))))

	http.HandleFunc("/hello/", sayHelloName)
	
	err := http.ListenAndServe(":10542", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
