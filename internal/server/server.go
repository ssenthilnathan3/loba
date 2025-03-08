package server

import (
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping)

}
