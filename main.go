package main

import (
	cnf "GenomeBustersBackend/configurationHandler"
	"GenomeBustersBackend/webserver"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	v, err := cnf.GetConfig()
	if err != nil {
		fmt.Printf("unable to parse config file %s", err)
		return
	}
	fmt.Println("Starting Busted")
	port := ":" + strconv.Itoa(v.GetInt("port"))
	apiport := ":" + strconv.Itoa(v.GetInt("apiPort"))

	fileServer := http.FileServer(http.Dir(v.GetString("serverRoot") + "/build"))
	// fh := http.Handle("/", fileServer)
	// http.HandleFunc("/api/gene_search/", webserver.GeneSearch)

	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	server := &http.Server{Addr: port, Handler: fileServer}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	apiServer := &http.Server{Addr: apiport, Handler: http.HandlerFunc(webserver.GeneSearch)}
	go func() {
		if err := apiServer.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	fmt.Printf("Server running on port %d, with api on port %d\n", v.GetInt("port"), v.GetInt("apiPort"))
	<-keyboardInterrupt
	fmt.Printf("\nShutting Down Server...\n")
	if err := server.Shutdown(nil); err != nil {
		fmt.Printf("Error: %s", err)
	}
	if err := apiServer.Shutdown(nil); err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Goodbye!\n")
}
