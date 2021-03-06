package restserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"blockchain/chain"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Setup() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/blocks", GetBlocksHandler)
	r.HandleFunc("/add", AddBlockHandler)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Printf("✔ REST Interface is UP on http://localhost:%s", port)
	http.ListenAndServe(":"+port, loggedRouter)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to MJS Blockchain. Check out https://mjs-blockchain.herokuapp.com/blocks for all block data. To add a block, run a curl command to /add. Ex: curl -X POST -d '{'Data':'hello world'}' http://localhost:8000/add")
}
func GetBlocksHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(chain.BC.Blocks)
}

type NewBlockData struct {
	Data string
}

func AddBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var blockData NewBlockData
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&blockData)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		chain.BC.AddBlock(blockData.Data)
		json.NewEncoder(w).Encode(chain.BC.Blocks[len(chain.BC.Blocks)-1])
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
