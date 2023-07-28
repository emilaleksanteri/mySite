package blogParser

import (
	"html/template"
	"testing"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
		expectedData string
	}{
		{"#This is a title#", "<h2 id='title'>This is a title</h2>", "This is a title"},
		{"#a title #", "<h2 id='title'>a title </h2>", "a title "},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if blog.Title != template.HTML(tt.expectedHTML) {
			t.Fatalf("Title html was not right, expected %q and got=%q", string(blog.Title), tt.expectedHTML)
		}

		if blog.TitleData != tt.expectedData {
			t.Fatalf("Title data is not correct expected=%q, got=%q", tt.expectedData, blog.TitleData)
		}
	}
}

func TestCover(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{"@./images/cat.png@", "<img id='cover' src=./images/cat.png alt='cover photo' />"},
		{"@https://google.com/cat.png@", "<img id='cover' src=https://google.com/cat.png alt='cover photo' />"},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if blog.CoverPhoto != template.HTML(tt.expectedHTML) {
			t.Fatalf("Cover photo does not match expected=%q got=%q", tt.expectedHTML, string(blog.CoverPhoto))
		}
	}
}

func TestPreview(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{"$This is a preview$", "<p id='preview'>This is a preview</p>"},
		{"$Preview this lel$", "<p id='preview'>Preview this lel</p>"},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if blog.Preview != template.HTML(tt.expectedHTML) {
			t.Fatalf("Blog preview did not match expected preview %q got=%q", tt.expectedHTML, string(blog.Preview))
		}
	}
}
