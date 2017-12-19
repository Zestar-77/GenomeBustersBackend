package webserver

import (
	"net/http"
)

// APIRequestHandler handles api server commands
func APIRequestHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/API/"):]
	switch url {
	case "scanForGenes":
		// TODO
		break
	}
}
