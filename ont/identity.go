package ont

import (
	"crypto/rand"
	"encoding/base64"

	jsoniter "github.com/json-iterator/go"
	"github.com/vinely/Identity/model"
	"github.com/vinely/dids"
	chain "github.com/vinely/ontchain"
)

// Identity - my type of Identity
type Identity struct {
	dids.DIDNode
	PublicKey []PublicKey `json:"publicKey,omitempty"`
}

func createPassword() string {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return ""
	}
	passwd := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
	base64.StdEncoding.Encode(passwd, buf[:])
	return string(buf[:])
}

// NewIdentity - create an Identity . also return an identity for new
func NewIdentity() (*Identity, *chain.ManagedIdentity, error) {
	id, err := dids.CreateID(OntMethod)
	if err != nil {
		return nil, nil, err
	}

	return IdentityFromID(id.ID())
}

// IdentityFromID - create Identity from ID(DIDs)
func IdentityFromID(id *dids.ID) (*Identity, *chain.ManagedIdentity, error) {
	passwd := createPassword()

	identity, err := chain.GetIdentityFromID(string(*id), []byte(passwd))
	if err != nil {
		return nil, nil, err
	}
	i, err := IdentityFromOntid(identity)
	if err != nil {
		return nil, nil, err
	}
	return i, identity, nil
}

// IdentityFromOntid - create Identity from managed identity
func IdentityFromOntid(id *chain.ManagedIdentity) (*Identity, error) {
	pk, err := getPublicKeyFromSDK(&id.Identity, 1)
	if err != nil {
		return nil, err
	}
	i := &Identity{}
	i.ID = dids.ID(id.ID)
	i.PublicKey = []PublicKey{*pk}
	return i, nil
}

// for model kvdb

// Value - interface for any  convert to/from value
// type Value interface {
// Marshal() ([]byte, error)
// Unmarshal(data []byte) error // this method need Unmarshal data to interface.
// Value() interface{}          // return struct value for this point
// }

// Marshal - convert  json publickey to string ([]byte)
func (i *Identity) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*i)
}

// Unmarshal - convert  json string ([]byte) to publickey
func (i *Identity) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	return nil
}

// Value - interface implement
func (i *Identity) Value() interface{} {
	return *i
}

// GetIdentityFromDB - get publickey from db
func GetIdentityFromDB(id string) (*Identity, error) {
	i := &Identity{}
	err := model.GetIdentity(i, id)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// SaveToDB - Save to database
func (i *Identity) SaveToDB() error {
	err := model.SetIdentity(string(i.ID), i)
	if err != nil {
		return err
	}
	for _, pk := range i.PublicKey {
		err := pk.SaveToDB()
		if err != nil {
			return err
		}
	}
	return nil
}
