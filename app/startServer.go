package app
import (
"github.com/gorilla/mux"
"net/http"
"fmt"
"GoMessages/app/controllers"
"GoMessages/app/messages"
)


func StartServer() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	topicsContainer := messages.NewTopicsContainer()
	addNewMessageHandler := controllers.NewAddNewMessageHandler(topicsContainer)
	messageRedistributor := controllers.NewMessageRedistributor(topicsContainer)

	r.Handle("/infocenter/{topic}", addNewMessageHandler).Methods("POST")
	r.Handle("/infocenter/{topic}", messageRedistributor).Methods("GET")

	r.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, req.URL.Path[1:])
	})

	http.Handle("/", r)
	err := http.ListenAndServe(":9090", nil);
	if err != nil {
		fmt.Println("Error while starting server!")
	}
}