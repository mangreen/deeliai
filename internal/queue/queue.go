package queue

// QueueProducer 是將任務發布到佇列的抽象介面
type QueueProducer interface {
	// Produce 方法接受一個任務字串，並將其發布
	Produce(message string) error
	// 也可以加入 Close 等方法來清理資源
	Close() error
}

// QueueConsumer 是處理任務的抽象介面
type QueueConsumer interface {
	// Produce 方法接受一個任務字串，並將其發布
	Consume(callback func(string))
}
