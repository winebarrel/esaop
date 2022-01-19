package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/winebarrel/esaop"
)

func main() {
	cfg := parseArgs()
	r := esaop.NewRouter(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	err := http.ListenAndServe(addr, r)

	if err != nil {
		log.Fatal(err)
	}
}
