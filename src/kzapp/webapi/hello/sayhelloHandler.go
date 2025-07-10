package hello

import (
	"fmt"
	"html/template" // 用於 HTML 樣板（可選但建議使用）
	"kzapp/webapi/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath" // For safe path joining
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"

	markdown "github.com/gomarkdown/markdown" // Markdown parser
	"github.com/gomarkdown/markdown/html"     // HTML renderer for markdown
	"github.com/gomarkdown/markdown/parser"   // Parser options for markdown
)

type GreetingHandler struct{}

// for compile time check
var _ pkg.Handler = (*GreetingHandler)(nil)

func (h GreetingHandler) InitServiceGin(router *gin.Engine) {
	router.GET("/", pkg.ChainGin(h.sayhello), gin.Logger())
	router.GET("/greet", pkg.ChainGin(h.greet), gin.Logger())
	router.GET("/img/{image}", pkg.ChainGin(h.giveMeCorgi), gin.Logger())
}

func (h GreetingHandler) InitService(router *mux.Router) {
	router.HandleFunc("/", pkg.Chain(h.sayhello, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/greet", pkg.Chain(h.greet, pkg.Method("GET"), pkg.Logging()))
	router.HandleFunc("/img/{image}", pkg.Chain(h.giveMeCorgi, pkg.Method("GET"), pkg.Logging()))

	// todo static file server
	// router.PathPrefix("/imgs/").Handler(h.giveMeImages())
	// router.HandleFunc("/imgs/", pkg.Chain(h.giveMeImages().ServeHTTP, pkg.Method("GET"), pkg.Logging()))
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

// todo static file server
func (h GreetingHandler) giveMeImages() http.Handler {
	// 設定靜態檔案目錄
	imageDir := "webapi/images"

	// 使用 FileServer 提供靜態檔案服務
	return http.StripPrefix("/imgs/", http.FileServer(http.Dir(imageDir)))
}

// example: http://localhost:80/img/corgi_laugh.jpg
func (h GreetingHandler) giveMeCorgi(w http.ResponseWriter, r *http.Request) {
	// 使用 FileServer 來提供靜態檔案服務
	// 適用於：
	// 1. 提供靜態資源（圖片、CSS、JS、HTML 檔案等）
	// 2. 簡單的檔案瀏覽功能
	// 3. 開發環境中的檔案伺服器
	// 4. 不需要複雜路由邏輯的檔案服務

	// 設定圖片目錄路徑
	imageDir := "images/"
	supportedFormats := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}

	// 從 URL 路徑中提取圖片檔名
	vars := mux.Vars(r)
	imageName := vars["image"]

	// 驗證檔名安全性（防止路徑遍歷攻擊）
	if imageName == "" || containsDotDot(imageName) {
		http.Error(w, "Invalid image name", http.StatusBadRequest)
		return
	}

	// 構建完整的檔案路徑
	imagePath := filepath.Join(imageDir, imageName)

	// 檢查檔案是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// 檢查檔案副檔名是否支援
	ext := filepath.Ext(imageName)
	contentType, supported := supportedFormats[ext]
	if !supported {
		http.Error(w, "Unsupported image format", http.StatusBadRequest)
		return
	}

	// 設定適當的 HTTP 標頭
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 快取一年
	w.Header().Set("Access-Control-Allow-Origin", "*")          // 允許跨域存取

	// 提供圖片檔案
	http.ServeFile(w, r, imagePath)

	log.Printf("Served image: %s", imageName)
}

// containsDotDot 檢查路徑是否包含 ".." 以防止路徑遍歷攻擊
func containsDotDot(path string) bool {
	return strings.Contains(path, "..")
}

func (h GreetingHandler) sayhello(w http.ResponseWriter, r *http.Request) {
	fileName := "README.md"
	currentDir, err := pkg.GetResourceDir()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	mdFilepath := filepath.Join(currentDir, fileName)

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
