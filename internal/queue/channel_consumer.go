package queue

import (
	"context"
	"deeliai/internal/interfaces"
	"deeliai/internal/service"
	"log"
	"time"
)

// channelConsumer 是 QueueConsumer 介面基於 Go Channel 的實現
type channelConsumer struct {
	queue         chan string
	scrapeService *service.ScrapeService
	workerCount   int
}

func NewChannelConsumer(q chan string, s *service.ScrapeService, cnt int) interfaces.QueueConsumer {
	return &channelConsumer{
		queue:         q,
		scrapeService: s,
		workerCount:   cnt,
	}
}

// Start 啟動 worker pool
func (cc *channelConsumer) Start() {
	log.Printf("Starting %d scrape workers...", cc.workerCount)

	// 啟動多個 worker goroutine
	for i := 0; i < cc.workerCount; i++ {
		log.Printf("Worker #%d started...", i)
		// 呼叫 cc.Consume 執行 cc.scrapeService.ProcessScrapeTask
		go func(id int) {
			cc.Consume()
			log.Printf("Worker #%d stopped.", id) // 真的結束時才印
		}(i)
	}
}

// Consume 是 channelConsumer 的執行邏輯
func (cc *channelConsumer) Consume() {
	// 確保每個單獨的爬取任務都有自己的超時控制
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 每個 callback 都是一個無窮迴圈，持續從 channel 中讀取任務
	for articleID := range cc.queue {
		cc.scrapeService.ProcessScrapeTask(ctx, articleID)
	}
}
