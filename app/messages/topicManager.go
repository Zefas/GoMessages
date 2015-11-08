package messages
import (
	"log"
)


func newTopicManager(topic string) *topicManager {
	addMessageCh := make(chan string)
	addSubscriberCh := make(chan chan MessageOutput)
	unSubscribeCh := make(chan chan MessageOutput)
	handle := topicManagerHandle{addMessageCh: addMessageCh, addSubscriberCh: addSubscriberCh, unSubscribeCh: unSubscribeCh}

	return &topicManager{
		topic: topic,
		listeners: make([] chan MessageOutput, 0),
		messageCounter: 0,
		addMessageCh: addMessageCh, addSubscriberCh: addSubscriberCh, unSubscribeCh: unSubscribeCh,
		handle: handle}
}

type topicManager struct {
	topic           string
	listeners       [] chan MessageOutput
	messageCounter  int

	addMessageCh    <-chan string
	addSubscriberCh <-chan chan MessageOutput
	unSubscribeCh   <-chan chan MessageOutput

	handle topicManagerHandle
}

func (this * topicManager) getHandle() topicManagerHandle {
	return this.handle
}

func (this *topicManager) startRunning() {
	select {
	case message := <-this.addMessageCh:
		this.addMessage(message)
	case newSubscriberCh := <-this.addSubscriberCh:
		this.subscribe(newSubscriberCh)
	case subscriberCh := <-this.unSubscribeCh:
		this.unSubscribe(subscriberCh)
	}
}

func (this *topicManager) subscribe(newSubscriberCh chan MessageOutput) {
	this.listeners = append(this.listeners, newSubscriberCh)
}

/*func (this *topicManager) subscribe() chan MessageOutput {
	topicCh := make(chan MessageOutput)
	this.listeners = append(this.listeners, topicCh)

	log.Printf("topicManager#Subscribe: adding chanel - '%v'\n", topicCh)
	return topicCh
}*/

func (this *topicManager) unSubscribe(removeCh <-chan MessageOutput) {
	var indexToRemove int = -1
	for index, ch := range this.listeners {
		if removeCh == ch {
			indexToRemove = index
		}
	}

	if indexToRemove >= 0 {
		log.Printf("topicManager#UnSubscribe: Removing channel - '%v'.\n", removeCh)
		this.listeners = append(this.listeners[:indexToRemove], this.listeners[indexToRemove + 1:]...)
	} else {
		log.Println("topicManager#UnSubscribe: Cannot find channel to remove.")
	}
}

func (this *topicManager) addMessage(message string) {
	log.Printf("topicManager#addMessage(): %s\n", message)

	// TODO This is not "Goroutine-Safe" I guess
	this.messageCounter++

	messageOutput := MessageOutput{this.messageCounter, message}

	if len(this.listeners) <= 0 {
		log.Println("topicManager#addMessage: No Listeneres to send to.");
	}

	for _, listenerCh := range this.listeners {
		if listenerCh != nil {
			log.Printf("topicManager#addMessage: sending to chanel - %v\n", listenerCh)
			listenerCh <- messageOutput
		}
	}
}



