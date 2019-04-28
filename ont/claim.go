package ont

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vinely/Identity/model"
)

var (
	seperator = ":"
)

// Claims - For a type to be a Claims object
// type Claims interface {
// 	Value
// 	ID() string
// 	Get(id string) error
// 	Set() error
// 	Hdr() string
// 	Cont() string
// }

// ClaimStatus - status of claim. maybe changed a lot
// Header and content can be rarely changeable after created
// when claim changes just change status.
// so status is the main part of claim
type ClaimStatus struct {
	Header     string `json:"hd"`  // ID of related ClaimHeader
	Content    string `json:"cnt"` // ID of related ClaimContent
	Status     bool   `json:"st"`  // status of claims. valid or other status
	IssuedAt   int64  `json:"iat,omitempty"`
	Audience   string `json:"aud,omitempty"`
	ExpiresAt  int64  `json:"exp,omitempty"`
	NotBefore  int64  `json:"nbf,omitempty"`
	SignMethod string `json:"smd,omitempty"` // not used now
	Signature  string `json:"sig,omitempty"`
}

// Sign - make signature
func (s *ClaimStatus) Sign() error {
	return nil
}

// Verify - verify the signature
func (s *ClaimStatus) Verify() bool {
	return true
}

// ID - get id from status
func (s *ClaimStatus) ID() string {
	return s.Hdr() + seperator + s.Cont()
}

// Hdr - get id from header
func (s *ClaimStatus) Hdr() string {
	return s.Header
}

// Cont - get id from content
func (s *ClaimStatus) Cont() string {
	return s.Content
}

// Marshal - convert  json status to string ([]byte)
func (s *ClaimStatus) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*s)
}

// Unmarshal - convert  json string ([]byte) to status
func (s *ClaimStatus) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}

// Get - get status from database
func (s *ClaimStatus) Get(id string) error {
	return model.GetClaims(s, id)
}

// Set - set status to database
func (s *ClaimStatus) Set() error {
	return model.SetClaims(s, s.ID())
}


