package main

import (
	"log"
	"log/slog"

	"github.com/hrvadl/coursework_db/pkg/server"
)

const port = ":8080"

func main() {
	srv := server.NewHTTP()

	slog.Info("Server starting...")
	log.Fatal(srv.ListenAndServe(port))
}
