package rpc

import (
	"context"
	"errors"
	"github.com/L-DENG/wallet-sign-go/leveldb"
	"github.com/L-DENG/wallet-sign-go/protobuf/wallet"
	"github.com/L-DENG/wallet-sign-go/ssm"
	"github.com/ethereum/go-ethereum/log"
	"strconv"
)

const MaxGenerateAddressNumber = 1000 * 10

func (s *RpcServer) GetSupportSignWay(ctx context.Context, in *wallet.SupportSignWayRequest) (*wallet.SupportSignWayResponse, error) {
	if in.Type == "ecdsa" || in.Type == "ed25519" {
		return &wallet.SupportSignWayResponse{
			Code:    strconv.Itoa(1),
			Msg:     "Support this sign way",
			Support: true,
		}, nil
	}
	return &wallet.SupportSignWayResponse{
		Code:    strconv.Itoa(0),
		Msg:     "Do not support this sign way",
		Support: false,
	}, nil
}

func (s *RpcServer) ExportPublicList(ctx context.Context, in *wallet.ExportPublicKeyRequest) (*wallet.ExportPublicKeyResponse, error) {
	if in.Number > MaxGenerateAddressNumber {
		return &wallet.ExportPublicKeyResponse{
			Code: strconv.Itoa(1),
			Msg:  "Number must be less than 10000",
		}, nil
	}
	var keyPairList []leveldb.KeyPair
	var pubKeyList []*wallet.PublicKey
	for count := 0; count < int(in.Number); count++ {
		priKeyStr, pubKeyStr, decPubKeyStr, err := ssm.CreateECDSAKeyPair()
		if err != nil {
			log.Error("create key pair failed", "err", err)
			return nil, err
		}
		KeyPair := leveldb.KeyPair{
			CompressPubkey: pubKeyStr,
			PrivateKey:     priKeyStr,
		}
		pukItem := &wallet.PublicKey{
			CompressPubkey:   pubKeyStr,
			DecompressPubkey: decPubKeyStr,
		}
		pubKeyList = append(pubKeyList, pukItem)
		keyPairList = append(keyPairList, KeyPair)
	}
	isOk := s.db.StorePrivateKeys(keyPairList)
	if !isOk {
		log.Error("store key pairs failed", "err", isOk)
		return nil, errors.New("store key pairs failed")
	}
	return &wallet.ExportPublicKeyResponse{
		Code:      strconv.Itoa(1),
		Msg:       "create keys success",
		PublicKey: pubKeyList,
	}, nil
}

func (s *RpcServer) SignTxMessage(ctx context.Context, in *wallet.SignTxMessageRequest) (*wallet.SignTxMessageResponse, error) {
	privKey, isOk := s.db.GetPrivateKey(in.PublicKey)
	if !isOk {
		return nil, errors.New("get private key failed")
	}
	signature, err := ssm.SignMessage(privKey, in.TxMessage)
	if err != nil {
		return nil, err
	}
	return &wallet.SignTxMessageResponse{
		Code:      strconv.Itoa(1),
		Msg:       "sign tx message success",
		Signature: signature,
	}, nil
}
