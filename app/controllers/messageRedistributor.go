package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"GoMessages/app/infrastructure/sse"
	"GoMessages/app/messages"
	"log"
)


func NewMessageRedistributor(topicsContainer messages.ITopicsContainer) *MessageRedistributor {
	if topicsContainer == nil {
		panic("parameter cannot be nil.")
	}

	result := MessageRedistributor{}
	result.topicsContainer = topicsContainer
	return &result
}

type MessageRedistributor struct {
	topicsContainer messages.ITopicsContainer
}

func (this *MessageRedistributor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	timeoutCh := time.After(30 * time.Second)
	sse.AddSSEHeaders(w)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	topic := mux.Vars(req)["topic"]
	ch := this.topicsContainer.GetTopicManager(topic).Subscribe(topic)

	messageWriter := messageWriter{&w, flusher}
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-timeoutCh:
			log.Println("MessageRedistributor#ServeHTTP: Timeout!")
			close(ch)
			messageWriter.writeTimeout(w, flusher)
			return
		case newMessage := <-ch:
			// TODO Implement, as nothing gets send right now
			if newMessage != nil {
				log.Println("MessageRedistributor#ServeHTTP: Received Message.")
				messageWriter.writeNewMessage(newMessage, &w, flusher)
			}
		default:
			log.Println("MessageRedistributor#ServeHTTP: Waiting.")
		}
	}
}
