package analyzer

import "strconv"

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
	"AUG": 'M',
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

type gene struct {
	start    int
	end      int
	label    string
	identity []rune
}

func getPermutations(genome []rune) (invert []rune, reverse []rune, inverse []rune) {
	invert = make([]rune, len(genome))
	reverse = make([]rune, len(genome))
	inverse = make([]rune, len(genome))
	for i, v := range genome {
		switch v {
		case 'T':
			invert[i] = 'A'
			inverse[len(genome)-1-i] = 'A'
			break
		case 'G':
			invert[i] = 'C'
			inverse[len(genome)-1-i] = 'C'
			break
		case 'A':
			invert[i] = 'T'
			inverse[len(genome)-1-i] = 'T'
			break
		case 'C':
			invert[i] = 'G'
			inverse[len(genome)-1-i] = 'G'
			break
		}
		reverse[len(genome)-1-i] = v
	}
	return
}

func thing(genome []rune) []gene {
	gen1 := make(chan []gene)
	gen2 := make(chan []gene)
	gen3 := make(chan []gene)
	gen4 := make(chan []gene)
	UnknownCounter := &concurrentCounter{}
	go count(genome, gen1, UnknownCounter)
	invert, reverse, inverse := getPermutations(genome)
	go count(invert, gen2, UnknownCounter)
	go count(reverse, gen3, UnknownCounter)
	go count(inverse, gen4, UnknownCounter)
	return append(append(append(<-gen1, <-gen2...), <-gen3...), <-gen4...)
}

func count(runeArray []rune, genes chan []gene, UnknownCounter *concurrentCounter) {
	geneStore := make([]gene, 0)
	inphase := false
	temp := '0'
	temp2 := '0'

	unk := UnknownCounter.addAndGetCount()
	current := gene{-1, -1, "unat" + strconv.Itoa(unk), nil}

	//3 is codon length this does not change, 1 and 2 are checking the entirety of the codon
	for i := 0; i < len(runeArray) || inphase; {
		temp = runeArray[(i+1)%len(runeArray)]
		temp2 = runeArray[(i+2)%len(runeArray)]
		if inphase && current.start%len(runeArray) != i%len(runeArray) {
			if runeArray[i%len(runeArray)] == 'T' && ((temp == 'A' && (temp2 == 'A' || temp2 == 'G')) || (temp == 'G' && temp2 == 'A')) {
				inphase = false
				current.end = i % len(runeArray)
				i = current.start + 1
				current.start = current.start % len(runeArray)
				geneStore = append(geneStore, current)

				//fmt.Println(current.start, " ", current.end)
				// TODO Get actual gene label
				unk := UnknownCounter.addAndGetCount()
				current = gene{-1, -1, "unat" + strconv.Itoa(unk), nil}
			} else {
				current.identity = append(current.identity, codonToAmino(runeArray, i))
				i += 3
			}
		} else if i%len(runeArray) == current.start%len(runeArray) {
			genes <- nil
			panic("Never ending gene")

		} else {
			if runeArray[i%len(runeArray)] == 'A' && temp == 'T' && temp2 == 'G' {
				inphase = true
				current.start = i
				current.identity = append(current.identity, codonToAmino(runeArray, i))
				i += 3
			} else {
				i++
			}
		}
	}
	genes <- geneStore

}
