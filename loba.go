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
		// TODO: implement consistent load balancing
	}
}
