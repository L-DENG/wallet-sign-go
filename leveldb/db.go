package leveldb

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	*leveldb.DB
}

func NewLevelStore(Path string) (*LevelDB, error) {
	handle, err := leveldb.OpenFile(Path, nil)
	if err != nil {
		log.Error("open leveldb err", "err", err)
		return nil, err
	}
	return &LevelDB{handle}, nil
}

func (db *LevelDB) Put(key []byte, value []byte) error {
	return db.DB.Put(key, value, nil)
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	return db.DB.Get(key, nil)
}

func (db *LevelDB) Delete(key []byte) error {
	return db.DB.Delete(key, nil)
}

func toByte(hexDataString string) []byte {
	dataBytes, _ := hex.DecodeString(hexDataString)
	return dataBytes
}

func toHexString(dataBytes []byte) string {
	return hex.EncodeToString(dataBytes)
}
