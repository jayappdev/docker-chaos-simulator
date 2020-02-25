package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/javdevapp/docker-chaos-simulator/cmdexec"
)

type ChaosKillHandler struct {
	// mu sync.Mutex // guards n
	// n  int
}

func (c *ChaosKillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// h.mu.Lock()
	// defer h.mu.Unlock()
	// h.n++
	// fmt.Fprintf(w, "count is %d\n", h.n)

	vars := mux.Vars(r)

	fmt.Println("Service ChaosKillHandler is called " + vars["containerid"])

	pause := cmdexec.CreateChaosKillSimulator(vars["containerid"], r.Body)
	err, _ := pause.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
}
