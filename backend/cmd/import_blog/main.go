package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotel-backend/internal/config"
	"hotel-backend/internal/database"
)

// StaticArticle represents the model in static JSON file
type StaticArticle struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Image    string `json:"image"`
	Content  string `json:"content"`
	AdLink   string `json:"adLink"`
}

func main() {
	log.Println("🔄 Starting blog articles Markdown to HTML DB import...")

	// 1. Load config and connect
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema to ensure posts table is created
	if err := db.InitSchema(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize schema: %v", err)
	}
	log.Println("✅ Database connected and schema verified")

	// 2. Locate articles.json
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "frontend"
	}
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "../../frontend"
	}

	articlesPath := filepath.Join(frontendDir, "data_json/blog/articles.json")
	log.Printf("📂 Reading articles from: %s", articlesPath)

	data, err := ioutil.ReadFile(articlesPath)
	if err != nil {
		log.Fatalf("❌ Failed to read articles.json: %v", err)
	}

	var staticArticles []StaticArticle
	if err := json.Unmarshal(data, &staticArticles); err != nil {
		log.Fatalf("❌ Failed to parse articles JSON: %v", err)
	}

	// 3. Clear existing posts to prevent duplicate seeding (optional but clean for seeding)
	_, err = db.Pool.Exec(ctx, `TRUNCATE TABLE posts RESTART IDENTITY`)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to truncate posts table: %v", err)
	}

	// 4. Convert and insert articles
	query := `
		INSERT INTO posts (
			id, title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO NOTHING`

	for _, sa := range staticArticles {
		// Convert Markdown content/adLink to HTML
		htmlContent := markdownToHTML(sa.Content)
		htmlAdLink := markdownToHTML(sa.AdLink)

		// Set default SEO TKD fields based on title and content excerpt
		seoTitle := sa.Title
		seoKeywords := strings.Join([]string{sa.Category, "飯店推薦", "汽車旅館休息"}, ",")
		
		excerpt := sa.Title
		if len(sa.Content) > 150 {
			excerpt = sa.Content[:150]
		}
		seoDescription := strings.ReplaceAll(excerpt, "\n", " ")

		// Tags JSON: Category is added as first tag
		tags := []string{sa.Category}
		tagsJSON, _ := json.Marshal(tags)

		// Parse date (format: 2026-01-05)
		parsedDate, errDate := time.Parse("2006-01-02", sa.Date)
		if errDate != nil {
			parsedDate = time.Now()
		}

		idNum := 1
		fmt.Sscanf(sa.ID, "%d", &idNum)

		_, err = db.Pool.Exec(ctx, query,
			idNum, sa.Title, string(tagsJSON), sa.Image, htmlContent, htmlAdLink,
			seoTitle, seoKeywords, seoDescription, parsedDate, parsedDate,
		)

		if err != nil {
			log.Printf("❌ Failed to import article %d: %v", idNum, err)
		} else {
			log.Printf("✅ Imported post %d: %s", idNum, sa.Title)
		}
	}

	log.Println("🎉 Seeding completed successfully!")
}

// Markdown to HTML simple parser
func markdownToHTML(md string) string {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	blocks := strings.Split(md, "\n\n")
	var htmlBlocks []string
	
	inList := false
	var listItems []string

	flushList := func() {
		if len(listItems) > 0 {
			htmlBlocks = append(htmlBlocks, "<ul>\n"+strings.Join(listItems, "\n")+"\n</ul>")
			listItems = nil
			inList = false
		}
	}

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		if strings.HasPrefix(block, "#### ") {
			flushList()
			title := strings.TrimPrefix(block, "#### ")
			htmlBlocks = append(htmlBlocks, "<h4>"+parseInlines(title)+"</h4>")
		} else if strings.HasPrefix(block, "### ") {
			flushList()
			title := strings.TrimPrefix(block, "### ")
			htmlBlocks = append(htmlBlocks, "<h3>"+parseInlines(title)+"</h3>")
		} else if strings.HasPrefix(block, "- ") || strings.HasPrefix(block, "* ") {
			inList = true
			lines := strings.Split(block, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "- ") {
					item := strings.TrimPrefix(line, "- ")
					listItems = append(listItems, "  <li>"+parseInlines(item)+"</li>")
				} else if strings.HasPrefix(line, "* ") {
					item := strings.TrimPrefix(line, "* ")
					listItems = append(listItems, "  <li>"+parseInlines(item)+"</li>")
				} else if inList && line != "" {
					if len(listItems) > 0 {
						lastIdx := len(listItems) - 1
						// Strip trailing </li> and append text
						listItems[lastIdx] = listItems[lastIdx][:len(listItems[lastIdx])-5] + " " + parseInlines(line) + "</li>"
					}
				}
			}
		} else {
			flushList()
			lines := strings.Split(block, "\n")
			var processedLines []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					processedLines = append(processedLines, parseInlines(line))
				}
			}
			htmlBlocks = append(htmlBlocks, "<p>"+strings.Join(processedLines, "<br/>")+"</p>")
		}
	}
	flushList()

	return strings.Join(htmlBlocks, "\n\n")
}

func parseInlines(text string) string {
	// 1. Parse images: ![alt](url) -> <img src="url" alt="alt" />
	for {
		imgStart := strings.Index(text, "![")
		if imgStart == -1 {
			break
		}
		imgMid := strings.Index(text[imgStart:], "](")
		if imgMid == -1 {
			break
		}
		imgEnd := strings.Index(text[imgStart+imgMid:], ")")
		if imgEnd == -1 {
			break
		}
		
		alt := text[imgStart+2 : imgStart+imgMid]
		url := text[imgStart+imgMid+2 : imgStart+imgMid+imgEnd]
		
		tag := fmt.Sprintf(`<img src="%s" alt="%s" />`, url, alt)
		text = text[:imgStart] + tag + text[imgStart+imgMid+imgEnd+1:]
	}

	// 2. Parse links: [text](url) -> <a href="url" target="_blank">text</a>
	for {
		lnkStart := strings.Index(text, "[")
		if lnkStart == -1 {
			break
		}
		lnkMid := strings.Index(text[lnkStart:], "](")
		if lnkMid == -1 {
			break
		}
		lnkEnd := strings.Index(text[lnkStart+lnkMid:], ")")
		if lnkEnd == -1 {
			break
		}
		
		val := text[lnkStart+1 : lnkStart+lnkMid]
		url := text[lnkStart+lnkMid+2 : lnkStart+lnkMid+lnkEnd]
		
		tag := fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, url, val)
		text = text[:lnkStart] + tag + text[lnkStart+lnkMid+lnkEnd+1:]
	}

	return text
}
