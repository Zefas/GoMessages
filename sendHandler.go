package main
import (
"net/http"
"fmt"
)


func handleSend(w http.ResponseWriter, req *http.Request) {
	if !validSendRequest(req) {
		fmt.Println("Error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Data."))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}


// TODO Validate Request
func validSendRequest(req *http.Request) bool {
	return false;
}
