package analyzer

import (
	"GenomeBustersBackend/genedatabase"
	"GenomeBustersBackend/global"
	"math"
	"strconv"
)

var minLength = 1

func codonToAmino(local []rune, si int) rune {
	var ret rune
	st := si % len(local)
	su := (si + 1) % len(local)
	sv := (si + 2) % len(local)
	var codon string
	codon = string(local[st]) + string(local[su]) + string(local[sv])
	ret = global.CodonMap[codon]
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

// Analyze analyzes the genome and returns found genes.
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
				if i-current.Start > minLength {
					current.End = i % len(runeArray)
					sum = sum + int(math.Abs(float64(i-current.Start)))

					current.Start = current.Start % len(runeArray)
					// TODO Get actual gene label
					label := genedatabase.GetGeneLabel([]byte(string(current.Identity)))
					if label == "" {
						global.Log.Println(label)
						label = "unat" + strconv.Itoa(UnknownCounter.addAndGetCount())
					}
					current.Label = label
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
