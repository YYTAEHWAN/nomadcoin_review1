package blockchain

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nomadcoders_review/utils"
	"github.com/nomadcoders_review/wallet"
)

const (
	minerReward int = 50
)

var ErrorNoMoney error = errors.New("not enough money")
var ErrorNotValid error = errors.New("Tx invalid")

type mempool struct {
	Txs map[string]*Tx
	M   sync.Mutex
}

var m *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
	Signature string   `json:"signature"`
}

type TxIn struct {
	TxId      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txID"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) makeId() {
	t.Id = utils.Hashing(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}

func (t *Tx) MakeTxTimestamp() {
	fmt.Print("\tTx 시간 설정\t")
	t.Timestamp = int(time.Now().Unix())
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, mtx := range Mempool().Txs {
		for _, input := range mtx.TxIns {
			if input.TxId == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{Address: address, Amount: minerReward},
	}
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.makeId()

	return tx
}

// func which is used when a new transaction is added
func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	fmt.Println("---AddTx함수 실행---")
	// 72a7ca6ab3cd7405040f1578f8b2d809de4c0e9c579a37ae89c06a2af48bff823996cb2e95032ac16adf590675d9fed3e9bfaf172e081b1d6198d9620b2b76f5
	fmt.Println("wallet 주소는 ", wallet.Wallet())
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		fmt.Println("AddTx에서 에러 발생!")
		fmt.Println("---AddTx함수 종료 1(err)---")
		if err == ErrorNotValid {
			return nil, ErrorNotValid
		}
		return nil, ErrorNoMoney
	}
	m.Txs[tx.Id] = tx
	fmt.Println("---AddTx함수 종료 2(normal)---")
	return tx, nil
}

func makeTx(from string, to string, amount int) (*Tx, error) {
	fmt.Println("---makeTx 함수 실행---")
	if TotalBalanceByAddress(from, Blockchain()) < amount {
		fmt.Println("makeTx함수에서  TotalBalance 부족으로 에러 발생")
		return nil, errors.New("not enough money-1")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	fmt.Println("makeTx의 UTxOutsByAddress 함수 사용")
	for i, uTxOut := range UTxOutsByAddress(from, Blockchain()) {
		if total >= amount {
			break
		}
		fmt.Print(i, "번째 uTxOut")
		fmt.Println(" ", uTxOut.TxID, uTxOut.Index, uTxOut.Amount)
		total += uTxOut.Amount
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
	}
	fmt.Println("makeTx의 UTxOutsByAddress 함수 종료")
	if change := total - amount; change != 0 {
		fmt.Println("잔돈은 ", change)
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	totxOut := &TxOut{to, amount}
	txOuts = append(txOuts, totxOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.makeId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		fmt.Println("여기서 에러?-209번째 줄")
		fmt.Println("---makeTx 함수 종료 1(err)---")
		return nil, ErrorNotValid
	}
	fmt.Println("---makeTx 함수 종료 2(normal)---")
	return tx, nil
}

func TotalBalanceByAddress(address string, b *blockchain) int {
	fmt.Println("---TotalBalanceByAddress 함수 시작---")
	var total int
	txOuts := UTxOutsByAddress(address, b)
	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	fmt.Println("---TotalBalanceByAddress 함수 종료---")
	return total
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {

	var uTxOuts []*UTxOut
	creatorIds := make(map[string]bool)
	fmt.Print(" test1 ")
	for _, block := range Blocks(b) { // 모든 블록을 다 본다.
		fmt.Print("test2 ")
		for _, Tx := range block.Transaction {
			fmt.Print("test3 ")
			for _, input := range Tx.TxIns {
				if input.Signature == "COINBASE" { // Tx.TxIn이 누구의 것이라고 볼 필요 없이 COINBASE이기 때문에 넘어간다.
					fmt.Print("COINBASE 여서 통과 ")
					break
				}
				fmt.Println(" - test4 - ")
				if FindTx(b, input.TxId).TxOuts[input.Index].Address == address {
					fmt.Println("Spent TxOutsByAddress found!")
					// 쓴 TxOut들의 TxId 기록을 통해 기록한다.
					creatorIds[input.TxId] = true // TxId만 넣어도 어떤 TxOut인지 알기 때문에 (2개 중에 하나라서) TxId만 찾아도 된다.
				}
			}
			fmt.Print("test999 ")
			// 이제 안 쓴 TxOut들 UTxOuts에다가 집어넣기
			for index, output := range Tx.TxOuts { // 모든 블록의 Tx들의 TxOuts에서
				fmt.Print("wow 1 ")
				if output.Address == address { // TxOut의 address가 확인하고자 하는 address 와 같으면
					fmt.Print("wow 2 ")
					if ok := creatorIds[Tx.Id]; !ok { // 또 그 TxOut이 사용되지 않았으면 (TxId가 기록되어있지 않으면)
						fmt.Print("wow 3 ")
						uTxOut := &UTxOut{Tx.Id, index, output.Amount} // 그럼 안쓴 TxOut이니까 UTxOut에 기록
						if !(isOnMempool(uTxOut)) {                    // 안쓴 TxOut 중에서도 Mempool에 있는 TxOut인지 확인하고
							fmt.Println("wow 4 ")
							uTxOuts = append(uTxOuts, uTxOut) // 진짜 안 쓴 TxOut이면 UTxOut 만들어 추가하기
						}
					}
				}
			}
		}
	}
	// 의문들 다 해결
	// 한 Tx에는 TxIns, TxOuts 가 있는데
	// TxIns 는 여러개의 TxOut을 합쳐 TxIns를 하나의 돈 덩어리로 보게끔 하는 것이고
	// TxOuts 는 오직 두개 이하의 TxOut만 존재할 수 있음  ex) 상대방 주소로 하나, 나에게로 잔돈 하나 이렇게 두개
	fmt.Print("uTxOuts의 간략한 출력", uTxOuts, "  ")
	fmt.Println(uTxOuts[0].TxID, uTxOuts[0].Index, uTxOuts[0].Amount)
	fmt.Println("UTXO End!")
	return uTxOuts
}

// 해당 Tx에 들어간 TxIn들이 유효한지를 검증하는 함수
// TxIn의 signature(서명)을 검증한다
// ex) TxIn.TxID, 서명, publicKey를 넣어서 사용된 TxIn이 그 address의 사람에 의해서 서명되었는지를 확인
func validate(t *Tx) bool {
	fmt.Println("---validate 함수 실행---")
	var valid = true
	fmt.Println("시작시에 valid는", valid)

	fmt.Println("t.TxIns는 ", t.TxIns)
	for i, txIn := range t.TxIns {
		fmt.Printf("%d번째 t.TxIns\t", i)
		prevTx := FindTx(Blockchain(), txIn.TxId)
		if prevTx == nil {
			valid = false
			fmt.Println("validate 함수 오류 발생")
			break
		}
		fmt.Println("아직까지는 valid가 ", valid)
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, t.Id, address)
	}
	if valid {
		fmt.Println("---validate 함수 종료---")
	} else {
		fmt.Println("--- valid = false인채로 validate 함수 종료---")
	}
	return valid
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.M.Lock()
	defer m.M.Unlock()
	fmt.Println("AddPeerTx에서 추가할 tx는 ", tx)
	m.Txs[tx.Id] = tx
}
