package analyzer

import (
	"GenomeBustersBackend/genedatabase"
	"GenomeBustersBackend/global"
	"GenomeBustersBackend/specialfilereaders"
	"bufio"
	"log"
	"os"
	"strconv"
	"testing"
)

// TestMain runs prior to the entire test suit, and is used for setup/tear down of the environment
func TestMain(m *testing.M) {
	file, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, 0666) // Ignore log messages by sending them to platform agnostic os.DevNull
	if err != nil {
		panic(err)
	}
	global.Log = log.New(file, "TESTING: ", log.Lmicroseconds)
	closer, err := genedatabase.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	file.Close()
	closer()
	os.Exit(exitCode)
}

func TestGetOneInSequenceStartOverEnd(T *testing.T) {
	test := [11]rune{'C', 'T', 'T', 'T', 'G', 'A', 'A', 'A', 'A', 'T', 'G'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 3)
	temp := <-gen1
	if temp[0].Start != 8 {
		T.Error("incorrect Start")
	}
	if temp[0].End != 3 {
		T.Error("incorrect End")
	}
}
func TestGetOneInSequenceLoopBack(T *testing.T) {
	test := [11]rune{'G', 'C', 'T', 'T', 'T', 'G', 'A', 'A', 'A', 'A', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 3)
	temp := <-gen1
	if temp[0].Start != 9 {
		T.Error("incorrect Start")
	}
	if temp[0].End != 4 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTGA(T *testing.T) {
	test := [11]rune{'G', 'A', 'T', 'G', 'A', 'A', 'A', 'T', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 3)
	temp := <-gen1
	if temp[0].Start != 1 {
		T.Error("incorrect Start")
	}
	if temp[0].End != 10 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTAG(T *testing.T) {
	test := [11]rune{'A', 'G', 'T', 'G', 'A', 'A', 'A', 'A', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 0)
	temp := <-gen1
	if temp[0].Start != 7 {
		T.Error("incorrect Start")
	}
	if temp[0].End != 10 {
		T.Error("incorrect End")
	}
}

func TestGetOneInSequenceEndOverEndTAA(T *testing.T) {
	test := [11]rune{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'T', 'G', 'T'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 0)
	temp := <-gen1
	if temp[0].Start != 7 {
		T.Error("incorrect Start")
	}
	if temp[0].End != 10 {
		T.Error("incorrect End")
	}
}

func TestNoGene(T *testing.T) {
	test := [11]rune{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 3)
	temp := <-gen1
	if len(temp) != 0 {
		T.Error("there should be no genes")
	}
}

func TestInfiniteGene(T *testing.T) {
	test := [11]rune{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A', 'T', 'G', 'G'}
	gen1 := make(chan []Gene)
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				switch x := r.(type) {
				case string:
					if x != "Never ending gene" {
						T.Error("Wrong Panic Message")
					}
					break
				default:
					T.Error("NonString Panic")
				}
				gen1 <- nil
			} else {
				T.Error("No Panic")
			}

		}()
		count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 3)

	}()
	<-gen1
}

func TestGeneCount(T *testing.T) {
	test := [34]rune{'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'A', 'A', 'A'}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 0)
	temp := <-gen1
	if len(temp) != 5 {
		T.Error("there should 5 genes")
	}
	for i := 0; i < 5; i++ {
		if i+1 != temp[i].UUID {
			T.Error("Inccorect UUID on Element " + strconv.Itoa(i) + " of " + strconv.Itoa(temp[i].UUID))
		}
	}
}
func TestGeneCountWithMinLength(T *testing.T) {
	test := [37]rune{'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
		'A', 'A', 'A', 'A'}
	gen1 := make(chan []Gene)
	cc1 := &concurrentCounter{}
	cc2 := &concurrentCounter{}
	go count(test[:], gen1, cc1, cc2, 5)
	temp := <-gen1
	if len(temp) != 1 {
		T.Error("there should 1 gene" + strconv.Itoa(len(temp)))
	}

}

func TestGenesInPhase(T *testing.T) {
	minLength = 0

	test := [34]rune{'A', 'T', 'G', 'T', 'A', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A', 'A',
		'A', 'T', 'G', 'T', 'A', 'A',
	}
	gen1 := make(chan []Gene)

	go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 0)
	temp := <-gen1
	if len(temp) != 3 {
		T.Error("there should 3 genes")
		go count(test[:], gen1, &concurrentCounter{}, &concurrentCounter{}, 0)
		temp := <-gen1
		if len(temp) != 3 {
			T.Error("there should 3 genes")
		}
		for i := 0; i < 3; i++ {
			if i+1 != temp[i].UUID {
				T.Error("Inccorect UUID on Element " + strconv.Itoa(i) + " of " + strconv.Itoa(temp[i].UUID))
			}
		}
	}
}
func BenchmarkThing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		file, err := os.Open("../Bs-916.gb")
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(file)
		testFile, err := specialfilereaders.NewGenebankFile(reader)
		testGenome := []rune(testFile.ReadGenome())

		this := Analyze(testGenome)
		this.GenesFound++
		this.GenesFound--
		j := 0
		j++

	}
}

func BenchmarkCount(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file, err := os.Open("../sequence.gb")
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(file)
		testFile, err := specialfilereaders.NewGenebankFile(reader)
		testGenome := []rune(testFile.ReadGenome())
		gen1 := make(chan []Gene)

		go count(testGenome, gen1, &concurrentCounter{}, &concurrentCounter{}, 3)
		<-gen1

	}
}
func BenchmarkGenBankIn(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file, err := os.Open("../sequence.gb")
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(file)
		testFile, err := specialfilereaders.NewGenebankFile(reader)
		testGenome := []rune(testFile.ReadGenome())
		testGenome[0] = 'A'

	}
}
