package messages
import "fmt"


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
	fmt.Printf("topicManager#AddMessage(): %s\n" + message)
}



