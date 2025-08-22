package queue

import "errors"

// channelProducer 是 QueueProducer 介面基於 Go Channel 的實現
type channelProducer struct {
	queue chan string
}

func NewChannelProducer(queue chan string) QueueProducer {
	return &channelProducer{queue: queue}
}

// Produce 將任務字串送入 Go Channel
func (p *channelProducer) Produce(message string) error {
	select {
	case p.queue <- message:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// Close 關閉 Go Channel
func (p *channelProducer) Close() error {
	close(p.queue)
	return nil
}
