package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"github.com/emilaleksanteri/mySite/tree/main/src/blogParser"
	"net/http"
	"os"
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
			LongDescription: "Currently building a hybrid clone thats based on WhatsApp and Discord for me and my friends to use. The project is written in TypeScript, using Next.js as the framework. The app uses MySQL as the database, tRPC for client to interact with the server in a type safe manner. I wanted to plau with some AI stuff with the chat application so currently the users can invite ChatGPT to be part of their conversation by tagging AI into the chat with @ai. The application connects to a separate WebSocket server I had made in TypeScript which is used to let users know who is typing and to invalidate the chats cache whenever a new message arrives in the chat. This also served as a gateway to play funny sounds to all chat participants with the press of a sound board button as part of the 'chatting experience'. I also built a WebSocket server that I deployed with fly.io, it is still very bare bones server, I am planning to redo it in Elixir or Rust once I have a bit of extra time.",
			Id:              "chatapp",
			Url:             "https://chat-app-tau-teal.vercel.app/",
		},
		{
			ProjectTitle:    "BlogApp",
			ShortSummary:    "A basic React CRUD app that uses a REST API to store and retrieve and share posts.",
			LongDescription: "Application where people can post links to their favourite blogs.Users can then like and comment these on the application. Application supports CRUD, the backend is built on express and uses MongoDB, the front end is fully in React. For testing I used Jest and Cypress.",
			Id:              "blogapp",
			Url:             "https://github.com/emilaleksanteri/BlogApp",
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

type Error struct {
	Error    int
	ErrorMsg string
}

func displayError(res http.ResponseWriter, req *http.Request, error *Error) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	template := template.Must(template.ParseFiles("pages/error.html"))
	responseData := map[string]interface{}{
		"Error":    error.Error,
		"ErrorMsg": error.ErrorMsg,
	}

	template.Execute(res, responseData)
}

func resume(res http.ResponseWriter, req *http.Request) {
	resume, err := os.Open("cv/Emil-Lystimaki-CV.pdf")
	if err != nil {
		fmt.Println("Could not get cv")
		error := Error{
			Error:    500,
			ErrorMsg: "Could not get the CV :(",
		}
		displayError(res, req, &error)
		return
	}
	defer resume.Close()

	res.Header().Set(
		"Content-Type",
		"application/pdf",
	)
	if _, err := io.Copy(res, resume); err != nil {
		fmt.Println("Could not serve resume bruh moment")
		error := Error{
			Error:    500,
			ErrorMsg: "Could not get the CV :(",
		}
		displayError(res, req, &error)
		return
	}
}

type Blog struct {
	Thumbnail template.HTML
	Title     template.HTML
	Preview   template.HTML
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
		fmt.Println("Not able to read file", fileName)
		error := Error{
			Error:    404,
			ErrorMsg: "Could not find the blog you were looking for... :(",
		}
		displayError(res, req, &error)
		return
	}

	defer openFile.Close()

	var text string
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		text += scanner.Text()
	}

	myBlog := blogParser.BlogText{}
	myBlog.ParseMarkDown(text)

	template := template.Must(template.ParseFiles("pages/blog.html"))
	responseData := map[string]interface{}{
		"Title":      myBlog.Title,
		"TitleText":  myBlog.TitleData,
		"Id":         id,
		"Thumbnail":  myBlog.CoverPhoto,
		"Paragraphs": myBlog.Paragraphs,
	}

	template.Execute(res, responseData)

}

func blogs(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	files, err := os.ReadDir("blogs")
	if err != nil {
		fmt.Println("No blogs in blogs dir")
		error := Error{
			Error:    404,
			ErrorMsg: "There are no blogs to display, come back later when there are.",
		}
		displayError(res, req, &error)
		return
	}

	var blogs []Blog
	for _, blog := range files {
		filePath := "blogs/" + blog.Name()
		openFile, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Not able to read file at", filePath)
			error := Error{
				Error:    500,
				ErrorMsg: "Something went terribly wrong when trying to display blogs!",
			}
			displayError(res, req, &error)
			return
		}
		defer openFile.Close()

		var text string
		scanner := bufio.NewScanner(openFile)
		for scanner.Scan() {
			text += scanner.Text()
		}

		myBlog := blogParser.BlogText{}
		myBlog.ParseMarkdownPreview(text)

		blo := Blog{
			Title:     myBlog.Title,
			Preview:   myBlog.Preview,
			Thumbnail: myBlog.CoverPhoto,
			Id:        blog.Name(),
		}

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

	ids := req.URL.Query()["id"]
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
			LongDescription: "Currently building a hybrid clone thats based on WhatsApp and Discord for me and my friends to use. The project is written in TypeScript, using Next.js as the framework. The app uses MySQL as the database, tRPC for client to interact with the server in a type safe manner. I wanted to plau with some AI stuff with the chat application so currently the users can invite ChatGPT to be part of their conversation by tagging AI into the chat with @ai. The application connects to a separate WebSocket server I had made in TypeScript which is used to let users know who is typing and to invalidate the chats cache whenever a new message arrives in the chat. This also served as a gateway to play funny sounds to all chat participants with the press of a sound board button as part of the 'chatting experience'. I also built a WebSocket server that I deployed with fly.io, it is still very bare bones server, I am planning to redo it in Elixir or Rust once I have a bit of extra time.",
			Id:              "chatapp",
			Url:             "https://chat-app-tau-teal.vercel.app/",
		},
		{
			ProjectTitle:    "BlogApp",
			ShortSummary:    "A basic React CRUD app that uses a REST API to store and retrieve and share posts.",
			LongDescription: "Application where people can post links to their favourite blogs.Users can then like and comment these on the application. Application supports CRUD, the backend is built on express and uses MongoDB, the front end is fully in React. For testing I used Jest and Cypress.",
			Id:              "blogapp",
			Url:             "https://github.com/emilaleksanteri/BlogApp",
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

func serve(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	switch path {
	case "/":
		hello(res, req)
	case "/seemore":
		seeMore(res, req)
	case "/seeless":
		seeLess(res, req)
	case "/blogs":
		blogs(res, req)
	case "/blog":
		blog(res, req)
	case "/resume":
		resume(res, req)
	default:
		error := Error{
			Error:    404,
			ErrorMsg: "Page not found, wonder where it went!",
		}
		displayError(res, req, &error)
	}
}

func main() {
	fileServe := http.FileServer(http.Dir("."))
	http.Handle("/pages", http.StripPrefix("/", fileServe))

	http.HandleFunc("/", serve)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
