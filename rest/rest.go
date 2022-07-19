package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nomadcoders_review/blockchain"
	"github.com/nomadcoders_review/p2p"
	"github.com/nomadcoders_review/utils"
	"github.com/nomadcoders_review/wallet"
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

type balanceResponse struct {
	Address      string `json:"address"`
	TotalBalance int    `json:"totalBalance"`
}

type errResponse struct {
	ErrMessage string `json:"errMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type myWalletResponse struct {
	Address string
}

type addPeerPayload struct {
	Address string
	Port    string
}

func (u urlDescription) String() string {
	return "Hello I'm URL Description"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentaion", // 설명
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
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "See balance for an address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to Web Socket",
		},
		{
			URL:         url("/peers"),
			Method:      "GET",
			Description: "see peers",
		},
		{
			URL:         url("/peers"),
			Method:      "POST",
			Description: "Add a peer",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b) 이거 3개를 대체하는 밑에 함수 한줄
	json.NewEncoder(rw).Encode(data) // 위에 3개와 같음
}

// 사용자가 전해준 변수를 받기 위해 사용하는 변수 추출 함수
func block(rw http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		// utils.HandleErr(json.NewDecoder(r.Body).Decode()) 아마 나중엔 Tx를 가져오지 않을까
		newBlock := blockchain.AddBlock(blockchain.Blockchain())
		p2p.BroadcastNewBlock(newBlock)
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
	blockchain.Status(blockchain.Blockchain(), rw, r)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	encoder := json.NewEncoder(rw)

	// totalBool := r.URL.Query().Get("total")
	// if totalBool == "true" {
	// 	total := blockchain.TotalBalanceByAddress(address)
	// 	utils.HandleErr(encoder.Encode(balanceResponse{address, total})) // 여러 개의 변수를 보내고 싶을 때
	// } else {
	// 	OwnedTxOuts := blockchain.Blockchain().BalanceByAddress(address)
	// 	utils.HandleErr(encoder.Encode(OwnedTxOuts))
	// }

	totlaResult := r.URL.Query().Get("total")
	switch totlaResult {
	case "true":
		amount := blockchain.TotalBalanceByAddress(address, blockchain.Blockchain())
		utils.HandleErr(encoder.Encode(balanceResponse{address, amount})) // 여러 개의 변수를 보내고 싶을 때
	default:
		OwnedTxOuts := blockchain.UTxOutsByAddress(address, blockchain.Blockchain())
		utils.HandleErr(encoder.Encode(OwnedTxOuts))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {

	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	fmt.Println(payload.To, "에게", payload.Amount, "만큼의 코인을 보냅니다.")
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount) // 문제는 여기다 // 문제는 여기다 // 문제는 여기다 // 문제는 여기다
	if err != nil {
		if tx == nil {
			fmt.Println("transactions함수에서 왜 안돼 err 안으로 들어옴")
		}
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errResponse{err.Error()})
		// 잔액이 부족합니다 or 유효하지 않은 트랜잭션의 의미
		return
	}
	p2p.BrodcastNewTx(tx) // Header가 Brodcast를 기다리지 않도록 go routine을 사용해줄 수 있어
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	utils.HandleErr(json.NewEncoder(rw).Encode(myWalletResponse{Address: address}))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Strat(aPort int) {
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/wallet/", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost:%d\n", aPort)
	port = fmt.Sprintf(":%d", aPort)
	log.Fatal(http.ListenAndServe(port, router))
}
