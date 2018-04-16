package genedatabase

import (
	"GenomeBustersBackend/global"
)

// Dump dumps out the contents of the database
func Dump() {
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		global.Log.Printf("Key: %v, Value: %v", k, v)
	}
}
