package analyzer

import "fmt"

/**
local:= array
si= start index
int le= array length for loop back
*/

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

func codonToAmino(local []rune, si int, le int) rune {
	var ret rune
	st := si % le
	su := (si + 1) % le
	sv := (si + 2) % le
	var codon string
	codon = string(local[st]) + string(local[su]) + string(local[sv])
	ret = codonMap[codon]
	return ret
}

type gene struct {
	start    int
	end      int
	identity []rune
}

func thing(genome []rune, arrayLength int) []gene {
	genome2 := make([]rune, arrayLength)
	for i := 0; i < arrayLength; i++ {
		switch genome[i] {
		case 'T':
			genome2[i] = 'A'
			break
		case 'G':
			genome2[i] = 'C'
			break
		case 'A':
			genome2[i] = 'T'
			break
		case 'C':
			genome2[i] = 'G'
			break
		}
	}
	gen1 := count(genome, arrayLength)
	gen2 := count(genome2, arrayLength)
	return append(gen1, gen2...)

}

func count(runeArray []rune, arrayLength int) []gene {
	//srand(time(nullptr))

	var geneStore []gene
	genePosition := 0
	inphase := false
	temp := '0'
	temp2 := '0'
	current := gene{0, 0, nil}
	//3 is codon length this does not change, 1 and 2 are checking the entirety of the codon
	for i := 0; i < arrayLength+3 || inphase; {
		temp = runeArray[i%arrayLength+1]
		temp2 = runeArray[i%arrayLength+2]
		if inphase {
			if runeArray[i%arrayLength] == 'T' && ((temp == 'A' && (temp2 == 'A' || temp2 == 'G')) || (temp == 'G' && temp2 == 'A')) {
				inphase = false
				current.end = i
				geneStore[genePosition] = current
				i = current.start + 1
				fmt.Println(current.start, " ", current.end)
				current = gene{0, 0, nil}
			} else {
				current.identity = append(current.identity, codonToAmino(runeArray, genePosition, arrayLength))
				i += 3
			}
		} else {
			if runeArray[i%arrayLength] == 'A' && temp == 'T' && temp2 == 'G' {
				inphase = true
				current.start = i
			} else {
				i++
			}
		}
	}
	return geneStore
}
