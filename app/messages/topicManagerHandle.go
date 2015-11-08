package messages


type topicManagerHandle struct {
	addMessageCh chan <- string
	addSubscriberCh chan <- chan MessageOutput
	unSubscribeCh chan <- chan MessageOutput
}

func (this *topicManagerHandle) subscribe() {
	topicCh := make(chan MessageOutput)
	this.addSubscriberCh <- topicCh
}

func (this *topicManagerHandle) unSubscribe(removeCh chan MessageOutput) {
	this.unSubscribeCh <- removeCh
}