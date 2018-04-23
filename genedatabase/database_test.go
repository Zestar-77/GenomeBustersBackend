package genedatabase

import (
	"GenomeBustersBackend/global"
	"log"
	"os"
	"testing"
)

// TestMain runs prior to the entire test suit, and is used for setup/tear down of the environment
func TestMain(m *testing.M) {
	file, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, 0666) // Ignore log messages by sending them to platform agnostic os.DevNull
	if err != nil {
		panic(err)
	}
	global.Log = log.New(file, "TESTING: ", log.Lmicroseconds)

	exitCode := m.Run()
	file.Close()
	os.Exit(exitCode)
}

func WipeDatabase(T *testing.T) {
	if err := os.RemoveAll("busted.db"); err != nil {
		T.Fatalf("Unable to delete database: %v", err)
	}
}

func TestDBAddAndRetrive(T *testing.T) {
	closer, err := InitializeDatabase()
	if err != nil {
		T.Error(err)
		return
	}
	defer WipeDatabase(T)
	defer closer()
	AddGenBank("sequence.gbtest")
	label := GetGeneLabel([]byte{77, 70, 70, 83, 73, 72, 82, 71, 89, 82, 81, 72, 70, 72, 73, 84, 78})
	if label != "TEST" {
		T.Error("Either unable to add or unable to retrive from db")
	}
}

func TestDBPersistance(T *testing.T) {
	closer, err := InitializeDatabase()
	if err != nil {
		T.Error(err)
		return
	}
	defer WipeDatabase(T)
	AddGenBank("sequence.gbtest")
	closer()
	closer, err = InitializeDatabase()
	if err != nil {
		T.Error(err)
		return
	}
	defer closer()
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	pass := false
	for iter.Next() {
		v := iter.Value()
		if string(v) == "TEST" {
			pass = true
		}
	}
	if !pass {
		T.Error("Unable to find TEST gene")
	}
}

func TestMultipleGenes(T *testing.T) {
	closer, err := InitializeDatabase()
	if err != nil {
		T.Error(err)
		return
	}
	defer WipeDatabase(T)
	defer closer()
	AddGenBank("sequence.1.gbtest")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	pass := 0
	for iter.Next() {
		v := iter.Value()
		if string(v) == "TEST1" {
			pass++
		} else if string(v) == "TEST2" {
			pass++
		}
	}
	if pass != 2 {
		T.Error("Unable to find all entered genes")
	}
}
