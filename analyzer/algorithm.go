package analyzer

import (
	"math"
	"strconv"
)

//import "fmt"

var codonMap = map[string]rune{
	"TTT": 'F',
	"TTC": 'F',
	"TTA": 'L',
	"TTG": 'L',
	"CTT": 'L',
	"CTC": 'L',
	"CTA": 'L',
	"CTG": 'L',
	"ATT": 'I',
	"ATC": 'I',
	"ATA": 'I',
	"ATG": 'M',
	"GTT": 'V',
	"GTC": 'V',
	"GTA": 'V',
	"GTG": 'V',

	"TCT": 'S',
	"TCC": 'S',
	"TCA": 'S',
	"TCG": 'S',
	"CCT": 'P',
	"CCC": 'P',
	"CCA": 'P',
	"CCG": 'P',
	"ACT": 'T',
	"ACC": 'T',
	"ACA": 'T',
	"ACG": 'T',
	"GCT": 'A',
	"GCC": 'A',
	"GCA": 'A',
	"GCG": 'A',

	"TAT": 'Y',
	"TAC": 'Y',
	"TAA": 0,
	"TAG": 0,
	"CAT": 'H',
	"CAC": 'H',
	"CAA": 'Q',
	"CAG": 'Q',
	"AAT": 'N',
	"AAC": 'N',
	"AAA": 'K',
	"AAG": 'K',
	"GAT": 'D',
	"GAC": 'D',
	"GAA": 'E',
	"GAG": 'E',

	"TGT": 'C',
	"TGC": 'C',
	"TGA": 0,
	"TGG": 'W',
	"CGT": 'R',
	"CGC": 'R',
	"CGA": 'R',
	"CGG": 'R',
	"AGT": 'S',
	"AGC": 'S',
	"AGA": 'R',
	"AGG": 'R',
	"GGT": 'G',
	"GGC": 'G',
	"GGA": 'G',
	"GGG": 'G',
}
var minLength =0
func codonToAmino(local []rune, si int) rune {
	var ret rune
	st := si % len(local)
	su := (si + 1) % len(local)
	sv := (si + 2) % len(local)
	var codon string
	codon = string(local[st]) + string(local[su]) + string(local[sv])
	ret = codonMap[codon]
	return ret
}

// Gene is a struct that holds information about a gene. Upmost Care must be taken when initializing this.
type Gene struct {
	UUID     int    `json:"id"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Label    string `json:"label"`
	Identity []rune `json:"sequence"`
}

// Genome is a list of genes, that is json serializable for what the frontend expects
type Genome struct {
	Genes          []Gene `json:"features"`
	GenesFound     int    `json:"features_found"`
	SequenceLength int    `json:"sequence_length"`
	Filename       string `json:"filename"`
}



// Thing analyzes the genome and returns found genes.
func Analyze(genome []rune) Genome {
	gen := make(chan []Gene)

	UnknownCounter := &concurrentCounter{}
	UUIDCounter := &concurrentCounter{}
	go count(genome, gen, UnknownCounter, UUIDCounter, 68)

	genes := <-gen
	return Genome{genes, len(genes), len(genome), ""}
}

func count(runeArray []rune, genes chan []Gene, UnknownCounter, UUIDCounter *concurrentCounter, minLength int) {
	geneStore := make([]Gene, 0)
	inphase := false
	temp := '0'
	temp2 := '0'
	sum := 0
	unk := UnknownCounter.count
	current := Gene{UUIDCounter.count, -1, -1, "unat" + strconv.Itoa(unk), nil}

	//3 is codon length this does not change, 1 and 2 are checking the entirety of the codon
	for i := 0; i < len(runeArray) || inphase; {
		temp = runeArray[(i+1)%len(runeArray)]
		temp2 = runeArray[(i+2)%len(runeArray)]
		if inphase && current.Start%len(runeArray) != i%len(runeArray) {
			if runeArray[i%len(runeArray)] == 'T' && ((temp == 'A' && (temp2 == 'A' || temp2 == 'G')) || (temp == 'G' && temp2 == 'A')) {
				inphase = false
				if i-current.Start>minLength {
					current.End = i % len(runeArray)
					sum = sum + int(math.Abs(float64(i-current.Start)))

					current.Start = current.Start % len(runeArray)
					// TODO Get actual gene label
					current.Label = "unat" + strconv.Itoa(UnknownCounter.addAndGetCount())
					current.UUID = UUIDCounter.addAndGetCount()
					geneStore = append(geneStore, current)
				} else {
					i = current.Start + 1
				}
				current = Gene{0, -1, -1, "", nil}
			} else {
				current.Identity = append(current.Identity, codonToAmino(runeArray, i))
				i += 3
			}
		} else if i%len(runeArray) == current.Start%len(runeArray) {
			genes <- nil
			panic("Never ending gene")

		} else {
			if runeArray[i%len(runeArray)] == 'A' && temp == 'T' && temp2 == 'G' {
				inphase = true
				current.Start = i
				current.Identity = append(current.Identity, codonToAmino(runeArray, i))
				i += 3
			} else {
				i++
				//looking for start codon
			}
		}
	}

	genes <- geneStore

}
