package service

import (
	"net/http"
)

func instantiateWebService() http.Handler {
	// create a new service
	fs := http.FileServer(http.Dir("/app/web"))
	http.Handle("/", fs)

	return fs
}
