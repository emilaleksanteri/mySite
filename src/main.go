package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func message(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	template := template.Must(template.ParseFiles("responses/text.html"))
	template.Execute(res, nil)
	
}

type TestData struct {
	Heading string
	Body    string
}

func hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	template := template.Must(template.ParseFiles("pages/hello.html"))
	testData := map[string][]TestData{
		"testData": {
			{Heading: "hello world", Body: "I am text"},	
			{Heading: "hello world", Body: "I am text"},	
			{Heading: "hello world", Body: "I am text"},	
		},
	}

	template.Execute(res, testData)

}

func main() {
	fileServe := http.FileServer(http.Dir("."))
	http.Handle("/", http.StripPrefix("/", fileServe))

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/test", message)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
