package loba

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ssenthilnathan3/loba/internal/config"
)

func BalanceLoads(cfg *config.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var selectedTarget *config.Target

		for _, resource := range cfg.Targets {
			capacity, err := strconv.Atoi(resource.Capacity)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid capacity: %v", err), http.StatusInternalServerError)
				return
			}

			consumptionRate, err := strconv.Atoi(resource.Consumption_rate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid consumption rate: %v", err), http.StatusInternalServerError)
				return
			}

			if consumptionRate < capacity {
				if selectedTarget == nil || selectedTarget.Consumption_rate < resource.Consumption_rate {
					selectedTarget = &resource
				}
			}
		}

		if selectedTarget == nil {
			http.Error(w, "No available target found", http.StatusServiceUnavailable)
			return
		}

		http.Redirect(w, r, selectedTarget.Host_url, http.StatusTemporaryRedirect)
	}
}

func ConsistentlyBalanceLoads(cfg *config.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var hashRing []uint32
		// hash host ips
		for _, target := range cfg.Targets {
			hashedURL, err := config.Hash(target.Host_url)
			if err != nil {
				fmt.Printf("Error parsing targetURL %s", err)
				return
			}

			hashRing = append(hashRing, hashedURL)
		}
		// maintain a hashring with those hashed ips
		if len(hashRing) == 0 {
			http.Error(w, "No available target found", http.StatusServiceUnavailable)
			return
		}

		// route all the incoming reqs to the first available ip
		for _, hashedURL := range hashRing {
			for _, target := range cfg.Targets {
				url, err := config.Hash(target.Host_url)
				if url == hashedURL && err == nil {
					http.Redirect(w, r, target.Host_url, http.StatusTemporaryRedirect)
					return
				}
			}
		}

		// if the ip is not available, throw error
		http.Error(w, "No available target found", http.StatusServiceUnavailable)
	}
}
