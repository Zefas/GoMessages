package main
import (
"net/http"
"GoMessages/sse"
"time"
"encoding/json"
"fmt"
)


func handleSSE(w http.ResponseWriter, req *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	sse.AddSSEHeaders(w)

	for seconds := 0; seconds <= 10; seconds++ {
		time.Sleep(1 * time.Second)

		data, _ := json.Marshal(SSEData{Name: "ABCDEFGH", Time: seconds})
		fmt.Println(string(data))

		fmt.Fprintf(w, "id: %s\n", "eventId01")
		fmt.Fprintf(w, "data: %s\n\n", string(data))
		flusher.Flush()
	}
}

type SSEData struct {
	Name string
	Time int
}
