package webserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"GenomeBustersBackend/analyzer"
	"GenomeBustersBackend/specialFileReaders"
)

var isFasta *regexp.Regexp
var isGenBank *regexp.Regexp

func init() {
	if reg, err := regexp.Compile(".*\\.[fF][aA][sS][tT]"); err == nil {
		isFasta = reg
	} else {
		panic(err)
	}
	if reg, err := regexp.Compile(".*\\.[gG][bG]"); err == nil {
		isGenBank = reg
	} else {
		panic(err)
	}
}

type dataWrapper struct {
	Sequence analyzer.Genome `json:"sequence"`
}

// GeneSearch handles an uploaded fasta file and handles searching.
func GeneSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Incoming file\n")
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Could not read uploaded file. Is it a FASTA file\n?")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("Filename is \"%s\"\n", header.Filename)
	var genes analyzer.Genome
	if !(isFasta.MatchString(header.Filename)) {
		if !(isGenBank.MatchString(header.Filename)) {
			fmt.Printf("Could not read uploaded file. Is it a FASTA file?\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("File is genbank\n")
		genFile, err := specialFileReaders.NewGenebankFile(file)
		if err != nil {
			fmt.Printf("could not read uploaded file\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Analyzing uploaded genbank\n")
		genes = analyzer.Thing([]rune(genFile.ReadGenome())) // Maybee convert all rune arrays to strings to prevent unneeded memory duplication
	} else {
		fmt.Printf("file is fasta")
		reader := bufio.NewReader(file)
		var genome []rune
		reader.ReadBytes('\n')
		for err == nil {
			var line string
			line, err = reader.ReadString('\n')
			if err == nil {
				genome = append(genome, []rune(strings.TrimSpace(line))...)
			}
		}
		if err != io.EOF {
			fmt.Printf("Could not read uploaded file. Is it a FASTA file?\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Analyzing uploaded fasfa\n")
		genes = analyzer.Thing(genome)
	}

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println("Returning file")
	genes.Filename = header.Filename
	template := dataWrapper{genes}

	jsonData, err := json.Marshal(template)
	if err != nil {
		fmt.Printf("Error in json marshelling! %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonData)
	fmt.Println("done")
}
