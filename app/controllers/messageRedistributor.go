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
	newMessagesCh := this.topicsContainer.Subscribe(topic)
	log.Printf("MessageRedistributor#ServeHTTP: subscribed via channel: %v\n", newMessagesCh)

	messageWriter := messageWriter{&w, flusher}
	for {
		select {
		case <-timeoutCh:
			log.Println("MessageRedistributor#ServeHTTP: Timeout!")
			this.topicsContainer.UnSubscribe(topic, newMessagesCh)
			messageWriter.writeTimeout(w, flusher)
			return
		case newMessage := <-newMessagesCh:
			log.Println("MessageRedistributor#ServeHTTP: Received Message.")
			messageWriter.writeNewMessage(&newMessage, &w, flusher)
		default:
			// Do Nothing
		}
	}
}
