package genedatabase

import "github.com/syndtr/goleveldb/leveldb"

var db *leveldb.DB

// InitializeDatabase initializes the database structure and returns the call to close the database
func InitializeDatabase() (func() error, error) {
	var err error
	db, err = leveldb.OpenFile("busted.db", nil)
	if err != nil {
		return nil, err
	}

	return db.Close, nil
}
