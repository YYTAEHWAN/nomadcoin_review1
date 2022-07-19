package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/nomadcoders_review/utils"
)

const (
	fileName = "nomadcoin.wallet"
)

/*
func reminder() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)

	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	fmt.Printf("privateKey : %x\n\n", keyAsBytes)

	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleErr(err)
	fmt.Printf("what these are? : \n r : %d\n s : %d\n", r, s)

	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("signature : %x\n\n", signature)

	/////////////////////////////
	priKeyAsBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)
	restoredKey, err := x509.ParseECPrivateKey(priKeyAsBytes)
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	// fmt.Printf("sigBytes : %d\n\n", sigBytes)
	// fmt.Printf("rBytes : %d\n\n", rBytes)
	// fmt.Printf("sBytes : %d\n\n", sBytes)

	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	ok := ecdsa.Verify(&restoredKey.PublicKey, hashAsBytes, &bigR, &bigS)
	fmt.Println(ok)
}
*/

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
	// no -> 존재 안함?
	// yes -> 존재함?
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistPrivKey(privKey *ecdsa.PrivateKey) {
	privKeyAsBytes, err := x509.MarshalECPrivateKey(privKey)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, privKeyAsBytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	privKeyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	privKey, err := x509.ParseECPrivateKey(privKeyAsBytes)
	utils.HandleErr(err)
	return privKey
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromk(privkey *ecdsa.PrivateKey) string {
	return encodeBigInts(privkey.X.Bytes(), privkey.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(signature string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	lastHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := &big.Int{}, &big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(lastHalfBytes)

	return bigA, bigB, nil
}

func Verify(signature, payload, address string) bool {
	fmt.Println("---wallet.go에 있는 Verify함수 실행---")
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)

	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&publicKey, payloadAsB, r, s)
	fmt.Println("Verify의 값 ok는 ", ok)
	fmt.Println("---wallet.go에 있는 Verify함수 종료---")
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		// has a wallet already?
		if hasWalletFile() {
			// yes -> restore from file
			w.privateKey = restoreKey()
		} else {
			// no -> create new privKey, save to file
			privKey := createPrivKey()
			persistPrivKey(privKey)
			w.privateKey = privKey
		}
		w.Address = aFromk(w.privateKey)
	}
	return w
}

func Start() {
	Wallet()
}
