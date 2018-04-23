package genedatabase

// GetGeneLabel takes the codon sequence and looks up in the database
// for a corresponding entry. If one exists, returns its label, else return nil
func GetGeneLabel(sequence []byte) string {
	v, err := db.Get(sequence, nil)
	if err != nil {
		return ""
	}
	return string(v)
}
