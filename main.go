package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("yt.html")
		t.Execute(w, nil)
	} else {
		//r.ParseForm()
		// logic part of log in
		var values = r.FormValue("message")
		fmt.Println("message:", values)
		//fmt.Println("form:", r.Form)
		valString := strings.Split(values, "\n")
		for key, value := range valString { // range over []string
			//fmt.Println(key, value)
			fmt.Println("Processing ", key, value)

			//timeout := time.Duration(5) * time.Second
			//transport := &http.Transport{
			//	ResponseHeaderTimeout: timeout,
			//	DisableKeepAlives:     true,
			//}
			//client := &http.Client{
			//	Transport: transport,
			//}

			//resp, err := client.Get(value)
			//if err != nil {
			//	fmt.Println(err)
			//}

			// process file

			//yt-dlp -x --audio-format mp3 --audio-quality 0 https://youtu.be/BS5N_lAIohQ
			cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "0", value)
			//if err := cmd.Run(); err != nil {
			//	fmt.Println("Error: ", err)
			//}
			out, err := cmd.CombinedOutput()
			if err != nil {
				//	log.Fatal(err)
				//  fmt.Println("Error: ", err)
				fmt.Fprintf(w, "Error: %d", err) // write data to response

			}
			fmt.Fprintf(w, "%s\n", out)
			// process file

			//defer resp.Body.Close()

			//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
			//w.Header().Set("Content-Disposition", "attachment")
			//w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			//w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

			//stream the body to the client without fully loading it into memory
			//io.Copy(w, resp.Body)

		}
	}
}

func main() {
	http.HandleFunc("/hello", sayhelloName) // setting router rule
	http.HandleFunc("/yt", yt)
	err := http.ListenAndServe(":10542", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
