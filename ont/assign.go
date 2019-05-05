package ont

import (
	"github.com/vinely/Identity/model"
)

const (
	// AssignClaimType -
	AssignClaimType = "asn"
)

// NewAssignedID -
func NewAssignedID() (*model.StandardClaim, error) {
	identity, mi, err := NewIdentity()
	if err != nil {
		return nil, err
	}
	err = SaveManagedIdentity(mi)
	if err != nil {
		return nil, err
	}
	err = identity.SaveToDB()
	if err != nil {
		return nil, err
	}
	return AssignID(identity.String())
}

// AssignID -
func AssignID(id string) (*model.StandardClaim, error) {
	h, hd, _ := model.NewHeader(AssignClaimType, id, adminID)
	h.Set()
	c, cc, _ := model.NewContent(nil, nil)
	c.Set()
	return NewClaim(hd, cc)
}
