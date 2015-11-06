package messages
import (
	"log"
)


func newTopicManager(topic string) *topicManager {
	return &topicManager{topic, make([] chan *MessageOutput, 10)}
}

type topicManager struct {
	topic string
	listeners [] chan *MessageOutput
}


func (this *topicManager) Subscribe(topic string) chan *MessageOutput {
	topicCh := make(chan *MessageOutput)

	this.listeners = append(this.listeners, topicCh)

	return topicCh
}

func (this *topicManager) addMessage(message string)  {
	// TODO Send to all subscribers
	log.Printf("topicManager#addMessage(): %s\n", message)
}



