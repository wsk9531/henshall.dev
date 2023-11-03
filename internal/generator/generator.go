package generator

import (
	"bytes"
	"io"
	"io/fs"
	"sort"

	"html/template"

	ui "github.com/wsk9531/henshall.dev/ui"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// reader.go
func NewPagesFromFS(filesystem fs.FS) ([]Page, error) {
	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}
	var pages []Page
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		page, err := getPage(filesystem, f.Name())
		if err != nil {
			return nil, err // todo: needs clarification, fail if one file fails or all ?
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func getPage(fileSystem fs.FS, fileName string) (Page, error) {
	pageFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Page{}, err
	}
	defer pageFile.Close()

	return newPage(pageFile)
}

// Renderer.go
type PageRenderer struct {
	templ    *template.Template
	mdParser goldmark.Markdown
}

func NewPageRenderer() (*PageRenderer, error) {
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
	return &PageRenderer{templ: templ, mdParser: parser}, nil
}

func (r *PageRenderer) Render(w io.Writer, p Page) error {
	return r.templ.ExecuteTemplate(w, "page.tmpl.html", newPageVM(p, r))
}

func (r *PageRenderer) RenderBlog(w io.Writer, p Page) error {
	return r.templ.ExecuteTemplate(w, "blog.tmpl.html", newPageVM(p, r))
}

func (r *PageRenderer) RenderIndex(w io.Writer, pages []Page) error {
	return r.templ.ExecuteTemplate(w, "index.tmpl.html", newIndexVM(pages, r))
}

type pageViewModel struct {
	Page
	HTMLBody template.HTML
}

func newPageVM(p Page, r *PageRenderer) pageViewModel {
	vm := pageViewModel{Page: p}
	var buf bytes.Buffer
	if err := r.mdParser.Convert([]byte(p.Body), &buf); err != nil {
		panic(err)
	}
	vm.HTMLBody = template.HTML(buf.Bytes())
	return vm
}

// An indexViewModel wraps a list all of the posts.
//
// This is a workaround for errors in shared templates like:
// executing "top" at <.Title>: can't evaluate field Title in type []generator.Page
// where we use the top template for both Page and Index Generation.
//
// TODO: Review this approach vs just using different partial templates.
type indexViewModel struct {
	Pages       []Page
	Title       string
	Description string
}

func newIndexVM(p []Page, r *PageRenderer) indexViewModel {
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Published.After(p[j].Published)
	})
	vm := indexViewModel{Pages: p}
	vm.Title = "Blog Index"
	vm.Description = "Blog Articles (sorted by publication date)"
	return vm
}
