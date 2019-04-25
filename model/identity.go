package model

import (
	kvdb "github.com/vinely/kvdb"
)

var (
	didDB *kvdb.BoltDB
)

func init() {
	d, _ := kvdb.NewKVDataBase("bolt://did.db/did?count=50")
	didDB = d.(*kvdb.BoltDB)
}

// Identity - interface for any identity and document convert to/from value
type Identity interface {
	Value
}

// GetIdentity - get publickey from database
func GetIdentity(i Identity, id string) error {
	return get(keyDB, i, id)
}

// SetIdentity - set publickey from database
func SetIdentity(id string, i Identity) error {
	return set(keyDB, id, i)
}
