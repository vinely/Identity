package ont

import (
	"fmt"
	"testing"

	"github.com/vinely/dids"
)

func TestSetAndGetKeyFromDB(t *testing.T) {
	pk := &PublicKey{}
	pk.ID = "test"
	pk.Type = "x"
	pk.Index = 1
	pk.PublicKeyHex = dids.PublicKeyHex{Value: "082732"}
	err := pk.SaveToDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	out, err := GetPublicKeyFromDB("test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(" out is : %v\n", *out)
}
