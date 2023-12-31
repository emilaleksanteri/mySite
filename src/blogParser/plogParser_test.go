package blogParser

import (
	"fmt"
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

func TestParagraph(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{"&Paragraph test&", "<p id='paragraph'>Paragraph test</p>"},
		{"&This is my essay paragraph&", "<p id='paragraph'>This is my essay paragraph</p>"},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 1 {
			t.Fatalf("Expectd blog paragrapsh to be length 1, got=%T", len(blog.Paragraphs))
		}
		paragraph := blog.Paragraphs[0].Text
		if paragraph != template.HTML(tt.expectedHTML) {
			t.Fatalf("paragraph did not match the expected %q got=%q", tt.expectedHTML, string(paragraph))
		}
	}
}

func TestParagraphs(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML []string
	}{
		{
			"&Paragraph test&&Paragraph test&",
			[]string{"<p id='paragraph'>Paragraph test</p>", "<p id='paragraph'>Paragraph test</p>"}},
		{
			"&This is my essay paragraph&&This is my essay paragraph&",
			[]string{"<p id='paragraph'>This is my essay paragraph</p>", "<p id='paragraph'>This is my essay paragraph</p>"}},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 2 {
			t.Fatalf("Expectd blog paragrapsh to be length 1, got=%T", len(blog.Paragraphs))
		}

		for i, para := range blog.Paragraphs {
			paragraph := para.Text
			if paragraph != template.HTML(tt.expectedHTML[i]) {
				t.Fatalf("paragraph did not match the expected %q got=%q", tt.expectedHTML[i], string(paragraph))
			}
		}
	}
}

func TestQuoteInParagraph(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{`&Paragraph "quote" test&`, "<p id='paragraph'>Paragraph <p id='quote'>quote</p> test</p>"},
		{`&This is my "quote" essay paragraph&`, "<p id='paragraph'>This is my <p id='quote'>quote</p> essay paragraph</p>"},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 1 {
			t.Errorf("Expected just a single paragraph got=%T", len(blog.Paragraphs))
		}

		paragraph := blog.Paragraphs[0]
		if paragraph.Text != template.HTML(tt.expectedHTML) {
			t.Errorf("HTML did not match expected %q and got %q", tt.expectedHTML, string(paragraph.Text))
		}
	}
}

func TestLinkInParagraph(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{
			`&Paragraph [link]{yomama.com} test&`,
			"<p id='paragraph'>Paragraph <a id='link' href=yomama.com>link</a> test</p>",
		},
		{
			`&This is my [yomama.com] essay paragraph&`,
			"<p id='paragraph'>This is my <a id='link' href=yomama.com>yomama.com</a> essay paragraph</p>",
		},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 1 {
			t.Fatalf("Expected to get a single paragraph got=%T", len(blog.Paragraphs))
		}

		paragraph := blog.Paragraphs[0]
		if paragraph.Text != template.HTML(tt.expectedHTML) {
			t.Fatalf("Blog html does not match expected %q got=%q", tt.expectedHTML, string(paragraph.Text))
		}
	}
}

func TestImageParagraph(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{
			`&Paragraph *cat.png* test&`,
			"<p id='paragraph'>Paragraph <img id='image' src=cat.png alt='inline image' /> test</p>",
		},
		{
			`&This is my *inline.jpg* essay paragraph&`,
			"<p id='paragraph'>This is my <img id='image' src=inline.jpg alt='inline image' /> essay paragraph</p>",
		},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 1 {
			t.Fatalf("Expected to get a single paragraph got=%T", len(blog.Paragraphs))
		}

		paragraph := blog.Paragraphs[0]
		if paragraph.Text != template.HTML(tt.expectedHTML) {
			t.Fatalf("Blog html does not match expected %q got=%q", tt.expectedHTML, string(paragraph.Text))
		}
	}
}

func TestImageWithTextParagraph(t *testing.T) {
	tests := []struct {
		input        string
		expectedHTML string
	}{
		{
			`&Paragraph +*cat.png*Image of a cat+ test&`,
			"<p id='paragraph'>Paragraph <div id='imageWithText'><img id='image' src=cat.png alt='inline image' /><p id='blockText'>Image of a cat</p></div> test</p>",
		},
		{
			`&This is my +Inline photo*inline.jpg*+ essay paragraph&`,
			"<p id='paragraph'>This is my <div id='imageWithText'><p id='blockText'>Inline photo</p><img id='image' src=inline.jpg alt='inline image' /></div> essay paragraph</p>",
		},
	}

	for _, tt := range tests {
		blog := BlogText{}
		blog.ParseMarkDown(tt.input)

		if len(blog.Paragraphs) != 1 {
			t.Fatalf("Expected to get a single paragraph got=%T", len(blog.Paragraphs))
		}

		paragraph := blog.Paragraphs[0]
		if paragraph.Text != template.HTML(tt.expectedHTML) {
			t.Fatalf("Blog html does not match expected %q got=%q", tt.expectedHTML, string(paragraph.Text))
		}
		fmt.Println(string(paragraph.Text))
	}
}
