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

func getQuestions(w http.ResponseWriter, r *http.Request) {
	//kintone
	app := &kintone.App{
		Domain:   "out8j59wrrrm.cybozu.com",
		ApiToken: token,
		AppId:    2,
	}

	fields := []string{"question", "category", "first", "second"}
	records, err := app.GetRecords(fields, "isUse in (\"使用\")")
	if err != nil {
		log.Fatal(err)
	}
	n := rand.Intn(len(records))
	question := records[n]
	fmt.Println(question.Fields)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(records[n])
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
