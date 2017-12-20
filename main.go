package main

import (
	"GenomeBustersBackend/webserver"
	"net/http"
)

// serverRoot is a constant string stating where the root of the filesystem is for serving files.
// TODO: Make this configurable. It should probably also default to something like '/opt/GenomeBusters/www/'
var serverRoot = "./GenomeBusters/polymorphs-frontend-master/build"

func main() {
	http.HandleFunc("/api/", webserver.APIRequestHandler)
	fileServer := http.FileServer(http.Dir(serverRoot))
	http.Handle("/", fileServer)
	http.ListenAndServe(":8080", nil)
}
