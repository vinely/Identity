package model

import (
	"encoding/base64"

	jsoniter "github.com/json-iterator/go"
	kvdb "github.com/vinely/kvdb"
	"golang.org/x/crypto/sha3"
)

var (
	claimHeaderDB  *kvdb.BoltDB
	claimStatusDB  *kvdb.BoltDB
	claimContentDB *kvdb.BoltDB
)

func init() {
	d, _ := kvdb.NewKVDataBase("bolt://claimheader.db/header?count=50")
	claimHeaderDB = d.(*kvdb.BoltDB)
	d, _ = kvdb.NewKVDataBase("bolt://claimstatus.db/status?count=50")
	claimStatusDB = d.(*kvdb.BoltDB)
	d, _ = kvdb.NewKVDataBase("bolt://claimcontent.db/content?count=50")
	claimContentDB = d.(*kvdb.BoltDB)
}

// ClaimDB - get claim status db
func ClaimDB() *kvdb.BoltDB {
	return claimStatusDB
}

// Claims - For a type to be a Claims object
type Claims interface {
	Value
	ID() string
	Get(id string) error
	Set() error
	Hdr() (*ClaimHeader, error)
	Cont() (*ClaimContent, error)
}

// GetClaims - get claims from database
func GetClaims(c Claims, id string) error {
	return get(claimStatusDB, c, id)
}

// SetClaims - set claims to database
func SetClaims(c Claims, id string) error {
	return set(claimStatusDB, id, c)
}

// ClaimHeader - header of claim.
// content of header identify the one claim to the others
// so ID of claim is symbol of this content
type ClaimHeader struct {
	Type      string `json:"typ"`
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	Arguments string `json:"arg,omitempty"`
	Algorithm string `json:"alg,omitempty"`
}

// ID - get id from header
func (hd *ClaimHeader) ID() string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(hd)
	if err != nil {
		return ""
	}
	// now simply using sha3.
	// TODO change hmac method with algorith
	id := sha3.Sum224(data)
	return base64.URLEncoding.EncodeToString(id[:])
}

// Marshal - convert  json claimheader to string ([]byte)
func (hd *ClaimHeader) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*hd)
}

// Unmarshal - convert  json string ([]byte) to claimheader
func (hd *ClaimHeader) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, hd)
	if err != nil {
		return err
	}
	return nil
}

// Get - get claimheader from database
func (hd *ClaimHeader) Get(id string) error {
	return get(claimHeaderDB, hd, id)
}

// Set - set claimheader from database
func (hd *ClaimHeader) Set() error {
	return set(claimHeaderDB, hd.ID(), hd)
}

// FindHeader - find some claim header
func FindHeader(page uint, check func(k, v []byte) (*ClaimHeader, error)) ([]*ClaimHeader, error) {
	db := claimHeaderDB
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
		cs := make([]*ClaimHeader, len(s))
		for k, v := range s {
			cs[k] = v.(*ClaimHeader)
		}
		return cs, nil
	}
	return nil, data
}

// ClaimContent - Content of Claim
type ClaimContent struct {
	Scope    map[string]interface{} `json:"scp,omitempty"`
	Contents map[string]interface{} `json:"cnt,omitempty"`
}

// ID - get id from content
func (c *ClaimContent) ID() string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	// now simply using sha3.
	// TODO change hmac method with algorith
	id := sha3.Sum224(data)
	return base64.URLEncoding.EncodeToString(id[:])
}

// Marshal - convert  json content to string ([]byte)
func (c *ClaimContent) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*c)
}

// Unmarshal - convert  json string ([]byte) to content
func (c *ClaimContent) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

// Get - get content from database
func (c *ClaimContent) Get(id string) error {
	return get(claimContentDB, c, id)
}

// Set - set content to database
func (c *ClaimContent) Set() error {
	return set(claimContentDB, c.ID(), c)
}

// FindContent - find some claim content
func FindContent(page uint, check func(k, v []byte) (*ClaimContent, error)) ([]*ClaimContent, error) {
	db := claimContentDB
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
		cs := make([]*ClaimContent, len(s))
		for k, v := range s {
			cs[k] = v.(*ClaimContent)
		}
		return cs, nil
	}
	return nil, data
}

// StandardClaim - Basic elements for claim
type StandardClaim struct {
	Header  *ClaimHeader  `json:"header"`
	Content *ClaimContent `json:"content"`
	Claim   Claims        `json:"status"`
}

// Marshal - convert  json status to string ([]byte)
func (s *StandardClaim) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(*s)
}

// Unmarshal - convert  json string ([]byte) to status
func (s *StandardClaim) Unmarshal(data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}

// Get - get status from database
func (s *StandardClaim) Get(c Claims, id string) error {
	if err := GetClaims(c, id); err != nil {
		return err
	}
	header, err := c.Hdr()
	if err != nil {
		return err
	}
	content, err := c.Cont()
	if err != nil {
		return err
	}
	s.Claim = c
	s.Header = header
	s.Content = content
	return nil
}

// Set - set status to database
func (s *StandardClaim) Set() error {
	if err := s.Header.Set(); err != nil {
		return err
	}
	if err := s.Claim.Set(); err != nil {
		return err
	}
	if err := s.Content.Set(); err != nil {
		return err
	}
	return nil
}


// NewHeader - Create a claim header
func NewHeader(t, subject, issuer string) (*ClaimHeader, string, error) {
	ch := &ClaimHeader{
		Type:    t,
		Subject: subject,
		Issuer:  issuer,
	}
	return ch, ch.ID(), ch.Set()
}

// NewContent - Create a claim content
func NewContent(scope, content map[string]interface{}) (*ClaimContent, string, error) {
	cc := &ClaimContent{
		Scope:    scope,
		Contents: content,
	}
	return cc, cc.ID(), cc.Set()
}