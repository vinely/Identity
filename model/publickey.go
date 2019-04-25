package model

import (
	kvdb "github.com/vinely/kvdb"
)

var (
	keyDB *kvdb.BoltDB
)

func init() {
	d, _ := kvdb.NewKVDataBase("bolt://key.db/key?count=50")
	keyDB = d.(*kvdb.BoltDB)
}

// PublicKey - interface for any publickey convert to/from value
type PublicKey interface {
	Value
}

// GetPublicKey - get publickey from database
func GetPublicKey(pk PublicKey, id string) error {
	return get(keyDB, pk, id)
}

// SetPublicKey - set publickey from database
func SetPublicKey(id string, pk PublicKey) error {
	return set(keyDB, id, pk)
}
