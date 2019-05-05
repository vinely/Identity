package ont

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	kvdb "github.com/vinely/kvdb"
	chain "github.com/vinely/ontchain"
)

var (
	// ManagedDB - store secret for client. so it is a high rank private database
	ManagedDB *kvdb.BoltDB

	adminID = "did:ont:TMhYqdkBGMYQ7qDyutNK2z7k345QUDjiDZ"

	// AdminIdentity -  Admin mangedidentity of this application
	AdminIdentity *chain.ManagedIdentity
	// Admin - admin identity for check
	Admin *Identity
)

func init() {
	d, _ := kvdb.NewKVDataBase("bolt://admin.db/managedidentity?count=50")
	ManagedDB = d.(*kvdb.BoltDB)
	AdminIdentity, err := AddRandomPasswordIdentity(adminID)
	if err != nil {
		panic(err.Error())
	}
	Admin, err = IdentityFromOntid(AdminIdentity)
	if err != nil {
		panic(err.Error())
	}
	Admin.SaveToDB()
}

// NewRandomPasswordIdentity - create an random password account
func NewRandomPasswordIdentity(id string) (*chain.ManagedIdentity, error) {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return nil, fmt.Errorf("generate ID error, %s", err)
	}
	passwd := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
	base64.StdEncoding.Encode(passwd, buf[:])
	return chain.GetIdentityFromID(id, passwd)
}

// SaveManagedIdentity - save to db
func SaveManagedIdentity(id *chain.ManagedIdentity) error {
	return saveToDB(id, ManagedDB)
}

// AddRandomPasswordIdentity - add admin identity if it wasn't existed
func AddRandomPasswordIdentity(id string) (*chain.ManagedIdentity, error) {
	return AddRandomPasswordAccount(id, ManagedDB, false)
}

func saveToDB(id *chain.ManagedIdentity, db kvdb.KVMethods) error {
	ret := db.SetData(id.ID, id)
	if ret.Result {
		return nil
	}
	return ret
}

// SaveIDToDB - Save managed Identity to DB
func SaveIDToDB(id *chain.ManagedIdentity, db kvdb.KVMethods, force bool) error {
	if db.Exists(id.ID) && !force {
		kr := db.Get(id.ID)
		if kr.Result {
			return errors.New("Already existed")
		}
	}
	return saveToDB(id, db)
}

// AddRandomPasswordAccount - Create Wallet to db
// id is the key
// force - true is write anyway. false is don't write if existed
func AddRandomPasswordAccount(id string, db kvdb.KVMethods, force bool) (*chain.ManagedIdentity, error) {
	if db.Exists(id) && !force {
		kr := db.Get(id)
		if kr.Result {
			fmt.Printf("%v\n", string(kr.Data.([]byte)))
			i := &chain.ManagedIdentity{}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err := json.Unmarshal(kr.Data.([]byte), i)
			if err != nil {

				return nil, err
			}
			return i, nil
		}
	}
	i, err := NewRandomPasswordIdentity(id)
	if err != nil {
		return nil, err
	}
	saveToDB(i, db)
	return i, nil
}
