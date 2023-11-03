package generator_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/wsk9531/henshall.dev/internal/generator"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `---
title: Post 1
description: Description 1
url: first
tags: 
  - tdd 
  - go
---
Let's get 
this M.O.N.E.Y`
		secondBody = `---
title: Post 2
description: Description 2
url: second
tags: 
  - tdd 
  - go
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

	posts, err := generator.NewPagesFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}
	assertPost(t, posts[0], generator.Page{
		Frontmatter: generator.Frontmatter{Title: "Post 1", Description: "Description 1", Url: "first", Tags: []string{"tdd", "go"}},
		Body: `Let's get 
this M.O.N.E.Y`,
	})
}

func assertPost(t *testing.T, got, want generator.Page) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestRender(t *testing.T) {
	var (
		aPost = generator.Page{
			Frontmatter: generator.Frontmatter{
				Title:       "First Post",
				Url:         "testpost",
				Description: "The first of many lovely blog posts",
				Tags:        []string{"go", "tdd"},
				Published:   time.Date(2023, 11, 3, 0, 0, 0, 0, time.UTC),
			},
			Body: `#Hello world! 
This is my *very cool* blog post`,
		}
		aSecondPost = generator.Page{
			Frontmatter: generator.Frontmatter{
				Title:     "Second Post",
				Url:       "testpost-2",
				Published: time.Date(1999, 25, 12, 0, 0, 0, 0, time.UTC),
			},
		}
	)

	postRenderer, err := generator.NewPageRenderer()
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
		posts := []generator.Page{aPost, aSecondPost}

		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var (
		aPost = generator.Page{
			Frontmatter: generator.Frontmatter{
				Title:       "xyz",
				Description: "zyzz",
			},
			Body: "xyz",
		}
	)

	postRenderer, err := generator.NewPageRenderer()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = postRenderer.Render(io.Discard, aPost)
	}
}
