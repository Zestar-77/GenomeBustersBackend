package genedatabase

import (
	"GenomeBustersBackend/global"
)

// Dump dumps out the contents of the database
// This exists for debugging reasons and is not critical.
func Dump() {
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		global.Log.Printf("Key: %v, Value: %v", k, v)
	}
}
