package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nomadcoders_review/blockchain"
	"github.com/nomadcoders_review/utils"
)

var port string

type url string

// 이 MarshalText() 함수는 Marshal 함수가 struct를
// json으로 변환할 때 자동으로 호출해주는 함수이다.
// 철자가 틀리면(시그니처가 틀리면) 절대 호출되지 않을 것임
func (u url) MarshalText() ([]byte, error) {
	// 여기서 URL이 json으로 어떻게 변할지도 정할 수 있음
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// 이렇게 멋있게 써놓으면 마치 실제 설명처럼 띄워줌
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"` // 설명
	Payload     string `json:"payload,omitempty"`
}

func (u urlDescription) String() string {
	return "Hello I'm URL Description"
}

type addBlockBody struct {
	Message string
}

type errResponse struct {
	ErrMessage string `json:"errMessage"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentaion",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Block's Difficulty",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b) 이거 3개를 대체하는 밑에 함수 한줄
	json.NewEncoder(rw).Encode(data) // 위에 3개와 같음
}

func block(rw http.ResponseWriter, r *http.Request) {
	// 사용자가 전해준 변수를 받기 위해 사용하는 변수 추출 함수
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

// 중간에 사용해서 middle ware라고 하네요
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func Strat(aPort int) {
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost:%d\n", aPort)
	port = fmt.Sprintf(":%d", aPort)
	log.Fatal(http.ListenAndServe(port, router))
}
