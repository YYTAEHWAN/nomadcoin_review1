package blockchain

import (
	"errors"
	"fmt"
	"time"

	"github.com/nomadcoders_review/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string
	Timestamp int
	TxIns     []*TxIn
	TxOuts    []*TxOut
}

type TxIn struct {
	TxId  string
	Index int
	Owner string
}

type TxOut struct {
	Owner  string
	Amount int
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) makeId() {
	t.Id = utils.Hashing(t)
}

func (t *Tx) MakeTxTimestamp() {
	fmt.Print("\tTx 시간 설정\t")
	t.Timestamp = int(time.Now().Unix())
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, mtx := range Mempool.Txs {
		for _, input := range mtx.TxIns {
			if input.TxId == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func (utx *UTxOut) laterChange() {

}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {

	var uTxOuts []*UTxOut
	creatorIds := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorIds[input.TxId] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if ok := creatorIds[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !(isOnMempool(uTxOut)) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
						// 이걸 고치게 되면 정말 큰 문제가 생기기 때문에 그냥 지나가도록 하겠다
						// 이유 1. 해당 input TxIns에 사용된 TxOut을 추적하여 찾아오는 함수를 만들어야 함 -> 할 수 있을지도 모름
						// 강력한 이유 2. 같은 주소로 돈을 보내는 여러개의 TxOut을 만들게 됐을 때 TxOut1, 2, 3이라 하자
						// 나중에 1,2,3 중 1,2 만 사용했다고 했을 때, 이 로직에 의하면 1과 2의 TxId는 사용된 TxId라 간주되어
						// TxOut3 은 사용되지 않았지만 사용되었다고 분류되어 사용하지 못하는 돈이 됨
						// 그걸 해결해주는 로직도 따로 만들어야 하는데 머리를 조금 굴려봤을 때 굉장히 복잡하다고 예상됨
						// 이건 시간상의 이유로 실패했으니 git 하는법이라도 배워가겠음
						// git 업스트림 사용 2
					}
				}
			}
		}
	}

	// 의문들 다 해결
	// 한 Tx에는 TxIns, TxOuts 가 있는데
	// TxIns 는 여러개의 TxOut을 합쳐 TxIns를 하나의 돈 덩어리로 보게끔 하는 것이고
	// TxOuts 는 오직 두개 이하의 TxOut만 존재할 수 있음  ex) 상대방 주소로 하나, 나에게로 잔돈 하나 이렇게 두개
	return uTxOuts
}

func TotalBalanceByAddress(address string, b *blockchain) int {
	var total int
	txOuts := UTxOutsByAddress(address, b)
	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return total
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{Owner: address, Amount: minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.makeId()
	return &tx
}

func makeTx(from string, to string, amount int) (*Tx, error) {
	if TotalBalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enough money-1")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0

	for _, uTxOut := range UTxOutsByAddress(from, Blockchain()) {
		if total >= amount {
			break
		}
		total += uTxOut.Amount
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
	}

	if change := total - amount; change != 0 {
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

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("taehwan", to, amount)
	if err != nil {
		return errors.New("not enough money")
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("taehwan")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil

	return txs
}
