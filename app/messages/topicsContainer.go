package messages


type ITopicsContainer interface {
	GetTopicManager(topic string) *topicManager
	AddMessage(messageData *MessageInput)
}


func NewTopicsContainer() ITopicsContainer {
	return &topicsContainer{make(map[string]*topicManager)}
}

type topicsContainer struct {
	topicManagers map[string]*topicManager
}

func (this *topicsContainer) GetTopicManager(topic string) *topicManager {
	manager, _ := this.topicManagers[topic]
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

