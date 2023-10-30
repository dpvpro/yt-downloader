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
	yt_path string = "/tmp/yt_downloader/"
	files   string = "/files/"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

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
		var values = r.FormValue("message")
		fmt.Println("message:", values)

		valString1 := strings.Split(values, "\r\n")
		valString2 := removeEmptyStrings(valString1)
		fmt.Println(valString1)
		fmt.Println(valString2)
		//fmt.Println(valString)

		process(valString2)

		//err := process(valString2)
		//if err != nil {
		//	fmt.Fprintf(w, "Error: %d", err) // write data to response
		//}

		http.Redirect(w, r, "/serve/", http.StatusSeeOther)

	}
}

func process(arr_clips []string) error {

	var err error
	var pwd string

	pwd, err = os.Getwd()
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

		// process file

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
			return err
		}
		fmt.Printf("%s\n", out)

	}

	return nil
}
func serve(w http.ResponseWriter, r *http.Request) {

	// list directory
	fmt.Println("Listing", yt_path, "directory")
	c, err := os.ReadDir(yt_path)
	check(err)
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	http.Redirect(w, r, files, http.StatusSeeOther)

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
