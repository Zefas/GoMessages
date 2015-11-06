package messages



type MessageInput struct {
	Topic string
	Message string
}

type MessageOutput struct {
	Id int
	Data string
}
