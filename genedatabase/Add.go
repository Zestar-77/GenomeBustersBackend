package genedatabase

import (
	"GenomeBustersBackend/global"
	"GenomeBustersBackend/specialFileReaders"
	"os"
)

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
	for e := gbfile.GetGenes().Front(); e != nil; e = e.Next() {
		gene := e.Value.(specialFileReaders.Gene)
		AddGene(gene.Name, gene.Sequence)
	}

	return nil
}

// AddGene adds the gene to the database
// With the key `sequence` and the value `label`
func AddGene(label string, sequence []byte) {
	db.Put(sequence, []byte(label), nil)
}
