package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nekoteoj/lab-cms/internal/pkg/config"
)

func main() {
	cfg := config.Load()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Lab CMS")
	})

	log.Printf("Server starting on port %s [%s]", cfg.Port, cfg.Env)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
