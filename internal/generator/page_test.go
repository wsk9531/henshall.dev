package generator

import (
	"io"
	"testing"
	"testing/fstest"
)

func TestFileContentSplitting(t *testing.T) {

	fs := SetupFiles()

	t.Run("it splits a valid frontmatter + markdown file", func(t *testing.T) {
		const wantedFrontmatter string = `---
title: Post 1
published: 2023-01-01
---`
		const wantedContent string = `Really Good Content`

		file, err := fs.Open("valid.md")
		assertNoErr(t, err)
		b, err := io.ReadAll(file)
		assertNoErr(t, err)

		got, err := separateFrontMatterAndContent(string(b))
		want := markdownDocument{wantedFrontmatter, wantedContent}

		assertNoErr(t, err)
		AssertEqual[markdownDocument](t, got, want)
	})

	t.Run("it errors when opening dashes are not at start of file", func(t *testing.T) {
		file, err := fs.Open("missing_opening_dashes.md")
		assertNoErr(t, err)
		b, err := io.ReadAll(file)
		assertNoErr(t, err)

		_, got := separateFrontMatterAndContent(string(b))
		want := ErrMalformedFrontmatterDashes

		assertError(t, got, want)
	})

	t.Run("it errors when closing dashes are not present in file", func(t *testing.T) {
		file, err := fs.Open("missing_closing_dashes.md")
		assertNoErr(t, err)
		b, err := io.ReadAll(file)
		assertNoErr(t, err)

		_, got := separateFrontMatterAndContent(string(b))
		want := ErrMalformedFrontmatterDashes

		assertError(t, got, want)
	})
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func SetupFiles() fstest.MapFS {
	const validMdFile string = `---
title: Post 1
published: 2023-01-01
---
Really Good Content`

	const missingOpeningDashes = `
title: Post 1
published: 2023-01-01
---
Really Good Content`

	const missingClosingDashes = `---
title: Post 1
published: 2023-01-01
Really Good Content`

	fs := fstest.MapFS{
		"valid.md":                  {Data: []byte(validMdFile)},
		"missing_opening_dashes.md": {Data: []byte(missingOpeningDashes)},
		"missing_closing_dashes.md": {Data: []byte(missingClosingDashes)},
	}

	return fs
}
