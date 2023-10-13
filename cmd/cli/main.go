package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "github.com/wsk9531/henshall.dev/internal/generator"
)

const (
	CONTENT_DIR string = "pages"
	OUTPUT_DIR string = "dist"
	DEFAULT_PORT string = "5000"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	var contentDir string
	generateCmd.StringVar(&contentDir, "input", CONTENT_DIR, "path to content for rendering")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	var port string
	serveCmd.StringVar(&port, "port", DEFAULT_PORT, "localhost port to serve content")

	if len(os.Args) < 2 {
		fmt.Println("expected 'generate' or 'serve' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "generate":
		generateCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'generate'")
		fmt.Println("  dir:", contentDir)
		fmt.Println("  tail:", generateCmd.Args())

		// TODO: Add generate cli subcommand
		// e.g.
		// posts, _ := generator.NewPostsFromFS(os.DirFS(contentDir))
		// log.Println(posts)
		//postRenderer, _ := generator.NewPostRenderer()
		// for _, page := range(posts) {
		// 	postRenderer.Render(page)
		// }


	case "serve":
		serveCmd.Parse(os.Args[2:])

		url, err := parseLocalURL(port)
		if err != nil {
			log.Fatal(err)
		}
		addr := url.Host

		log.Printf("Starting server on: %s", addr)
		log.Fatal(http.ListenAndServe(addr, http.FileServer(http.FS(os.DirFS(OUTPUT_DIR)))))

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
