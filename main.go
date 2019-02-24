package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseMessage represents the message sent by a server
type ResponseMessage struct {
	Message string `json:"message"`
}

func getJSONMessage(message string) ([]byte, error) {
	responseMessage := ResponseMessage{
		Message: message,
	}
	jsonResponse, err := json.Marshal(responseMessage)

	if err != nil {
		log.Printf("error occurred while json marshalling pong response: %v", err)
		return nil, err
	}

	return jsonResponse, nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		jsonResponse, err := getJSONMessage("pong")

		if err != nil {
			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusOK)
		w.Write(jsonResponse)
	})

	log.Printf("server is starting ...")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatalf("failed to start the server : %v", err)
	}

	log.Printf("started listening at :8080")
}
