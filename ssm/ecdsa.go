package ssm

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateECDSAKeyPair() (string, string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Error("generate key fail", "err", err)
		return "0x00", "0x00", "0x00", err
	}
	privateKeyString := hex.EncodeToString(crypto.FromECDSA(privateKey))
	pubKeyString := hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey))
	compressedAddress := hex.EncodeToString(crypto.CompressPubkey(&privateKey.PublicKey))
	return privateKeyString, pubKeyString, compressedAddress, nil
}

func SignMessage(privateKey, txMsgHex string) (string, error) {
	privateKeyEcdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error("hex string to ecdsa key fail", "err", err)
		return "0x00", err
	}
	txMsg, _ := hex.DecodeString(txMsgHex)
	sig, err := crypto.Sign(txMsg, privateKeyEcdsa)
	if err != nil {
		log.Error("sign transaction fail", "err", err)
		return "0x00", err
	}
	return hex.EncodeToString(sig), nil
}

func VerifySign(pubkey, msgHash, sig string) bool {
	publicKey, _ := hex.DecodeString(pubkey)
	MessageHash, _ := hex.DecodeString(msgHash)
	signature, _ := hex.DecodeString(sig)
	fmt.Println("MessageHash len:", len(MessageHash))
	fmt.Println("signature len:", len(signature))
	return crypto.VerifySignature(publicKey, MessageHash, signature)
}

//
//func PrivateKeyToStringWithRecover(privateKey *ecdsa.PrivateKey) (string, error) {
//	// 转换为字符串
//	privateKeyBytes := crypto.FromECDSA(privateKey)
//	privateKeyStr := hex.EncodeToString(privateKeyBytes)
//
//	// 验证可以正确恢复
//	recoveredKey, err := crypto.HexToECDSA(privateKeyStr)
//	if err != nil {
//		return "", fmt.Errorf("key conversion verification failed: %w", err)
//	}
//
//	// 验证私钥是否相同
//	if privateKey.D.Cmp(recoveredKey.D) != 0 {
//		return "", fmt.Errorf("key recovery check failed")
//	}
//
//	return privateKeyStr, nil
//}

//func main() {
//	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//	PrivateKeyToStringWithRecover(privateKey)
//
//}
