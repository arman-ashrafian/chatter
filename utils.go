package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func shutdownServer() {
	fmt.Println("Server shutdown")
	os.Exit(0)
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
