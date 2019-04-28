package model

import (
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

// Claims - For a type to be a Claims object
type Claims interface {
	Value
	ID() string
	Get(id string) error
	Set() error
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
	id := sha3.Sum224(data)
	return string(id[:])
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

// ClaimStatus - status of claim. maybe changed a lot
type ClaimStatus struct {
	Owner     string `json:"owner"`
	Status    bool   `json:"st"` // status of claims. valid or other status
	IssuedAt  int64  `json:"iat,omitempty"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
}

// ID - get id from status
func (s *ClaimStatus) ID() string {
	// return s.Owner + "#st"
	// now using header id. don't need fragment
	return s.Owner
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
	return get(claimStatusDB, s, id)
}

// Set - set status to database
func (s *ClaimStatus) Set() error {
	return set(claimStatusDB, s.ID(), s)
}

// ClaimContent - Content of Claim
type ClaimContent struct {
	Owner    string                 `json:"owner"`
	Scope    map[string]interface{} `json:"scp,omitempty"`
	Contents map[string]interface{} `json:"cnt,omitempty"`
}

// ID - get id from content
func (c *ClaimContent) ID() string {
	// return c.Owner + "#cnt"
	// now using header id. don't need fragment
	return c.Owner
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

// StandardClaim - Basic elements for claim
type StandardClaim struct {
	Header  *ClaimHeader  `json:"header"`
	Status  *ClaimStatus  `json:"status"`
	Content *ClaimContent `json:"content"`
}

// ID - get id from content
func (s *StandardClaim) ID() string {
	return s.Status.ID()
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
func (s *StandardClaim) Get(id string) error {
	header := &ClaimHeader{}
	if err := header.Get(id); err != nil {
		return err
	}
	status := &ClaimStatus{}
	if err := status.Get(id); err != nil {
		return err
	}
	content := &ClaimContent{}
	if err := content.Get(id); err != nil {
		return err
	}
	s.Header = header
	s.Status = status
	s.Content = content
	return nil
}

// Set - set status to database
func (s *StandardClaim) Set() error {
	if err := s.Header.Set(); err != nil {
		return err
	}
	if err := s.Status.Set(); err != nil {
		return err
	}
	if err := s.Content.Set(); err != nil {
		return err
	}
	return nil
}
