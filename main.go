package main

import (
	cnf "GenomeBustersBackend/configurationHandler"
	"GenomeBustersBackend/genedatabase"
	"GenomeBustersBackend/interactive"
	"GenomeBustersBackend/webserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	v := cnf.GetConfig()
	log.Println("Starting Busted")
	port := ":" + strconv.Itoa(v.GetInt("port"))
	apiport := ":" + strconv.Itoa(v.GetInt("apiPort"))

	fileServer := http.FileServer(http.Dir(v.GetString("serverRoot") + "/build"))
	// fh := http.Handle("/", fileServer)
	// http.HandleFunc("/api/gene_search/", webserver.GeneSearch)

	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	closeDB, err := geneDatabase.InitializeDatabase()
	if err != nil {
		log.Printf("Unable to open gene database: %v\nAll genes will be marked unat\n", err)
	} else {
		defer closeDB()
	}

	server := &http.Server{Addr: port, Handler: fileServer}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	apiServer := &http.Server{Addr: apiport, Handler: http.HandlerFunc(webserver.GeneSearch)}
	go func() {
		if err := apiServer.ListenAndServe(); err != nil {
			log.Printf("Error: %s\n", err)
			keyboardInterrupt <- nil
		}
	}()

	log.Printf("Server running on port %d, with api on port %d\n", v.GetInt("port"), v.GetInt("apiPort"))

	if v.GetBool("LogToConsole") {
		if err := interactive.RunTui(keyboardInterrupt); err != nil {
			<-keyboardInterrupt
		}
	}

	log.Printf("\nShutting Down Server...\n")
	if err := server.Shutdown(nil); err != nil {
		log.Printf("Error: %s", err)
	}
	if err := apiServer.Shutdown(nil); err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Goodbye!\n")
}
