package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type VocabEntry struct {
	ID          string `json:"id"`
	Language    string `json:"language"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	ImageURL    string `json:"imageUrl"`
}

var (
	vocabStore = make(map[string][]VocabEntry)
	mu         sync.Mutex
)

func addVocab(w http.ResponseWriter, r *http.Request) {
	var entry VocabEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	vocabStore[entry.Language] = append(vocabStore[entry.Language], entry)
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
}

func getVocab(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	mu.Lock()
	entries := vocabStore[lang]
	mu.Unlock()
	json.NewEncoder(w).Encode(entries)
}

func main() {
	http.HandleFunc("/vocab", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addVocab(w, r)
		} else if r.Method == http.MethodGet {
			getVocab(w, r)
		}
	})

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
