package controllers
import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"errors"
	"GoMessages/app/infrastructure/httphelper"
	"GoMessages/app/messages"
)


func NewAddNewMessageHandler(topicsContainer messages.ITopicsContainer) *AddNewMessageHandler {
	result := AddNewMessageHandler{}
	result.topicsContainer = topicsContainer
	return &result
}

type AddNewMessageHandler struct {
	topicsContainer messages.ITopicsContainer
}

func (this *AddNewMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data, err := parseRequest(req)
	if err != nil {
		fmt.Println("Error: Cannot read message.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Data."))
		return
	}
	w.WriteHeader(http.StatusNoContent)

	this.topicsContainer.AddMessage(data)
}

func parseRequest(req *http.Request) (*messages.MessageInput, error) {
	message, err := httphelper.ReadBodyAsString(req)
	if err != nil {
		return nil, errors.New("Invalid Request Parameters");
	}
	return &messages.MessageInput{mux.Vars(req)["topic"], message}, nil
}
