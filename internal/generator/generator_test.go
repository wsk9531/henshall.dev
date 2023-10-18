package generator_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"testing/fstest"

	approvals "github.com/approvals/go-approval-tests"
	generator "github.com/wsk9531/henshall.dev/internal/generator"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
URL: first.html
Tags: tdd, go
---
Let's get 
this M.O.N.E.Y`
		secondBody = `Title: Post 2
Description: Description 2
URL: second.html
Tags: tdd, go
---
Ok
I 
Pull 
Up`
	)

	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := generator.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}
	assertPost(t, posts[0], generator.Post{
		Title:       "Post 1",
		Description: "Description 1",
		URL:         "first.html",
		Tags:        []string{"tdd", "go"},
		Body: `Let's get 
this M.O.N.E.Y`,
	})
}

func assertPost(t *testing.T, got, want generator.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestRender(t *testing.T) {
	var (
		aPost = generator.Post{
			Title: "First Post",
			Body: `#Hello world! 
This is my *very cool* blog post`,
			URL:         "testpost",
			Description: "The first of many lovely blog posts",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := generator.NewPostRenderer()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		if err := postRenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})

	t.Run("it renders an index of posts", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []generator.Post{aPost, {Title: "Post 2"}, {Title: "Post 3"}}

		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var (
		aPost = generator.Post{
			Title:       "hello world",
			Body:        "this is a post",
			URL:         "testpage",
			Description: "this is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := generator.NewPostRenderer()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = postRenderer.Render(io.Discard, aPost)
	}
}
