package blogParser

import (
	"html/template"
)

type Paragraph struct {
	Text     template.HTML
	TextData string
}

type Link struct {
	LinkText string
	Link     string
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
	LinkTextStart = '['
	LinkTextEnd   = ']'
	LinkStart     = '{'
	LinkEnd       = '}'
	Image         = '*'
	Para          = '&'
	Quote         = '"'
	CoverPhoto    = '@'
	Preview       = '$'
	EOF           = 0
)

func (b *BlogText) ParseMarkDown(input string) {
	b.Input = input
	if len(b.Input) < 3 {
		return
	}

	b.CurrPos = 0
	b.Curr = b.Input[b.CurrPos]

	for b.Curr != EOF {
		switch b.Curr {
		case Title:
			b.parseTitle()
		case CoverPhoto:
			b.parseCoverPhoto()
		case Preview:
			b.parsePreview()
		case Para:
			b.parseParagraph()
		default:
			// don't wanna parse non markdown things
			b.nextByte()
		}
	}
}

func (b *BlogText) ParseMarkdownPreview(input string) {
	b.Input = input
	if len(b.Input) < 3 {
		return
	}

	b.CurrPos = 0
	b.Curr = b.Input[b.CurrPos]

	for b.Curr != EOF {
		switch b.Curr {
		case Title:
			b.parseTitle()
		case CoverPhoto:
			b.parseCoverPhoto()
		case Preview:
			b.parsePreview()
		default:
			b.nextByte()
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

func (b *BlogText) parseParagraph() {
	b.nextByte()
	paragraph := Paragraph{TextData: "<p id='paragraph'>"}
	for b.Curr != Para {
		switch b.Curr {
		case Quote:
			b.parseQuote(&paragraph)
		case LinkTextStart:
			b.parseLink(&paragraph)
		case Image:
			b.parseImage(&paragraph)
		default:
			paragraph.TextData += string(b.Curr)
			b.nextByte()
		}
	}

	b.nextByte()
	paragraph.TextData += "</p>"
	paragraph.Text = template.HTML(paragraph.TextData)
	b.Paragraphs = append(b.Paragraphs, paragraph)
}

func (b *BlogText) parseLink(paragraph *Paragraph) {
	link := Link{}
	b.nextByte()

	for b.Curr != LinkTextEnd {
		link.LinkText += string(b.Curr)
		b.nextByte()
	}

	b.nextByte()
	if b.Curr != LinkStart {
		link.Link = link.LinkText
	} else {
		b.nextByte()
		for b.Curr != LinkEnd {
			link.Link += string(b.Curr)
			b.nextByte()
		}
		b.nextByte()
	}

	element := "<a id='link' href=" + link.Link + ">" + link.LinkText + "</a>"
	paragraph.TextData += element
}

func (b *BlogText) parseQuote(paragraph *Paragraph) {
	b.nextByte()
	quote := "<p id='quote'>"
	for b.Curr != Quote {
		quote += string(b.Curr)
		b.nextByte()
	}

	b.nextByte()
	quote += "</p>"
	paragraph.TextData += quote
}

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

func (b *BlogText) parseImage(paragraph *Paragraph) {
	b.nextByte()
	var url string
	for b.Curr != Image {
		url += string(b.Curr)
		b.nextByte()
	}

	b.nextByte()
	imgElement := "<img id='image' src=" + url + " alt='inline image' />"
	paragraph.TextData += imgElement
}

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
