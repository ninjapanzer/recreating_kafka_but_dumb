package internal

type Topic struct {
	Name      string
	Compacted bool
}

type Consumer struct {
	Position  uint64
	TopicName string
}

type Producer struct {
	TopicName string
}
