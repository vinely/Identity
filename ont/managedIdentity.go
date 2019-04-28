package ont

import (
	kvdb "github.com/vinely/kvdb"
	chain "github.com/vinely/ontchain"
)

var (
	// ManagedDB - store secret for client. so it is a high rank private database
	ManagedDB *kvdb.BoltDB
)

func init() {
	d, _ := kvdb.NewKVDataBase("bolt://admin.db/managedidentity?count=50")
	ManagedDB = d.(*kvdb.BoltDB)
}

// SaveManagedIdentity - save to db
func SaveManagedIdentity(id *chain.ManagedIdentity) error {
	ret := ManagedDB.SetData(id.ID, *id)
	if ret.Result {
		return ret
	}
	return nil
}
