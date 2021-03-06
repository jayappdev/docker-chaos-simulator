package service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/javdevapp/docker-chaos-simulator/cmdexec"
)

type ChaosNetworkHandler struct {
	// mu sync.Mutex // guards n
	// n  int
}

func (c *ChaosNetworkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// h.mu.Lock()
	// defer h.mu.Unlock()
	// h.n++
	// fmt.Fprintf(w, "count is %d\n", h.n)

	vars := mux.Vars(r)

	fmt.Println("Service ChaosNetworkEmulatorSimulator is called " + vars["containerid"])

	pause := cmdexec.CreateChaosNetworkEmulatorSimulator(vars["containerid"], r.Body)
	err, _, errorReader := pause.Execute()
	if err != nil {
		errorLines, _ := ioutil.ReadAll(errorReader)
		w.Write([]byte("Error while applying Network Chaos " + "\n" + string(errorLines)))

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
}
