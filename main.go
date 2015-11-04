package main
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/receive", handleReceive)

	err := http.ListenAndServe(":9090", nil);
	if err != nil {
		fmt.Println("Error while starting server!")
	}

}




