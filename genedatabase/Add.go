package genedatabase

import (
	"GenomeBustersBackend/global"
	"GenomeBustersBackend/specialFileReaders"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// AddGenBank adds the genbank file
func AddGenBank(path string) error {
	global.Log.Printf("Adding genbank files genes from file: %s\n", path)
	file, err := os.Open(path)
	if err != nil {
		global.Log.Printf("Unable to add genes from genbank: %v\n", err)
		return err
	}

	gbfile, err := specialFileReaders.NewGenebankFile(file)
	if err != nil {
		global.Log.Printf("Unable to add genes from genbank: %v\n", err)
		return err
	}

	global.Log.Printf("Adding %d genes from %s\n", gbfile.GetGenes().Len(), path)
	batch := new(leveldb.Batch)
	counter := 0
	for e := gbfile.GetGenes().Front(); e != nil; e = e.Next() {
		counter++
		global.Log.Printf("Count %d, Label: %s, sequence: %v", counter, e.Value.(specialFileReaders.Gene).Name, e.Value.(specialFileReaders.Gene).Sequence)
		gene := e.Value.(specialFileReaders.Gene)
		batch.Put([]byte(gene.Sequence), []byte(gene.Name))
	}

	err = db.Write(batch, &opt.WriteOptions{NoWriteMerge: false, Sync: true})
	if err != nil {
		panic(err)
	}

	return nil
}

// AddGene adds the gene to the database
// With the key `sequence` and the value `label`
func AddGene(label string, sequence []byte) {
	db.Put(sequence, []byte(label), nil)
}
