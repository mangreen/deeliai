package service

import (
	"context"
	"deeliai/internal/interfaces"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

type ScrapeService struct {
	articleRepo interfaces.ArticleRepository
}

// NewScrapeService 接受 worker 數量作為參數
func NewScrapeService(repo interfaces.ArticleRepository) *ScrapeService {
	return &ScrapeService{
		articleRepo: repo,
	}
}

// ProcessScrapeTask 處理單個爬取任務
func (w *ScrapeService) ProcessScrapeTask(ctx context.Context, articleID string) {
	id, err := uuid.Parse(articleID)
	if err != nil {
		log.Printf("Invalid article ID in queue: %s", articleID)
		return
	}

	// 這裡需要從資料庫取回 URL，由於我們直接傳 ID，所以需要先查詢一次
	article, err := w.articleRepo.FindByID(ctx, id)
	if err != nil {
		log.Printf("Article not found for ID: %s", articleID)
		return
	}

	title, desc, img, err := w.scrapeMetadata(article.URL)
	if err != nil {
		log.Printf("Failed to scrape URL %s: %v", article.URL, err)
		// 爬取失敗，標記為失敗並增加重試次數
		if err := w.articleRepo.MarkScrapeFailed(ctx, id); err != nil {
			log.Printf("Failed to mark scrape as failed: %v", err)
		}
		return
	}

	// 爬取成功，更新資料庫
	if err := w.articleRepo.UpdateMetadata(ctx, id, title, desc, img); err != nil {
		log.Printf("Failed to update article metadata: %v", err)
	} else {
		log.Printf("Successfully scraped and updated article ID: %s", articleID)
	}
}

// scrapeMetadata 實際的爬取邏輯，使用 goquery
func (w *ScrapeService) scrapeMetadata(url string) (title, description, imageURL string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	// 優先抓取 OpenGraph Metadata
	title = doc.Find("meta[property='og:title']").AttrOr("content", "")
	description = doc.Find("meta[property='og:description']").AttrOr("content", "")
	imageURL = doc.Find("meta[property='og:image']").AttrOr("content", "")

	// 若 OpenGraph 找不到，退回抓取一般 HTML 標籤
	if title == "" {
		title = doc.Find("title").Text()
	}
	if description == "" {
		description = doc.Find("meta[name='description']").AttrOr("content", "")
	}
	// Image 暫不退回，因為一般 img 標籤可能不適合作為預覽圖

	return title, description, imageURL, nil
}
