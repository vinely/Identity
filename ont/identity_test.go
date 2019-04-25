package ont

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func printoutdoc(doc *Identity) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	d, err := json.Marshal(doc)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("doc:%s\n", d)
	fmt.Println("PublicKey:[")
	for _, v := range doc.PublicKey {
		fmt.Printf("%v,\n", v)
	}
	fmt.Println("]")
	return nil
}

func testnewdoc() (*Identity, error) {
	doc, id, err := NewIdentity()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = printoutdoc(doc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("id:%v\n", id)
	return doc, nil
}

func TestNew(t *testing.T) {
	testnewdoc()

}

func TestIdentitySetAndGetKeyFromDB(t *testing.T) {
	doc, err := testnewdoc()
	if err != nil {
		fmt.Println(err)
	}
	err = doc.SaveToDB()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("doc ID is:" + doc.ID)
	fmt.Println("Output")
	out, err := GetIdentityFromDB(string(doc.ID))
	if err != nil {
		fmt.Println(err)
	}
	printoutdoc(out)

}
