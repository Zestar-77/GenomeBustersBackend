package genedatabase

// GetGeneLabel takes the codon sequence
func GetGeneLabel(sequence []byte) string {
	v, err := db.Get(sequence, nil)
	if err != nil {
		return ""
	}
	return string(v)
}
