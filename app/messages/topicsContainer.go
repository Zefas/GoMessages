package messages
import "log"


type ITopicsContainer interface {
	AddMessage(messageData *MessageInput)
	Subscribe(topic string) <-chan MessageOutput
	UnSubscribe(topic string, removeCh <-chan MessageOutput)
}

func NewTopicsContainer() ITopicsContainer {
	return &topicsContainer{
		topicManagers: *NewConcurrentTopicMap(),
		topicExitAnnouncementCh: make(chan unSubscribeResult)}
}

type topicsContainer struct {
	topicManagers           concurrentTopicMap
	topicExitAnnouncementCh chan unSubscribeResult
	topicManagersStarted    int
}

func (this *topicsContainer) findOrCreateTopicManager(topic string) *topicManager {
	topicManager := this.topicManagers.GetOrAdd(topic, func()*topicManager {
		topicManager := newTopicManager(topic, this.topicExitAnnouncementCh)

		log.Printf("topicsContainer#findOrCreateTopicManager: Starting new TopicManager for topic '%s'.\n", topic)
		go topicManager.startRunning()
		this.topicManagersStarted += 1

		return topicManager
	})

	return topicManager
}

func (this *topicsContainer) AddMessage(messageData *MessageInput) {
	this.topicManagers.IfPresentOrNot(
		messageData.Topic,
		func(topicManager *topicManager) {
			topicManager.getHandle().addMessageCh <- messageData.Message
		},
		func() {
			log.Printf(
				"topicsContainer#AddMessage: TopicManager '%s' is not running as there are no listeners, message '%s' will be ignored.\n",
				messageData.Topic, messageData.Message);
		});
}


func (this *topicsContainer) Subscribe(topic string) <- chan MessageOutput {
	topicManager := this.findOrCreateTopicManager(topic)

	listenToMessagesCh := make(chan MessageOutput)
	topicManager.getHandle().addSubscriberCh <- listenToMessagesCh
	return listenToMessagesCh
}

func (this *topicsContainer) UnSubscribe(topic string, removeCh <-chan MessageOutput) {
	this.topicManagers.IfPresent(topic, func(topicManager *topicManager) {
		topicManager.getHandle().unSubscribeCh <- removeCh

		answer := <- this.topicExitAnnouncementCh

		if(answer.stopped) {
			this.topicManagers.Delete(topic)
		}
	})
}

type unSubscribeResult struct {
	topic string
	stopped bool
}


