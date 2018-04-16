package main

import (
	cnf "GenomeBustersBackend/configurationHandler"
	"GenomeBustersBackend/genedatabase"
	"GenomeBustersBackend/global"
	"GenomeBustersBackend/interactive"
	"GenomeBustersBackend/webserver"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	v := cnf.GetConfig()
	global.Log.Println("Starting Busted")
	port := ":" + strconv.Itoa(v.GetInt("port"))
	apiport := ":" + strconv.Itoa(v.GetInt("apiPort"))

	fileServer := http.FileServer(http.Dir(v.GetString("serverRoot") + "/build"))
	// fh := http.Handle("/", fileServer)
	// http.HandleFunc("/api/gene_search/", webserver.GeneSearch)

	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	closeDB, err := genedatabase.InitializeDatabase()
	if err != nil {
		global.Log.Printf("Unable to open gene database: %v\nAll genes will be marked unat\n", err)
	} else {
		defer closeDB()
	}

	server := &http.Server{Addr: port, Handler: fileServer}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			global.Log.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	apiServer := &http.Server{Addr: apiport, Handler: http.HandlerFunc(webserver.GeneSearch)}
	go func() {
		if err := apiServer.ListenAndServe(); err != nil {
			global.Log.Printf("Error: %s\n", err)
			keyboardInterrupt <- nil
		}
	}()

	global.Log.Printf("Server running on port %d, with api on port %d\n", v.GetInt("port"), v.GetInt("apiPort"))

	if !v.GetBool("LogToConsole") {
		if err := interactive.RunTui(keyboardInterrupt); err != nil {
			<-keyboardInterrupt
		}
	} else {
		<-keyboardInterrupt
	}

	global.Log.Printf("Shutting Down Server...\n")
	if err := server.Shutdown(nil); err != nil {
		global.Log.Printf("%s", err)
	}
	if err := apiServer.Shutdown(nil); err != nil {
		global.Log.Printf("%s", err)
	}
	global.Log.Printf("Goodbye!\n")
}
