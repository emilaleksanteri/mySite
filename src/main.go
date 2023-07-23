package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Project struct {
	ProjectTitle    string
	ShortSummary    string
	LongDescription string
	Id              string
	Url             string
}

func getProjectById(id string) Project {
	projects := []Project{
		{
			ProjectTitle:    "chatApp",
			ShortSummary:    "Live chat application built with the T3 stack",
			LongDescription: "Currently building a hybrid clone thats based on WhatsApp and Discord for me and my friends to use. The project is written in TypeScript, using Next.js as the framework. The app uses MySQL as the database, tRPC for client to interact with the server in a type safe manner. I wanted to plau with some AI stuff with the chat application so currently the users can invite ChatGPT to be part of their conversation by tagging AI into the chat with @ai. The application connects to a separate WebSocket server I had made in TypeScript which is used to let users know who is typing and to invalidate the chats cache whenever a new message arrives in the chat. This also served as a gateway to play funny sounds to all chat participants with the press of a sound board button as part of the 'chatting experience'.",
			Id:              "chatapp",
			Url:             "https://chat-app-tau-teal.vercel.app/",
		},

		{
			ProjectTitle:    "HaikuForu",
			ShortSummary:    "A website that generates Haikus with the use of AI",
			LongDescription: "A web app I made with Next.js and deployed on Vercel, written in TypeScript using React. I wanted to use OpenAI's API so the application uses this as a haiku generator for  any topic the user inputs",
			Id:              "haiku",
			Url:             "https://haiku-for-u.vercel.app/",
		},
	}

	project := Project{}
	for _, proj := range projects {
		if proj.Id == id {
			project = proj
		}
	}

	return project
}

func resume(res http.ResponseWriter, req *http.Request) {
	resume, err := os.Open("cv/emil-lystimaki-cv.pdf")
	if err != nil {
		fmt.Println("Could not get cv")
	}
	defer resume.Close()

	res.Header().Set(
		"Content-Type",
		"application/pdf",
	)
	if _, err := io.Copy(res, resume); err != nil {
		fmt.Println("Could not serve resume bruh moment")
		res.WriteHeader(500)
	}
}

type Blog struct {
	Thumbnail string
	Title     string
	Preview   string
	Text      string
	Id        string
}

func blog(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	ids := req.URL.Query()["id"]

	id := ids[0]
	fileName := "blogs/" + id
	openFile, err := os.Open(fileName)
	if err != nil {
		fmt.Errorf("Not able to read file at=%T", fileName)
	}

	defer openFile.Close()

	blo := Blog{}
	blo.Id = id

	var text []string
	scanner := bufio.NewScanner(openFile)
	lineNum := 0
	for scanner.Scan() {
		switch lineNum {
		case 0:
			blo.Thumbnail = scanner.Text()
		case 1:
			blo.Title = scanner.Text()
		case 2:
			blo.Preview = scanner.Text()
		default:
			text = append(text, scanner.Text())
		}
		lineNum += 1
	}
	blo.Text = strings.Join(text, " ")

	template := template.Must(template.ParseFiles("pages/blog.html"))
	responseData := map[string]interface{}{
		"Title": blo.Title,
		"Id": blo.Id,
		"Text": blo.Text,
		"Thumbnail": blo.Thumbnail,
	}

	template.Execute(res, responseData)

}

func blogs(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	files, err := ioutil.ReadDir("blogs")
	if err != nil {
		fmt.Errorf("No blogs in blogs dir")
	}

	var blogs []Blog
	for _, blog := range files {
		blo := Blog{}
		blo.Id = blog.Name()

		filePath := "blogs/" + blog.Name()
		openFile, err := os.Open(filePath)
		if err != nil {
			fmt.Errorf("Not able to read file at=%T", filePath)
		}
		defer openFile.Close()

		var text []string
		scanner := bufio.NewScanner(openFile)
		lineNum := 0
		for scanner.Scan() {
			switch lineNum {
			case 0:
				blo.Thumbnail = scanner.Text()
			case 1:
				blo.Title = scanner.Text()
			case 2:
				blo.Preview = scanner.Text()
			default:
				text = append(text, scanner.Text())
			}
			lineNum += 1
		}
		blo.Text = strings.Join(text, " ")
		blogs = append(blogs, blo)
	}

	template := template.Must(template.ParseFiles("pages/blogs.html"))
	responseData := map[string][]Blog{
		"blogs": blogs,
	}

	template.Execute(res, responseData)
}

func seeLess(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	ids := req.URL.Query()["id"]

	id := ids[0]
	project := getProjectById(id)
	template := template.Must(template.ParseFiles("responses/seeLess.html"))
	responseData := map[string]interface{}{
		"ProjectTitle": project.ProjectTitle,
		"ShortSummary": project.ShortSummary,
		"Id":           project.Id,
		"Url":          project.Url,
	}

	template.Execute(res, responseData)
}

func seeMore(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	ids, err := req.URL.Query()["id"]
	if err {
		fmt.Errorf("Id missing from url params!")
	}

	id := ids[0]

	project := getProjectById(id)
	template := template.Must(template.ParseFiles("responses/seeMore.html"))
	responseData := map[string]interface{}{
		"ProjectTitle":    project.ProjectTitle,
		"LongDescription": project.LongDescription,
		"Id":              project.Id,
		"Url":             project.Url,
	}

	template.Execute(res, responseData)

}

func hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	projects := []Project{
		{
			ProjectTitle:    "chatApp",
			ShortSummary:    "Live chat application built with the T3 stack",
			LongDescription: "Currently building a hybrid clone thats based on WhatsApp and Discord for me and my friends to use. The project is written in TypeScript, using Next.js as the framework. The app uses MySQL as the database, tRPC for client to interact with the server in a type safe manner. I wanted to plau with some AI stuff with the chat application so currently the users can invite ChatGPT to be part of their conversation by tagging AI into the chat with @ai. The application connects to a separate WebSocket server I had made in TypeScript which is used to let users know who is typing and to invalidate the chats cache whenever a new message arrives in the chat. This also served as a gateway to play funny sounds to all chat participants with the press of a sound board button as part of the 'chatting experience'.",
			Id:              "chatapp",
			Url:             "https://chat-app-tau-teal.vercel.app/",
		},

		{
			ProjectTitle:    "HaikuForu",
			ShortSummary:    "A website that generates Haikus with the use of AI",
			LongDescription: "A web app I made with Next.js and deployed on Vercel, written in TypeScript using React. I wanted to use OpenAI's API so the application uses this as a haiku generator for  any topic the user inputs",
			Id:              "haiku",
			Url:             "https://haiku-for-u.vercel.app/",
		},
	}

	template := template.Must(template.ParseFiles("pages/hello.html"))
	mainPgData := map[string][]Project{
		"projects": projects,
	}

	template.Execute(res, mainPgData)

}

func main() {
	// TOTO PLAN
	// 1. markdown to html parser
	// 2. pass blog markdown files through parser to get HTML
	// 3. serve html into page

	fileServe := http.FileServer(http.Dir("."))
	http.Handle("/pages", http.StripPrefix("/", fileServe))

	http.HandleFunc("/", hello)
	http.HandleFunc("/seemore", seeMore)
	http.HandleFunc("/seeless", seeLess)
	http.HandleFunc("/blogs", blogs)
	http.HandleFunc("/resume", resume)
	http.HandleFunc("/blog", blog)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
