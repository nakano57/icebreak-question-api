package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Question struct {
	QUESTION string   `json:"question"`
	CATEGORY string   `json:"category"`
	CHOICES  *Choices `json:"choices"`
}

type Choices struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

var q []Question

func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(q)
}

func main() {
	// ルーターのイニシャライズ
	r := mux.NewRouter()

	//テストデータ
	q = append(q, Question{QUESTION: "元気ですかーッ！", CATEGORY: "体調", CHOICES: &Choices{First: "元気です", Second: "めっちゃ元気です!!!!!!"}})

	// ルート(エンドポイント)
	r.HandleFunc("/api/random", getQuestions).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
