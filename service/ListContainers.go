package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/javdevapp/docker-chaos-simulator/dockerinfo"
)

func ListContainers(w http.ResponseWriter, r *http.Request) { // h.mu.Lock()
	// defer h.mu.Unlock()
	// h.n++
	// fmt.Fprintf(w, "count is %d\n", h.n)

	containers, err := dockerinfo.ListContainer(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(containers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
