package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/wsk9531/henshall.dev/internal/generator"
)

const (
	DEFAULT_CONTENT_DIR string = "pages"
	DEFAULT_OUTPUT_DIR  string = "dist"
	DEFAULT_PORT        string = "5000"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	contentDir := *generateCmd.String("content", DEFAULT_CONTENT_DIR, "path to content for rendering")
	outputDir := *generateCmd.String("output", DEFAULT_OUTPUT_DIR, "path to rendered site files")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveDir := *serveCmd.String("dir", DEFAULT_OUTPUT_DIR, "path to rendered site files")
	port := *serveCmd.String("port", DEFAULT_PORT, "localhost port to serve content")

	if len(os.Args) < 2 {
		fmt.Println("expected 'generate' or 'serve' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "generate":
		generateCmd.Parse(os.Args[2:]) //nolint:errcheck // this error is handled by the Flags' parser fail behaviour for now

		fmt.Println("subcommand 'generate'")
		fmt.Println(" content dir:", contentDir)
		fmt.Println(" output dir:", outputDir)

		err := os.Mkdir(outputDir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		posts, err := generator.NewPostsFromFS(os.DirFS(contentDir))
		if err != nil {
			log.Fatal(err)
		}
		postRenderer, err := generator.NewPostRenderer()
		if err != nil {
			log.Fatal(err)
		}

		for _, page := range posts {
			// Render content into buffer
			buf := bytes.Buffer{}
			if err := postRenderer.Render(&buf, page); err != nil {
				log.Fatal(err)
			}

			// write buffer to file
			path := outputDir + "/" + page.SanitisedTitle() + ".html"
			err = os.WriteFile(path, buf.Bytes(), 0660)
			if err != nil {
				log.Fatal(err)
			}
		}

	case "serve":
		serveCmd.Parse(os.Args[2:]) //nolint:errcheck // this error is handled by the Flags' parser fail behaviour for now
		fmt.Println("subcommand 'serve'")
		fmt.Println(" site content dir:", serveDir)
		fmt.Println(" port:", port)
		url, err := parseLocalURL(port)
		if err != nil {
			log.Fatal(err)
		}

		addr := url.Host
		fs := http.FS(os.DirFS(serveDir))

		log.Printf("Starting server on: %s", addr)

		log.Fatal(http.ListenAndServe(addr, http.FileServer(fs)))

	default:
		fmt.Println("expected 'generate' or 'serve' subcommands")
		os.Exit(1)
	}
}

func parseLocalURL(port string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%s", port))
	if err != nil {
		fmt.Println(err)
	}
	return u, nil
}
