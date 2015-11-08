package messages


type ITopicsContainer interface {
	AddMessage(messageData *MessageInput)
	Subscribe(topic string) chan MessageOutput
	UnSubscribe(topic string, removeCh chan MessageOutput)
}

func NewTopicsContainer() ITopicsContainer {
	return &topicsContainer{make(map[string]*topicManager)}
}

type topicsContainer struct {
	topicManagers map[string]*topicManager
}

func (this *topicsContainer) getTopicManager(topic string) *topicManager {
	manager, exists := this.topicManagers[topic]
	if !exists {
		manager = newTopicManager(topic)
		this.topicManagers[topic] = manager
	}

	return manager;
}


func (this *topicsContainer) AddMessage(messageData *MessageInput)  {
	topicManager, exists := this.topicManagers[messageData.Topic]
	if !exists {
		topicManager = newTopicManager(messageData.Topic)
		this.topicManagers[messageData.Topic] = topicManager
	}

	topicManager.addMessage(messageData.Message)
}

func (this *topicsContainer) Subscribe(topic string) chan MessageOutput {
	topicManager := this.getTopicManager(topic)
	return topicManager.subscribe()
}

func (this *topicsContainer) UnSubscribe(topic string, removeCh chan MessageOutput) {
	topicManager := this.getTopicManager(topic)
	topicManager.unSubscribe(removeCh)
}

