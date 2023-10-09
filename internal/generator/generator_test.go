package generator_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"testing/fstest"

	approvals "github.com/approvals/go-approval-tests"
	ssg "github.com/wsk9531/henshall.dev/internal/generator"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Let's get 
this M.O.N.E.Y`
		secondBody = `Title: Post 2
Description: Description 2
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

	posts, err := ssg.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}
	assertPost(t, posts[0], ssg.Post{
		Title: "Post 1", 
		Description: "Description 1", 
		Tags: []string{"tdd", "go"},
		Body: `Let's get 
this M.O.N.E.Y`,
	})
}

func assertPost(t *testing.T, got, want ssg.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestRender(t *testing.T) {
	var ( 
		aPost = ssg.Post{
			Title: "First Post", 
			Body: `#Hello world! 
This is my *very cool* blog post`,
			Description: "The first of many lovely blog posts",
			Tags: []string{"go", "tdd"},
		}
	)

	postRenderer, err := ssg.NewPostRenderer()
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
		posts := []ssg.Post{aPost, {Title: "Post 2"}, {Title: "Post 3"}}
	
		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var ( 
		aPost = ssg.Post{
			Title: "hello world", 
			Body: "this is a post",
			Description: "this is a description",
			Tags: []string{"go", "tdd"},
		}
	)

	postRenderer, err := ssg.NewPostRenderer()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postRenderer.Render(io.Discard, aPost)
	}
}

