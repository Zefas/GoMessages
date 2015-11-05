package main
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/static/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, req.URL.Path[1:])
	})
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/receive", handleReceive)
	http.HandleFunc("/sse", handleSSE)

	err := http.ListenAndServe(":9090", nil);
	if err != nil {
		fmt.Println("Error while starting server!")
	}
}
