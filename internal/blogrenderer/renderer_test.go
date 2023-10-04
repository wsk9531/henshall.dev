package blogrenderer_test

import (
	"bytes"
	"io"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/wsk9531/henshall.dev/internal/blogrenderer"
)

func TestRender(t *testing.T) {
	var ( 
		aPost = blogrenderer.Post{
			Title: "First Post", 
			Body: `#Hello world! 
This is my *very cool* blog post`,
			Description: "The first of many lovely blog posts",
			Tags: []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T){
		buf := bytes.Buffer{}
		if err := postRenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})

	t.Run("it renders an index of posts", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []blogrenderer.Post{aPost, {Title: "Post 2"}, {Title: "Post 3"}}
	
		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var ( 
		aPost = blogrenderer.Post{
			Title: "hello world", 
			Body: "this is a post",
			Description: "this is a description",
			Tags: []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postRenderer.Render(io.Discard, aPost)
	}
}

