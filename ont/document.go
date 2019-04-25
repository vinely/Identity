package ont

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vinely/Identity/model"
)

// Document - my type of Document
type Document struct {
	Identity
}

// for model kvdb

// Value - interface for any  convert to/from value
// type Value interface {
// 	Marshal() ([]byte, error)
// 	Unmarshal([]byte) error
// }

// Marshal - convert  json publickey to string ([]byte)
func (doc *Document) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*doc)
}

// Unmarshal - convert  json string ([]byte) to publickey
func (doc *Document) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, doc)
	if err != nil {
		return err
	}
	return nil
}

// GetDocumentFromDB - get publickey from db
// document is identity with ddo
// document is one kind of identity
func GetDocumentFromDB(id string) (*Document, error) {
	doc := &Document{}
	err := model.GetIdentity(doc, id)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// SaveToDB - Save to database
// document is identity with ddo
// document is one kind of identity
func (doc *Document) SaveToDB() error {
	err := model.SetIdentity(string(doc.ID), doc)
	if err != nil {
		return err
	}
	for _, pk := range doc.PublicKey {
		err := pk.SaveToDB()
		if err != nil {
			return err
		}
	}
	return nil
}
