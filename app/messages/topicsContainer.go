package messages
import "log"


type ITopicsContainer interface {
	AddMessage(messageData *MessageInput)
	Subscribe(topic string) <-chan MessageOutput
	UnSubscribe(topic string, removeCh <-chan MessageOutput)
}

func NewTopicsContainer() ITopicsContainer {
	return &topicsContainer{
		topicManagers: make(map[string]*topicManager),
		topicExitAnnouncementCh: make(chan unSubscribeResult)}
}

type topicsContainer struct {
	topicManagers           map[string]*topicManager
	topicExitAnnouncementCh chan unSubscribeResult
}

func (this *topicsContainer) findOrCreateTopicManager(topic string) *topicManager {
	topicManager, exists := this.topicManagers[topic]
	if !exists {
		topicManager = newTopicManager(topic, this.topicExitAnnouncementCh)
		this.topicManagers[topic] = topicManager

		log.Printf("topicsContainer#findOrCreateTopicManager: Starting new TopicManager '%v'.\n", topicManager)
		go topicManager.startRunning()
	}

	return topicManager
}

func (this *topicsContainer) AddMessage(messageData *MessageInput)  {
	topicManager, exists := this.topicManagers[messageData.Topic]
	if !exists {
		log.Printf(
			"topicsContainer#AddMessage: TopicManager '%s' is not running as there are no listeners, message '%s' will be ignored.\n",
			messageData.Topic, messageData.Message);
		return
	}

	log.Printf("topicsContainer#AddMessage: Adding message to topic - '%s'.\n", messageData.Topic)
	topicManager.getHandle().addMessageCh <- messageData.Message
}


func (this *topicsContainer) Subscribe(topic string) <- chan MessageOutput {
	topicManager := this.findOrCreateTopicManager(topic)
	log.Printf("topicsContainer#Subscribe: TopicManager - '%v'.\n", topicManager)

	listenToMessagesCh := make(chan MessageOutput)
	topicManager.getHandle().addSubscriberCh <- listenToMessagesCh
	return listenToMessagesCh
}

func (this *topicsContainer) UnSubscribe(topic string, removeCh <-chan MessageOutput) {
	topicManager, exists := this.topicManagers[topic]
	if exists {
		topicManager.getHandle().unSubscribeCh <- removeCh

		answer := <- this.topicExitAnnouncementCh

		if(answer.stopped) {
			delete(this.topicManagers, topic)
		}
	}
}

type unSubscribeResult struct {
	topic string
	stopped bool
}


