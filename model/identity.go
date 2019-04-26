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

// // ListIdentity - list identity
// too complicated
// func ListIdentity(page uint, check func(k, v []byte) (Identity, error)) ([]Identity, error) {
// 	data := didDB.List(uint(page), func(k, v []byte) *kvdb.KVResult {
// 		d, err := check(k, v)
// 		if err != nil {
// 			return &kvdb.KVResult{
// 				Result: false,
// 				Info:   err.Error(),
// 			}
// 		}
// 		return &kvdb.KVResult{
// 			Data:   d,
// 			Result: true,
// 			Info:   "",
// 		}
// 	})
// 	if data.Result {
// 		return data.Data.([]Identity), nil
// 	}
// 	return nil, data
// }

// IdentityDataBase - kv database of identity
func IdentityDataBase() kvdb.KVMethods {
	return keyDB
}
