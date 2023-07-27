package blogParser

import (
	"fmt"
	"html/template"
)

type Paragraph struct {
	Text template.HTML
}

type BlogText struct {
	Title      template.HTML
	TitleData  string
	CoverPhoto template.HTML
	Paragraphs []Paragraph
	Input      string
	Curr       int
	Next       int
}

const (
	Title         = '#'
	LinkStart     = '['
	LinkEnd       = ']'
	LinkNameStart = '{'
	LinkNameEnd   = '}'
	Image         = '*'
	Para          = '&'
	Quote         = '$'
	CoverPhoto    = '@'
	EOF           = ';'
)

func (b *BlogText) ParseMarkDown(input string) {
	b.Input = input
	if len(b.Input) < 2 {
		return
	}

	b.Curr = 0
	for b.Input[b.Curr] != EOF {
		switch b.Input[b.Curr] {
		case '#':
			b.parseTitle()
		case '@':
			b.parseCoverPhoto()
		}
	}

}

func (b *BlogText) parseTitle() {
	text := "<h2>"
	var titleData string
	b.Curr += 1
	for b.Input[b.Curr] != Title {
		text += string(b.Input[b.Curr])
		titleData += string(b.Input[b.Curr])
		b.Curr += 1
	}
	b.Curr += 1
	text += "</h2>"
	b.Title = template.HTML(text)
	b.TitleData = titleData
}

func (b *BlogText) parseParagraph() {}

func (b *BlogText) parseCoverPhoto() {
	b.Curr += 1
	var url string
	for b.Input[b.Curr] != CoverPhoto {
		url += string(b.Input[b.Curr])
		b.Curr += 1
	}
	b.Curr += 1

	imgElement := "<img id='cover' src=" + url + " alt='cover photo' />"
	b.CoverPhoto = template.HTML(imgElement)
}

func (b *BlogText) parseImage() {}

// Goal:
// Parse input line by line
// Have case for each type of markup token
// Depending on token type, add correct HTML with style to Title or make a Paragraph object
// If paragraph contains a link, place an <a> tag into the paragraph
// Final output should be a title in HTML
// And an array of paragrapsh in HTML

// 1. Input titles into HTML succesfully
// 2. Parse input paragraphs into a paragraph struct succesfully
// 3. Combine
