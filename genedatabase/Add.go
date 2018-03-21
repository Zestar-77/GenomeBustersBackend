package genedatabase

import (
	"GenomeBustersBackend/specialFileReaders"
	"os"
)

func addGenBank(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	gbfile, err := specialFileReaders.NewGenebankFile(file)
	if err != nil {
		return err
	}

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
