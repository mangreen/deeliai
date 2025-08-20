package scraper

import (
	"context"
	"log"
	"time"

	"deeliai/internal/queue"
	"deeliai/internal/repository"
)

// ScrapeScheduler 定時檢查失敗的爬取任務並重新排入佇列
type ScrapeScheduler struct {
	articleRepo repository.ArticleRepository
	producer    queue.QueueProducer
}

func NewScrapeScheduler(repo repository.ArticleRepository, producer queue.QueueProducer) *ScrapeScheduler {
	return &ScrapeScheduler{
		articleRepo: repo,
		producer:    producer,
	}
}

// Start 啟動排程器，每隔 5 分鐘檢查一次
func (s *ScrapeScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Scrape Scheduler started...")

	// 立即執行一次檢查
	s.checkAndRequeue(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Scrape Scheduler shutting down...")
			return
		case <-ticker.C:
			s.checkAndRequeue(ctx)
		}
	}
}

func (s *ScrapeScheduler) checkAndRequeue(ctx context.Context) {
	log.Println("Checking for failed scrape tasks...")
	articles, err := s.articleRepo.FindFailedScrapes(ctx)
	if err != nil {
		log.Printf("Error finding failed scrapes: %v", err)
		return
	}

	for _, article := range articles {
		log.Printf("Re-queuing failed scrape task for article ID: %s", article.ID.String())
		if err := s.producer.Produce(article.ID.String()); err != nil {
			log.Printf("Failed to requeue article %s: %v", article.ID.String(), err)
		}
	}
}
