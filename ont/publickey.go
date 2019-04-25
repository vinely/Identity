package ont

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vinely/Identity/model"
	"github.com/vinely/dids"
	sdk "github.com/vinely/ontchain/ontsdk"
)

const (
	// PublicKeyFragmentHeader - publickey label will be fragment of did scheme.
	// The header will be prefix of the fragment
	PublicKeyFragmentHeader = "Key_"
)

// PublicKey for ont
type PublicKey struct {
	dids.BasePublicKey
	Index int
	dids.PublicKeyHex
}

func getPublicKeyFromSDK(id *sdk.Identity, i int) (*PublicKey, error) {
	c, err := id.GetControllerDataByIndex(i)
	if err != nil {
		return nil, err
	}
	pk := &PublicKey{}
	pk.ID = dids.ID(id.ID + "#" + PublicKeyFragmentHeader + c.ID)
	pk.Controller = dids.ID(id.ID)
	pk.Type = ""
	pk.Index = i
	pk.PublicKeyHex = dids.PublicKeyHex{Value: dids.PublicKeyValue(c.Public)}
	return pk, nil
}

// Key - key value of public key
func (pk *PublicKey) Key() string {
	return string(pk.PublicKeyHex.Value)
}

// Valid - for interface PublicKey
func (pk *PublicKey) Valid() bool {
	// TODO not implement
	return true
}

// for model kvdb

// Value - interface for any  convert to/from value
// type Value interface {
// 	Marshal() ([]byte, error)
// 	Unmarshal([]byte) error
// }

// Marshal - convert  json publickey to string ([]byte)
func (pk *PublicKey) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*pk)
}

// Unmarshal - convert  json string ([]byte) to publickey
func (pk *PublicKey) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, pk)
	if err != nil {
		return err
	}
	return nil
}

// Value - interface implement
func (pk *PublicKey) Value() interface{} {
	return *pk
}

// GetPublicKeyFromDB - get publickey from db
func GetPublicKeyFromDB(id string) (*PublicKey, error) {
	pk := &PublicKey{}
	err := model.GetPublicKey(pk, id)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

// SaveToDB - Save to database
func (pk *PublicKey) SaveToDB() error {
	return model.SetPublicKey(string(pk.ID), pk)
}
