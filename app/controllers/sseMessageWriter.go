package controllers
import (
	"fmt"
	"net/http"
	"GoMessages/app/messages"
)

type messageWriter struct {
	w *http.ResponseWriter
	flusher http.Flusher
}

func (*messageWriter) writeNewMessage(newMessage *messages.MessageOutput, w *http.ResponseWriter, flusher http.Flusher) {
	fmt.Fprintf(*w, "id: %d\n", newMessage.Id)
	fmt.Fprintf(*w, "event: %s\n", "msg")
	fmt.Fprintf(*w, "data: %s\n\n", newMessage.Data)
	flusher.Flush()
}

func (*messageWriter) writeTimeout(w http.ResponseWriter, flusher http.Flusher) {
	fmt.Fprintf(w, "event: %s\n", "timeout")
	fmt.Fprintf(w, "data: %s\n\n", "30s")
	flusher.Flush()
}