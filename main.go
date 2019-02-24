package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/badger"
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
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("crash"))

			if err == badger.ErrKeyNotFound {
				jsonResponse, err := getJSONMessage("pong")

				if err != nil {
					return err
				}

				log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusOK)
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResponse)
				return nil
			}

			if err != nil {
				return err
			}

			crash, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			if bytes.Compare(crash, []byte("yes")) == 0 {
				jsonResponse, err := getJSONMessage("crashing")

				if err != nil {
					return err
				}

				log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write(jsonResponse)
				return nil
			}

			jsonResponse, err := getJSONMessage("pong")

			if err != nil {
				return err
			}

			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusOK)
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
			return nil
		})

		if err != nil {
			log.Printf("error occurred: %v", err)
			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/crashping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ttlTime := r.URL.Query().Get("time")
		crash := r.URL.Query().Get("crash")

		ttl, err := time.ParseDuration(ttlTime)

		if err != nil {
			log.Printf("error occurred while parsing time: %v", err)
			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.Update(func(txn *badger.Txn) error {
			err := txn.SetWithTTL([]byte("crash"), []byte(crash), ttl)
			return err
		})

		if err != nil {
			log.Printf("error occurred while updating 'crash' key in db: %v", err)
			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonResponse, err := getJSONMessage(fmt.Sprintf("crash set to %v with ttl %v !", crash, ttlTime))

		if err != nil {
			log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("RequestMethod:%v RequestUrl:%v ResponseCode:%v", r.Method, r.URL.Path, http.StatusOK)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})

	log.Printf("server is starting ...")
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatalf("failed to start the server : %v", err)
	}

	log.Printf("started listening at :8080")
}
