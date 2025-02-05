package main

import (
	"fmt"
	"net/http"
	"strings"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Fprintf(w, "%T, %+v\n", r.Form, r.Form) // print information on server side.
	fmt.Fprintln(w, "path", r.URL.Path)
	fmt.Fprintln(w, "scheme", r.URL.Scheme)
	fmt.Fprintln(w, r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintln(w, "Hello astaxie!") // write data to response
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
