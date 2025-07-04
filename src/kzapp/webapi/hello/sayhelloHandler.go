package hello

import (
	"fmt"
	"html/template" // For HTML templates (optional but good practice)
	"kzapp/webapi/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath" // For safe path joining

	"github.com/gorilla/mux"

	markdown "github.com/gomarkdown/markdown" // Markdown parser
	"github.com/gomarkdown/markdown/html"     // HTML renderer for markdown
	"github.com/gomarkdown/markdown/parser"   // Parser options for markdown
)

type GreetingHandler struct{}

func (h GreetingHandler) InitService(router *mux.Router) {
	router.HandleFunc("/", pkg.Chain(h.sayhello, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/greet", pkg.Chain(h.greet, pkg.Method("GET"), pkg.Logging()))
}

func (h GreetingHandler) greet(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	fmt.Printf("userName: %s\n", userName)
	greeting := "Welcome Back!"

	response := GreetingResponse{
		Name:    userName,
		Message: &greeting,
	}

	pkg.JsonResponse(w, response)
}

func (h GreetingHandler) sayhello(w http.ResponseWriter, r *http.Request) {

	mdFilepath := "./../webapi/hello/README.md"

	// Read the Markdown file content
	markdownBytes, err := os.ReadFile(mdFilepath)
	if err != nil {
		log.Printf("Error reading markdown file %s: %v", mdFilepath, err)
		http.Error(w, "Could not read markdown file", http.StatusInternalServerError)
		return
	}

	// Configure Markdown parser options
	// Extensions like Autolink, FencedCode, and Strikethrough are commonly useful
	extensions := parser.CommonExtensions | parser.NoIntraEmphasis
	p := parser.NewWithExtensions(extensions)

	// Configure HTML renderer options
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert Markdown to HTML
	htmlBytes := markdown.ToHTML(markdownBytes, p, renderer)

	// Prepare data for the HTML template
	data := TemplateData{
		Title:   filepath.Base(mdFilepath), // Use filename as title
		Content: template.HTML(htmlBytes),  // The converted HTML
	}

	// Parse and execute the HTML template
	tmpl, err := template.New("markdownPage").Parse(htmlTemplate)
	if err != nil {
		log.Printf("Error parsing HTML template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing HTML template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

// A simple HTML template to wrap our Markdown-converted HTML
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: sans-serif; margin: 2em; line-height: 1.6; }
        pre { background-color: #eee; padding: 1em; overflow-x: auto; }
        code { font-family: monospace; }
        h1, h2, h3 { color: #eee; }
        a { color:rgb(192, 169, 39); text-decoration: none; }
        a:hover { color: #FFD700; text-decoration: underline; }
    </style>
</head>
<body style="background-color: #000;">	
    {{.Content}}
</body>
</html>
`

// TemplateData struct to pass data to the HTML template
type TemplateData struct {
	Title   string
	Content template.HTML // Use template.HTML to prevent Go's HTML escaping
}
