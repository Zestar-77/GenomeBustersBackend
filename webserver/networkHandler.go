package webserver

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"GenomeBustersBackend/analyzer"
	"GenomeBustersBackend/specialFileReaders"
)

var (
	isFasta   *regexp.Regexp
	isGenBank *regexp.Regexp
)

// init compiles necessary gene files
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
// Returns the request with a json formatted list of genes
func GeneSearch(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming file\n")
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Could not read uploaded file. Is it a FASTA file\n?")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Printf("Filename is \"%s\"\n", header.Filename)
	var genes analyzer.Genome
	if !(isFasta.MatchString(header.Filename)) {
		if !(isGenBank.MatchString(header.Filename)) {
			log.Printf("Could not read uploaded file. Is it a FASTA file?\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("File is genbank\n")
		genFile, err := specialFileReaders.NewGenebankFile(file)
		if err != nil {
			log.Printf("could not read uploaded file\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Analyzing uploaded genbank\n")
		genes = analyzer.Analyze([]rune(genFile.ReadGenome())) // Maybee convert all rune arrays to strings to prevent unneeded memory duplication
	} else {
		log.Printf("file is fasta")
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
			log.Printf("Could not read uploaded file. Is it a FASTA file?\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Analyzing uploaded fasfa\n")
		genes = analyzer.Analyze(genome)
	}

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	log.Println("Returning file")
	genes.Filename = header.Filename
	template := dataWrapper{genes}

	jsonData, err := json.Marshal(template)
	if err != nil {
		log.Printf("Error in json marshelling! %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonData)
	log.Println("done")
}
