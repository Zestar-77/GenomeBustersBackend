package analyzer

import (
	"testing"
	"strconv"
)

	func TestGetPermutations(T *testing.T){
	test := [4]rune{'T', 'G', 'A', 'C'}

	invert, reverse, inverse := getPermutations(test[:])
		if invert[0]=='A'&&invert[1]=='C'&&invert[2]=='T'&&invert[3]=='G' {

	}else{
		T.Error("Invert failed")
	}
	if reverse[0]=='C'&&reverse[1]=='A'&&reverse[2]=='G'&&reverse[3]=='T' {

	}else{
		T.Error("Invert failed")
	}
	if inverse[0]=='G'&&inverse[1]=='T'&&inverse[2]=='C'&&inverse[3]=='A' {

	}else{
		T.Error("Invert failed")
	}
	}
	func TestGetOneInSequence(T *testing.T){
		test := [11]rune{'A', 'T', 'G', 'C', 'T', 'T', 'T', 'G', 'A', 'A','A'}
		gen1 := make(chan []Gene)

		go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
		temp:= <- gen1
		if temp[0].Start!=0 {
			T.Error("incorrect Start")
		}
		if temp[0].End!=6 {
			T.Error("incorrect End")
		}
	}
func TestGetOneInSequenceStartOverEnd(T *testing.T){
	test := [11]rune{ 'C', 'T', 'T', 'T', 'G', 'A', 'A','A', 'A', 'T', 'G'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if temp[0].Start!=8 {
		T.Error("incorrect Start")
	}
	if temp[0].End!=3 {
		T.Error("incorrect End")
	}
}
func TestGetOneInSequenceLoopBack(T *testing.T){
	test := [11]rune{ 'G','C', 'T', 'T', 'T', 'G', 'A', 'A','A', 'A', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if temp[0].Start!=9 {
		T.Error("incorrect Start")
	}
	if temp[0].End!=4 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTGA(T *testing.T){
	test := [11]rune{ 'G', 'A', 'T', 'G', 'A', 'A','A', 'T', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if temp[0].Start!=1 {
		T.Error("incorrect Start")
	}
	if temp[0].End!=10 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTAG(T *testing.T){
	test := [11]rune{ 'A', 'G', 'T', 'G', 'A', 'A','A', 'A', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if temp[0].Start!=7 {
		T.Error("incorrect Start")
	}
	if temp[0].End!=10 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTAA(T *testing.T){
	test := [11]rune{ 'A', 'A', 'A', 'A', 'A', 'A','A', 'A', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if temp[0].Start!=7 {
		T.Error("incorrect Start")
	}
	if temp[0].End!=10 {
		T.Error("incorrect End")
	}
}

func TestNoGene(T *testing.T){
	test := [11]rune{ 'A', 'A', 'A', 'A', 'A', 'A','A', 'A', 'A', 'A', 'A'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if len(temp)!=0{
		T.Error("there should be no genes")
	}
}

func TestInfiniteGene(T *testing.T){
	test := [11]rune{ 'A', 'A', 'A', 'A', 'A', 'A','A', 'A', 'T', 'G', 'G'}
	gen1 := make(chan []Gene)
	go func() {
		defer func() {
			r := recover()
			if r!= nil {
				switch x:= r.(type) {
				case string:
					if x!="Never ending gene" {
						T.Error("Wrong Panic Message")
					}
					break
				default:
					T.Error("NonString Panic")
				}
				gen1<-nil
			}else{
				T.Error("No Panic")
			}

		}()
		count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{})

	}()
	<- gen1
}

func TestGeneCount(T *testing.T){
	test := [34]rune{ 'A', 'T', 'G', 'T', 'A', 'A',
	                  'A', 'T', 'G', 'T', 'A', 'A',
	                  'A', 'T', 'G', 'T', 'A', 'A',
	                  'A', 'T', 'G', 'T', 'A', 'A',
	                  'A', 'T', 'G', 'T', 'A', 'A',
	                  'A', 'A', 'A', 'A'}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if len(temp)!=5{
		T.Error("there should 5 genes")
	}
	for i := 0; i<5; i++ {
		if i+1!= temp[i].UUID{
			T.Error("Inccorect UUID on Element " + strconv.Itoa(i) + " of "+ strconv.Itoa(temp[i].UUID))
		}
	}
}

func TestGenesInPhase(T *testing.T){
	test := [34]rune{'A', 'T', 'G', 'T', 'A', 'A', 'A',
					 'A', 'T', 'G', 'T', 'A', 'A', 'A',
					 'A', 'T', 'G', 'T', 'A', 'A',
	}
	gen1 := make(chan []Gene)

	go count(test[:],gen1, &concurrentCounter{}, &concurrentCounter{})
	temp:= <- gen1
	if len(temp)!=3{
		T.Error("there should 5 genes")
	}
	for i := 0; i<3; i++ {
		if i+1!= temp[i].UUID{
			T.Error("Inccorect UUID on Element " + strconv.Itoa(i) + " of "+ strconv.Itoa(temp[i].UUID))
		}
	}
}