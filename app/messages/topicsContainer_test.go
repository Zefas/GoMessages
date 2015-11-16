package messages
import (
	"testing"
	"time"
	"fmt"
)

func TestTopicManagerShouldNotStartWithoutSubscribersWhenMessageIsSent(t *testing.T) {
	c := NewTopicsContainer()

	c.AddMessage(&MessageInput{"topic01", "Message01"})

	cast, _ := c.(*topicsContainer)
	_, exists := cast.topicManagers["topic01"]
	if exists {
		t.Error("Topic Manager should not exist without subscribers.")
	}
}

func TestOneManagerWorks(t *testing.T) {
	c := NewTopicsContainer()
	s1 := c.Subscribe("topic01")
	s2 := c.Subscribe("topic01")

	c.AddMessage(&MessageInput{"topic01", "Message01"})

	assertSubscriberReceivedMessage("Message01", s1, t)
	assertSubscriberReceivedMessage("Message01", s2, t)


	c.UnSubscribe("topic01", s1)
	c.UnSubscribe("topic01", s2)

	cast, _ := c.(*topicsContainer)
	_, exists := cast.topicManagers["topic01"]
	if exists {
		t.Error("Topic Manager should be closed after last subscriber.")
	}
}


func TestMultipleManagerWorks(t *testing.T) {
	c := NewTopicsContainer()
	s1 := c.Subscribe("topic01")
	s2 := c.Subscribe("topic02")
	s3 := c.Subscribe("topic03")

	c.AddMessage(&MessageInput{"topic01", "Message01"})
	c.AddMessage(&MessageInput{"topic02", "Message02"})
	c.AddMessage(&MessageInput{"topic03", "Message03"})

	assertSubscriberReceivedMessage("Message01", s1, t)
	assertSubscriberReceivedMessage("Message02", s2, t)
	assertSubscriberReceivedMessage("Message03", s3, t)

	c.UnSubscribe("topic01", s1)
	c.UnSubscribe("topic02", s2)
	c.UnSubscribe("topic03", s3)

	cast, _ := c.(*topicsContainer)
	_, exists := cast.topicManagers["topic01"]
	if exists {
		t.Error("Topic Manager should be closed after last subscriber.")
	}
	cast, _ = c.(*topicsContainer)
	_, exists = cast.topicManagers["topic02"]
	if exists {
		t.Error("Topic Manager should be closed after last subscriber.")
	}
	cast, _ = c.(*topicsContainer)
	_, exists = cast.topicManagers["topic03"]
	if exists {
		t.Error("Topic Manager should be closed after last subscriber.")
	}
}

func TestTopicManagerCreationShouldBeSynchronized(t *testing.T) {
	c := NewTopicsContainer()

	iterationCount := 30
	for i := 0; i < iterationCount; i++ {
		go subscribeSimultaneously(c, fmt.Sprintf("topic0%d", i), 200 * time.Millisecond)
	}

	time.Sleep(600 * time.Millisecond)

	cast, _ := c.(*topicsContainer)
	if cast.topicManagersStarted != iterationCount {
		t.Fatalf("Race Condition - same topic manager were started at the same time. Expected - %d, got - %d", iterationCount, cast.topicManagersStarted)
	}
}

func subscribeSimultaneously(c ITopicsContainer, topic string, duration time.Duration) {
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
	time.AfterFunc(duration, func() {
		c.Subscribe(topic)
	})
}
















