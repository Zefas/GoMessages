package messages
import "log"


type ITopicsContainer interface {
	AddMessage(messageData *MessageInput)
	Subscribe(topic string) <-chan MessageOutput
	UnSubscribe(topic string, removeCh <-chan MessageOutput)
}

func NewTopicsContainer() ITopicsContainer {
	return &topicsContainer{make(map[string]*topicManager)}
}

type topicsContainer struct {
	topicManagers map[string]*topicManager
}

func (this *topicsContainer) findTopicManager(topic string) *topicManager {
	manager, exists := this.topicManagers[topic]
	if !exists {
		manager = newTopicManager(topic)
		this.topicManagers[topic] = manager
	}

	return manager;
}

func (this *topicsContainer) findOrCreateTopicManager(topic string) *topicManager {
	topicManager, exists := this.topicManagers[topic]
	if !exists {
		topicManager = newTopicManager(topic)
		this.topicManagers[topic] = topicManager

		log.Printf("topicsContainer#findOrCreateTopicManager: Starting new TopicManager '%v'.\n", topicManager)
		go topicManager.startRunning()
	}

	return topicManager
}

func (this *topicsContainer) AddMessage(messageData *MessageInput)  {
	topicManager := this.findOrCreateTopicManager(messageData.Topic)
	log.Printf("topicsContainer#AddMessage: TopicManager - '%v'.\n", topicManager)

	log.Printf("topicsContainer#AddMessage: Adding message via channel - '%v'.\n", topicManager.getHandle().addMessageCh)
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
	topicManager := this.findTopicManager(topic)
	topicManager.unSubscribe(removeCh)
}

