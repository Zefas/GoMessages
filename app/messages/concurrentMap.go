package messages
import (
	"sync"
)



func NewConcurrentTopicMap() *concurrentTopicMap {
	return &concurrentTopicMap{m: make(map[string]*topicManager)}
}


type concurrentTopicMap struct {
	sync.RWMutex
	m map[string]*topicManager
}

func (this *concurrentTopicMap) Get(key string) (*topicManager, bool) {
	this.RLock()
	topicManager, exists := this.m[key]
	this.RUnlock()

	return topicManager, exists
}

func (this *concurrentTopicMap) IfPresent(key string, consume func(*topicManager)) {
	this.IfPresentOrNot(key, consume, nil)
}

func (this *concurrentTopicMap) IfPresentOrNot(key string, consume func(*topicManager), notPresent func()) {
	topicManager, exists := this.Get(key)
	if exists {
		consume(topicManager);
	} else if notPresent != nil {
		notPresent()
	}
}

func  (this *concurrentTopicMap) Delete(key string) {
	this.Lock()
	delete(this.m, key)
	this.Unlock()
}

func (this *concurrentTopicMap) GetOrAdd(key string, createFunc func()*topicManager ) *topicManager {
	this.Lock()
	topicManager, exists := this.m[key]
	if !exists {
		topicManager = createFunc()
		this.m[key] = topicManager
	}
	this.Unlock()
	return topicManager
}
