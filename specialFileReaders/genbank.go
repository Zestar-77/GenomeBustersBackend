package specialFileReaders

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var numbers *regexp.Regexp

func init() {
	if reg, err := regexp.Compile("[\\s\\d]+"); err == nil {
		numbers = reg
	} else {
		panic(err)
	}
}

// GenebankFile is a reader for reading genebank files. Currently, it does not do much more that ignore all
// Information prior to the actual genemap (should that exist. Otherwise, there is a problem)
type GenebankFile struct {
	genome string
}

// NewGenebankFile creates a new GenbankFile struct for reading a genbank file from an io.Reader
func NewGenebankFile(reader io.Reader) (*GenebankFile, error) {
	s := bufio.NewScanner(reader)
	file := GenebankFile{}
	reachedGenome := false
	for s.Scan() {
		if !reachedGenome && strings.TrimSpace(s.Text()) != "ORIGIN" {
			continue
		} else if !reachedGenome {
			reachedGenome = true
			continue
		} else if strings.TrimSpace(s.Text()) == "//" {
			reachedGenome = false
			continue // End of file
		}
		file.genome += strings.TrimSpace(numbers.ReplaceAllString(s.Text(), "")) // Not actually sure if trim space is needed
	}
	return &file, nil
}

// ReadGenome Reads the part of the genebank file containing the genome
func (g *GenebankFile) ReadGenome() string {
	return g.genome
}
