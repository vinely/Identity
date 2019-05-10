package ont

import (
	"github.com/vinely/Identity/model"
)

const (
	// GrantClaimType -
	GrantClaimType = "grant"
)

var (
	// DefaultGrantContent - all scope allowed grant
	DefaultGrantContent *model.ClaimContent
	// DefaultGrantContentString -
	DefaultGrantContentString string
)

func init() {
	DefaultGrantContent, DefaultGrantContentString, _ = model.NewContent(model.M{"scope": "all"}, model.M{"type": "grant"})
	DefaultGrantContent.Set()
}

// DefaultGrantID  - client is id of client
func DefaultGrantID(client string) (*model.StandardClaim, error) {
	h, hd, _ := model.NewHeader(GrantClaimType, client, adminID)
	h.Set()
	return NewClaim(hd, DefaultGrantContentString)
}
