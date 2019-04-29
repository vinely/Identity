package ont

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/vinely/Identity/model"
	kvdb "github.com/vinely/kvdb"
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
// Hdr() (*ClaimHeader, error)
// Cont() (*ClaimContent, error)
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
	return s.Header + seperator + s.Content
}

// Hdr - get id from header
func (s *ClaimStatus) Hdr() (*model.ClaimHeader, error) {
	h := &model.ClaimHeader{}
	err := h.Get(s.Header)
	return h, err
}

// Cont - get id from content
func (s *ClaimStatus) Cont() (*model.ClaimContent, error) {
	c := &model.ClaimContent{}
	err := c.Get(s.Content)
	return c, err
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

// ClaimList - list claim
func ClaimList(page uint, check func(k, v []byte) (*ClaimStatus, error)) ([]*ClaimStatus, error) {
	db := model.ClaimDB()
	data := db.List(uint(page), func(k, v []byte) *kvdb.KVResult {
		d, err := check(k, v)
		if err != nil {
			return &kvdb.KVResult{
				Result: false,
				Info:   err.Error(),
			}
		}
		return &kvdb.KVResult{
			Data:   d,
			Result: true,
			Info:   "",
		}
	})
	if data.Result {
		s := data.Data.([]interface{})
		cs := make([]*ClaimStatus, len(s))
		for k, v := range s {
			cs[k] = v.(*ClaimStatus)
		}
		return cs, nil
	}
	return nil, data

}

// NewHeader - Create a claim header
func NewHeader(t, subject, issuer string) (*model.ClaimHeader, string, error) {
	ch := &model.ClaimHeader{
		Type:    t,
		Subject: subject,
		Issuer:  issuer,
	}
	return ch, ch.ID(), ch.Set()
}

// NewContent - Create a claim content
func NewContent(scope, content map[string]interface{}) (*model.ClaimContent, string, error) {
	cc := &model.ClaimContent{
		Scope:    scope,
		Contents: content,
	}
	return cc, cc.ID(), cc.Set()
}

// NewClaim - create a claim
func NewClaim(header, content string) (*model.StandardClaim, error) {
	c := &ClaimStatus{}
	c.Header = header
	c.Content = content
	c.Status = true
	c.IssuedAt = time.Now().Unix()
	c.ExpiresAt = time.Now().Add(time.Hour).Unix()
	hd, err := c.Hdr()
	if err != nil {
		return nil, err
	}
	cnt, err := c.Cont()
	if err != nil {
		return nil, err
	}
	s := &model.StandardClaim{}
	s.Claim = c
	s.Header = hd
	s.Content = cnt
	return s, nil
}
