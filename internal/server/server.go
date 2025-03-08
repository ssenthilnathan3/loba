package server

import (
	"fmt"
	"net/http"

	"github.com/ssenthilnathan3/loba"
	"github.com/ssenthilnathan3/loba/internal/config"
)

func Run() error {
	config, err := config.NewConfiguration()
	if err != nil {
		return fmt.Errorf("error creating configuration %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", loba.BalanceLoads(config))
	mux.HandleFunc("/ping", ping)

	if err := http.ListenAndServe("localhost:8000", mux); err != nil {
		return fmt.Errorf("could not start the HTTP server %s", err)
	}

	return nil
}
