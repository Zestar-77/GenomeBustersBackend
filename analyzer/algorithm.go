package analyzer

import "fmt"
/**
local:= array
si= start index
int le= array length for loop back
 */
func codonToAmino(local []rune,  si int, le int) rune{
	var ret rune
	st := si%le
	su := (si+1)%le
	sv := (si+2)%le
	switch (local[st]){
	case 'T':
		switch (local[su]) {
		case 'T':
			switch (local[su]) {
			case 'T':
			case 'C':
				ret='F'
				break
			case 'A':
			case 'G':
				ret='L'
				break
			default:

			}
			break
		case 'C':
			ret='S'
			break
		case 'A':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='Y'
				break
			case 'A':
			case 'G':
				break

			}
			break
		case 'G':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='C'
				break
			case 'A':
				break
			case 'G':
				ret='W'
				break
			}
			break
		}
		break
	case 'C':
		switch (local[su]) {
		case 'T':
			ret='L'
			break
		case 'C':
			ret='P'
			break
		case 'A':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='H'
				break
			case 'A':
			case 'G':
				ret='Q'
				break

			}
			break
		case 'G':
			ret='R'
			break
		}
		break
	case 'A':
		switch (local[su]) {
		case 'T':
			switch (local[sv]) {
			case 'T':
			case 'C':
			case 'A':
				ret='I'
				break
			case 'G':
				ret='M'
				break
			}
			break
		case 'C':
			ret='T'
			break
		case 'A':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='N'
				break
			case 'A':
			case 'G':
				ret='K'
				break
			}
			break
		case 'G':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='S'
				break
			case 'A':
			case 'G':
				ret='R'
				break

			}
			break
		default:
		}
		break
	case 'G':
		switch (local[su]) {
		case 'T':
			ret='V'
			break
		case 'C':
			ret='A'
			break
		case 'A':
			switch (local[sv]) {
			case 'T':
			case 'C':
				ret='D'
				break
			case 'A':
			case 'G':
				ret='E'
				break

			}
			break
		case 'G':
			ret='G'
			break

		}
		break

	}
	return ret
}
type gene struct{
	start int
	end int
	identity []rune
}
func thing( genome []rune, arrayLength int) []gene{
	genome2 := make([]rune,arrayLength)
	for i:=0; i<arrayLength; i++{
		switch(genome[i]){
		case 'T':
			genome2[i]='A'
			break
		case 'G':
			genome2[i]='C'
			break
		case 'A':
			genome2[i]='T'
			break
		case 'C':
			genome2[i]='G'
			break
		}
	}
	gen1:= count(genome, arrayLength);
	gen2:= count(genome2, arrayLength);
	return append(gen1,gen2...);

}
func count(runeArray []rune, arrayLength int) []gene {
	//srand(time(nullptr))

	var geneStore []gene
	genePosition:=0
	inphase:=false
	temp:='0'
	temp2:='0'
	current:=gene{0,0,nil}
	//3 is codon length this does not change, 1 and 2 are checking the entirety of the codon
	for i:=0; i<arrayLength+3||inphase;{
		temp=runeArray[i%arrayLength+1]
		temp2=runeArray[i%arrayLength+2]
		if inphase {
			if runeArray[i%arrayLength]=='T'&&((temp=='A'&&(temp2=='A'||temp2=='G'))||(temp=='G'&&temp2=='A')){
				inphase = false
				current.end=i
				geneStore[genePosition]=current
				i=current.start+1
				fmt.Println(current.start," ", current.end)
				current=gene{0,0, nil}
			}else{
				append(current.identity, codonToAmino(runeArray,genePosition,arrayLength)
				i+=3
			}
		}else{
			if runeArray[i%arrayLength]=='A'&&temp=='T'&&temp2=='G' {
				inphase = true
				current.start=i
			}else{
				i++
			}
		}
	}
	return geneStore
}