package queue

// channelConsumer 是 QueueConsumer 介面基於 Go Channel 的實現
type channelConsumer struct {
	queue chan string
}

func NewChannelConsumer(q chan string) QueueConsumer {
	return &channelConsumer{
		queue: q,
	}
}

// Consume 是 channelConsumer 的執行邏輯
func (cc *channelConsumer) Consume(callback func(string)) {
	// 每個 callback 都是一個無窮迴圈，持續從 channel 中讀取任務
	for articleID := range cc.queue {
		callback(articleID)
	}
}
