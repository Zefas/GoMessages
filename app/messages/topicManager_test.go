package messages
import (
	"testing"
	"time"
)


func TestShouldNotReceiveMessagesWhenNoneWasSent(t *testing.T) {
	// Topic
	unsubscribeCh := make(chan unSubscribeResult)
	t1 := newTopicManager("topic01", unsubscribeCh)
	go t1.startRunning()

	// Subscribers
	var s1 chan MessageOutput = make(chan MessageOutput)
	t1.getHandle().addSubscriberCh <- s1

	forLoop:
	for {
		select {
		case <-s1:
			t.Error("Should not receive any messages when none were sent.")
		case <-time.After(1 * time.Second):
			break forLoop
		}
	}

	t1.getHandle().unSubscribeCh <- s1

	closed := <- unsubscribeCh
	if !closed.stopped {
		t.Error("Topic Manager should be stopped.")
	}
}

func TestMultipleSubscribersReceiveMessagesWhenOneWasSent(t *testing.T) {
	// Topic
	unsubscribeCh := make(chan unSubscribeResult)
	t1 := newTopicManager("topic01", unsubscribeCh)
	go t1.startRunning()

	// Subscribers
	var s1 chan MessageOutput = make(chan MessageOutput)
	t1.getHandle().addSubscriberCh <- s1
	var s2 chan MessageOutput = make(chan MessageOutput)
	t1.getHandle().addSubscriberCh <- s2

	// Send Message
	t1.getHandle().addMessageCh <- "message01"

	assertSubscriberReceivedMessage("message01", s1, t)
	assertSubscriberReceivedMessage("message01", s2, t)

	t1.getHandle().unSubscribeCh <- s1
	closed := <- unsubscribeCh
	if closed.stopped {
		t.Error("Topic Manager should not stop when subscribers still remain.")
	}

	t1.getHandle().unSubscribeCh <- s2
	closed = <- unsubscribeCh
	if !closed.stopped {
		t.Error("Topic Manager should be stopped.")
	}
}


func TestExitsAfterLastSubscriberUnsubscribed(t *testing.T) {
	// Topic
	unsubscribeCh := make(chan unSubscribeResult)
	t1 := newTopicManager("topic01", unsubscribeCh)
	go t1.startRunning()

	// Subscribe
	var s1 chan MessageOutput = make(chan MessageOutput)
	t1.getHandle().addSubscriberCh <- s1
	var s2 chan MessageOutput = make(chan MessageOutput)
	t1.getHandle().addSubscriberCh <- s2

	// UnSubscribe
	t1.getHandle().unSubscribeCh <- s1
	closed := <- unsubscribeCh
	if closed.stopped {
		t.Error("Topic Manager should not stop when subscribers still remain.")
	}

	t1.getHandle().unSubscribeCh <- s2
	closed = <- unsubscribeCh
	if !closed.stopped {
		t.Error("Topic Manager should be stopped.")
	}

	// Attempt to send
	forLoop:
	for {
		select {
		case t1.getHandle().addSubscriberCh <- s1:
			t.Error("Sending should not succeed when Topic Manager is finished")
		case <-time.After(1 * time.Second):
			break forLoop
		}
	}
}


func assertSubscriberReceivedMessage(expected string, subscriberCh <-chan MessageOutput, t *testing.T) {
	msg := <-subscriberCh
	if msg.Data != expected {
		t.Errorf("Subscriber haven't received correct message: expected - '%s', actual - '%s'", expected, msg.Data)
	}
}

