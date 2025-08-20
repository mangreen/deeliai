package queue

import "errors"

// ChannelProducer 是 QueueProducer 介面基於 Go Channel 的實現
type ChannelProducer struct {
	queue chan string
}

func NewChannelProducer(queue chan string) *ChannelProducer {
	return &ChannelProducer{queue: queue}
}

// Produce 將任務字串送入 Go Channel
func (p *ChannelProducer) Produce(message string) error {
	select {
	case p.queue <- message:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// Close 關閉 Go Channel
func (p *ChannelProducer) Close() error {
	close(p.queue)
	return nil
}
