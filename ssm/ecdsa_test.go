package ssm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
)

func TestCreateECDSAKeyPair(t *testing.T) {
	privKey, pubKey, cpubKey, _ := CreateECDSAKeyPair()
	fmt.Println("privKey=", privKey)
	fmt.Println("pubKey=", pubKey)
	fmt.Println("cpubKey=", cpubKey)
}

func TestSignMessage(t *testing.T) {
	privKey := "117ff4db9619a76e1a43c96cbfbc33db1d88138fdb25100192bb5cdf31d07e74"
	tx := types.NewTransaction(
		0, // nonce
		common.HexToAddress("0x0000"),
		big.NewInt(1000000000000000000), // 1 ETH
		21000,                           // gas limit
		big.NewInt(20000000000),         // gas price
		nil,                             // data
	)
	txHash := tx.Hash()
	fmt.Println("msgHash=", txHash.Bytes())
	fmt.Println("msgHashLen=", len(txHash))
	signature, err := SignMessage(privKey, txHash.Bytes())
	if err != nil {
		fmt.Println("sign tx fail")
	}
	fmt.Println("Signature: ", signature)
	fmt.Println("Signature len: ", len(signature))
}

func TestVerifySign(t *testing.T) {
	msgHash := "0775b8c5fbfee9950d678f9a8369eb02847799d63b7a48a1b67e52ed31cc7892"
	// 04fd13e3e6895a2d9d38821e47f60f6c701da7bf243cc1f0a21a15ff8b09322ceda16bbc6d07df0e408dc19ffb756a75bf0c7cbf6836ef89594e0d97cbfeff319c
	publicKey := "04d20a0b30fd29e121d3f267817d2886fb18cbcb5e6a59baf6183a7762bf89ef7c94d55ef0870d03e510b495d02d2f4fcb4b1a5546139bfed3ecac951674621869"
	signature := "ecbc9b61fa0f77c670895eaaa9b5c6f8014aab70e79a4b28c4b3f1359c3d8bee300ebcb17aed8bc3ee339cd43521a06054253d0cdcd6ecd1713f69c55927d1be01"
	ok := VerifySign(publicKey, msgHash, signature)
	fmt.Println("ok==", ok)
}
