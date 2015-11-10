package messages
import (
	"log"
)


func newTopicManager(topic string, informAboutExitCh chan<- unSubscribeResult) *topicManager {
	addMessageCh := make(chan string)
	addSubscriberCh := make(chan chan MessageOutput)
	unSubscribeCh := make(chan (<-chan MessageOutput))
	handle := topicManagerHandle{addMessageCh: addMessageCh, addSubscriberCh: addSubscriberCh, unSubscribeCh: unSubscribeCh}

	return &topicManager{
		topic: topic,
		listeners: make([] chan MessageOutput, 0),
		messageCounter: 0,
		informAboutExitCh: informAboutExitCh,
		addMessageCh: addMessageCh, addSubscriberCh: addSubscriberCh, unSubscribeCh: unSubscribeCh,
		handle: handle}
}

type topicManager struct {
	topic           string
	listeners       [] chan MessageOutput
	messageCounter  int

	informAboutExitCh chan<- unSubscribeResult
	addMessageCh    <-chan string
	addSubscriberCh <-chan (chan MessageOutput)
	unSubscribeCh   <-chan (<-chan MessageOutput)

	handle topicManagerHandle
}

func (this * topicManager) getHandle() topicManagerHandle {
	return this.handle
}

func (this *topicManager) startRunning() {
	for {
		select {
		case message := <-this.addMessageCh:
			log.Printf("topicManager#startRunning: Received message - '%s'\n", message)
			this.addMessage(message)
		case listenerCh := <-this.addSubscriberCh:
			log.Printf("topicManager#startRunning: New Subscriber '%v'.\n", listenerCh)
			this.subscribe(listenerCh)
		case listenerCh := <-this.unSubscribeCh:
			log.Printf("topicManager#startRunning: Unsubscribing '%v'.\n", listenerCh)
			this.unSubscribe(listenerCh)

			// Stop running if no listeners
			if len(this.listeners) == 0 {
				log.Printf("topicManager#startRunning: No Listeners - Exit.\n")
				this.informAboutExitCh <- unSubscribeResult{topic: this.topic, stopped: true}
				return
			} else {
				this.informAboutExitCh <- unSubscribeResult{topic: this.topic, stopped: false}
			}
		}
	}

	panic("topicManager#startRunning: Should never reach this!.\n")
}

func (this *topicManager) subscribe(listenToMessagesCh chan MessageOutput) {
	this.listeners = append(this.listeners, listenToMessagesCh)
}

func (this *topicManager) unSubscribe(listenToMessagesCh <-chan MessageOutput) {
	var indexToRemove int = -1
	for index, ch := range this.listeners {
		if listenToMessagesCh == ch {
			indexToRemove = index
		}
	}

	if indexToRemove >= 0 {
		this.listeners = append(this.listeners[:indexToRemove], this.listeners[indexToRemove + 1:]...)
	} else {
		log.Println("topicManager#UnSubscribe: Cannot find channel to remove.")
	}
}

func (this *topicManager) addMessage(message string) {

	// TODO This is not "Goroutine-Safe" I guess
	this.messageCounter++

	messageOutput := MessageOutput{this.messageCounter, message}

	if len(this.listeners) <= 0 {
		log.Println("topicManager#addMessage: No Listeneres to send to.");
	}

	for _, listenerCh := range this.listeners {
		if listenerCh != nil {
			listenerCh <- messageOutput
		}
	}
}



