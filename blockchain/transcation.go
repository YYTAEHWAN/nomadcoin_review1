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

func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {

	var uTxOuts []*UTxOut
	creatorIds := make(map[string]bool)

	for _, block := range Blockchain().Blocks() {
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
						uTxOuts = append(uTxOuts, uTxOut)
						// 이럴 경우 내가 걱정되는 건
						// input에 사용된 txId는 output의 txId라는 건데
						// 한 블록에 같은 주소로 여러개의 output Txs가 생성되었을 때
						// 6개 TxOut 중 4개만 쓰이고 2개는 안쓰였다면?
						// 그래도 이 로직에 따르면 그 2개는 spentTxOut 으로 분류될텐데 이 로직을 어떻게 고쳐야 할까
					}
				}
			}
		}
	}
	// 의문 2 TxIn 은 여러개의 TxOut이 모여서 만들어진다고 했는데
	// 그럼 하나의 TxIn은 하나의 TxId, Index를 만드는데 어떻게 여러개의 TxOut을 합칠 수 있지?

	// 의문 3
	// 그럼 어떻게 해결하지?

	return uTxOuts
}

func (b *blockchain) TotalBalanceByAddress(address string) int {
	var total int
	txOuts := b.UTxOutsByAddress(address)
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
	if Blockchain().TotalBalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}

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
