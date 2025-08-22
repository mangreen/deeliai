package scraper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"deeliai/internal/queue"
	"deeliai/internal/repository"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

type ScrapeWorker struct {
	articleRepo repository.ArticleRepository
	consumer    queue.QueueConsumer
	workerCount int // worker 數量
}

// NewScrapeWorker 接受 worker 數量作為參數
func NewScrapeWorker(repo repository.ArticleRepository, consumer queue.QueueConsumer, count int) *ScrapeWorker {
	return &ScrapeWorker{
		articleRepo: repo,
		consumer:    consumer,
		workerCount: count,
	}
}

// Start 啟動 worker pool
func (w *ScrapeWorker) Start() {
	log.Printf("Starting %d scrape workers...", w.workerCount)

	// 啟動多個 worker goroutine
	for i := 0; i < w.workerCount; i++ {
		log.Printf("Worker #%d started...", i)
		// 呼叫 QueueConsumer 執行 processScrapeTask
		w.consumer.Consume(w.processScrapeTask)
		log.Printf("Worker #%d stopped.", i)
	}
}

// processScrapeTask 處理單個爬取任務
func (w *ScrapeWorker) processScrapeTask(articleID string) {
	// 確保每個單獨的爬取任務都有自己的超時控制
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
func (w *ScrapeWorker) scrapeMetadata(url string) (title, description, imageURL string, err error) {
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
