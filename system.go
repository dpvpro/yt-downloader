package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net/http"
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
			return value, err
		}
		fmt.Printf("%s\n", out)

	}

	return "", nil
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
