package webserver

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/GenomeBustersBackend/analyzer"
	"github.com/GenomeBustersBackend/specialFileReaders"
)

var isFasta *regexp.Regexp
var isGenBank *regexp.Regexp

func init() {
	if reg, err := regexp.Compile("*.\\.[fF][aA][sS][tT]"); err == nil {
		isFasta = reg
	} else {
		panic(err)
	}
	if reg, err := regexp.Compile("*.\\.[gG][bG]"); err == nil {
		isFasta = reg
	} else {
		panic(err)
	}
}

// GeneSearch handles an uploaded fasta file and handles searching.
func GeneSearch(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Could not read uploaded file. Is it a FASTA file?")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	if !(isFasta.MatchString(header.Filename)) {
		if !(isGenBank.MatchString(header.Filename)) {
			fmt.Printf("Could not read uploaded file. Is it a FASTA file?")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		genFile, err := specialFileReaders.NewGenebankFile(file)
		if err != nil {
			fmt.Printf("could not read uploaded file\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		analyzer.Thing([]rune(genFile.ReadGenome())) // Maybee convert all rune arrays to strings to prevent unneeded memory duplication
	} else {
		reader := bufio.NewReader(file)
		var genome []rune
		reader.ReadBytes('\n')
		for err == nil {
			var line string
			line, err = reader.ReadString('\n')
			if err != nil {
				genome = append(genome, []rune(strings.TrimSpace(line))...)
			}
		}
		if err != io.EOF {
			fmt.Printf("Could not read uploaded file. Is it a FASTA file?")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		analyzer.Thing(genome)
	}

	// TODO Upload response

}
