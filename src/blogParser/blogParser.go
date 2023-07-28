package blogParser

import (
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
	CurrPos    int
	Curr       byte
	Next       int
	Preview    template.HTML
}

const (
	Title         = '#'
	LinkStart     = '['
	LinkEnd       = ']'
	LinkNameStart = '{'
	LinkNameEnd   = '}'
	Image         = '*'
	Para          = '&'
	Quote         = '"'
	CoverPhoto    = '@'
	Preview       = '$'
	EOF           = 0
)

func (b *BlogText) ParseMarkDown(input string) {
	b.Input = input
	if len(b.Input) < 2 {
		return
	}

	b.CurrPos = 0
	b.Curr = b.Input[b.CurrPos]

	for b.Curr != EOF {
		switch b.Curr {
		case '#':
			b.parseTitle()
		case '@':
			b.parseCoverPhoto()
		case '$':
			b.parsePreview()
		}
	}
}

func (b *BlogText) nextByte() {
	b.CurrPos += 1
	if b.CurrPos >= len(b.Input) {
		b.Curr = 0
		return
	}
	b.Curr = b.Input[b.CurrPos]
}

func (b *BlogText) parseTitle() {
	text := "<h2 id='title'>"
	var titleData string
	b.nextByte()
	for b.Curr != Title {
		text += string(b.Curr)
		titleData += string(b.Curr)
		b.nextByte()
	}
	b.nextByte()
	text += "</h2>"
	b.Title = template.HTML(text)
	b.TitleData = titleData
}

func (b *BlogText) parsePreview() {
	b.nextByte()
	preview := "<p id='preview'>"
	for b.Curr != Preview {
		preview += string(b.Curr)
		b.nextByte()
	}
	b.nextByte()
	preview += "</p>"
	b.Preview = template.HTML(preview)
}

func (b *BlogText) parseParagraph() {}

func (b *BlogText) parseCoverPhoto() {
	b.nextByte()
	var url string
	for b.Curr != CoverPhoto {
		url += string(b.Curr)
		b.nextByte()
	}
	b.nextByte()

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
