package model

import (
	kvdb "github.com/vinely/kvdb"
)

// Value - interface for any  convert to/from value
// Warning - Value interface need be a point implement because Unmarshal will return data into interface object
type Value interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error // this method need Unmarshal data to interface.
	Value() interface{}          // return struct value for this point
}

func get(db kvdb.KVMethods, v Value, id string) error {
	info := db.Get(id)
	if !info.Result {
		return info
	}
	return v.Unmarshal(info.Data.([]byte))
}

func set(db kvdb.KVMethods, id string, v Value) error {
	if info := db.SetData(id, v); !info.Result {
		return info
	}
	return nil
}

func del(db kvdb.KVMethods, id string) error {
	info := db.Delete(id)
	if !info.Result {
		return info
	}
	return nil
}

func list(db kvdb.KVMethods, page uint, check func(k, v []byte) (Value, error)) ([]Value, error) {
	data := didDB.List(uint(page), func(k, v []byte) *kvdb.KVResult {
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
		return data.Data.([]Value), nil
	}
	return nil, data
}
