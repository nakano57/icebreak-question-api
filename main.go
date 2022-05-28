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
)

var token string

type Records struct {
	Records []*Question `json:"records"`
}

type Question struct {
	Question string   `json:"question"`
	Category string   `json:"category"`
	Choices  *Choices `json:"choices"`
}

type Choices struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

var qs interface{}
var q Records

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
	fmt.Println("レスポンス")
	fmt.Println(string(response))

	//うまくいかない
	if err := json.Unmarshal([]byte(response), &qs); err != nil {
		panic(err)
	}

	//fmt.Println(qs["records"])
	mapData := qs.(map[string]interface{})
	delete(mapData, "totalCount")

	jsonStr2, err := json.Marshal(mapData)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Totalcount削除")
	fmt.Println(string(jsonStr2))

	qs3 := make([]*Question, 0)

	//うまくいかない
	if err := json.Unmarshal([]byte(jsonStr2), &qs3); err != nil {
		panic(err)
	}

	fmt.Println(qs3)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(qs3)
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
	//root := Question{Question: "元気ですかーッ！", Category: "体調", Choices: &Choices{First: "元気です", Second: "めっちゃ元気です!!!!!!"}}
	//q = append(q, root)

	// ルート(エンドポイント)
	r.HandleFunc("/api/random", getQuestions).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
