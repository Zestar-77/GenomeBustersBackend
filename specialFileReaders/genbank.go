package specialFileReaders

import (
	"GenomeBustersBackend/global"
	"bufio"
	"bytes"
	"container/list"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var numbers *regexp.Regexp
var geneIndexRange *regexp.Regexp
var geneName *regexp.Regexp
var locusTag *regexp.Regexp

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// This init function will compile needed regex expressions
func init() {
	reg, err := regexp.Compile("[\\s\\d]+")
	panicIfError(err)
	numbers = reg
	reg, err = regexp.Compile(`(\d+)..(\d+)`)
	panicIfError(err)
	geneIndexRange = reg
	reg, err = regexp.Compile(`/gene="(\w+)"`)
	panicIfError(err)
	geneName = reg
	reg, err = regexp.Compile(`/locus_tag="(\w+)"`)
	panicIfError(err)
	locusTag = reg
}

// GenebankFile is a reader for reading genebank files. Currently, it does not do much more that ignore all
// Information prior to the actual genemap (should that exist. Otherwise, there is a problem)
type GenebankFile struct {
	genome string
	genes  *list.List
}

// Represents a gene
type Gene struct {
	Name       string
	start, end uint
	Sequence   []byte
}

// NewGenebankFile creates a new GenbankFile struct for reading a genbank file from an io.Reader
func NewGenebankFile(reader io.Reader) (*GenebankFile, error) {
	s := bufio.NewScanner(reader)
	file := GenebankFile{}
	reachedGenome := false
	file.genes = list.New()
	var gene []string
	for s.Scan() {
		if line := strings.TrimSpace(s.Text()); !reachedGenome && line != "ORIGIN" {
			if strings.HasPrefix(line, "gene") {
				if gene != nil {
					file.finalizeGene(gene)
				}
				gene = make([]string, 14)
			} else if len(gene) > 0 {

			}
			continue
		} else if !reachedGenome {
			reachedGenome = true
			continue
		} else if line == "//" {
			reachedGenome = false
			continue // End of file
		}
		file.genome += strings.ToUpper(strings.TrimSpace(numbers.ReplaceAllString(s.Text(), ""))) // Not actually sure if trim space is needed
	}
	file.finalizeGeneSequence()
	return &file, nil
}

func (gf *GenebankFile) finalizeGene(geneData []string) {
	r := geneIndexRange.FindStringSubmatch(geneData[0])
	if len(r) != 3 {
		log.Printf("Error, gene format invalid: %s", geneData[0])
		return
	}

	start, err := strconv.ParseUint(r[1], 10, 0)
	if err != nil {
		log.Printf("Error, gene format invalid: %s", geneData[0])
		return
	}
	end, err := strconv.ParseUint(r[2], 10, 0)
	if err != nil {
		log.Printf("Error, gene format invalid: %s", geneData[0])
		return
	}

	name := ""
	locus := ""
	for _, line := range geneData {
		n := geneName.FindStringSubmatch(line)
		if len(n) == 2 {
			name = n[1]
		}
		l := locusTag.FindStringSubmatch(line)
		if len(l) == 2 {
			locus = l[1]
		}
	}
	if name == "" {
		name = locus
	}
	if name == "" {
		return
	}

	g := Gene{Name: name, start: uint(start), end: uint(end)}
	gf.genes.PushBack(g)
}

// 410 -1750

func (gf *GenebankFile) finalizeGeneSequence() {
	for e := gf.genes.Front(); e != nil; e = e.Next() {
		g := e.Value.(Gene)
		if g.start > 0 && g.end > g.start && g.end < (uint)(len(gf.genome)) {
			g.Sequence = codonToAmino(gf.genome[g.start : g.end+1])
		} else {
			log.Printf("GB had an gene that went past the end of the genome.\n")

		}
	}
}

func codonToAmino(s string) []byte {
	buffer := bytes.NewBuffer(make([]byte, len(s)/3+1))
	for i := 0; i < len(s); i += 3 {
		buffer.WriteRune(global.CodonMap[s[i:i+3]])
	}
	return buffer.Bytes()
}

// ReadGenome Reads the part of the genebank file containing the genome
func (gf *GenebankFile) ReadGenome() string {
	return gf.genome
}

// GetGenes Returns a list of all genes in this genome,
// Already parsed to single letter codons.
func (gf *GenebankFile) GetGenes() *list.List {
	return gf.genes
}