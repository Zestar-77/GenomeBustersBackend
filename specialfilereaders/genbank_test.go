package specialfilereaders

import (
	"GenomeBustersBackend/global"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func TestRegex(T *testing.T) {
	r := geneIndexRange.FindStringSubmatch("gene 5..6")
	if len(r) != 3 {
		T.Error("gene index range regex is broken")
	} else {
		start, err1 := strconv.ParseUint(r[1], 10, 0)
		end, err2 := strconv.ParseUint(r[2], 10, 0)
		if err1 != nil || err2 != nil || start != 5 || end != 6 {
			T.Error("gene index range regex is broken")
		}
	}
	r = geneName.FindStringSubmatch(`/gene="TestGene"`)
	if len(r) != 2 || r[0] != `/gene="TestGene"` || r[1] != "TestGene" {
		T.Error("Gene name regex is broken")
	}
	r = locusTag.FindStringSubmatch(`/locus_tag="TestLocus"`)
	if len(r) != 2 || r[0] != `/locus_tag="TestLocus"` || r[1] != "TestLocus" {
		T.Error("Locus Tag regex is broken")
	}
	s := numbers.ReplaceAllString(`	156 TGBADTISNVWIOR`, "")
	if s != "TGBADTISNVWIOR" {
		T.Error("numbers regex is broken")
	}
}

func TestCodonToAmino(T *testing.T) {
	// FLIVSYW
	b := codonToAmino("TTTTTAATCGTGTCTTATTGG", false)
	if bytes.Compare(b, []byte("FLIVSYW")) != 0 {
		T.Fail()
	}
	// KNYHRIT
	// AAAAATTACCACAGAATAACC
	b = codonToAmino("TTTTTAATGGTGTCTTATTGG", true)
	if bytes.Compare(b, []byte("KNYHRIT")) != 0 {
		fmt.Printf("%v", string(b))
		T.Error(string(b))
	}
}

func TestReadingFile(T *testing.T) {
	file, err := os.Open("sequence.gbtest")
	if err != nil {
		T.FailNow()
	}
	defer file.Close()
	gf, err := NewGenebankFile(file)
	if err != nil {
		T.FailNow()
	}
	if gf.genes.Len() != 2 {
		T.Errorf("List of genes is incomplete: %v", gf.genes)
		return
	}
	if gf.genome != strings.ToUpper("atctttttcgatgttttttagtatccacagaggttatcgacaacattttcacattaccaacccctgttaacaaggttttttcaacaggttgtccgctttgtggataagattgtgacaaccattgcaagctctcgtttattttggtattatatttgtgttttaactcttgattactaatcctacctttcctctttatccacaaagtgtggataagttgtggattgatttcacacagcttgtgtagaaggttgtccacaagttgtgaaatttgtcgaaaagctatttatctactatattata") {
		T.Errorf("Did not actually read in genome, instead read:\n %s", gf.genome)
	}
}
