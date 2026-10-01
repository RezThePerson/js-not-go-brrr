package main

import (
	"log/slog"
	"os"
)

func main() {


	
	if err := api.Start(); err != nil {
		slog.Error("starting api server", "err", err)
		os.Exit(1)
	}
}
