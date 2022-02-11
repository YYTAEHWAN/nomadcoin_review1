package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/nomadcoders_review/utils"
)

const (
	message       string = "i love you"
	privateKey    string = "3077020101042067790203cc8b6089947d592f11273330ebcb87ff4e2b532bd8a02d91c7c2a263a00a06082a8648ce3d030107a144034200047f6419b8b1e72c9709fbe31eb948a3162f8083f9ce93edc5a77abe2e90522f32447c19557811e487cb6b255957f5db4aebfd65548fe98a7e9a505c38b2b423cd"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	signature     string = "4ad1fcd38eacc6f284974bda47722700499b33124067325797332d80df9df915e17ce5d2aa2dd1ccc0bfc85f61cf1467e50dc9d17f9b52990d48b0f2dd38b7d6"
)

func reminder() {
	/*privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
	*/

}

func Start() {

	priKeyAsBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)
	_, err = x509.ParseECPrivateKey(priKeyAsBytes)
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	fmt.Printf("sigBytes : %d\n\n", sigBytes)
	fmt.Printf("rBytes : %d\n\n", rBytes)
	fmt.Printf("sBytes : %d\n\n", sBytes)

	var bigR, bigS = big.Int{}, big.Int{}
	r := bigR.SetBytes(rBytes)
	s := bigS.SetBytes(sBytes)

	fmt.Printf("%x\n\n%x\n\n", r, s)
}
