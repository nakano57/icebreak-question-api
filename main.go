package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/koron/go-dproxy"
)

var token string

type Root struct {
	Records []Question
}

type Question struct {
	Question string
	Category string
	Choices  *Choices
}

type Choices struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

var q interface{}

type Parameters struct {
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	//kintone
	jsonStr := `{
		"app": 2,
		"id": 1,
		"query": "isUse in (\"使用\")",
		"fields": ["question","category", "first", "second"]}`
	url := "https://out8j59wrrrm.cybozu.com/k/v1/records.json"
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("X-Cybozu-API-Token", token)
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(response))

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//json.Unmarshal(response, &q)
	//うまくいかない
	if err := json.Unmarshal(response, &q); err != nil {
		panic(err)
	}
	qs := dproxy.New(q).M("records")
	fmt.Println(qs)
	//fmt.Println(err)
	json.NewEncoder(w).Encode(qs)
}

func main() {
	//Token周り
	err := godotenv.Load(".env")

	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("ファイルが存在しません: %v", err)
	}

	token = os.Getenv("TOKEN")

	// ルーターのイニシャライズ
	r := mux.NewRouter()

	//テストデータ
	//q = append(q, Question{QUESTION: "元気ですかーッ！", CATEGORY: "体調", CHOICES: &Choices{First: "元気です", Second: "めっちゃ元気です!!!!!!"}})

	// ルート(エンドポイント)
	r.HandleFunc("/api/random", getQuestions).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
