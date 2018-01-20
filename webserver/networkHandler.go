package webserver

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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
	if err != nil || isFasta.MatchString(header.Filename) {
		fmt.Printf("Could not read uploaded file. Is it a FASTA file?")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var buffer bytes.Buffer
	reader.ReadBytes('\n')
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			buffer.WriteString(strings.TrimSpace(line))
		}
	}
	if err != io.EOF {
		fmt.Printf("Could not read uploaded file. Is it a FASTA file?")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// analyzer.AnalyzeFastaData(buffer.Bytes())
}
