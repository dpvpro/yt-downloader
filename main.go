package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var (
	ytPath         string = "/tmp/yt_downloader/"
	fileUrl        string = "/mp3s/"
	dowloadedItems *[]string
	pwd            string
	err            error
	workerLimit    int = 2
)

func FilterUrlStrings(s []string) []string {
	// filter empty strings and strings that begins with http or https prefix
	var r []string
	regExFilter, _ := regexp.Compile("^https?")
	for _, str := range s {
		if str != "" && regExFilter.MatchString(str) {
			r = append(r, str)
		}
	}
	return r
}

func Process(arr_clips []string) (item string, error error) {

	err = os.RemoveAll(ytPath)
	Check(err)
	err = os.Mkdir(ytPath, 0755)
	Check(err)
	err = os.Chdir(ytPath)
	Check(err)

	// На основе ratelim.go из курса на Степике от Василия Романова

	wg := &sync.WaitGroup{}
	quotaChanel := make(chan bool, workerLimit) // ratelim.go

	for index, clip := range arr_clips { // range over []string
		wg.Add(1)
		// go startWorker(i, wg, quotaCh)
		go func(indexGorutine int, clipGorutine string) {
			quotaChanel <- true // ratelim.go, берём свободный слот
			defer wg.Done()
			defer func() {
				<-quotaChanel
			}()

			fmt.Println("Processing ", indexGorutine, clipGorutine)

			//yt-dlp -x --audio-format mp3 --audio-quality 0 https://youtu.be/BS5N_lAIohQ
			cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "0", clipGorutine)
			//if err := cmd.Run(); err != nil {
			//	fmt.Println("Error: ", err)
			//}
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error: ", err)
				// time.Sleep(5 * time.Second)
				// log.Fatal(err)
				// return value, err
				return
			}
			fmt.Printf("%s\n", out)

		}(index, clip)
	}

	wg.Wait()

	// сформировать переменную перед циклом для ошибок и после проверять
	// https://pkg.go.dev/errors#Join
	// использовать мьютексы что бы использовать в разных горутинах общую
	// переменную которую я сформировал выше

	err = os.Chdir(pwd)
	Check(err)

	return "", nil
}

func Yt(w http.ResponseWriter, r *http.Request) {

	err := os.Chdir(pwd)
	Check(err)

	fmt.Println("method:", r.Method) //get request method
	fmt.Println("url:", r.URL)       //get request method

	html, err := template.ParseFiles("yt.html")
	Check(err)
	err = html.Execute(w, nil)
	Check(err)

}

func Waiting(w http.ResponseWriter, r *http.Request) {

	err := os.Chdir(pwd)
	Check(err)

	r.ParseForm()
	var values string = r.FormValue("message")
	fmt.Println("page message:")
	fmt.Println(values)
	fmt.Println(":end page message")
	splitValues := strings.Split(values, "\r\n")
	filterValues := FilterUrlStrings(splitValues)
	fmt.Println("-------")
	fmt.Println(splitValues)
	fmt.Println("-------")
	fmt.Println(filterValues)
	fmt.Println("-------")

	dowloadedItems = &filterValues

	html, err := template.ParseFiles("waiting.html")
	Check(err)
	err = html.Execute(w, nil)
	Check(err)

}

func Download(w http.ResponseWriter, r *http.Request) {

	fmt.Println(*dowloadedItems)
	item, err := Process(*dowloadedItems)
	if err != nil {
		fmt.Fprintf(w, "Ошибка: '%s' со ссылкой: '%s'", err, item) // write data to response
		return
	}

	http.Redirect(w, r, "/serve/", http.StatusFound)

}

func Serve(w http.ResponseWriter, r *http.Request) {

	// list directory
	fmt.Println("Listing", ytPath, "directory")
	c, err := os.ReadDir(ytPath)
	Check(err)
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// redirect to directory listing
	w.Header().Set("Content-Type", "audio/mpeg")
	http.Redirect(w, r, fileUrl, http.StatusSeeOther)

}

func main() {

	pwd, err = os.Getwd()
	Check(err)
	// исходим из того что у нас используется отдельный домен или поддомен
	http.HandleFunc("/", Yt)
	http.HandleFunc("/waiting/", Waiting)
	http.HandleFunc("/download/", Download)
	http.HandleFunc("/serve/", Serve)

	http.Handle(fileUrl,
		http.StripPrefix(fileUrl,
			http.FileServer(
				http.Dir(ytPath))))

	http.HandleFunc("/hello/", SayHelloName)

	err = http.ListenAndServe(":10542", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
