package main
import (
"net/http"
"fmt"
)

func handleReceive(w http.ResponseWriter, req *http.Request) {
	if !validReceiveRequest(req) {
		fmt.Println("Error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Data."))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// TODO Validate Request
func validReceiveRequest(req *http.Request) bool {
	return false;
}
