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

type idAndIndexSlice struct {
	InputTxId string
	Index     []int
}

func (i idAndIndexSlice) appendInt(index int) {
	i.Index = append(i.Index, index)
}

// type testVar struct {
// 	v []*idAndIndexSlice
// }

func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {

	var uTxOuts []*UTxOut
	blocks := b.Blocks()
	creatorTx := make(map[string]bool)
	creatorTxIndex := make(map[string]int)

	var testVar []*idAndIndexSlice

	for _, block := range blocks {
		for _, tx := range block.Transaction {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTx[input.TxId] = true
					creatorTxIndex[input.TxId] = input.Index

					idAndIndex := &idAndIndexSlice{input.TxId, nil}
					idAndIndex.appendInt(input.Index)
					testVar = append(testVar, idAndIndex)
				}
			}
		}
		for _, tx := range block.Transaction {
			for index, output := range tx.TxOuts {
				if output.Owner == address { // 주소가 같으면 일단 요구하는 사람의 TxOut들이고
					// 이제 그 TxOut이 사용됐는지 안됐는지를 확인
					if ok := creatorTx[tx.Id]; !ok { // TxIn으로 사용된 TxOut이 있는 tx.Id 리스트를 사용하여 사용되었을 가능성이 있는지 확인
						// 만약 이 TxOut이 존재하는 tx.Id가  사용된 TxOut tx.Id 리스트에 존재하지 않는다면 무조건 UTxOuts에 들어가야함
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						uTxOuts = append(uTxOuts, uTxOut)
					} else { // 한 tx.Id 에는 같은 address로 돈을 보내는 많은 txOut이 있을 수 있고 그러므로 인덱스를 확인해야 함
						if !(creatorTxIndex[tx.Id] == index) {
							// 인덱스가 같지 않으면 사용되지 않은 것이므로 추가해주는 것이 맞는데
							// 여기서 문제가 발생함 // 여기서 문제가 발생함 // 여기서 문제가 발생함 // 여기서 문제가 발생함 // 여기서 문제가 발생함 // 여기서 문제가 발생함
							// 말했듯 한 tx.Id라도 같은 address로 돈을 보내는 많은 tx.Out이 있을 수 있기 때문에
							// 한 tx.Id에서 사용된 txOut, 사용되지 않은 txOut을 구분해야하는데
							// 그 구분은 index로 하면 된다고 하지만 그 인덱스를 어디에 저장하여야 할지 모르겠음
							uTxOut := &UTxOut{tx.Id, index, output.Amount}
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
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
