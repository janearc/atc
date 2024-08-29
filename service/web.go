package service

import (
	"net/http"
)

func instantiateWebService() http.Handler {
	// create a new service (remember: this is in docker, so the path is absolute in the container)
	fs := http.FileServer(http.Dir("/app/web"))
	http.Handle("/", fs)

	return fs
}
