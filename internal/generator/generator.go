package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"html/template"

	ui "github.com/wsk9531/henshall.dev/ui"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// post.go
type Post struct {
	Title       string
	Description string
	Body        string
	Tags        []string
}

func (p Post) SanitisedTitle() string {
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagSeperator         = "Tags: "
	bodySeparator        = "---"
)

func NewPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{Title: readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagSeperator), ", "),
		Body:        readBody(scanner)}, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

// reader.go
func NewPostsFromFS(filesystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(filesystem, f.Name())
		if err != nil {
			return nil, err // todo: needs clarification, fail if one file fails or all ?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	return NewPost(postFile)
}

// Renderer.go
type PostRenderer struct {
	templ    *template.Template
	mdParser goldmark.Markdown
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(ui.Files, "templates/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	parser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("nord"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
					chromahtml.WithClasses(false),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	return &PostRenderer{templ: templ, mdParser: parser}, nil
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.templ.ExecuteTemplate(w, "blog.tmpl.html", newPostVM(p, r))
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return r.templ.ExecuteTemplate(w, "index.tmpl.html", posts)
}

type postViewModel struct {
	Post
	HTMLBody template.HTML
}

func newPostVM(p Post, r *PostRenderer) postViewModel {
	vm := postViewModel{Post: p}
	var buf bytes.Buffer
	if err := r.mdParser.Convert([]byte(p.Body), &buf); err != nil {
		panic(err)
	}
	vm.HTMLBody = template.HTML(buf.Bytes())
	return vm
}
