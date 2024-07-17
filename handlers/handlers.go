package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"testtask/db"

	"github.com/gorilla/mux"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type NotAllowedHandler struct{}

// ServeHTTP implements http.Handler.
func (h NotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(w, r)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Default Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	w.WriteHeader(http.StatusNotFound)
	Body := r.URL.Path + " is not supported. Thanks for visiting!"
	fmt.Fprintf(w, "%s", Body)
}

// MethodNotAllowedHandler is executed when the HTTP method is incorrect
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, "with method")
	w.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed!\n"
	fmt.Fprintf(w, "%s", Body)
}

// AddHandler is for adding a new message
func AddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AddHandler Serving:", r.URL.Path, "from", r.Host)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if len(d) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No input!")
		return
	}

	var messages = []db.Message{}
	err = json.Unmarshal(d, &messages)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(messages)

	result := db.InsertMessage(messages[1])
	if !result {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetAllHandler is for getting all data from the messages database
func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllHandler Serving:", r.URL.Path, "from", r.Host)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if len(d) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No input!")
		return
	}

	var user = db.Message{}
	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		return
	}

	err = SliceToJSON(db.ListAllMessages(), w)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetMessageDataHandler + GET returns the full record of a messages
func GetMessageDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetMessageDataHandler Serving:", r.URL.Path, "from", r.Host)
	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("ID value not set!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := db.FindMessageID(intID)
	if t.ID != 0 {
		err := t.ToJSON(w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		return
	}

	log.Println("User not found:", id)
	w.WriteHeader(http.StatusBadRequest)
}

// DeleteHandler is for deleting an existing user
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteHandler Serving:", r.URL.Path, "from", r.Body)

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("ID value not set!")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var user = db.Message{}
	err := user.FromJSON(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("id:", err)
		return
	}

	t := db.FindMessageID(intID)
	if t.Text != "" {
		log.Println("About to delete:", t)
		deleted := db.DeleteMessageID(intID)
		if deleted {
			log.Println("Message deleted:", t)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			log.Println("Message ID not found:", id)
			w.WriteHeader(http.StatusNotFound)
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}
