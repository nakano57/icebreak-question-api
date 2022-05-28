package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kintone-labs/go-kintone"
)

var token string

type Question struct {
	Question string
	Choices  Choices
	Category []string
}

type Choices struct {
	First  string
	Second string
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	//kintone
	app := &kintone.App{
		Domain:   "out8j59wrrrm.cybozu.com",
		ApiToken: token,
		AppId:    2,
	}

	fields := []string{"question", "category", "choice", "first", "second"}
	records, err := app.GetRecords(fields, "isUse in (\"使用\")")
	if err != nil {
		log.Fatal(err)
	}
	n := rand.Intn(len(records))
	question := records[n]

	q := new(Question)
	q.Question = string(question.Fields["question"].(kintone.SingleLineTextField))
	q.Category = question.Fields["category"].(kintone.MultiSelectField)
	q.Choices.First = string(question.Fields["first"].(kintone.SingleLineTextField))
	q.Choices.Second = string(question.Fields["second"].(kintone.SingleLineTextField))

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	json.NewEncoder(w).Encode(q)
}

func main() {
	//Token周り
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("ファイルが存在しません: %v", err)
	}
	token = os.Getenv("TOKEN")

	// ルーターのイニシャライズ
	r := mux.NewRouter()

	//テストデータ
	//root := Question{Question: "元気ですかーッ！", Category: "体調", Choices: &Choices{First: "元気です", Second: "めっちゃ元気です!!!!!!"}}
	//q = append(q, root)

	// ルート(エンドポイント)
	r.HandleFunc("/api/random", getQuestions).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
