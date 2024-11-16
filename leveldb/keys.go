package leveldb

import "github.com/ethereum/go-ethereum/log"

type Keys struct {
	db *LevelDB
}

func NewKeys(path string) (*Keys, error) {
	db, err := NewLevelStore(path)
	if err != nil {
		log.Error("Error opening keys database", "path", path, "err", err)
		return nil, err
	}
	return &Keys{db: db}, nil
}

func (k *Keys) GetPrivateKey(publicKey string) (string, bool) {
	key := []byte(publicKey)
	data, err := k.db.Get(key)
	if err != nil {
		return "0x00", false
	}
	bStr := toHexString(data)
	return bStr, true
}

func (k *Keys) StorePrivateKeys(keyPairList []KeyPair) bool {
	for _, pair := range keyPairList {
		privateKey := toByte(pair.PrivateKey)
		pubKey := []byte(pair.CompressPubkey)
		err := k.db.Put(pubKey, privateKey)
		if err != nil {
			log.Error("Error storing private key", "err", err, "public key", pubKey, "private key", privateKey)
			return false
		}
	}
	return true
}
