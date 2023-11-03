package generator

import (
	"errors"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const YAMLSeparator string = "---"

type Page struct {
	Frontmatter
	Body string
}

func newPage(r io.Reader) (Page, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return Page{}, err
	}
	doc, err := separateFrontMatterAndContent(string(b))
	if err != nil {
		return Page{}, err
	}

	fm, err := newFrontmatter(doc.FrontContent)
	if err != nil {
		return Page{}, err
	}

	page := Page{Frontmatter: fm, Body: doc.BodyContent}
	return page, nil
}

func newFrontmatter(s string) (Frontmatter, error) {
	front := Frontmatter{}

	err := yaml.Unmarshal([]byte(s), &front)
	if err != nil {
		return Frontmatter{}, err
	}

	return front, nil
}

type Frontmatter struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Url         string    `yaml:"url"`
	Tags        []string  `yaml:"tags"`
	Published   time.Time `yaml:"published"`
	Updated     time.Time `yaml:"updated"`
}

// internal representation of a file with yaml frontmatter and text content,
// before unmarshal'ing into our Frontmatter struct
type markdownDocument struct {
	FrontContent string
	BodyContent  string
}

var ErrMalformedFrontmatterDashes = errors.New("malformed frontmatter opening/closing dashes in file")

// separateFrontMatterAndContent returns a struct of type markdownDocument,
// which is an intermediary to hold frontmatter content before marshaling
func separateFrontMatterAndContent(data string) (markdownDocument, error) {
	// Check opening --- is present
	if !strings.HasPrefix(data, YAMLSeparator) {
		return markdownDocument{}, ErrMalformedFrontmatterDashes
	}

	// Check for closing --- is present
	closingDashesIdx := strings.Index(data[len(YAMLSeparator):], YAMLSeparator)
	if closingDashesIdx == -1 {
		return markdownDocument{}, ErrMalformedFrontmatterDashes
	}

	// idx = opening dashes dashes + frontmatter content + end of closing dashes
	contentStartIdx := closingDashesIdx + 2*len(YAMLSeparator)

	frontmatter := strings.TrimSpace(data[:contentStartIdx])
	body := strings.TrimSpace(data[contentStartIdx:])

	// Return a MarkdownDocument struct with front matter and content.
	return markdownDocument{FrontContent: frontmatter, BodyContent: body}, nil
}
